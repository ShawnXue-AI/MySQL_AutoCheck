<template>
  <div class="benchmarks-view">
    <div class="bm-toolbar">
      <div class="bm-search-wrap">
        <svg viewBox="0 0 24 24" class="bm-search-icon"><path d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z" fill="currentColor"/></svg>
        <input type="text" class="bm-search-input" v-model="searchText" placeholder="grep 巡检项名称或说明..." />
      </div>
      <div class="bm-tag-group">
        <button class="bm-tag" :class="{ active: activeCategory === '' }" @click="activeCategory = ''; page = 1">*</button>
        <button class="bm-tag" v-for="cat in categories" :key="cat" :class="{ active: activeCategory === cat }" @click="activeCategory = cat; page = 1">{{ cat }}</button>
      </div>
      <span class="bm-count">
        <span class="terminal-prompt">$</span> {{ filteredItems.length }} 结果
      </span>
    </div>

    <div v-if="pagedItems.length === 0" class="empty-state">
      <p class="empty-text">// 没有匹配的巡检项</p>
    </div>

    <div v-else class="bm-table-wrap">
      <table class="bm-table">
        <thead>
          <tr>
            <th class="col-name">check_name</th>
            <th class="col-cat">category</th>
            <th class="col-baseline">baseline</th>
            <th class="col-ok">OK</th>
            <th class="col-warn">WARN</th>
            <th class="col-err">ERR</th>
            <th class="col-desc">description</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(row, i) in pagedItems" :key="i">
            <td class="col-name"><code>{{ row.name }}</code></td>
            <td class="col-cat"><span class="cat-chip" :style="{ borderColor: catColors[row.category] || 'var(--border-hover)', color: catColors[row.category] || 'var(--text-tertiary)' }">{{ row.category }}</span></td>
            <td class="col-baseline">{{ row.baseline }}</td>
            <td class="col-ok">{{ row.ok }}</td>
            <td class="col-warn">{{ row.warn }}</td>
            <td class="col-err">{{ row.error }}</td>
            <td class="col-desc">{{ row.desc }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="bm-pager" v-if="totalPages > 1">
      <button class="bm-page-btn" :disabled="page <= 1" @click="page--">prev</button>
      <span class="bm-page-info">Page {{ page }} / {{ totalPages }} · {{ filteredItems.length }} items</span>
      <button class="bm-page-btn" :disabled="page >= totalPages" @click="page++">next</button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'BenchmarksView',
  data() {
    return {
      searchText: '',
      activeCategory: '',
      page: 1,
      pageSize: 25,
      benchmarks: [
        { name: 'DATABASE_VERSION', category: '数据库信息', baseline: '-', ok: 'MySQL 8.0+ / MySQL 5.7', warn: 'MySQL 5.6', error: 'MySQL 5.5 及以下', desc: '数据库版本检查' },
        { name: 'UPTIME', category: '数据库信息', baseline: '运行时间', ok: '> 864000 秒 (10天)', warn: '3600 ~ 864000 秒', error: '< 3600 秒 (1小时)', desc: '检查是否刚重启' },
        { name: 'MAX_CONNECTIONS', category: '数据库信息', baseline: '最大连接数', ok: '>= 200', warn: '100 ~ 199', error: '< 100', desc: '最大连接数设置' },
        { name: 'INNODB_BUFFER_POOL_SIZE', category: '数据库信息', baseline: '缓冲池 / 总内存', ok: '>= 50%', warn: '20% ~ 50%', error: '< 20%', desc: 'InnoDB 缓冲池大小，推荐为内存的 70%' },
        { name: 'OS_VERSION', category: '数据库信息', baseline: '操作系统版本', ok: 'CentOS 7+ / Ubuntu 20+', warn: 'CentOS 6 / Ubuntu 18', error: 'CentOS 5 / Ubuntu 16 以下', desc: '操作系统版本检查' },
        { name: 'INNODB_LOG_FILE_SIZE', category: '数据库信息', baseline: '日志文件大小', ok: '>= 2 GB', warn: '1 ~ 2 GB', error: '< 1 GB', desc: 'Redo log 大小' },
        { name: 'INNODB_LOG_FILES_IN_GROUP', category: '数据库信息', baseline: '日志文件数量', ok: '>= 2', warn: 'N/A', error: '< 2', desc: '日志组成员数' },
        { name: 'INNODB_FLUSH_LOG_AT_TRX_COMMIT', category: '数据库信息', baseline: '日志刷新策略', ok: '1', warn: 'N/A', error: '!= 1', desc: '推荐值为 1，保证 ACID' },
        { name: 'SYNC_BINLOG', category: '数据库信息', baseline: 'binlog 同步策略', ok: '1', warn: 'N/A', error: '!= 1', desc: '推荐值为 1' },
        { name: 'BINLOG_FORMAT', category: '数据库信息', baseline: '二进制日志格式', ok: 'ROW', warn: 'MIXED', error: 'STATEMENT', desc: '推荐使用 ROW 格式' },
        { name: 'TRANSACTION_ISOLATION', category: '数据库信息', baseline: '事务隔离级别', ok: 'READ-COMMITTED', warn: 'REPEATABLE-READ', error: 'SERIALIZABLE / READ-UNCOMMITTED', desc: '推荐 RC' },
        { name: 'INNODB_FLUSH_METHOD', category: '数据库信息', baseline: 'InnoDB 刷新方法', ok: 'O_DIRECT', warn: 'N/A', error: 'N/A', desc: '推荐使用 O_DIRECT' },
        { name: 'INNODB_FILE_PER_TABLE', category: '数据库信息', baseline: '独立表空间', ok: 'ON', warn: 'N/A', error: 'OFF', desc: '建议开启' },
        { name: 'OPEN_FILES_LIMIT', category: '数据库信息', baseline: '打开文件限制', ok: '>= 65535', warn: '20000 ~ 65535', error: '< 20000', desc: '文件句柄数' },
        { name: 'TABLE_OPEN_CACHE', category: '数据库信息', baseline: '表打开缓存', ok: '>= 4000', warn: '2000 ~ 3999', error: '< 2000', desc: '表缓存大小' },
        { name: 'MAX_ALLOWED_PACKET', category: '数据库信息', baseline: '最大允许包大小', ok: '>= 64 MB', warn: '16 ~ 64 MB', error: '< 16 MB', desc: '最大数据包' },
        { name: 'WAIT_TIMEOUT', category: '数据库信息', baseline: '等待超时', ok: '28800 (8小时)', warn: '600 ~ 28800', error: '< 600', desc: '非交互超时' },
        { name: 'INTERACTIVE_TIMEOUT', category: '数据库信息', baseline: '交互超时', ok: '28800 (8小时)', warn: '600 ~ 28800', error: '< 600', desc: '交互连接超时' },
        { name: 'CURRENT_CONNECTIONS', category: '性能指标', baseline: '连接数 / 最大连接数', ok: '< 80%', warn: '80% ~ 90%', error: '> 90%', desc: '连接数使用率' },
        { name: 'SLOW_QUERIES', category: '性能指标', baseline: '慢查询数', ok: '< 10 / 天', warn: '10 ~ 100 / 天', error: '> 100 / 天', desc: '慢查询数量' },
        { name: 'TABLE_LOCKS_WAITED', category: '性能指标', baseline: '表锁等待数', ok: '0', warn: '1 ~ 100', error: '> 100', desc: '表锁等待次数' },
        { name: 'QUERY_CACHE_HITS', category: '性能指标', baseline: '查询缓存命中率', ok: '> 80%', warn: '50% ~ 80%', error: '< 50%', desc: 'MySQL 5.7 已废弃' },
        { name: 'INNODB_ROW_LOCK_WAITS', category: '性能指标', baseline: '行锁等待次数', ok: '0', warn: '1 ~ 100', error: '> 100', desc: '行锁等待' },
        { name: 'INNODB_ROW_LOCK_TIME_AVG', category: '性能指标', baseline: '平均行锁等待时间', ok: '< 100 ms', warn: '100 ~ 1000 ms', error: '> 1000 ms', desc: '行锁平均等待' },
        { name: 'CREATED_TMP_DISK_TABLES', category: '性能指标', baseline: '磁盘临时表创建数', ok: '0', warn: '< 1000 / 天', error: '> 1000 / 天', desc: '磁盘临时表过多需优化' },
        { name: 'ABORTED_CLIENTS', category: '性能指标', baseline: '中断客户端数', ok: '0', warn: '< 10 / 天', error: '> 10 / 天', desc: '客户端断连' },
        { name: 'ABORTED_CONNECTS', category: '性能指标', baseline: '中断连接数', ok: '0', warn: '< 10 / 天', error: '> 10 / 天', desc: '连接失败次数' },
        { name: 'EMPTY_PASSWORD_USERS', category: '安全配置', baseline: '空密码用户', ok: '0', warn: 'N/A', error: '> 0', desc: '存在空密码用户即异常' },
        { name: 'REMOTE_ROOT_USERS', category: '安全配置', baseline: '远程 root 用户', ok: '0', warn: 'N/A', error: '> 0', desc: '不允许远程 root' },
        { name: 'SSL_STATUS', category: '安全配置', baseline: 'SSL 状态', ok: '已启用', warn: 'N/A', error: '未启用', desc: '生产环境建议开启 SSL' },
        { name: 'REMOTE_USERS', category: '安全配置', baseline: '远程用户数', ok: '< 5', warn: '5 ~ 10', error: '> 10', desc: '远程访问用户数量' },
        { name: 'DATABASE_SIZES', category: '存储信息', baseline: '数据库总大小', ok: '< 500 GB', warn: '500 GB ~ 1 TB', error: '> 1 TB', desc: '数据库容量' },
        { name: 'TABLE_FRAGMENTATION_DETAILED', category: '存储信息', baseline: '表碎片率', ok: '碎片 < 20% (碎片 > 2GB)', warn: '碎片 20% ~ 50% (碎片 > 2GB)', error: '碎片 > 50% (碎片 > 2GB)', desc: '碎片表数量' },
        { name: 'TABLES_WITHOUT_PRIMARY_KEY', category: '存储信息', baseline: '无主键表', ok: '0', warn: '1 ~ 10', error: '> 10', desc: '无主键表数量' },
        { name: 'PARTITIONED_TABLES', category: '存储信息', baseline: '分区表状态', ok: '分区策略合理', warn: '需补充新分区', error: '无分区且数据量 > 50GB', desc: '分区表检查' },
        { name: 'SYSTEM_CPU_USAGE', category: '系统资源', baseline: 'CPU 使用率', ok: '< 50%', warn: '50% ~ 80%', error: '> 80%', desc: 'CPU 利用率' },
        { name: 'SYSTEM_MEMORY_USAGE', category: '系统资源', baseline: '内存使用率', ok: '< 60%', warn: '60% ~ 80%', error: '> 80%', desc: '内存利用率' },
        { name: 'SYSTEM_SWAP_USAGE', category: '系统资源', baseline: 'SWAP 使用率', ok: '< 20%', warn: '20% ~ 50%', error: '> 50%', desc: 'SWAP 使用率' },
        { name: 'SYSTEM_DISK_USAGE', category: '系统资源', baseline: '磁盘使用率', ok: '< 70%', warn: '70% ~ 85%', error: '> 85%', desc: '磁盘空间利用率' },
        { name: 'MYSQL_PROCESS', category: '系统资源', baseline: 'MySQL 进程数', ok: '< 50', warn: '50 ~ 100', error: '> 100', desc: 'MySQL 线程数' },
        { name: 'FILE_DESCRIPTOR_LIMIT', category: '系统资源', baseline: '文件描述符限制', ok: '>= 65535', warn: '20000 ~ 65535', error: '< 20000', desc: '系统文件描述符' },
        { name: 'NETWORK_CONNECTIONS', category: '系统资源', baseline: '网络连接数', ok: '< 200', warn: '200 ~ 500', error: '> 500', desc: 'TCP 连接数' },
        { name: 'MYSQL_ERROR_LOG', category: '系统资源', baseline: '错误日志', ok: '无错误', warn: '少量警告', error: '有 ERROR 级别日志', desc: '错误日志检查' },
        { name: 'MASTER_STATUS', category: '主从复制', baseline: '主库状态', ok: '运行正常', warn: 'N/A', error: '异常', desc: '主库状态' },
        { name: 'SLAVE_STATUS', category: '主从复制', baseline: '从库状态', ok: 'IO 和 SQL 线程均为 Yes', warn: 'Seconds_Behind_Master > 30', error: '任一线程非 YES', desc: '从库同步状态' },
        { name: 'BINARY_LOGS', category: '主从复制', baseline: '二进制日志', ok: '已开启且保留 > 7 天', warn: '已开启但保留 < 7 天', error: '未开启', desc: 'binlog 配置' },
        { name: 'GTID_MODE', category: '主从复制', baseline: 'GTID 模式', ok: 'ON', warn: 'N/A', error: 'OFF', desc: '推荐开启 GTID' },
        { name: 'GTID_CONSISTENCY', category: '主从复制', baseline: 'GTID 一致性', ok: 'ON', warn: 'N/A', error: 'OFF', desc: 'GTID 一致性检查' },
        { name: 'MGR_GROUP_MEMBERS', category: '集群信息', baseline: 'MGR 组成员', ok: '所有成员 ONLINE', warn: '部分成员 ONLINE', error: '有成员 OFFLINE', desc: 'MGR 集群状态' },
        { name: 'MGR_MEMBER_STATS', category: '集群信息', baseline: 'MGR 成员统计', ok: '队列空', warn: '队列较小', error: '队列较大', desc: 'MGR 事务队列' },
        { name: 'PXC_WSREP_STATUS', category: '集群信息', baseline: 'PXC 集群状态', ok: 'Primary 且所有节点同步', warn: '非 Primary', error: '节点数 < 预期', desc: 'PXC 状态检查' }
      ],
      catColors: {
        '数据库信息': '#58a6ff',
        '性能指标': '#3fb950',
        '安全配置': '#f0883e',
        '存储信息': '#bc8cff',
        '系统资源': '#e5534b',
        '主从复制': '#539bf5',
        '集群信息': '#db61a2'
      }
    }
  },
  computed: {
    categories() {
      const cats = new Set(this.benchmarks.map(b => b.category))
      return [...cats]
    },
    filteredItems() {
      let items = this.benchmarks
      if (this.activeCategory) {
        items = items.filter(b => b.category === this.activeCategory)
      }
      if (this.searchText.trim()) {
        const kw = this.searchText.trim().toLowerCase()
        items = items.filter(b =>
          b.name.toLowerCase().includes(kw) ||
          b.desc.toLowerCase().includes(kw) ||
          b.baseline.toLowerCase().includes(kw) ||
          b.ok.toLowerCase().includes(kw) ||
          b.warn.toLowerCase().includes(kw) ||
          b.error.toLowerCase().includes(kw)
        )
      }
      return items
    },
    totalPages() {
      return Math.max(1, Math.ceil(this.filteredItems.length / this.pageSize))
    },
    pagedItems() {
      const start = (this.page - 1) * this.pageSize
      return this.filteredItems.slice(start, start + this.pageSize)
    }
  },
  watch: {
    searchText() { this.page = 1 },
    activeCategory() { this.page = 1 }
  }
}
</script>

<style scoped>
.benchmarks-view { max-width: 1300px; margin: 0 auto; }

.bm-toolbar {
  display: flex; align-items: center; gap: 14px; margin-bottom: 16px; flex-wrap: wrap;
}

.bm-search-wrap {
  position: relative; flex: 0 0 300px;
}
.bm-search-icon {
  position: absolute; left: 10px; top: 50%; transform: translateY(-50%);
  width: 16px; height: 16px; color: var(--text-muted);
}
.bm-search-input {
  width: 100%; padding: 8px 12px 8px 34px;
  background: var(--bg-input); border: 1px solid var(--border-hover);
  border-radius: var(--radius-sm); color: var(--text-secondary);
  font-size: 13px; font-family: var(--font-code); outline: none;
  transition: border-color var(--transition-normal);
}
.bm-search-input:focus { border-color: var(--accent-blue); box-shadow: var(--shadow-glow-blue); }
.bm-search-input::placeholder { color: var(--text-muted); font-family: var(--font-code); font-size: 11.5px; }

.bm-tag-group { display: flex; gap: 5px; flex-wrap: wrap; }
.bm-tag {
  padding: 4px 12px; border: 1px solid var(--border-hover); border-radius: 20px;
  background: transparent; color: var(--text-tertiary);
  font-size: 12px; font-family: var(--font-code); cursor: pointer;
  transition: all var(--transition-fast); white-space: nowrap;
}
.bm-tag:hover { border-color: var(--text-muted); color: var(--text-secondary); }
.bm-tag.active { background: var(--accent-blue-subtle); border-color: var(--accent-blue); color: var(--accent-blue); }

.bm-count {
  font-size: 12px; color: var(--text-muted); margin-left: auto; white-space: nowrap;
  font-family: var(--font-code); display: flex; align-items: center; gap: 5px;
}
.terminal-prompt { color: var(--accent-green); font-weight: 500; }

.bm-table-wrap {
  overflow-x: auto; border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
}

.bm-table {
  width: 100%; border-collapse: collapse; font-size: 13px;
}
.bm-table thead {
  background: var(--bg-primary); position: sticky; top: 0; z-index: 1;
}
.bm-table th {
  padding: 10px 12px; text-align: left; font-weight: 600; color: var(--text-tertiary);
  font-family: var(--font-code); font-size: 10.5px; text-transform: uppercase;
  letter-spacing: 0.4px; border-bottom: 1px solid var(--border-default); white-space: nowrap;
}
.bm-table td {
  padding: 10px 12px; border-bottom: 1px solid var(--border-default);
  color: var(--text-secondary); vertical-align: middle;
}
.bm-table tbody tr { background: var(--bg-secondary); transition: background 0.1s; }
.bm-table tbody tr:hover { background: var(--bg-tertiary); }
.bm-table tbody tr:last-child td { border-bottom: none; }

.col-name { width: 230px; }
.col-name code {
  font-family: var(--font-code); font-size: 11.5px;
  color: var(--accent-blue); background: var(--accent-blue-subtle);
  padding: 2px 6px; border-radius: 4px;
}
.col-cat { width: 86px; white-space: nowrap; }
.col-baseline { min-width: 130px; font-size: 12.5px; }
.col-ok { min-width: 140px; font-size: 12.5px; color: var(--accent-green); }
.col-warn { min-width: 110px; font-size: 12.5px; color: var(--accent-amber); }
.col-err { min-width: 120px; font-size: 12.5px; color: var(--accent-red); }
.col-desc { min-width: 160px; font-size: 12.5px; color: var(--text-tertiary); line-height: 1.4; }

.cat-chip {
  display: inline-block; padding: 1px 8px; border: 1px solid; border-radius: 10px;
  font-size: 11px; font-weight: 500; font-family: var(--font-code);
}

.bm-pager {
  display: flex; align-items: center; justify-content: center; gap: 14px;
  margin-top: 16px; padding: 8px 0;
}
.bm-page-btn {
  padding: 6px 18px; border: 1px solid var(--border-hover); border-radius: var(--radius-sm);
  background: transparent; color: var(--text-secondary);
  font-family: var(--font-code); font-size: 12px; cursor: pointer;
  transition: all var(--transition-fast);
}
.bm-page-btn:hover:not(:disabled) { border-color: var(--accent-blue); color: var(--accent-blue); background: var(--accent-blue-subtle); }
.bm-page-btn:disabled { opacity: 0.3; cursor: not-allowed; }
.bm-page-info { font-size: 12px; color: var(--text-tertiary); font-family: var(--font-code); }
</style>
