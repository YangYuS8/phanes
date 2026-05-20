package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/YangYuS8/phanes/internal/contracts"
)

type Config struct {
	WorkspaceDir string        `json:"workspace_dir"`
	Runtime      RuntimeConfig `json:"runtime"`
	Cache        CacheConfig   `json:"cache"`
	Save         SaveConfig    `json:"save"`
}

type RuntimeConfig struct {
	BindHost string `json:"bind_host"`
	HTTPPort uint32 `json:"http_port"`
	Mode     string `json:"mode"`
}

type CacheConfig struct {
	Dir string `json:"dir"`
}

type SaveConfig struct {
	DBPath string `json:"db_path"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) Validate() error {
	bindHost, err := contracts.NormalizeBindHost(c.Runtime.BindHost)
	if err != nil {
		return err
	}
	c.Runtime.BindHost = bindHost

	if c.Cache.Dir == "" {
		return errors.New("cache.dir is required")
	}
	if c.Save.DBPath == "" {
		return errors.New("save.db_path is required")
	}
	if err := contracts.CacheCleanAllowed(c.Cache.Dir, c.Save.DBPath); err != nil {
		return fmt.Errorf("invalid cache/save paths: %w", err)
	}
	return nil
}
