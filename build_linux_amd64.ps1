# 设置目标平台
$env:GOOS = "linux"
$env:GOARCH = "amd64"

# 读取 VERSION 文件
$VERSION = Get-Content -Raw -Path "VERSION" | ForEach-Object { $_.Trim() }

# 获取当前 UTC 时间（ISO 8601 格式）
#$BUILD_TIME = (Get-Date -AsUTC).ToString("yyyy-MM-ddTHH:mm:ssZ")
$BUILD_TIME = ([datetime]::UtcNow).ToString("yyyy-MM-ddTHH:mm:ssZ")

# 获取 Git 提交短哈希
$GIT_COMMIT = (git rev-parse --short HEAD 2>$null) -replace "`n", ""

# 设置输出目录
$BIN_PATH = Join-Path -Path (Get-Location) -ChildPath "bin"

# 打印构建信息
Write-Host "VERSION=$VERSION"
Write-Host "BUILD_TIME=$BUILD_TIME"
Write-Host "GIT_COMMIT=$GIT_COMMIT"
Write-Host "BIN_PATH=$BIN_PATH"

# 确保输出目录存在
New-Item -ItemType Directory -Force -Path "$BIN_PATH\bin_linux_amd64" | Out-Null

# 构建 nwpc_data_client
Write-Host "Building nwpc_data_client..."
go build `
    -ldflags "-X github.com/cemc-oper/nwpc-data-client/common.Version=$VERSION `
             -X github.com/cemc-oper/nwpc-data-client/common.BuildTime=$BUILD_TIME `
             -X github.com/cemc-oper/nwpc-data-client/common.GitCommit=$GIT_COMMIT" `
    -o "$BIN_PATH\bin_linux_amd64\nwpc_data_client" `
    "data_client/main.go"
Write-Host "Build complete: $BIN_PATH\bin_linux_amd64\nwpc_data_client"

# 构建 nwpc_data_checker
Write-Host "Building nwpc_data_checker..."
go build `
    -ldflags "-X github.com/cemc-oper/nwpc-data-client/common.Version=$VERSION `
             -X github.com/cemc-oper/nwpc-data-client/common.BuildTime=$BUILD_TIME `
             -X github.com/cemc-oper/nwpc-data-client/common.GitCommit=$GIT_COMMIT" `
    -o "$BIN_PATH\bin_linux_amd64\nwpc_data_checker" `
    "data_checker/main.go"

Write-Host "Build complete: $BIN_PATH\bin_linux_amd64\nwpc_data_checker"