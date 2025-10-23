package configs

import (
	"log"
	"os"

	"path/filepath"

	"gopkg.in/yaml.v2"
)

type ConfigProvider[T any] interface {
	DefaultConfig() *T
}

func LoadConfig[T any, P ConfigProvider[T]](path string, provider P) *T {
	defaultCfg := provider.DefaultConfig()
	fileCfg := readConfigFile(path, provider)
	mergedCfg := mergeStructs(defaultCfg, fileCfg)
	writeConfigFile(path, mergedCfg)
	return mergedCfg
}

func readConfigFile[T any, P ConfigProvider[T]](path string, provider P) *T {
	defaultCfg := provider.DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			createConfigFile[T](path, defaultCfg)
			return defaultCfg
		}
		log.Fatalf("cannot read config file %s: %v", path, err)
	}

	var fileCfg *T
	if err := yaml.Unmarshal(data, fileCfg); err != nil {
		log.Fatalf("yaml file %s unmarshal failed: %v", path, err)
	}
	return fileCfg
}

func createConfigFile[T any](path string, cfg *T) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("failed to create directory %s: %v", dir, err)
	}

	writeConfigFile(path, cfg)
}

func writeConfigFile[T any](path string, cfg *T) {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatalf("failed to marshal yaml %s: %v", path, err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		log.Fatalf("failed to write file %s: %v", path, err)
	}
}

func mergeStructs[T any](dst, src *T) *T {
	data, _ := yaml.Marshal(src)
	_ = yaml.Unmarshal(data, dst)
	return dst
}
