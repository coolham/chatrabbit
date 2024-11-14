package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Config interface {
	Get(key string) interface{}
	GetInt(key string) int
	GetFloat64(key string) float64
	GetBool(key string) bool
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringSlice(key string) []string
	SetString(key string, value string)
	SetDefault(key string, value interface{})
}

var (
	configInstance Config
	lock           sync.Mutex
)

// 获取项目根目录
func getProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("project root not found")
		}
		dir = parent
	}
}

// 获取配置实例
func GetConfig() (Config, error) {
	if configInstance == nil {
		// 在不持有锁的情况下调用 InitConfig
		if err := InitConfig(""); err != nil {
			return nil, err
		}
	}

	lock.Lock()
	defer lock.Unlock()

	return configInstance, nil
}

// 显式初始化函数
func InitConfig(file string) error {
	lock.Lock()
	defer lock.Unlock()

	if configInstance != nil {
		return nil
	}

	if file == "" {
		file = os.Getenv("RABBIT_CONFIG_FILE")
		if file == "" {
			root, err := getProjectRoot()
			if err != nil {
				return err
			}
			file = filepath.Join(root, "conf", "config.yaml")
		}
	}

	// check file exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return fmt.Errorf("config file does not exist: %s", file)
	}
	configInstance = NewViperConfig(file)
	return nil
}
