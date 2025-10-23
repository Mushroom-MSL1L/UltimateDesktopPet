package configs

import (
	"os"

	"path/filepath"

	pp "UltimateDesktopPet/pkg/print"

	"gopkg.in/yaml.v2"
)

type ConfigProvider[T any] interface {
	DefaultConfig() *T
}

func LoadConfig[T any, P ConfigProvider[T]](path string, provider P) *T {
	pp.Info(pp.Config, "loading config %s", path)
	defaultCfg := provider.DefaultConfig()
	fileCfg := readConfigFile(path, provider)
	mergedCfg := mergeStructs(defaultCfg, fileCfg)
	writeConfigFile(path, mergedCfg)
	pp.Assert(pp.Config, "load config success %s", path)
	return mergedCfg
}

func readConfigFile[T any, P ConfigProvider[T]](path string, provider P) *T {
	defaultCfg := provider.DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			pp.Warn(pp.Config, "config file %s not exist, create new one with default", path)
			createConfigFile[T](path, defaultCfg)
			return defaultCfg
		}
		pp.Fatal(pp.Config, "cannot read config file %s: %v", path, err)
	}

	fileCfg := provider.DefaultConfig()
	if err := yaml.Unmarshal(data, fileCfg); err != nil {
		pp.Fatal(pp.Config, "yaml file %s unmarshal failed: %v", path, err)
	}
	return fileCfg
}

func createConfigFile[T any](path string, cfg *T) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		pp.Fatal(pp.Config, "failed to create directory %s: %v", dir, err)
	}

	writeConfigFile(path, cfg)
}

func writeConfigFile[T any](path string, cfg *T) {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		pp.Fatal(pp.Config, "failed to marshal yaml %s: %v", path, err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		pp.Fatal(pp.Config, "failed to write file %s: %v", path, err)
	}
}

func mergeStructs[T any](dst, src *T) *T {
	data, _ := yaml.Marshal(src)
	_ = yaml.Unmarshal(data, dst)
	return dst
}
