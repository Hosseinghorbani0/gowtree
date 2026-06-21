#Requires -Version 5.1
<#
.SYNOPSIS
    Install gowtree to %USERPROFILE%\bin and update PATH.
.DESCRIPTION
    Interactive Windows installer with optional GUI progress.
#>
param(
    [switch]$Silent,
    [switch]$NoPath,
    [string]$Source
)

$ErrorActionPreference = "Stop"

function Write-Step($msg) {
    if (-not $Silent) { Write-Host "  → $msg" -ForegroundColor Cyan }
}

function Show-Banner {
    if ($Silent) { return }
    Write-Host ""
    Write-Host "  ╔══════════════════════════════════════╗" -ForegroundColor Green
    Write-Host "  ║   🌳  gowtree — Windows Installer    ║" -ForegroundColor Green
    Write-Host "  ╚══════════════════════════════════════╝" -ForegroundColor Green
    Write-Host ""
}

Show-Banner

if (-not $Source) {
    $Source = Join-Path $PSScriptRoot 'gowtree.exe'
    if (-not (Test-Path $Source)) {
        $built = Join-Path $PSScriptRoot '..\gowtree.exe'
        if (Test-Path $built) { $Source = $built }
    }
}

if (-not (Test-Path $Source)) {
    Write-Host "  ✗ gowtree.exe not found." -ForegroundColor Red
    Write-Host "    Build first:  go build -o gowtree.exe ./cmd/gowtree" -ForegroundColor DarkGray
    Write-Host "    Or run:       .\scripts\build-installer.ps1" -ForegroundColor DarkGray
    exit 1
}

$installDir = Join-Path $env:USERPROFILE 'bin'
Write-Step "Install directory: $installDir"

if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    Write-Step "Created $installDir"
}

$dest = Join-Path $installDir 'gowtree.exe'
Copy-Item -Path $Source -Destination $dest -Force
Write-Step "Copied gowtree.exe"

if (-not $NoPath) {
    $pathDirs = ($env:Path -split ';') | ForEach-Object { $_.Trim() } | Where-Object { $_ -ne '' }
    if ($pathDirs -contains $installDir) {
        Write-Step "PATH already contains $installDir"
    } else {
        Write-Step "Adding to user PATH..."
        $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
        if ([string]::IsNullOrEmpty($userPath)) {
            $newPath = $installDir
        } else {
            $newPath = $userPath.TrimEnd(';') + ';' + $installDir
        }
        [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
        $env:Path = $env:Path + ';' + $installDir
        Write-Step "PATH updated (user scope)"
    }
}

if (-not $Silent) {
    Write-Host ""
    Write-Host "  ✅ gowtree installed successfully!" -ForegroundColor Green
    Write-Host ""
    Write-Host "  Open a NEW terminal and try:" -ForegroundColor White
    Write-Host "    gowtree -h" -ForegroundColor Yellow
    Write-Host "    gowtree -a -s -L 2 --icons" -ForegroundColor Yellow
    Write-Host ""
}

# Verify
$version = & $dest -version 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Step "Verified: $version"
}

exit 0
