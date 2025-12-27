@echo off
setlocal enabledelayedexpansion

:: get short commit hash
for /f %%i in ('git rev-parse HEAD') do set COMMIT=%%i
set COMMIT=%COMMIT:~0,8%

:: get tag for this commit (if any)
for /f "tokens=*" %%i in ('git describe --exact-match --abbrev=0 --tags %COMMIT% 2^>nul') do set TAG=%%i

if "%TAG%"=="" (
    set VERSION=%COMMIT%
) else (
    set VERSION=%TAG%
)

:: check for uncommitted changes
for /f "tokens=*" %%i in ('git diff --shortstat 2^>nul') do set DIFF=%%i

if not "%DIFF%"=="" (
    set VERSION=%VERSION%-dirty
)

echo %VERSION%
endlocal
