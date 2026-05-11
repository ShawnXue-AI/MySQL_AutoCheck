#!/bin/bash
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
LOG_DIR="$SCRIPT_DIR/logs"

stop_by_pid() {
    local pid_file="$1"
    local name="$2"
    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        if kill -0 "$pid" 2>/dev/null; then
            kill "$pid" 2>/dev/null
            echo "[OK] $name 已停止 (PID: $pid)"
            rm -f "$pid_file"
        else
            echo "[提示] $name 进程不存在，清理 PID 文件"
            rm -f "$pid_file"
        fi
    else
        echo "[提示] 未找到 $name PID 文件"
    fi
}

stop_by_pid "$LOG_DIR/backend.pid" "后端服务"
stop_by_pid "$LOG_DIR/frontend.pid" "前端服务"

echo "所有服务已停止。"
