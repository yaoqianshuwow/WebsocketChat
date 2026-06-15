<#
.SYNOPSIS
    WsChat-zero 双窗口聊天测试启动器
.DESCRIPTION
    同时启动两个聊天窗口（u1 和 u2），支持单聊、群聊和全自动测试模式。
.PARAMETER Mode
    模式：single（单聊，默认）、group（群聊）、auto（全自动双向测试）
.PARAMETER User1
    第一个用户名（默认 u1）
.PARAMETER User2
    第二个用户名（默认 u2）
.PARAMETER Password
    登录密码（默认 111111）
.PARAMETER Peer1
    第一个用户的聊天对象（单聊模式，默认 u2）
.PARAMETER Peer2
    第二个用户的聊天对象（单聊模式，默认 u1）
.PARAMETER GroupId
    群聊模式的群组 ID
.PARAMETER BaseUrl
    HTTP API 基础地址（默认 http://localhost:8888）
.PARAMETER WsUrl
    WebSocket 基础地址（默认 ws://localhost:8888）
.EXAMPLE
    .\test\start-chat-pair.ps1 -Mode single
    # 启动单聊 u1<->u2

.EXAMPLE
    .\test\start-chat-pair.ps1 -Mode group -GroupId 1
    # 启动群聊，两个用户在同个群组

.EXAMPLE
    .\test\start-chat-pair.ps1 -Mode auto
    # 全自动双向测试（单聊+群聊一体化验证）
#>

param(
    [ValidateSet("single", "group", "auto")]
    [string]$Mode = "single",
    [string]$User1 = "u1",
    [string]$User2 = "u2",
    [string]$Password = "111111",
    [string]$Peer1 = "u2",
    [string]$Peer2 = "u1",
    [long]$GroupId = 0,
    [string]$BaseUrl = "http://localhost:8888",
    [string]$WsUrl = "ws://localhost:8888"
)

$repoRoot = Split-Path -Parent $PSScriptRoot
$cliPath = Join-Path $PSScriptRoot "chat_cli.go"
$autoPath = Join-Path $PSScriptRoot "chat_pair_test.go"

if ($Mode -eq "group" -and $GroupId -le 0) {
    Write-Error "群聊模式必须传入 -GroupId"
    exit 1
}

function Start-ChatWindow {
    param(
        [string]$Title,
        [string]$Command
    )
    Start-Process powershell -ArgumentList @(
        "-NoExit",
        "-Command",
        "`$Host.UI.RawUI.WindowTitle = '$Title'; Set-Location '$repoRoot'; Write-Host '=== $Title ==='; $Command"
    ) -WindowStyle Normal | Out-Null
}

switch ($Mode) {
    "single" {
        Write-Host "启动单聊测试：$User1 <-> $User2"
        $cmd1 = "go run `"$cliPath`" -mode single -user $User1 -pass $Password -peer $Peer1 -base $BaseUrl -ws $WsUrl"
        $cmd2 = "go run `"$cliPath`" -mode single -user $User2 -pass $Password -peer $Peer2 -base $BaseUrl -ws $WsUrl"
        Start-ChatWindow -Title "$User1-single" -Command $cmd1
        Start-ChatWindow -Title "$User2-single" -Command $cmd2
    }
    "group" {
        Write-Host "启动群聊测试：$User1 和 $User2 在群组 $GroupId 中聊天"
        $cmd1 = "go run `"$cliPath`" -mode group -user $User1 -pass $Password -group-id $GroupId -base $BaseUrl -ws $WsUrl"
        $cmd2 = "go run `"$cliPath`" -mode group -user $User2 -pass $Password -group-id $GroupId -base $BaseUrl -ws $WsUrl"
        Start-ChatWindow -Title "$User1-group-$GroupId" -Command $cmd1
        Start-ChatWindow -Title "$User2-group-$GroupId" -Command $cmd2
    }
    "auto" {
        Write-Host "启动全自动双向测试（单聊+群聊）"
        $cmd = "go run `"$autoPath`""
        Start-ChatWindow -Title "auto-test" -Command $cmd
    }
}

Write-Host ""
Write-Host "已启动聊天窗口，输入 /quit 退出"
Write-Host ""
Write-Host "示例："
Write-Host "  .\test\start-chat-pair.ps1 -Mode single"
Write-Host "  .\test\start-chat-pair.ps1 -Mode group -GroupId 123"
Write-Host "  .\test\start-chat-pair.ps1 -Mode auto"
