package app

import (
	"embed"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "UltimateDesktopPet/internal/pet"
	_ "UltimateDesktopPet/internal/system"
)

var embeddedPetAssets embed.FS

var petSpriteDataURI string

func loadDefaultSprite() {
	defPath := filepath.ToSlash("assets/petImages/default/cat_move.gif")
	if data, err := embeddedPetAssets.ReadFile(defPath); err == nil && len(data) > 0 {
		petSpriteDataURI = "data:image/gif;base64," + base64.StdEncoding.EncodeToString(data)
		return
	}

	relPath := filepath.Join("..", "..", "assets", "petImages", "default", "cat_move.gif")
	data, err := os.ReadFile(relPath)
	if err != nil {
		if exe, eerr := os.Executable(); eerr == nil {
			exeDir := filepath.Dir(exe)
			tryPath := filepath.Join(exeDir, "..", "..", "assets", "petImages", "default", "cat_move.gif")
			if d2, e2 := os.ReadFile(tryPath); e2 == nil {
				data = d2
				err = nil
			}
		}
	}

	if err == nil && len(data) > 0 {
		petSpriteDataURI = "data:image/gif;base64," + base64.StdEncoding.EncodeToString(data)
	} else {
		petSpriteDataURI = ""
	}
}

func loadSpriteByName(name string) (string, error) {
	// normalize name
	name = filepath.ToSlash(strings.TrimPrefix(name, "/"))

	// try embedded first; embedded paths use the same relative path as in main.go embed
	embeddedPath := filepath.ToSlash("assets/petImages/" + name)
	if data, err := embeddedPetAssets.ReadFile(embeddedPath); err == nil && len(data) > 0 {
		return "data:image/gif;base64," + base64.StdEncoding.EncodeToString(data), nil
	}

	// fallback filesystem (relative to repo)
	fsPath := filepath.Join("..", "..", "assets", "petImages", filepath.FromSlash(name))
	if data, err := os.ReadFile(fsPath); err == nil && len(data) > 0 {
		return "data:image/gif;base64," + base64.StdEncoding.EncodeToString(data), nil
	}

	// fallback relative to executable
	if exe, eerr := os.Executable(); eerr == nil {
		exeDir := filepath.Dir(exe)
		tryPath := filepath.Join(exeDir, "..", "..", "assets", "petImages", filepath.FromSlash(name))
		if data, err := os.ReadFile(tryPath); err == nil && len(data) > 0 {
			return "data:image/gif;base64," + base64.StdEncoding.EncodeToString(data), nil
		}
	}

	return "", fmt.Errorf("sprite not found: %s", name)
}
