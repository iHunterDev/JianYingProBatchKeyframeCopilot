package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

type Config struct {
	Port          string   `json:"port"`            // http服务端口
	CsrfDomains   []string `json:"csrf_domain"`     // CSRF 域名列表
	DraftRootPath string   `json:"draft_root_path"` // 草稿根目录
}

func DefaultConfig() Config {
	defaultConfig := Config{
		Port:          ":9507",
		CsrfDomains:   []string{"keyframeai.top"},
		DraftRootPath: "",
	}

	return defaultConfig
}

type ConfigStore struct {
	configPath string
}

func NewConfigStore() (*ConfigStore, error) {
	configFilePath, err := xdg.ConfigFile("jianyingpro-batch-keyframe-copilot/config.json")
	if err != nil {
		return nil, fmt.Errorf("could not resolve path for config file: %w", err)
	}
	return &ConfigStore{
		configPath: configFilePath,
	}, nil
}

func (s *ConfigStore) Config() (Config, error) {
	_, err := os.Stat(s.configPath)
	if os.IsNotExist(err) {
		s.SaveConfig(DefaultConfig())
		return DefaultConfig(), nil
	}
	dir, fileName := filepath.Split(s.configPath)
	if len(dir) == 0 {
		dir = "."
	}
	buf, err := fs.ReadFile(os.DirFS(dir), fileName)
	if err != nil {
		return Config{}, fmt.Errorf("could not read the configuration file: %w", err)
	}
	if len(buf) == 0 {
		return DefaultConfig(), nil
	}
	cfg := Config{}
	if err := json.Unmarshal(buf, &cfg); err != nil {
		return Config{}, fmt.Errorf("configuration file does not have a valid format: %w", err)
	}
	return cfg, nil
}

func (s *ConfigStore) SaveConfig(cfg Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("could not marshal config data: %w", err)
	}
	dir, _ := filepath.Split(s.configPath)
	if len(dir) == 0 {
		dir = "."
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("could not create config directory: %w", err)
	}
	err = os.WriteFile(s.configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("could not write config file: %w", err)
	}
	return nil
}

// func main() {
// 	store, err := NewConfigStore()
// 	if err != nil {
// 		fmt.Printf("could not initialize the config store: %v\n", err)
// 		return
// 	}
// 	fmt.Println(store.configPath)
// 	cfg, err := store.Config()
// 	if err != nil {
// 		fmt.Printf("could not retrieve the configuration: %v\n", err)
// 		return
// 	}
// 	fmt.Printf("config: %v\n", cfg)
// }
