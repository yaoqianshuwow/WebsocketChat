@echo off
chcp 65001 >nul
REM ============================================
REM WsChat-zero 双窗口聊天测试启动器 (CMD版本)
REM 支持单聊和群聊两种模式
REM ============================================

set MODE=%1
if "%MODE%"=="" set MODE=single

set USER1=u1
set USER2=u2
set PASSWORD=111111
set BASE_URL=http://localhost:8888
set WS_URL=ws://localhost:8888

if /i "%MODE%"=="single" (
    echo 启动单聊测试：%USER1% ^<->^ %USER2%
    start "u1-single" cmd /c "cd /d %~dp0.. && go run test\chat_cli.go -mode single -user %USER1% -pass %PASSWORD% -peer %USER2% -base %BASE_URL% -ws %WS_URL%"
    timeout /t 1 >nul
    start "u2-single" cmd /c "cd /d %~dp0.. && go run test\chat_cli.go -mode single -user %USER2% -pass %PASSWORD% -peer %USER1% -base %BASE_URL% -ws %WS_URL%"
) else if /i "%MODE%"=="group" (
    set GROUP_ID=%2
    if "%GROUP_ID%"=="" (
        echo 用法: start-chat-pair.bat group ^<group_id^>
        echo 例如: start-chat-pair.bat group 1
        pause
        exit /b 1
    )
    echo 启动群聊测试：%USER1% 和 %USER2% 在群组 %GROUP_ID% 中聊天
    start "u1-group" cmd /c "cd /d %~dp0.. && go run test\chat_cli.go -mode group -user %USER1% -pass %PASSWORD% -group-id %GROUP_ID% -base %BASE_URL% -ws %WS_URL%"
    timeout /t 1 >nul
    start "u2-group" cmd /c "cd /d %~dp0.. && go run test\chat_cli.go -mode group -user %USER2% -pass %PASSWORD% -group-id %GROUP_ID% -base %BASE_URL% -ws %WS_URL%"
) else if /i "%MODE%"=="auto" (
    echo 启动全自动双向测试（单聊+群聊）
    start "auto-test" cmd /c "cd /d %~dp0.. && go run test\chat_pair_test.go"
) else (
    echo 用法: start-chat-pair.bat [single^|group ^<group_id^>^|auto]
    echo   single     - 单聊模式（默认）
    echo   group ^<id^> - 群聊模式
    echo   auto       - 全自动双向测试
)

echo.
echo 已启动聊天窗口，输入 /quit 退出
