<#
.SYNOPSIS
    Cross-compile nwpc-data-client binaries for a target Linux architecture.
.DESCRIPTION
    Builds nwpc_data_client and nwpc_data_checker for the specified GOARCH.
    Intended as a Windows fallback when GoReleaser is not available.
.PARAMETER Arch
    Target architecture: amd64 or arm64. Default: amd64.
#>
param(
    [ValidateSet("amd64", "arm64")]
    [string]$Arch = "amd64"
)

$env:GOOS = "linux"
$env:GOARCH = $Arch
$env:CGO_ENABLED = "0"

$VERSION = ((git describe --tags --always --dirty 2>$null) -replace "`n", "")
if (-not $VERSION) { $VERSION = (Get-Content -Raw -Path "VERSION").Trim() }
if (-not $VERSION) { $VERSION = "dev" }
$BUILD_TIME = ([datetime]::UtcNow).ToString("yyyy-MM-ddTHH:mm:ssZ")
$GIT_COMMIT = ((git rev-parse --short HEAD 2>$null) -replace "`n", "")
if (-not $GIT_COMMIT) { $GIT_COMMIT = "unknown" }

$BIN_PATH = Join-Path (Get-Location) "bin"
$OUT_DIR = Join-Path $BIN_PATH "bin_linux_$Arch"
New-Item -ItemType Directory -Force -Path $OUT_DIR | Out-Null

$Binaries = @(
    @{ Name = "nwpc_data_client";  Source = "data_client/main.go" },
    @{ Name = "nwpc_data_checker"; Source = "data_checker/main.go" }
)

foreach ($bin in $Binaries) {
    $outFile = Join-Path $OUT_DIR $bin.Name
    Write-Host "Building $($bin.Name) for linux/$Arch ..."
    go build `
        -ldflags "-X github.com/cemc-oper/nwpc-data-client/common.Version=$VERSION -X github.com/cemc-oper/nwpc-data-client/common.BuildTime=$BUILD_TIME -X github.com/cemc-oper/nwpc-data-client/common.GitCommit=$GIT_COMMIT" `
        -o $outFile `
        $bin.Source
    Write-Host "  -> $outFile"
}

Write-Host "Done. Output: $OUT_DIR"
