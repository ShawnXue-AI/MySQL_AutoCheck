#!/bin/bash

# MySQL数据库自动化巡检脚本
# 生成标准格式的巡检日志，便于后续Python分析

# 配置参数
MYSQL_HOST="localhost"
MYSQL_PORT="3306"
MYSQL_USER="root"
MYSQL_PASS=""
MYSQL_SOCKET=""
OUTPUT_DIR="./inspection_logs"

# 功能开关
RDS_CHECK=false
MGR_CHECK=false
PXC_CHECK=false

# 获取当前实例的IP地址
get_instance_ip() {
    # 方法1: 通过 hostname -I 获取
    local ip=$(hostname -I 2>/dev/null | awk '{print $1}')
    if [ -n "$ip" ] && [ "$ip" != "127.0.0.1" ]; then
        echo "$ip"
        return
    fi
    
    # 方法2: 通过 ip addr 获取
    ip=$(ip addr show 2>/dev/null | grep 'inet ' | grep -v '127.0.0.1' | awk '{print $2}' | cut -d/ -f1 | head -1)
    if [ -n "$ip" ]; then
        echo "$ip"
        return
    fi
    
    # 方法3: 通过 ifconfig 获取
    ip=$(ifconfig 2>/dev/null | grep 'inet ' | grep -v '127.0.0.1' | awk '{print $2}' | head -1)
    if [ -n "$ip" ]; then
        echo "$ip"
        return
    fi
    
    # 如果都获取不到，使用默认值
    echo "127.0.0.1"
}

# 获取IP地址
INSTANCE_IP=$(get_instance_ip)

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数（注意：LOG_PATH在参数解析后定义）
log_info() {
    if [ -n "$LOG_PATH" ]; then
        echo -e "${BLUE}[INFO]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_PATH"
    else
        echo -e "${BLUE}[INFO]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
    fi
}

log_warning() {
    if [ -n "$LOG_PATH" ]; then
        echo -e "${YELLOW}[WARNING]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_PATH"
    else
        echo -e "${YELLOW}[WARNING]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
    fi
}

log_error() {
    if [ -n "$LOG_PATH" ]; then
        echo -e "${RED}[ERROR]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_PATH"
    else
        echo -e "${RED}[ERROR]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
    fi
}

log_success() {
    if [ -n "$LOG_PATH" ]; then
        echo -e "${GREEN}[SUCCESS]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$LOG_PATH"
    else
        echo -e "${GREEN}[SUCCESS]${NC} $(date '+%Y-%m-%d %H:%M:%S') - $1"
    fi
}

# 检查MySQL连接
check_mysql_connection() {
    log_info "检查MySQL连接..."
    
    if command -v mysql &> /dev/null; then
        # 使用数组直接构建命令，更安全
        local cmd=("mysql")
        if [ -n "$MYSQL_SOCKET" ]; then
            cmd+=("-S" "$MYSQL_SOCKET")
        else
            cmd+=("-h" "$MYSQL_HOST" "-P" "$MYSQL_PORT")
        fi
        cmd+=("-u" "$MYSQL_USER")
        if [ -n "$MYSQL_PASS" ]; then
            cmd+=("-p$MYSQL_PASS")
        fi
        cmd+=("-e" "SELECT 1;")
        
        # 直接执行数组命令，不需要eval
        "${cmd[@]}" &> /dev/null
        
        if [ $? -eq 0 ]; then
            log_success "MySQL连接成功"
            return 0
        else
            log_error "MySQL连接失败"
            return 1
        fi
    else
        log_error "mysql命令未找到，请安装MySQL客户端"
        return 1
    fi
}

# 执行MySQL查询并记录结果
execute_mysql_query() {
    local query="$1"
    local section="$2"
    
    echo "=== $section ===" >> "$LOG_PATH"
    
    # 使用数组直接构建命令，更安全
    local cmd=("mysql")
    if [ -n "$MYSQL_SOCKET" ]; then
        cmd+=("-S" "$MYSQL_SOCKET")
    else
        cmd+=("-h" "$MYSQL_HOST" "-P" "$MYSQL_PORT")
    fi
    cmd+=("-u" "$MYSQL_USER")
    if [ -n "$MYSQL_PASS" ]; then
        cmd+=("-p$MYSQL_PASS")
    fi
    cmd+=("-N" "-e" "$query")
    
    # 直接执行数组命令，不需要eval
    "${cmd[@]}" 2>/dev/null >> "$LOG_PATH"
    
    echo "" >> "$LOG_PATH"
}

# 检查数据库基本信息
check_server_info() {
    log_info "检查数据库基本信息..."
    
    execute_mysql_query "SELECT VERSION();" "DATABASE_VERSION"
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Uptime';" "UPTIME"
    execute_mysql_query "SHOW GLOBAL VARIABLES LIKE 'max_connections';" "MAX_CONNECTIONS"
    execute_mysql_query "SHOW GLOBAL VARIABLES LIKE 'innodb_buffer_pool_size';" "INNODB_BUFFER_POOL_SIZE"
    execute_mysql_query "SHOW GLOBAL VARIABLES LIKE 'version_compile_os';" "OS_VERSION"
    
    # 新增服务器配置项
    execute_mysql_query "SHOW GLOBAL VARIABLES LIKE 'innodb_log_file_size';" "INNODB_LOG_FILE_SIZE"
    execute_mysql_query "SHOW GLOBAL VARIABLES LIKE 'innodb_log_files_in_group';" "INNODB_LOG_FILES_IN_GROUP"
    execute_mysql_query "SHOW GLOBAL VARIABLES LIKE 'innodb_flush_log_at_trx_commit';" "INNODB_FLUSH_LOG_AT_TRX_COMMIT"
    execute_mysql_query "SHOW GLOBAL VARIABLES LIKE 'sync_binlog';" "SYNC_BINLOG"
    execute_mysql_query "SHOW GLOBAL VARIABLES LIKE 'binlog_format';" "BINLOG_FORMAT"
    execute_mysql_query "SHOW GLOBAL VARIABLES LIKE 'transaction_isolation';" "TRANSACTION_ISOLATION"
    execute_mysql_query "SHOW GLOBAL VARIABLES LIKE 'innodb_flush_method';" "INNODB_FLUSH_METHOD"
    execute_mysql_query "SHOW GLOBAL VARIABLES LIKE 'innodb_file_per_table';" "INNODB_FILE_PER_TABLE"
    execute_mysql_query "SHOW GLOBAL VARIABLES LIKE 'open_files_limit';" "OPEN_FILES_LIMIT"
    execute_mysql_query "SHOW GLOBAL VARIABLES LIKE 'table_open_cache';" "TABLE_OPEN_CACHE"
    execute_mysql_query "SHOW GLOBAL VARIABLES LIKE 'max_allowed_packet';" "MAX_ALLOWED_PACKET"
    execute_mysql_query "SHOW GLOBAL VARIABLES LIKE 'wait_timeout';" "WAIT_TIMEOUT"
    execute_mysql_query "SHOW GLOBAL VARIABLES LIKE 'interactive_timeout';" "INTERACTIVE_TIMEOUT"
}

# 检查性能指标
check_performance() {
    log_info "检查性能指标..."
    
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Threads_connected';" "CURRENT_CONNECTIONS"
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Qcache_hits';" "QUERY_CACHE_HITS"
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Qcache_inserts';" "QUERY_CACHE_INSERTS"
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Innodb_buffer_pool_reads';" "INNODB_READS"
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Innodb_buffer_pool_read_requests';" "INNODB_READ_REQUESTS"
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Slow_queries';" "SLOW_QUERIES"
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Table_locks_waited';" "TABLE_LOCKS_WAITED"
    
    # 新增性能指标
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Questions';" "QUESTIONS"
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Queries';" "QUERIES"
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Com_select';" "COM_SELECT"
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Com_insert';" "COM_INSERT"
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Com_update';" "COM_UPDATE"
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Com_delete';" "COM_DELETE"
    
    # InnoDB锁相关
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Innodb_row_lock_waits';" "INNODB_ROW_LOCK_WAITS"
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Innodb_row_lock_time_avg';" "INNODB_ROW_LOCK_TIME_AVG"
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Innodb_row_lock_time_max';" "INNODB_ROW_LOCK_TIME_MAX"
    
    # 死锁检测
    execute_mysql_query "SHOW ENGINE INNODB STATUS;" "INNODB_ENGINE_STATUS"
    
    # 临时表使用
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Created_tmp_tables';" "CREATED_TMP_TABLES"
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Created_tmp_disk_tables';" "CREATED_TMP_DISK_TABLES"
    
    # 连接相关
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Aborted_clients';" "ABORTED_CLIENTS"
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Aborted_connects';" "ABORTED_CONNECTS"
    
    # 事务提交
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Com_commit';" "COM_COMMIT"
    execute_mysql_query "SHOW GLOBAL STATUS LIKE 'Com_rollback';" "COM_ROLLBACK"
}

# 检查安全配置
check_security() {
    log_info "检查安全配置..."
    
    execute_mysql_query "SELECT user, host FROM mysql.user WHERE authentication_string = '' OR authentication_string IS NULL;" "EMPTY_PASSWORD_USERS"
    execute_mysql_query "SELECT user, host FROM mysql.user WHERE user = 'root' AND host NOT IN ('localhost', '127.0.0.1');" "REMOTE_ROOT_USERS"
    execute_mysql_query "SHOW VARIABLES LIKE 'have_ssl';" "SSL_STATUS"
    execute_mysql_query "SELECT user, host, authentication_string FROM mysql.user WHERE host = '%';" "REMOTE_USERS"
}

# 检查存储信息
check_storage() {
    log_info "检查存储信息..."
    
    execute_mysql_query "SELECT table_schema as database_name, ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) as size_mb FROM information_schema.tables WHERE table_schema NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys') GROUP BY table_schema ORDER BY size_mb DESC;" "DATABASE_SIZES"
    
    # 详细的表碎片检查（只显示碎片率>10%的）
    execute_mysql_query "SELECT table_schema, table_name, ROUND((data_length + index_length) / 1024 / 1024, 2) as total_size_mb, ROUND(data_length / 1024 / 1024, 2) as data_size_mb, ROUND(index_length / 1024 / 1024, 2) as index_size_mb, ROUND(data_free / 1024 / 1024, 2) as free_size_mb, ROUND((data_free / (data_length + index_length)) * 100, 2) as fragmentation_percent FROM information_schema.tables WHERE data_length > 0 AND table_schema NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys') AND (data_free / (data_length + index_length)) * 100 > 10 ORDER BY fragmentation_percent DESC;" "TABLE_FRAGMENTATION_DETAILED"
    
    # 简单的碎片化表检查（碎片率>10%）
    execute_mysql_query "SELECT table_schema, table_name, ROUND((data_free / (data_length + index_length)) * 100, 2) as fragmentation_percent FROM information_schema.tables WHERE data_length > 0 AND table_schema NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys') HAVING fragmentation_percent > 10;" "FRAGMENTED_TABLES"
    
    # 检查无主键表
    execute_mysql_query "SELECT table_schema, table_name FROM information_schema.tables WHERE table_schema NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys') AND table_type = 'BASE TABLE' AND table_name NOT IN (SELECT DISTINCT table_name FROM information_schema.statistics WHERE table_schema NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys') AND index_name = 'PRIMARY');" "TABLES_WITHOUT_PRIMARY_KEY"
    
    # 检查分区表
    execute_mysql_query "SELECT table_schema, table_name, partition_name, partition_expression, partition_description FROM information_schema.partitions WHERE table_schema NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys') AND partition_name IS NOT NULL;" "PARTITIONED_TABLES"
    
    # 索引使用率分析（检查重复索引和未使用索引）
    execute_mysql_query "SELECT table_schema, table_name, index_name, GROUP_CONCAT(column_name ORDER BY seq_in_index SEPARATOR ', ') as columns FROM information_schema.statistics WHERE table_schema NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys') GROUP BY table_schema, table_name, index_name ORDER BY table_schema, table_name, index_name;" "INDEX_DETAILS"
    
    # 索引统计（只显示索引数量>=5的表）
    execute_mysql_query "SELECT table_schema, table_name, COUNT(DISTINCT index_name) as index_count FROM information_schema.statistics WHERE table_schema NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys') GROUP BY table_schema, table_name HAVING index_count >= 5 ORDER BY index_count DESC;" "INDEX_STATISTICS"
}

# 检查备份和复制状态
check_backup_replication() {
    log_info "检查备份和复制状态..."
    
    # 获取MySQL主版本号
    local mysql_major_version=$(echo "$MYSQL_VERSION" | cut -d'.' -f1)
    local mysql_minor_version=$(echo "$MYSQL_VERSION" | cut -d'.' -f2)
    local mysql_version_float=$(echo "$mysql_major_version.$mysql_minor_version" | bc -l 2>/dev/null || echo "0")
    
    log_info "检测到MySQL主版本: $mysql_major_version.$mysql_minor_version"
    
    # 使用数组构建命令，更安全
    local cmd=("mysql")
    if [ -n "$MYSQL_SOCKET" ]; then
        cmd+=("-S" "$MYSQL_SOCKET")
    else
        cmd+=("-h" "$MYSQL_HOST" "-P" "$MYSQL_PORT")
    fi
    cmd+=("-u" "$MYSQL_USER")
    if [ -n "$MYSQL_PASS" ]; then
        cmd+=("-p$MYSQL_PASS")
    fi
    
    # 主库状态 - 根据MySQL版本选择命令
    echo "=== MASTER_STATUS ===" >> "$LOG_PATH"
    if (( $(echo "$mysql_version_float >= 8.4" | bc -l 2>/dev/null) )); then
        # MySQL 8.4+ 使用 SHOW BINARY LOG STATUS
        "${cmd[@]}" -e "SHOW BINARY LOG STATUS;" 2>/dev/null >> "$LOG_PATH"
    else
        # 旧版本使用 SHOW MASTER STATUS
        "${cmd[@]}" -e "SHOW MASTER STATUS;" 2>/dev/null >> "$LOG_PATH"
    fi
    echo "" >> "$LOG_PATH"
    
    # 从库状态（详细格式）- 根据MySQL版本选择命令
    echo "=== SLAVE_STATUS ===" >> "$LOG_PATH"
    if (( $(echo "$mysql_version_float >= 8.4" | bc -l 2>/dev/null) )); then
        # MySQL 8.4+ 使用 SHOW REPLICA STATUS
        "${cmd[@]}" -e "SHOW REPLICA STATUS\G" 2>/dev/null >> "$LOG_PATH"
    else
        # 旧版本使用 SHOW SLAVE STATUS
        "${cmd[@]}" -e "SHOW SLAVE STATUS\G" 2>/dev/null >> "$LOG_PATH"
    fi
    echo "" >> "$LOG_PATH"
    
    # 二进制日志 - 保留表头
    echo "=== BINARY_LOGS ===" >> "$LOG_PATH"
    "${cmd[@]}" -e "SHOW BINARY LOGS;" 2>/dev/null >> "$LOG_PATH"
    echo "" >> "$LOG_PATH"
    
    # 检查GTID状态
    execute_mysql_query "SHOW GLOBAL VARIABLES LIKE 'gtid_mode';" "GTID_MODE"
    execute_mysql_query "SHOW GLOBAL VARIABLES LIKE 'enforce_gtid_consistency';" "GTID_CONSISTENCY"
    execute_mysql_query "SELECT @@GLOBAL.GTID_EXECUTED as gtid_executed;" "GTID_EXECUTED"
    
    # 详细的复制状态检查（表格格式）- 根据MySQL版本选择命令
    echo "=== SLAVE_STATUS_DETAILED ===" >> "$LOG_PATH"
    if (( $(echo "$mysql_version_float >= 8.4" | bc -l 2>/dev/null) )); then
        "${cmd[@]}" -e "SHOW REPLICA STATUS;" 2>/dev/null >> "$LOG_PATH"
    else
        "${cmd[@]}" -e "SHOW SLAVE STATUS;" 2>/dev/null >> "$LOG_PATH"
    fi
    echo "" >> "$LOG_PATH"
}

# 检查数据库详细信息
check_database_details() {
    log_info "检查数据库详细信息..."
    
    execute_mysql_query "SHOW DATABASES;" "DATABASE_LIST"
    
    # 使用数组构建命令，更安全
    local cmd=("mysql")
    if [ -n "$MYSQL_SOCKET" ]; then
        cmd+=("-S" "$MYSQL_SOCKET")
    else
        cmd+=("-h" "$MYSQL_HOST" "-P" "$MYSQL_PORT")
    fi
    cmd+=("-u" "$MYSQL_USER")
    if [ -n "$MYSQL_PASS" ]; then
        cmd+=("-p$MYSQL_PASS")
    fi
    cmd+=("-N" "-e" "SELECT schema_name FROM information_schema.schemata WHERE schema_name NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys');")
    
    local databases=$("${cmd[@]}" 2>/dev/null)
    
    for db in $databases; do
        execute_mysql_query "USE \`$db\`; SHOW TABLE STATUS;" "DATABASE_${db}_TABLES"
    done
}

# 检查系统资源
check_system_resources() {
    log_info "检查系统资源..."
    
    echo "=== SYSTEM_CPU_USAGE ===" >> "$LOG_PATH"
    top -bn1 | grep "Cpu(s)" | sed "s/.*, *\([0-9.]*\)%* id.*/\1/" | awk '{print 100 - $1"%"}' >> "$LOG_PATH"
    echo "" >> "$LOG_PATH"
    
    echo "=== SYSTEM_MEMORY_USAGE ===" >> "$LOG_PATH"
    free -m | grep Mem | awk '{printf "%.1f%%", $3/$2*100}' >> "$LOG_PATH"
    echo "" >> "$LOG_PATH"
    
    echo "=== SYSTEM_SWAP_USAGE ===" >> "$LOG_PATH"
    free -m | grep Swap | awk '{if ($2 == 0) print "0%"; else printf "%.1f%%", $3/$2*100}' >> "$LOG_PATH"
    echo "" >> "$LOG_PATH"
    
    echo "=== SYSTEM_DISK_USAGE ===" >> "$LOG_PATH"
    df -h / | awk 'NR==2 {print $5}' >> "$LOG_PATH"
    echo "" >> "$LOG_PATH"
    
    echo "=== MYSQL_PROCESS ===" >> "$LOG_PATH"
    ps aux | grep mysql | grep -v grep | head -1 >> "$LOG_PATH"
    echo "" >> "$LOG_PATH"
    
    # 新增：文件描述符限制
    echo "=== FILE_DESCRIPTOR_LIMIT ===" >> "$LOG_PATH"
    local mysql_pid=$(ps aux | grep mysqld | grep -v grep | awk '{print $2}' | head -1)
    if [ -n "$mysql_pid" ]; then
        cat /proc/$mysql_pid/limits | grep "open files" >> "$LOG_PATH"
    else
        echo "无法获取MySQL进程PID" >> "$LOG_PATH"
    fi
    echo "" >> "$LOG_PATH"
    
    # 新增：网络连接数
    echo "=== NETWORK_CONNECTIONS ===" >> "$LOG_PATH"
    netstat -an | grep :$MYSQL_PORT | grep ESTABLISHED | wc -l >> "$LOG_PATH"
    echo "" >> "$LOG_PATH"
    
    # 新增：MySQL错误日志检查
    echo "=== MYSQL_ERROR_LOG ===" >> "$LOG_PATH"
    local cmd=("mysql")
    if [ -n "$MYSQL_SOCKET" ]; then
        cmd+=("-S" "$MYSQL_SOCKET")
    else
        cmd+=("-h" "$MYSQL_HOST" "-P" "$MYSQL_PORT")
    fi
    cmd+=("-u" "$MYSQL_USER")
    if [ -n "$MYSQL_PASS" ]; then
        cmd+=("-p$MYSQL_PASS")
    fi
    cmd+=("-N" "-e" "SHOW VARIABLES LIKE 'log_error';")
    
    local error_log_path=$("${cmd[@]}" 2>/dev/null | awk '{print $2}')
    
    if [ -n "$error_log_path" ] && [ -f "$error_log_path" ]; then
        echo "错误日志路径: $error_log_path" >> "$LOG_PATH"
        echo "最近2000条错误日志:" >> "$LOG_PATH"
        tail -2000 "$error_log_path" 2>/dev/null >> "$LOG_PATH"
    else
        echo "未找到错误日志文件或路径为空" >> "$LOG_PATH"
    fi
    echo "" >> "$LOG_PATH"
}

# 生成巡检摘要
generate_summary() {
    log_info "生成巡检摘要..."
    
    echo "=== INSPECTION_SUMMARY ===" >> "$LOG_PATH"
    echo "巡检时间: $(date '+%Y-%m-%d %H:%M:%S')" >> "$LOG_PATH"
    echo "实例IP: $INSTANCE_IP" >> "$LOG_PATH"
    echo "实例端口: $MYSQL_PORT" >> "$LOG_PATH"
    if [ -n "$MYSQL_SOCKET" ]; then
        echo "MySQL连接: socket=$MYSQL_SOCKET" >> "$LOG_PATH"
    else
        echo "MySQL主机: $MYSQL_HOST:$MYSQL_PORT" >> "$LOG_PATH"
    fi
    echo "巡检用户: $MYSQL_USER" >> "$LOG_PATH"
    echo "日志文件: $LOG_FILE" >> "$LOG_PATH"
    echo "" >> "$LOG_PATH"
}

# 获取MySQL版本
get_mysql_version() {
    local cmd=("mysql")
    if [ -n "$MYSQL_SOCKET" ]; then
        cmd+=("-S" "$MYSQL_SOCKET")
    else
        cmd+=("-h" "$MYSQL_HOST" "-P" "$MYSQL_PORT")
    fi
    cmd+=("-u" "$MYSQL_USER")
    if [ -n "$MYSQL_PASS" ]; then
        cmd+=("-p$MYSQL_PASS")
    fi
    cmd+=("-N" "-e" "SELECT VERSION();")
    
    "${cmd[@]}" 2>/dev/null | cut -d'-' -f1
}

# 检查MGR集群状态
check_mgr_status() {
    log_info "检查MGR集群状态..."
    
    execute_mysql_query "SELECT * FROM performance_schema.replication_group_members;" "MGR_GROUP_MEMBERS"
    execute_mysql_query "SELECT * FROM performance_schema.replication_group_member_stats;" "MGR_MEMBER_STATS"
}

# 检查PXC集群状态
check_pxc_status() {
    log_info "检查PXC集群状态..."
    
    execute_mysql_query "SHOW STATUS LIKE 'wsrep%';" "PXC_WSREP_STATUS"
    execute_mysql_query "SHOW VARIABLES LIKE 'wsrep%';" "PXC_WSREP_VARIABLES"
}

# 主函数
main() {
    echo "========================================"
    echo "    MySQL数据库自动化巡检脚本"
    echo "========================================"
    echo ""
    
    # 创建输出目录
    mkdir -p "$OUTPUT_DIR"
    
    # 生成日志文件名（统一使用IP+端口号+日期的方式）
    LOG_FILE="mysql_inspection_${INSTANCE_IP}_${MYSQL_PORT}_$(date +%Y%m%d_%H%M%S).log"
    
    # 日志文件路径
    LOG_PATH="$OUTPUT_DIR/$LOG_FILE"
    
    # 显示巡检模式
    if [ "$RDS_CHECK" = true ]; then
        echo "巡检模式: RDS环境（跳过系统资源检查）"
    fi
    if [ "$MGR_CHECK" = true ]; then
        echo "巡检模式: MGR集群检查"
    fi
    if [ "$PXC_CHECK" = true ]; then
        echo "巡检模式: PXC集群检查"
    fi
    echo ""
    
    # 开始巡检
    log_info "开始MySQL数据库巡检..."
    
    # 检查MySQL连接
    if ! check_mysql_connection; then
        log_error "无法连接到MySQL数据库，巡检终止"
        exit 1
    fi
    
    # 获取并记录MySQL版本
    MYSQL_VERSION=$(get_mysql_version)
    log_info "检测到MySQL版本: $MYSQL_VERSION"
    
    # 执行各项检查
    generate_summary
    check_server_info
    check_performance
    check_security
    check_storage
    check_backup_replication
    check_database_details
    
    # 检查是否需要执行系统资源检查
    if [ "$RDS_CHECK" != true ]; then
        check_system_resources
    else
        log_info "RDS环境，跳过系统资源检查"
    fi
    
    # 检查是否需要执行集群检查
    if [ "$MGR_CHECK" = true ]; then
        check_mgr_status
    fi
    
    if [ "$PXC_CHECK" = true ]; then
        check_pxc_status
    fi
    
    # 巡检完成
    log_success "MySQL数据库巡检完成"
    log_info "巡检日志已保存至: $LOG_PATH"
    
    echo ""
    echo "========================================"
    #echo "巡检完成！请将日志文件上传至云服务器进行分析"
    #echo "日志文件: $LOG_PATH"
    echo "========================================"
}

# 解析命令行参数
parse_args() {
    while [ $# -gt 0 ]; do
        case "$1" in
            --help)
                show_help
                exit 0
                ;;
            -h)
                # 如果只有 -h 且没有下一个参数，或下一个参数以 - 开头，则认为是 --help
                if [ $# -eq 1 ] || [[ "$2" == -* ]]; then
                    show_help
                    exit 0
                else
                    # 否则认为是 --host
                    MYSQL_HOST="$2"
                    shift 2
                fi
                ;;
            --host)
                MYSQL_HOST="$2"
                shift 2
                ;;
            --port|-P)
                MYSQL_PORT="$2"
                shift 2
                ;;
            --user|-u)
                MYSQL_USER="$2"
                shift 2
                ;;
            --password|-p)
                MYSQL_PASS="$2"
                shift 2
                ;;
            --socket|-S)
                MYSQL_SOCKET="$2"
                shift 2
                ;;
            --output-dir|-o)
                OUTPUT_DIR="$2"
                shift 2
                ;;
            --rds-check)
                RDS_CHECK=true
                shift
                ;;
            --mgr-check)
                MGR_CHECK=true
                shift
                ;;
            --pxc-check)
                PXC_CHECK=true
                shift
                ;;
            *)
                # 假设第一个非选项参数是MySQL密码（保持向后兼容）
                if [ -z "$MYSQL_PASS" ]; then
                    MYSQL_PASS="$1"
                fi
                shift
                ;;
        esac
    done
}

# 显示帮助信息
show_help() {
    echo "用法: $0 [选项] [MySQL密码]"
    echo ""
    echo "连接选项:"
    echo "  --host HOST           MySQL服务器主机名或IP地址（默认: localhost）"
    echo "  -h HOST               同上（当不接参数时显示帮助）"
    echo "  -P, --port PORT       MySQL服务器端口号（默认: 3306）"
    echo "  -u, --user USER       连接MySQL的用户名（默认: root）"
    echo "  -p, --password PASS   连接MySQL的密码"
    echo "  -S, --socket SOCK     MySQL socket文件路径（优先级高于host+port）"
    echo ""
    echo "输出选项:"
    echo "  -o, --output-dir DIR  巡检日志输出目录（默认: ./inspection_logs）"
    echo ""
    echo "功能开关:"
    echo "  --help                显示帮助信息"
    echo "  -h                    显示帮助信息"
    echo "  --rds-check           RDS环境检查（跳过系统资源检查）"
    echo "  --mgr-check           检查MGR集群（MySQL Group Replication）"
    echo "  --pxc-check           检查PXC集群（Percona XtraDB Cluster）"
    echo ""
    echo "示例:"
    echo "  # 基本用法（使用默认参数）"
    echo "  $0"
    echo ""
    echo "  # 指定主机和端口"
    echo "  $0 --host 192.168.1.100 --port 3307"
    echo ""
    echo "  # 指定用户名和密码"
    echo "  $0 --user admin --password mypassword"
    echo ""
    echo "  # 使用socket连接"
    echo "  $0 --socket /var/run/mysqld/mysqld.sock"
    echo ""
    echo "  # 组合使用多个参数"
    echo "  $0 --host 192.168.1.100 --port 3307 --user admin --password mypassword --output-dir /data/inspection_logs"
    echo ""
    echo "  # 配合其他功能开关使用"
    echo "  $0 --host 192.168.1.100 --port 3307 --rds-check --mgr-check"
    echo ""
    echo "注意: 请确保MySQL客户端已安装并配置正确"
}

# 解析参数
parse_args "$@"

# 执行主函数
main "$@"