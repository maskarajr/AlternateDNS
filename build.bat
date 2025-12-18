@echo off
REM Build script for AlternateDNS
REM Usage: build.bat [version] [commit] [date]
REM Example: build.bat 1.0.0 abc123 2025-01-01

setlocal

REM Get version info from arguments or use defaults
set VERSION=%1
if "%VERSION%"=="" set VERSION=dev

set GIT_COMMIT=%2
if "%GIT_COMMIT%"=="" set GIT_COMMIT=unknown

set BUILD_DATE=%3
if "%BUILD_DATE%"=="" (
    REM Use PowerShell to get date in YYYY-MM-DD format
    for /f "delims=" %%I in ('powershell -Command "Get-Date -Format 'yyyy-MM-dd'"') do set BUILD_DATE=%%I
)

echo Building AlternateDNS for Windows (portable)...
echo Version: %VERSION%
echo Commit: %GIT_COMMIT%
echo Build Date: %BUILD_DATE%
echo.
echo Checking for C compiler (required for Fyne GUI)...
echo.

gcc --version >nul 2>&1
if errorlevel 1 goto :no_gcc

echo C compiler found!
echo.

REM Check for windres (resource compiler) and compile icon
where windres >nul 2>&1
if errorlevel 1 (
    echo Warning: windres not found. Executable will build without Explorer icon.
    echo To fix: windres comes with MinGW-w64. Make sure it's in your PATH.
    echo The app will still work, just without the icon in Explorer.
    goto :skip_icon
)

REM Check if icon files exist
if not exist "icon.ico" (
    echo Warning: icon.ico not found. Executable will build without Explorer icon.
    goto :skip_icon
)
if not exist "icon.rc" (
    echo Warning: icon.rc not found. Executable will build without Explorer icon.
    goto :skip_icon
)

echo Compiling icon resource...
windres -o icon.syso icon.rc
if errorlevel 1 (
    echo Warning: Failed to compile icon resource. Building without icon.
    if exist "icon.syso" del "icon.syso"
    goto :skip_icon
)

if not exist "icon.syso" (
    echo Warning: icon.syso was not created. Building without icon.
    goto :skip_icon
)

echo Icon resource compiled successfully.

:skip_icon

echo.
echo Starting build (this may take a minute)...
echo.

REM Create dist folder if it doesn't exist
if not exist "dist" mkdir dist

set CGO_ENABLED=1
go build -ldflags="-s -w -H windowsgui -X main.Version=%VERSION% -X main.GitCommit=%GIT_COMMIT% -X main.BuildDate=%BUILD_DATE%" -o dist\AlternateDNS.exe

REM Check if build succeeded
if errorlevel 1 goto :build_failed

REM Clean up resource file after build
if exist "icon.syso" (
    del "icon.syso"
    echo Cleaned up icon.syso
)

REM Copy default_config.yaml to dist folder for reference
if exist "default_config.yaml" (
    copy /Y "default_config.yaml" "dist\default_config.yaml" >nul
    echo Copied default_config.yaml to dist folder for reference
)

echo.
echo Build successful! dist\AlternateDNS.exe created.
echo File size:
dir dist\AlternateDNS.exe | findstr AlternateDNS.exe
echo.
pause
exit /b 0

:no_gcc
echo ERROR: C compiler (gcc) not found!
echo.
echo Fyne requires CGO which needs a C compiler.
echo.
echo To fix this, install MinGW-w64:
echo   1. Download from: https://www.mingw-w64.org/downloads/
echo   2. Or use MSYS2: https://www.msys2.org/
echo   3. Or use Chocolatey: choco install mingw
echo.
echo After installing, add gcc to your PATH and try again.
echo.
pause
exit /b 1

:build_failed
echo Build failed!
echo.
echo Check the error messages above.
echo.
pause
exit /b 1
