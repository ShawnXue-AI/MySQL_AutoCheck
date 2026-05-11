#!/bin/bash

# MySQL巡检配置脚本
# 用于配置Shell巡检脚本的连接参数

CONFIG_FILE="./mysql_inspection.conf"

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 显示当前配置
show_config() {
    if [ -f "$CONFIG_FILE" ]; then
        echo -e "${BLUE}当前配置:${NC}"
        cat "$CONFIG_FILE"
        echo ""
    else
        echo "配置文件不存在，使用默认配置"
    fi
}

# 交互式配置
interactive_config() {
    echo -e "${GREEN}MySQL巡检配置向导${NC}"
    echo "========================================"
    echo ""
    
    # 读取当前配置或使用默认值
    if [ -f "$CONFIG_FILE" ]; then
        source "$CONFIG_FILE"
    else
        MYSQL_HOST="localhost"
        MYSQL_PORT="3306"
        MYSQL_USER="root"
        MYSQL_PASS=""
    fi
    
    # 主机地址
    read -p "MySQL主机地址 [当前: $MYSQL_HOST]: " new_host
    if [ ! -z "$new_host" ]; then
        MYSQL_HOST="$new_host"
    fi
    
    # 端口
    read -p "MySQL端口 [当前: $MYSQL_PORT]: " new_port
    if [ ! -z "$new_port" ]; then
        MYSQL_PORT="$new_port"
    fi
    
    # 用户名
    read -p "MySQL用户名 [当前: $MYSQL_USER]: " new_user
    if [ ! -z "$new_user" ]; then
        MYSQL_USER="$new_user"
    fi
    
    # 密码（安全输入）
    echo -n "MySQL密码 [留空保持当前]: "
    read -s new_pass
    echo ""
    if [ ! -z "$new_pass" ]; then
        MYSQL_PASS="$new_pass"
    fi
    
    # 保存配置
    cat > "$CONFIG_FILE" << EOF
# MySQL巡检配置
# 生成时间: $(date '+%Y-%m-%d %H:%M:%S')

MYSQL_HOST="$MYSQL_HOST"
MYSQL_PORT="$MYSQL_PORT"
MYSQL_USER="$MYSQL_USER"
MYSQL_PASS="$MYSQL_PASS"
EOF
    
    echo -e "${GREEN}配置已保存到 $CONFIG_FILE${NC}"
    echo ""
}

# 测试连接
test_connection() {
    if [ -f "$CONFIG_FILE" ]; then
        source "$CONFIG_FILE"
        
        echo "测试MySQL连接..."
        
        if [ -z "$MYSQL_PASS" ]; then
            mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -e "SELECT VERSION();" 2>/dev/null
        else
            mysql -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_PASS" -e "SELECT VERSION();" 2>/dev/null
        fi
        
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}连接测试成功${NC}"
        else
            echo -e "${RED}连接测试失败${NC}"
        fi
    else
        echo "请先运行配置向导"
    fi
}

# 生成巡检脚本
generate_inspection_script() {
    if [ -f "$CONFIG_FILE" ]; then
        source "$CONFIG_FILE"
        
        # 创建带配置的巡检脚本
        cat > "mysql_inspection_with_config.sh" << 'EOF'
#!/bin/bash

# 自动生成的MySQL巡检脚本（带配置）
# 配置来源: mysql_inspection.conf

# 加载配置
if [ -f "./mysql_inspection.conf" ]; then
    source "./mysql_inspection.conf"
else
    echo "错误: 配置文件不存在"
    exit 1
fi

# 执行巡检
./mysql_inspection.sh "$MYSQL_PASS"
EOF
        
        chmod +x "mysql_inspection_with_config.sh"
        echo -e "${GREEN}已生成带配置的巡检脚本: mysql_inspection_with_config.sh${NC}"
        echo "使用方法: ./mysql_inspection_with_config.sh"
    else
        echo "请先运行配置向导"
    fi
}

# 显示帮助
show_help() {
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -s, --show          显示当前配置"
    echo "  -c, --config        交互式配置"
    echo "  -t, --test          测试MySQL连接"
    echo "  -g, --generate      生成带配置的巡检脚本"
    echo "  -h, --help          显示此帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 --config          # 交互式配置"
    echo "  $0 --test            # 测试连接"
    echo "  $0 --generate        # 生成巡检脚本"
}

# 主函数
main() {
    case "$1" in
        -s|--show)
            show_config
            ;;
        -c|--config)
            interactive_config
            ;;
        -t|--test)
            test_connection
            ;;
        -g|--generate)
            generate_inspection_script
            ;;
        -h|--help)
            show_help
            ;;
        *)
            if [ $# -eq 0 ]; then
                show_config
                echo ""
                echo "使用 '$0 --help' 查看可用选项"
            else
                echo "未知选项: $1"
                show_help
                exit 1
            fi
            ;;
    esac
}

# 脚本入口
main "$@"