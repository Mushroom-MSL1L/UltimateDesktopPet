package configLoader

import (
	"fmt"
	"os"

	"path/filepath"

	pp "github.com/Mushroom-MSL1L/UltimateDesktopPet/pkg/print"

	"gopkg.in/yaml.v3"
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

func UpdateConfig[T any, P ConfigProvider[T], L any](path string, provider P, key string, value L) {
	pp.Info(pp.Config, "updating config %s", path)

	data, err := os.ReadFile(path)
	if err != nil {
		pp.Fatal(pp.Config, "read config failed: %v", err)
	}

	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		pp.Fatal(pp.Config, "yaml unmarshal failed: %v", err)
	}
	if len(node.Content) == 0 || node.Content[0].Kind != yaml.MappingNode {
		pp.Fatal(pp.Config, "invalid yaml structure")
	}
	mapping := node.Content[0]

	found := false
	for i := 0; i < len(mapping.Content); i += 2 {
		k := mapping.Content[i]
		v := mapping.Content[i+1]
		if k.Value == key {
			v.Value = fmt.Sprint(value)
			v.Kind = yaml.ScalarNode
			found = true
			break
		}
	}
	if !found {
		mapping.Content = append(mapping.Content,
			&yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: key,
			},
			&yaml.Node{
				Kind:  yaml.ScalarNode,
				Value: fmt.Sprint(value),
			},
		)
	}

	out, err := yaml.Marshal(&node)
	if err != nil {
		pp.Fatal(pp.Config, "yaml marshal failed: %v", err)
	}

	if err := os.WriteFile(path, out, 0644); err != nil {
		pp.Fatal(pp.Config, "failed to write file %s: %v", path, err)
	}
	pp.Assert(pp.Config, "update config success %s", path)
}
