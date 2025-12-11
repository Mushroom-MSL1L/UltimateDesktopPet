#!/usr/bin/env bash

set -e

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"

DB_DIR="$ROOT_DIR/assets/db"
BUILD_BIN="$ROOT_DIR/build/bin"
FRONTEND_DIST="$ROOT_DIR/frontend/dist"
WAILSJS="$ROOT_DIR/frontend/wailsjs"
GO_MOD="$ROOT_DIR/go.mod"
GO_SUM="$ROOT_DIR/go.sum"

echo "=== Reset 開始 ==="

###################################################
# 1. 刪除 udp.db
###################################################
if [ -f "$DB_DIR/udp.db" ]; then
    echo "刪除：$DB_DIR/udp.db"
    rm "$DB_DIR/udp.db"
fi

###################################################
# 2. 刪除 static_assets.db
###################################################
if [ -f "$DB_DIR/static_assets.db" ]; then
    echo "刪除：$DB_DIR/static_assets.db"
    rm "$DB_DIR/static_assets.db"
fi

###################################################
# 4. 刪除 desktop_pet/build/bin
###################################################
if [ -d "$BUILD_BIN" ]; then
    echo "刪除資料夾：$BUILD_BIN"
    rm -rf "$BUILD_BIN"
fi

###################################################
# 5. 刪除 frontend/dist
###################################################
if [ -d "$FRONTEND_DIST" ]; then
    echo "刪除資料夾：$FRONTEND_DIST"
    rm -rf "$FRONTEND_DIST"
fi

###################################################
# 6. 刪除 frontend/wailsjs
###################################################
if [ -d "$WAILSJS" ]; then
    echo "刪除資料夾：$WAILSJS"
    rm -rf "$WAILSJS"
fi

# ###################################################
# # 7. 刪除 go.mod / go.sum
# ###################################################
# if [ -f "$GO_MOD" ]; then
#     echo "刪除：go.mod"
#     rm "$GO_MOD"
# fi

# if [ -f "$GO_SUM" ]; then
#     echo "刪除：go.sum"
#     rm "$GO_SUM"
# fi

# ###################################################
# # 8. 在 desktop_pet/ 重新 init module
# ###################################################
# cd "$ROOT_DIR"

# echo "重新初始化 Go Module：go mod init github.com/Mushroom-MSL1L/UltimateDesktopPet/desktop_pet"
# go mod init github.com/Mushroom-MSL1L/UltimateDesktopPet/desktop_pet
# go mod tidy

# echo "=== Reset 完成 ==="
