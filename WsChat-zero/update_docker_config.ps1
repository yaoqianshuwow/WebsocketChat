$configPath = "$env:USERPROFILE\.docker\daemon.json"
$config = @{
    builder = @{
        gc = @{
            defaultKeepStorage = "20GB"
            enabled = $true
        }
    }
    experimental = $false
    "fixed-cidr-v6" = ""
    ipv6 = $false
    "registry-mirrors" = @()
}
$json = $config | ConvertTo-Json
Set-Content -Path $configPath -Value $json -Encoding utf8
Write-Host "daemon.json updated with no mirrors"
