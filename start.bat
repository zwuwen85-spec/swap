@echo off
echo ========================================
echo   校园闲置物品交换平台 - 启动脚本
echo ========================================
echo.

echo [1/3] 检查数据库配置...
echo.

REM 检查MySQL服务
sc query MySQL80 >nul 2>&1
if %errorlevel% neq 0 (
    echo [警告] MySQL服务未运行，正在尝试启动...
    net start MySQL80
    timeout /t 3 >nul
)

REM 检查Redis服务
redis-cli ping >nul 2>&1
if %errorlevel% neq 0 (
    echo [警告] Redis服务未运行，正在尝试启动...
    redis-server --service-start
    timeout /t 2 >nul
)

echo [提示] 请确保MySQL已创建数据库 campus_swap
echo [提示] 如未创建，请执行： CREATE DATABASE campus_swap CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
echo.

echo [2/3] 启动后端服务...
echo.
cd backend
echo 正在启动后端服务...
start "后端服务" cmd /k "server.exe"
echo.

echo [3/3] 启动前端服务...
echo.
cd frontend
echo 正在安装依赖（首次运行）...
if not exist "node_modules" (
    call npm install
)
echo.
echo 正在启动前端开发服务器...
start "前端服务" cmd /k "npm run dev"
echo.

echo ========================================
echo   服务启动完成！
echo ========================================
echo.
echo 后端服务: http://localhost:8080
echo 前端服务: http://localhost:5173
echo.
echo 打开浏览器访问: http://localhost:5173
echo.
echo 按任意键关闭此窗口...
pause >nul
