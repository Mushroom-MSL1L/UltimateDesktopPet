package file

import (
	"encoding/base64"
	"fmt"

	pp "github.com/Mushroom-MSL1L/UltimateDesktopPet/pkg/print"

	"os"
	"path/filepath"
	"sort"
	"strings"
)

type SpriteTool struct {
	StaticAssetPath      string
	DefaultImageFolder   string
	SpecifiedImageFolder string
	URISchemePrefix      string
	Suffix               string
}

func NewSpriteTool(s SpriteTool) SpriteTool {
	ss := SpriteTool{
		StaticAssetPath:      "",
		DefaultImageFolder:   "default",
		SpecifiedImageFolder: "",
		URISchemePrefix:      "data:image/png;base64,",
		Suffix:               ".png",
	}
	if s.DefaultImageFolder == "" {
		s.DefaultImageFolder = ss.DefaultImageFolder
	}
	if s.URISchemePrefix == "" {
		s.URISchemePrefix = ss.URISchemePrefix
	}
	if s.Suffix == "" {
		s.Suffix = ss.Suffix
	}
	return s
}

func (s SpriteTool) LoadFrameFromDir(file string) (string, error) {
	fsPath := filepath.ToSlash(filepath.Join(s.StaticAssetPath, s.SpecifiedImageFolder))
	defaultPath := filepath.ToSlash(filepath.Join(s.StaticAssetPath, s.DefaultImageFolder))
	var frame string
	var err error
	/* try specified path */
	if frame, err = s.readFrameFromDisk(fsPath, file); err == nil {
		return frame, nil
	}
	pp.Warn(pp.File, "LoadFrameFromDir: load specified frame %s failed: %v", fsPath, err)
	/* try default path */
	if frame, err = s.readFrameFromDisk(defaultPath, file); err == nil {
		return frame, nil
	}

	return "", fmt.Errorf("LoadFrameFromDir: default path %s doesn't contain assets %s: %v", defaultPath, file, err)
}

func (s SpriteTool) readFrameFromDisk(path string, file string) (string, error) {
	targetFile := file + s.Suffix
	files, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), s.Suffix) || (f.Name() != targetFile) {
			continue
		}
		data, err := os.ReadFile(filepath.Join(path, f.Name()))
		if err != nil {
			continue
		}
		return (s.URISchemePrefix + base64.StdEncoding.EncodeToString(data)), nil
	}
	return "", fmt.Errorf("")
}

func (s SpriteTool) LoadFramesFromDir(animationType string) ([]string, error) {
	fsPath := filepath.ToSlash(filepath.Join(s.StaticAssetPath, s.SpecifiedImageFolder, animationType))
	defaultPath := filepath.ToSlash(filepath.Join(s.StaticAssetPath, s.DefaultImageFolder, animationType))
	var localFrames []string
	var err error
	/* try specified path */
	if localFrames, err = s.readFramesFromDisk(fsPath); err == nil && len(localFrames) > 0 {
		return localFrames, nil
	}
	pp.Warn(pp.File, "LoadFramesFromDir: load specified frames %s failed: %v", fsPath, err)
	/* try default path */
	if localFrames, err := s.readFramesFromDisk(defaultPath); err == nil && len(localFrames) > 0 {
		return localFrames, nil
	}

	return nil, fmt.Errorf("LoadFramesFromDir: default path doesn't contain animationType %s assets under path %s: %v", animationType, defaultPath, err)
}

func (s SpriteTool) readFramesFromDisk(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var frames []string
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), s.Suffix) {
			continue
		}
		data, err := os.ReadFile(filepath.Join(path, f.Name()))
		if err != nil {
			continue
		}
		frames = append(frames, s.URISchemePrefix+base64.StdEncoding.EncodeToString(data))
	}
	sortFrames(&frames)
	return frames, nil
}

func sortFrames(frames *[]string) {
	sort.SliceStable(*frames, func(i, j int) bool {
		return i < j
	})
}
