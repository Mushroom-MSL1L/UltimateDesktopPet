package app

import (
	"encoding/base64"
	"fmt"

	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "UltimateDesktopPet/internal/pet"
)

const staticPetAssetPath = "assets/petImages/"
const defaultPetImageFolder = "default"
const defaultPetAnimationType = "stand"
const pngSuffix = ".png"
const uriSchemePrefix = "data:image/png;base64,"

func loadDefaultFrames(petPath string) ([]string, error) {
	frames, err := loadFramesFromDir(petPath, defaultPetAnimationType)
	if err == nil {
		return frames, nil
	}
	return nil, err
}

func loadFramesFromDir(petPath string, animationType string) ([]string, error) {
	name := filepath.ToSlash(filepath.Join(petPath, animationType))
	errMessage := ""
	fsPath := filepath.Join(staticPetAssetPath, filepath.FromSlash(name))
	if localFrames, err := readFramesFromDisk(fsPath); err == nil && len(localFrames) > 0 {
		return localFrames, nil
	}
	errMessage += fmt.Sprintf("no png frames in embedded or fs path: %s; ", fsPath)

	return nil, fmt.Errorf("loadFramesFromDir(%s): %s", name, errMessage)
}

func readFramesFromDisk(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var frames []string
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), pngSuffix) {
			continue
		}
		data, err := os.ReadFile(filepath.Join(path, f.Name()))
		if err != nil {
			continue
		}
		frames = append(frames, uriSchemePrefix+base64.StdEncoding.EncodeToString(data))
	}
	sortFrames(&frames)
	return frames, nil
}

func sortFrames(frames *[]string) {
	sort.SliceStable(*frames, func(i, j int) bool {
		return i < j
	})
}
