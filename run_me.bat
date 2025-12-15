@echo off
REM === UltimateDesktopPet runner for Windows ===

set "ROOT_DIR=%~dp0"
set "DESKTOP_PET_DIR=%ROOT_DIR%app\desktop_pet"
set "DESKTOP_PET_CONFIG_DIR=%DESKTOP_PET_DIR%\configs\system.yaml"
set "SYNC_SERVER_DIR=%ROOT_DIR%app\sync_server"
set "SYNC_SERVER_CONFIG_DIR=%SYNC_SERVER_DIR%\configs\server.yaml"

REM Go path (If not include in $PATH, you can uncomment the following line and set it manually)
REM set PATH=C:\Program Files\Go\bin;%PATH%

REM ---------- CLI ----------
if "%~1"=="" goto :usage
if /I "%~1"=="update" call :update_deps & goto :eof
if /I "%~1"=="dev" call :desktop_pet_dev & goto :eof
if /I "%~1"=="build" call :desktop_pet_build & goto :eof
if /I "%~1"=="server" call :sync_server_run & goto :eof
if /I "%~1"=="all" (
    call :desktop_pet_build
    call :sync_server_run
    goto :eof
)
goto :usage

REM ---------- Update Go dependencies ----------
:update_deps
echo === Updating Go dependencies ===
cd /d "%ROOT_DIR%"
if not exist "go.mod" (
    echo No go.mod found, initializing module...
    go mod init github.com/Mushroom-MSL1L/UltimateDesktopPet
    go mod edit -go=1.24
)
call go get -u ./...
call go mod tidy
goto :eof

REM ---------- Desktop Pet ----------
:desktop_pet_dev
echo === Desktop Pet: wails dev ===
cd /d "%DESKTOP_PET_DIR%"

REM Check wails existence
where wails >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo wails not found, installing...
    call go install github.com/wailsapp/wails/v2/cmd/wails@latest
)

call wails dev
goto :eof

:desktop_pet_build
echo === Desktop Pet: wails build ^& create shortcut ===
cd /d "%DESKTOP_PET_DIR%"

REM Check wails existence
where wails >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo wails not found, installing...
    call go install github.com/wailsapp/wails/v2/cmd/wails@latest
)

set SHORTCUT_PATH=%ROOT_DIR%UltimateDesktopPet.bat
(
    echo @echo off
    if not exist "%DESKTOP_PET_CONFIG_DIR%" (
        echo Config file not found: %DESKTOP_PET_CONFIG_DIR%
        exit /b 1
    )
    echo set "DESKTOP_PET_CONFIG_DIR=%DESKTOP_PET_CONFIG_DIR%"
    echo cd /d "%DESKTOP_PET_DIR%"
    echo call wails dev
    echo pause
) > "%SHORTCUT_PATH%"

echo Shortcut created at %SHORTCUT_PATH%
goto :eof

REM ---------- Sync Server ----------
:sync_server_run
echo === Sync Server: swag + run ===
cd /d "%SYNC_SERVER_DIR%"
REM check swag existence
where swag >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo swag not found, installing...
    call go install github.com/swaggo/swag/cmd/swag@latest
)
call swag fmt
call swag init
call go run main.go -config="%SYNC_SERVER_CONFIG_DIR%"
goto :eof

:usage
echo Usage:
echo   run_me.bat update   - update all resources
echo   run_me.bat dev      - run wails dev
echo   run_me.bat build    - wails build + shortcut
echo   run_me.bat server   - run sync server
echo   run_me.bat all      - build desktop_pet + run server
exit /b 1
