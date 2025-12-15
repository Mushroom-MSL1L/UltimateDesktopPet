package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	pp "github.com/Mushroom-MSL1L/UltimateDesktopPet/pkg/print"
	"gopkg.in/yaml.v2"
)

// ConfigService exposes helpers to read and write the system config file.
// It avoids exiting the app on validation errors so the UI can surface them
// without crashing the process.
type ConfigService struct {
	path string
	mu   sync.Mutex
}

func NewConfigService(path string) *ConfigService {
	return &ConfigService{
		path: path,
	}
}

func (c *ConfigService) LoadSystemConfig() (System, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cfg := (&System{}).DefaultConfig()

	data, err := os.ReadFile(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			return *cfg, c.writeSystemConfig(cfg)
		}
		return System{}, fmt.Errorf("read system config: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return System{}, fmt.Errorf("parse system config: %w", err)
	}

	return *cfg, nil
}

func (c *ConfigService) SaveSystemConfig(cfg System) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.writeSystemConfig(&cfg); err != nil {
		return err
	}

	pp.Info(pp.Config, "system config saved to %s", c.path)
	return nil
}

func (c *ConfigService) writeSystemConfig(cfg *System) error {
	if err := os.MkdirAll(filepath.Dir(c.path), 0755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("encode system config: %w", err)
	}

	if err := os.WriteFile(c.path, data, 0644); err != nil {
		return fmt.Errorf("write system config: %w", err)
	}

	return nil
}
