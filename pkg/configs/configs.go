package configs

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadConfig[T any](path string) *T {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("cannot read config file %s: %v", path, err)
	}

	var cfg T
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalln("yaml file %s unmarshal failed: %v", path, err)
	}

	return &cfg
}
