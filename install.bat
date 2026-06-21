@echo off
setlocal EnableExtensions
title gowtree Installer

set "SCRIPT_DIR=%~dp0"
set "PS_SCRIPT=%SCRIPT_DIR%install.ps1"

echo.
echo   ========================================
echo     gowtree - Windows Installer
echo   ========================================
echo.

if not exist "%PS_SCRIPT%" (
  echo [ERROR] install.ps1 not found in %SCRIPT_DIR%
  pause
  exit /b 1
)

where gowtree.exe >nul 2>&1
if exist "%SCRIPT_DIR%gowtree.exe" goto :run
if exist "%SCRIPT_DIR%..\gowtree.exe" goto :run

echo [INFO] gowtree.exe not found — building from source...
where go >nul 2>&1
if errorlevel 1 (
  echo [ERROR] Go is not installed. Download from https://go.dev/dl/
  echo         Or download the installer from GitHub Releases.
  pause
  exit /b 1
)
pushd "%SCRIPT_DIR%"
go build -o gowtree.exe ./cmd/gowtree
if errorlevel 1 (
  echo [ERROR] Build failed.
  popd
  pause
  exit /b 1
)
popd
echo [OK] Built gowtree.exe
echo.

:run
powershell -NoProfile -ExecutionPolicy Bypass -File "%PS_SCRIPT%"
set "RC=%ERRORLEVEL%"
if %RC% neq 0 (
  echo.
  echo [ERROR] Installation failed.
  pause
  exit /b %RC%
)
echo.
pause
exit /b 0
