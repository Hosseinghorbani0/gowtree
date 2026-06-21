#Requires -Version 5.1
<#
.SYNOPSIS
    Build all Windows release artifacts for gowtree.
.OUTPUTS
    release/gowtree.exe
    release/gowtree-setup-1.4.0.exe
    release/gowtree-portable-1.4.0.zip
#>
param(
    [switch]$SkipInstaller
)

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
Set-Location $Root

$Version = "1.4.0"
$ReleaseDir = Join-Path $Root "release"
$PortableExe = Join-Path $ReleaseDir "gowtree.exe"
$SetupExe = Join-Path $ReleaseDir "gowtree-setup-$Version.exe"
$ZipFile = Join-Path $ReleaseDir "gowtree-portable-$Version.zip"

Write-Host ""
Write-Host "  gowtree release build v$Version" -ForegroundColor Green
Write-Host "  --------------------------------" -ForegroundColor DarkGray
Write-Host ""

New-Item -ItemType Directory -Force -Path $ReleaseDir | Out-Null

Write-Host "[1/4] Building gowtree.exe..." -ForegroundColor Cyan
$ldflags = "-s -w -X github.com/hosseinghorbani0/gowtree/internal/metadata.Version=v$Version"
go build -ldflags $ldflags -o $PortableExe ./cmd/gowtree
if ($LASTEXITCODE -ne 0) { throw "go build failed" }
Copy-Item $PortableExe (Join-Path $Root "gowtree.exe") -Force
$sizeKB = [math]::Round((Get-Item $PortableExe).Length / 1KB)
Write-Host "      -> release\gowtree.exe ($sizeKB KB)" -ForegroundColor DarkGreen

Write-Host "[2/4] Running tests..." -ForegroundColor Cyan
go test ./...
if ($LASTEXITCODE -ne 0) { throw "tests failed" }
Write-Host "      -> all tests passed" -ForegroundColor DarkGreen

Write-Host "[3/4] Creating portable ZIP..." -ForegroundColor Cyan
$zipStage = Join-Path $env:TEMP "gowtree-portable-$Version"
if (Test-Path $zipStage) { Remove-Item $zipStage -Recurse -Force }
New-Item -ItemType Directory -Force -Path $zipStage | Out-Null
Copy-Item $PortableExe $zipStage
Copy-Item (Join-Path $Root "install.bat") $zipStage
Copy-Item (Join-Path $Root "install.ps1") $zipStage
Copy-Item (Join-Path $Root "LICENSE") $zipStage
if (Test-Path $ZipFile) { Remove-Item $ZipFile -Force }
Compress-Archive -Path (Join-Path $zipStage "*") -DestinationPath $ZipFile -Force
Remove-Item $zipStage -Recurse -Force
Write-Host "      -> release\gowtree-portable-$Version.zip" -ForegroundColor DarkGreen

if ($SkipInstaller) {
    Write-Host "[4/4] Skipped installer (-SkipInstaller)" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "  Done! Files in: $ReleaseDir" -ForegroundColor Green
    Get-ChildItem $ReleaseDir | Format-Table Name, @{N='Size';E={"{0:N0} KB" -f ($_.Length/1KB)}} -AutoSize
    exit 0
}

Write-Host "[4/4] Building setup wizard (Inno Setup)..." -ForegroundColor Cyan
$iscc = @(
    "C:\InnoSetup6\ISCC.exe",
    "${env:ProgramFiles(x86)}\Inno Setup 6\ISCC.exe",
    "$env:ProgramFiles\Inno Setup 6\ISCC.exe"
) | Where-Object { Test-Path $_ } | Select-Object -First 1

if (-not $iscc) {
    throw "Inno Setup not found. Run: winget install JRSoftware.InnoSetup"
}

& $iscc "installer\gowtree.iss"
if ($LASTEXITCODE -ne 0) { throw "Inno Setup build failed" }

$builtSetup = Get-ChildItem "installer\output\gowtree-setup-*.exe" | Sort-Object LastWriteTime -Descending | Select-Object -First 1
if (-not $builtSetup) { throw "Setup exe not found in installer\output" }
Copy-Item $builtSetup.FullName $SetupExe -Force
Write-Host "      -> release\gowtree-setup-$Version.exe" -ForegroundColor DarkGreen

Write-Host ""
Write-Host "  Done! Test the installer:" -ForegroundColor Green
Write-Host "    .\release\gowtree-setup-$Version.exe" -ForegroundColor Yellow
Write-Host ""
Get-ChildItem $ReleaseDir | Format-Table Name, @{N='Size';E={"{0:N0} KB" -f ($_.Length/1KB)}} -AutoSize
