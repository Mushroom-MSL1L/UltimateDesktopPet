package file

import (
	"os"
	"path/filepath"

	pp "github.com/Mushroom-MSL1L/UltimateDesktopPet/desktop_pet/pkg/print"
)

func SafeCreateFile(filePath string) {
	info, err := os.Stat(filePath)
	if err == nil {
		if !info.IsDir() {
			return
		}
		pp.Fatal(pp.File, "path %s exists but is a directory", filePath)
	}
	if !os.IsNotExist(err) {
		pp.Fatal(pp.File, "failed to check file %s: %v", filePath, err)
	}

	dir := filepath.Dir(filePath)
	SafeCreateDir(dir)
	pp.Warn(pp.File, "file not exist, so safe create %s", filePath)

	f, err := os.Create(filePath)
	if err != nil {
		pp.Fatal(pp.File, "failed to create file %s: %v", filePath, err)
	}
	f.Close()
}

func SafeCreateDir(dir string) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		pp.Fatal(pp.File, "failed to create directory %s: %v", dir, err)
	}
}

func SafeOpenFile(filePath string) *os.File {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			pp.Fatal(pp.File, "file %s does not exist: %v", filePath, err)
		}
		pp.Fatal(pp.File, "failed to check file %s: %v", filePath, err)
	}

	f, err := os.Open(filePath)
	if err != nil {
		pp.Fatal(pp.File, "failed to open file %s: %v", filePath, err)
	}
	return f
}
