#!/usr/bin/env bash
set -e

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"

GO_VERSION="1.24"
DESKTOP_PET_DIR="$ROOT_DIR/app/desktop_pet"
DESKTOP_PET_CONFIG_DIR="$DESKTOP_PET_DIR/configs/system.yaml"
SYNC_SERVER_DIR="$ROOT_DIR/app/sync_server"
SYNC_SERVER_CONFIG_DIR="$SYNC_SERVER_DIR/configs/server.yaml"

echo "=== UltimateDesktopPet runner ==="

# ---------- Shared Utilities ----------
run_or_reset() {
    local cmd="$1"
    local reset_cmd="$2"

    echo ">> Running: $cmd"
    if ! eval "$cmd"; then
        echo "!! Command failed, running reset"
        eval "$reset_cmd"
        echo ">> Retry: $cmd"
        eval "$cmd"
    fi
}

# ---------- Update Go dependencies ----------
update_deps() {
    echo "=== Updating Go dependencies ==="
    cd "$ROOT_DIR"

    if [ ! -f "go.mod" ]; then
        echo "No go.mod found, initializing module with go version $GO_VERSION..."
        go mod init github.com/Mushroom-MSL1L/UltimateDesktopPet 
        go mod edit -go="$GO_VERSION"
    fi

    go get -u ./...
    go mod tidy
}

# ---------- Desktop Pet ----------
desktop_pet_dev() {
    echo "=== Desktop Pet: wails dev ==="
    cd "$DESKTOP_PET_DIR"
    run_or_reset "DESKTOP_PET_CONFIG_DIR=\"$DESKTOP_PET_CONFIG_DIR\" wails dev" "./reset.sh"
}

desktop_pet_build() {
    echo "=== Desktop Pet: wails build & create shortcut ==="

    OS_TYPE="$(uname -s)"
    case "$OS_TYPE" in
        Darwin)
            # macOS: .command
            SHORTCUT_PATH="$ROOT_DIR/UltimateDesktopPet.command"
            cat > "$SHORTCUT_PATH" <<EOF
#!/bin/bash
DESKTOP_PET_CONFIG_DIR="$DESKTOP_PET_CONFIG_DIR"
cd "$DESKTOP_PET_DIR" || exit 1
wails dev
EOF
            chmod 777 "$SHORTCUT_PATH"
            echo "macOS shortcut created at $SHORTCUT_PATH"
            ;;
        Linux)
            # Linux: .desktop
            SHORTCUT_PATH="$HOME/Desktop/UltimateDesktopPet.desktop"
            cat > "$SHORTCUT_PATH" <<EOF
[Desktop Entry]
Name=UltimateDesktopPet
Comment=Run Desktop Pet
Exec=bash -c 'DESKTOP_PET_CONFIG_DIR="$DESKTOP_PET_CONFIG_DIR" cd "$DESKTOP_PET_DIR" && wails dev'
Icon=
Terminal=true
Type=Application
EOF
            chmod +x "$SHORTCUT_PATH"
            echo "Linux desktop shortcut created at $SHORTCUT_PATH"
            ;;
        CYGWIN*|MINGW*|MSYS*|Windows_NT)
            # Windows: .bat
            SHORTCUT_PATH="$ROOT_DIR/UltimateDesktopPet.bat"
            cat > "$SHORTCUT_PATH" <<EOF
@echo off
set DESKTOP_PET_CONFIG_DIR="$DESKTOP_PET_CONFIG_DIR"
cd /d "$DESKTOP_PET_DIR"
wails dev
pause
EOF
            echo "Windows batch shortcut created at $SHORTCUT_PATH"
            ;;
        *)
            echo "Unsupported OS: $OS_TYPE"
            return 1
            ;;
    esac
}

# ---------- Sync Server ----------
sync_server_run() {
    echo "=== Sync Server: swagger + run ==="
    cd "$SYNC_SERVER_DIR"

    echo ">> swag fmt"
    swag fmt

    echo ">> swag init"
    swag init

    echo ">> go run main.go -config=$SYNC_SERVER_CONFIG_DIR"
    go run main.go -config=$(SYNC_SERVER_CONFIG_DIR)
}

# ---------- CLI ----------
case "$1" in
    update)
        update_deps
        ;;
    dev)
        desktop_pet_dev
        ;;
    build)
        desktop_pet_build
        ;;
    server)
        sync_server_run
        ;;
    all)
        desktop_pet_build
        sync_server_run
        ;;
    *)
        echo "Usage:"
        echo "  ./run_me.sh update   # update all related resources" 
        echo "  ./run_me.sh dev      # wails dev (auto reset if fail)"
        echo "  ./run_me.sh build    # wails build with a shortcut (auto reset if fail)"
        echo "  ./run_me.sh server   # swag + go run sync server"
        echo "  ./run_me.sh all      # build desktop_pet + run server"
        exit 1
        ;;
esac
