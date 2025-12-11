@echo off
echo Building AlternateDNS for Windows (portable)...
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
) else (
    echo Compiling icon resource...
    windres -o icon.syso icon.rc
    if errorlevel 1 (
        echo Warning: Failed to compile icon resource. Building without icon.
    ) else (
        echo Icon resource compiled successfully.
    )
)

echo.
echo Starting build (this may take a minute)...
echo.

REM Create dist folder if it doesn't exist
if not exist "dist" mkdir dist

set CGO_ENABLED=1
go build -ldflags="-s -w -H windowsgui" -o dist\AlternateDNS.exe

REM Clean up resource file after build
if exist "icon.syso" del "icon.syso"
if errorlevel 1 goto :build_failed

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
