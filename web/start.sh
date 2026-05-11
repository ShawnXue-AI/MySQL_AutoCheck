#!/bin/bash
# MySQL巡检报告系统 - 启动脚本 (Linux/Mac)

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
LOG_DIR="$SCRIPT_DIR/logs"
mkdir -p "$LOG_DIR"

echo "========================================"
echo " MySQL巡检报告系统 - 启动中..."
echo "========================================"
echo ""

# 检查 Node.js
if ! command -v node &> /dev/null; then
    echo "[错误] 未找到 Node.js，请先安装 Node.js"
    exit 1
fi

# 检查 Go
if ! command -v go &> /dev/null; then
    echo "[错误] 未找到 Go，请先安装 Go"
    exit 1
fi

# 安装前端依赖
echo "[1/4] 安装前端依赖..."
cd "$SCRIPT_DIR/frontend"
npm install --silent
echo "[OK] 前端依赖安装完成"

# 下载Go依赖
echo "[2/4] 下载后端依赖..."
cd "$SCRIPT_DIR/backend"
go mod tidy
echo "[OK] 后端依赖下载完成"

# 先停掉旧进程
if [ -f "$LOG_DIR/backend.pid" ]; then
    kill "$(cat "$LOG_DIR/backend.pid")" 2>/dev/null
    rm -f "$LOG_DIR/backend.pid"
fi
if [ -f "$LOG_DIR/frontend.pid" ]; then
    kill "$(cat "$LOG_DIR/frontend.pid")" 2>/dev/null
    rm -f "$LOG_DIR/frontend.pid"
fi

# 启动后端（nohup 守护）
echo "[3/4] 启动后端服务..."
cd "$SCRIPT_DIR/backend"
nohup go run . > "$LOG_DIR/backend.log" 2>&1 &
BACKEND_PID=$!
echo $BACKEND_PID > "$LOG_DIR/backend.pid"
echo "[OK] 后端服务已启动 (PID: $BACKEND_PID, 端口 8080)"

sleep 2

# 启动前端（nohup 守护）
echo "[4/4] 启动前端服务..."
cd "$SCRIPT_DIR/frontend"
nohup npm run dev > "$LOG_DIR/frontend.log" 2>&1 &
FRONTEND_PID=$!
echo $FRONTEND_PID > "$LOG_DIR/frontend.pid"
echo "[OK] 前端服务已启动 (PID: $FRONTEND_PID, 端口 3000)"

echo ""
echo "========================================"
echo " 启动完成！"
echo ""
echo " 前端地址: http://localhost:3000"
echo " 后端API:  http://localhost:8080/api"
echo ""
echo " 关闭终端不影响服务运行。"
echo " 查看后端日志: tail -f $LOG_DIR/backend.log"
echo " 查看前端日志: tail -f $LOG_DIR/frontend.log"
echo ""
echo " 停止服务请运行: web/stop.sh"
echo "========================================"
