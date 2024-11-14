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

func NewViperConfig(file string) Config {
	newViper := viper.New()

	newViper.SetConfigFile(file)
	newViper.SetConfigType("yaml")

	if err := newViper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(errors.Wrap(err, "viper can not find config file"))
		} else {
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

func (v *ViperConfig) GetInt(key string) int {
	return v.viper.GetInt(key)
}

func (v *ViperConfig) GetFloat64(key string) float64 {
	return v.viper.GetFloat64(key)
}

func (v *ViperConfig) GetBool(key string) bool {
	return v.viper.GetBool(key)
}

func (v *ViperConfig) GetString(key string) string {
	return v.viper.GetString(key)
}

func (v *ViperConfig) GetStringMap(key string) map[string]interface{} {
	return v.viper.GetStringMap(key)
}

func (v *ViperConfig) GetStringSlice(key string) []string {
	return v.viper.GetStringSlice(key)
}

func (v *ViperConfig) SetString(key string, value string) {
	v.viper.Set(key, value)
}

func (v *ViperConfig) SetDefault(key string, value interface{}) {
	v.viper.SetDefault(key, value)
}
