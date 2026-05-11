@echo off
chcp 65001 >nul
title MySQL巡检报告系统

echo ========================================
echo  MySQL巡检报告系统 - 启动中...
echo ========================================
echo.

:: 检查 Node.js
where node >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [错误] 未找到 Node.js，请先安装 Node.js (https://nodejs.org)
    echo.
    pause
    exit /b 1
)

:: 检查 Go
where go >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [错误] 未找到 Go，请先安装 Go (https://go.dev/dl/)
    echo.
    pause
    exit /b 1
)

:: 切换到项目根目录
cd /d "%~dp0.."

:: 安装前端依赖
echo [1/4] 安装前端依赖...
cd /d "%~dp0frontend"
call npm install 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [错误] 前端依赖安装失败
    pause
    exit /b 1
)
echo [OK] 前端依赖安装完成

:: 下载Go依赖
echo [2/4] 下载后端依赖...
cd /d "%~dp0backend"
go mod tidy 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [错误] Go依赖下载失败
    echo 请确认网络通畅后重试，或手动执行: cd web\backend ^&^& go mod tidy
    pause
    exit /b 1
)
echo [OK] 后端依赖下载完成

:: 启动后端（新窗口，可以看到日志和报错）
echo [3/4] 启动后端服务...
cd /d "%~dp0backend"
start "MySQL巡检-后端" cmd /k "echo 后端服务启动中... && go run . && pause"
echo 等待后端启动(5秒)...
timeout /t 5 /nobreak >nul

:: 测试后端是否启动成功
curl -s http://localhost:8080/api/health >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [警告] 后端服务可能未正常启动，请检查后端窗口的错误信息
    echo 你也可以手动打开后端窗口查看报错
) else (
    echo [OK] 后端服务已启动 (端口 8080)
)

:: 启动前端
echo [4/4] 启动前端服务...
cd /d "%~dp0frontend"
start "MySQL巡检-前端" cmd /k "echo 前端服务启动中... && npm run dev"
echo [OK] 前端服务已启动 (端口 3000)

echo.
echo ========================================
echo  启动完成！
echo.
echo  前端地址: http://localhost:3000
echo  后端API:  http://localhost:8080/api
echo.
echo  提示: 如果无法访问，请检查防火墙是否放行 3000/8080 端口
echo  注意: 关闭本窗口不会停止服务，请手动关闭两个新窗口
echo ========================================
echo.
pause
