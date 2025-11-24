@echo off
echo Building AlternateDNS for Windows (portable)...
echo.
echo Checking for C compiler (required for Fyne GUI)...
echo.

gcc --version >nul 2>&1
if errorlevel 1 goto :no_gcc

echo C compiler found!
echo.
echo Starting build (this may take a minute)...
echo.
set CGO_ENABLED=1
go build -ldflags="-s -w -H windowsgui" -o AlternateDNS.exe
if errorlevel 1 goto :build_failed

echo Build successful! AlternateDNS.exe created.
echo File size:
dir AlternateDNS.exe | findstr AlternateDNS.exe
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