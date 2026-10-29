#Requires -Version 5.1
<#
.SYNOPSIS
    Build gowtree Windows installer (single release artifact).
.OUTPUTS
    release/gowtree-setup-1.4.0.exe
#>
param(
    [switch]$SkipTests,
    [string]$Version = "1.4.0"
)

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
Set-Location $Root

$BuildDir = Join-Path $Root "build"
$ReleaseDir = Join-Path $Root "release"
$SetupExe = Join-Path $ReleaseDir "gowtree-setup-$Version.exe"
$InternalExe = Join-Path $BuildDir "gowtree.exe"

Write-Host ""
Write-Host "  gowtree installer build v$Version" -ForegroundColor Green
Write-Host "  ---------------------------------" -ForegroundColor DarkGray
Write-Host ""

New-Item -ItemType Directory -Force -Path $BuildDir, $ReleaseDir | Out-Null

Write-Host "[1/5] Generating icon and wizard images..." -ForegroundColor Cyan
& (Join-Path $Root "scripts\generate-assets.ps1")

Write-Host "[2/5] Embedding Windows icon and version info..." -ForegroundColor Cyan
# go install github.com/tc-hib/go-winres@latest
# if ($LASTEXITCODE -ne 0) { throw "go-winres install failed" }
# $syso = Join-Path $Root "cmd\gowtree\rsrc_windows_amd64.syso"
# if (Test-Path $syso) { Remove-Item $syso -Force }
# go-winres make --in winres/winres.json --out cmd/gowtree/
# if ($LASTEXITCODE -ne 0) { throw "go-winres make failed" }

Write-Host "[3/5] Building gowtree.exe (embedded in installer)..." -ForegroundColor Cyan
$ldflags = "-s -w -X github.com/hosseinghorbani0/gowtree/internal/metadata.Version=v$Version"
go build -ldflags $ldflags -o $InternalExe ./cmd/gowtree
if ($LASTEXITCODE -ne 0) { throw "go build failed" }
# Remove-Item $syso -Force -ErrorAction SilentlyContinue
Write-Host "      -> build\gowtree.exe (internal, not published)" -ForegroundColor DarkGray

if (-not $SkipTests) {
    Write-Host "[4/5] Running tests..." -ForegroundColor Cyan
    go test ./...
    if ($LASTEXITCODE -ne 0) { throw "tests failed" }
    Write-Host "      -> all tests passed" -ForegroundColor DarkGreen
} else {
    Write-Host "[4/5] Tests skipped" -ForegroundColor Yellow
}

Write-Host "[5/5] Building setup wizard..." -ForegroundColor Cyan
$iscc = @(
    "C:\InnoSetup6\ISCC.exe",
    "${env:ProgramFiles(x86)}\Inno Setup 6\ISCC.exe",
    "$env:ProgramFiles\Inno Setup 6\ISCC.exe"
) | Where-Object { Test-Path $_ } | Select-Object -First 1
if (-not $iscc) { throw "Inno Setup not found. Install: winget install JRSoftware.InnoSetup" }

& $iscc "/DMyAppVersion=$Version" "installer\gowtree.iss"
if ($LASTEXITCODE -ne 0) { throw "Inno Setup build failed" }

$builtSetup = Get-ChildItem "installer\output\gowtree-setup-*.exe" | Sort-Object LastWriteTime -Descending | Select-Object -First 1
if (-not $builtSetup) { throw "Setup exe not found" }

Get-ChildItem $ReleaseDir -Filter "gowtree*" | Remove-Item -Force -ErrorAction SilentlyContinue
Copy-Item $builtSetup.FullName $SetupExe -Force

Write-Host ""
Write-Host "  Done! Single release file:" -ForegroundColor Green
Write-Host "    $SetupExe" -ForegroundColor Yellow
Write-Host ""
Get-Item $SetupExe | Format-Table Name, @{N='MB';E={[math]::Round($_.Length/1MB,2)}} -AutoSize
Write-Host "  Test: double-click the setup exe and click Next through the wizard." -ForegroundColor Cyan
