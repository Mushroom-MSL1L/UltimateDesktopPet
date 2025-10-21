package file

import (
	"log"
	"os"
	"path/filepath"
)

func SafeCreateFile(filePath string) {
	info, err := os.Stat(filePath)
	if err == nil {
		if !info.IsDir() {
			return
		}
		log.Fatalf("path %s exists but is a directory", filePath)
	}
	if !os.IsNotExist(err) {
		log.Fatalf("failed to check file %s: %v", filePath, err)
	}

	dir := filepath.Dir(filePath)
	SafeCreateDir(dir)

	f, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("failed to create file %s: %v", filePath, err)
	}
	f.Close()
}

func SafeCreateDir(dir string) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalln("failed to create directory %s: %v", dir, err)
	}
}

func SafeOpenFile(filePath string) *os.File {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("file %s does not exist: %v", filePath, err)
		}
		log.Fatalf("failed to check file %s: %v", filePath, err)
	}

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file %s: %v", filePath, err)
	}
	return f
}
