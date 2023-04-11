package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type ViperConfig struct {
	viper *viper.Viper
}

// NewFileObj 加载配置文件
func NewViperConfig(file string) *ViperConfig {
	newViper := viper.New()

	newViper.SetConfigFile(file)
	newViper.SetConfigType("yaml")

	if err := newViper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误
			panic(errors.Wrap(err, "viper can not find config file"))
		} else {
			// 配置文件被找到，但产生了另外的错误
			panic(errors.Wrap(err, "viper read config file error"))
		}
	}

	newViper.WatchConfig()
	newViper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("config file has changed: name=%s, type=%s", e.Name, e.Op.String())
	})
	return &ViperConfig{
		viper: newViper,
	}

}

func (v *ViperConfig) Get(key string) interface{} {
	return v.viper.Get(key)
}

// GetInt for get config int
func (v *ViperConfig) GetInt(key string) int {
	return v.viper.GetInt(key)
}

// GetInt for get config int
func (v *ViperConfig) GetFloat64(key string) float64 {
	return v.viper.GetFloat64(key)
}

func (v *ViperConfig) GetBool(key string) bool {
	return v.viper.GetBool(key)
}

// GetStringMap for get config map[string]interface
func (v *ViperConfig) GetStringMap(key string) map[string]interface{} {
	return v.viper.GetStringMap(key)
}

// GetString for get config string
func (v *ViperConfig) GetString(key string) string {
	return v.viper.GetString(key)
}

// GetStringSlice for get config []string
func (v *ViperConfig) GetStringSlice(key string) []string {
	return v.viper.GetStringSlice(key)
}

func (v *ViperConfig) SetDefault(key string, value interface{}) {
	v.viper.SetDefault(key, value)
}
