#Requires -Version 5.1
<#
.SYNOPSIS
    Build gowtree.exe and optionally the Windows installer (Inno Setup).
.EXAMPLE
    .\scripts\build-installer.ps1
    .\scripts\build-installer.ps1 -SkipInstaller
#>
param(
    [switch]$SkipInstaller,
    [string]$Output = "gowtree.exe"
)

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
Set-Location $Root

Write-Host ""
Write-Host "  gowtree build" -ForegroundColor Green
Write-Host "  -----------------" -ForegroundColor DarkGray
Write-Host ""

Write-Host "[1/3] Building release binary..." -ForegroundColor Cyan
$ldflags = "-s -w -X github.com/hosseinghorbani0/gowtree/internal/metadata.Version=v1.4.0"
go build -ldflags $ldflags -o $Output ./cmd/gowtree
if ($LASTEXITCODE -ne 0) { throw "go build failed" }
$sizeKB = [math]::Round((Get-Item $Output).Length / 1KB)
Write-Host "      -> $Output ($sizeKB KB)" -ForegroundColor DarkGreen

Write-Host "[2/3] Running tests..." -ForegroundColor Cyan
go test ./...
if ($LASTEXITCODE -ne 0) { throw "tests failed" }
Write-Host "      -> all tests passed" -ForegroundColor DarkGreen

if ($SkipInstaller) {
    Write-Host "[3/3] Skipped installer (-SkipInstaller)" -ForegroundColor Yellow
    exit 0
}

Write-Host "[3/3] Building Windows installer..." -ForegroundColor Cyan
$iscc = @(
    "${env:ProgramFiles(x86)}\Inno Setup 6\ISCC.exe",
    "$env:ProgramFiles\Inno Setup 6\ISCC.exe",
    "iscc"
) | Where-Object { Test-Path $_ -ErrorAction SilentlyContinue } | Select-Object -First 1

if (-not $iscc) {
    $isccCmd = Get-Command iscc -ErrorAction SilentlyContinue
    if ($isccCmd) { $iscc = $isccCmd.Source }
}

if (-not $iscc) {
    Write-Host ""
    Write-Host "  WARNING: Inno Setup not found." -ForegroundColor Yellow
    Write-Host "    Install from: https://jrsoftware.org/isinfo.php" -ForegroundColor DarkGray
    Write-Host "    Binary is ready: $Output" -ForegroundColor DarkGray
    Write-Host "    Or run: .\install.bat" -ForegroundColor DarkGray
    exit 0
}

& $iscc "installer\gowtree.iss"
if ($LASTEXITCODE -ne 0) { throw "Inno Setup build failed" }

$setup = Get-ChildItem "installer\output\gowtree-setup-*.exe" | Sort-Object LastWriteTime -Descending | Select-Object -First 1
Write-Host ""
Write-Host "  Done!" -ForegroundColor Green
Write-Host "     Binary:   $Output" -ForegroundColor White
if ($setup) {
    Write-Host "     Installer: $($setup.FullName)" -ForegroundColor White
}
Write-Host ""
