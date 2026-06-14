param(
  [string]$ServerIp = '127.0.0.1'
)

$ErrorActionPreference = 'Stop'

$root = Resolve-Path (Join-Path $PSScriptRoot '..')
Set-Location $root

$releaseRoot = Join-Path $root 'release'
$buildRoot = Join-Path $releaseRoot 'build'
$binRoot = Join-Path $buildRoot 'bin'

function Ensure-CleanDir {
  param([string]$Path)
  if (Test-Path $Path) {
    Remove-Item $Path -Recurse -Force
  }
  New-Item -ItemType Directory -Force -Path $Path | Out-Null
}

function Build-GoBinary {
  param(
    [string]$WorkDir,
    [string]$Target,
    [string]$OutputPath
  )

  Push-Location $WorkDir
  try {
    $env:GOOS = 'linux'
    $env:GOARCH = 'amd64'
    $env:CGO_ENABLED = '0'
    $env:GOPROXY = 'https://goproxy.cn,direct'
    $env:GOSUMDB = 'off'
    go build -trimpath -ldflags "-s -w" -o $OutputPath $Target
  } finally {
    Pop-Location
    Remove-Item Env:GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
    Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue
    Remove-Item Env:GOPROXY -ErrorAction SilentlyContinue
    Remove-Item Env:GOSUMDB -ErrorAction SilentlyContinue
  }
}

function New-GoContext {
  param(
    [string]$Name,
    [string]$BinaryName,
    [string]$ConfigFile
  )

  $contextDir = Join-Path $buildRoot $Name
  $contextBin = Join-Path $contextDir 'bin'
  $contextEtc = Join-Path $contextDir 'etc'
  Ensure-CleanDir $contextDir
  New-Item -ItemType Directory -Force -Path $contextBin, $contextEtc | Out-Null
  Copy-Item -Path (Join-Path $binRoot $BinaryName) -Destination (Join-Path $contextBin 'service')

  $configContent = Get-Content (Join-Path $root $ConfigFile) -Raw
  $configContent = $configContent.Replace('__SERVER_IP__', $ServerIp)
  Set-Content -Path (Join-Path $contextEtc 'service.yaml') -Value $configContent -Encoding utf8

  return $contextDir
}

function New-WebContext {
  $contextDir = Join-Path $buildRoot 'web'
  $contextBin = Join-Path $contextDir 'bin'
  $contextDist = Join-Path $contextDir 'web-dist'
  Ensure-CleanDir $contextDir
  New-Item -ItemType Directory -Force -Path $contextBin | Out-Null
  Copy-Item -Path (Join-Path $binRoot 'web-server') -Destination (Join-Path $contextBin 'web-server')
  Copy-Item -Path (Join-Path $root 'web/dist/*') -Destination $contextDist -Recurse
  return $contextDir
}

try {
  Ensure-CleanDir $releaseRoot
  Ensure-CleanDir $buildRoot
  Ensure-CleanDir $binRoot

  Write-Host '==> Build frontend dist'
  $env:VITE_API_BASE_URL = "http://$ServerIp"
  $env:VITE_WS_BASE_URL = "ws://$ServerIp"
  Push-Location (Join-Path $root 'web')
  try {
    npm run build
  } finally {
    Pop-Location
    Remove-Item Env:VITE_API_BASE_URL -ErrorAction SilentlyContinue
    Remove-Item Env:VITE_WS_BASE_URL -ErrorAction SilentlyContinue
  }

  Write-Host '==> Build linux binaries'
  Build-GoBinary -WorkDir (Join-Path $root 'app/gateway') -Target '.' -OutputPath (Join-Path $binRoot 'gateway')
  Build-GoBinary -WorkDir (Join-Path $root 'app/user') -Target '.' -OutputPath (Join-Path $binRoot 'user')
  Build-GoBinary -WorkDir (Join-Path $root 'app/friend') -Target '.' -OutputPath (Join-Path $binRoot 'friend')
  Build-GoBinary -WorkDir (Join-Path $root 'app/file') -Target '.' -OutputPath (Join-Path $binRoot 'file')
  Build-GoBinary -WorkDir (Join-Path $root 'app/voice') -Target '.' -OutputPath (Join-Path $binRoot 'voice')
  Build-GoBinary -WorkDir (Join-Path $root 'app/msg-forward') -Target '.' -OutputPath (Join-Path $binRoot 'msg-forward')
  Build-GoBinary -WorkDir (Join-Path $root 'app/msg-store') -Target './cmd' -OutputPath (Join-Path $binRoot 'msg-store')
  Push-Location (Join-Path $root 'deploy/web-server')
  try {
    $env:GOWORK = 'off'
    $env:GOOS = 'linux'
    $env:GOARCH = 'amd64'
    $env:CGO_ENABLED = '0'
    $env:GOPROXY = 'https://goproxy.cn,direct'
    $env:GOSUMDB = 'off'
    go build -trimpath -ldflags "-s -w" -o (Join-Path $binRoot 'web-server') .
  } finally {
    Pop-Location
    Remove-Item Env:GOWORK -ErrorAction SilentlyContinue
    Remove-Item Env:GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:GOARCH -ErrorAction SilentlyContinue
    Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue
    Remove-Item Env:GOPROXY -ErrorAction SilentlyContinue
    Remove-Item Env:GOSUMDB -ErrorAction SilentlyContinue
  }

  $images = @(
    @{ Name = 'wschat-gateway:latest'; Dockerfile = 'deploy/Dockerfile.release-go'; Context = New-GoContext -Name 'gateway' -BinaryName 'gateway' -ConfigFile 'deploy/docker-config/gateway-api.yaml' },
    @{ Name = 'wschat-user:latest'; Dockerfile = 'deploy/Dockerfile.release-go'; Context = New-GoContext -Name 'user' -BinaryName 'user' -ConfigFile 'deploy/docker-config/user.yaml' },
    @{ Name = 'wschat-friend:latest'; Dockerfile = 'deploy/Dockerfile.release-go'; Context = New-GoContext -Name 'friend' -BinaryName 'friend' -ConfigFile 'deploy/docker-config/friend.yaml' },
    @{ Name = 'wschat-file:latest'; Dockerfile = 'deploy/Dockerfile.release-go'; Context = New-GoContext -Name 'file' -BinaryName 'file' -ConfigFile 'deploy/docker-config/file.yaml' },
    @{ Name = 'wschat-voice:latest'; Dockerfile = 'deploy/Dockerfile.release-go'; Context = New-GoContext -Name 'voice' -BinaryName 'voice' -ConfigFile 'deploy/docker-config/voice.yaml' },
    @{ Name = 'wschat-msg-forward:latest'; Dockerfile = 'deploy/Dockerfile.release-go'; Context = New-GoContext -Name 'msg-forward' -BinaryName 'msg-forward' -ConfigFile 'deploy/docker-config/msg.yaml' },
    @{ Name = 'wschat-msg-store:latest'; Dockerfile = 'deploy/Dockerfile.release-go'; Context = New-GoContext -Name 'msg-store' -BinaryName 'msg-store' -ConfigFile 'deploy/docker-config/msg-store.yaml' },
    @{ Name = 'wschat-web:latest'; Dockerfile = 'deploy/Dockerfile.release-web'; Context = New-WebContext }
  )

  Write-Host '==> Build release images'
  foreach ($image in $images) {
    docker build -t $image.Name -f $image.Dockerfile $image.Context
  }

  Write-Host '==> Export image tarball'
  $tarPath = Join-Path $releaseRoot 'wschat-images.tar'
  if (Test-Path $tarPath) {
    Remove-Item $tarPath -Force
  }
  docker save -o $tarPath @($images | ForEach-Object { $_.Name })

  Write-Host '==> Copy frontend static files'
  $webDist = Join-Path $releaseRoot 'web-dist'
  if (Test-Path $webDist) {
    Remove-Item $webDist -Recurse -Force
  }
  Copy-Item -Path (Join-Path $root 'web/dist/*') -Destination $webDist -Recurse

  $releaseConfigDir = Join-Path $releaseRoot 'config'
  Ensure-CleanDir $releaseConfigDir
  Copy-Item -Path (Join-Path $root 'deploy/docker-config/*') -Destination $releaseConfigDir -Recurse
  (Get-Content (Join-Path $releaseConfigDir 'file.yaml') -Raw).Replace('__SERVER_IP__', $ServerIp) | Set-Content -Path (Join-Path $releaseConfigDir 'file.yaml') -Encoding utf8

  Copy-Item -Path (Join-Path $root 'deploy/docker-compose.release.yml') -Destination (Join-Path $releaseRoot 'docker-compose.release.yml') -Force

  Write-Host "Done. Release directory: $releaseRoot"
} finally {
  Remove-Item Env:VITE_API_BASE_URL -ErrorAction SilentlyContinue
  Remove-Item Env:VITE_WS_BASE_URL -ErrorAction SilentlyContinue
  Remove-Item Env:GOPROXY -ErrorAction SilentlyContinue
  Remove-Item Env:GOSUMDB -ErrorAction SilentlyContinue
}
