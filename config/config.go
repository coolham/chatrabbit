package config

import (
	"chatrabbit/config/common"
	"fmt"
	"log"
	"os"
	"sync"
)

// type ConfigServe struct {
// }

var (
	viperServe *ViperConfig
	configFile string
)

// 配置文件地址
var lock sync.Mutex

func GetConfigServe() *ViperConfig {
	if nil != viperServe {
		return viperServe
	}

	lock.Lock()
	defer lock.Unlock()

	InitConfig(configFile)
	return viperServe
}

// 获取配置文件
func InitConfig(cFile string) bool {
	if cFile != "" {
		configFile = cFile
		fmt.Printf("user configFile=%s", configFile)
		ok := LoadConfig(configFile)
		if !ok {
			fmt.Printf("load config file=%s failed\n", configFile)
		} else {
			return ok
		}
	}

	configFile := os.Getenv(common.CONFIGFILE)
	if len(configFile) == 0 {
		msg := fmt.Errorf("can not find env %s, please check it", common.CONFIGFILE)
		panic(msg)
	}

	log.Printf("configFile=%s", configFile)
	ok := LoadConfig(configFile)
	if !ok {
		panic("load config file failed")
	}
	return ok
}

func LoadConfig(configFile string) bool {
	fileInfo, err := os.Stat(configFile)
	if err != nil {
		fmt.Printf("can not find config file: %s\n", configFile)
		return false
	}

	if fileInfo.IsDir() {
		fmt.Printf("can not read config file: %s\n", configFile)
		return false
	}
	viperServe = NewViperConfig(configFile)
	return viperServe != nil
}

func GetInt(key string) int {
	return GetConfigServe().GetInt(key)
}

func GetBool(key string) bool {
	return GetConfigServe().GetBool(key)
}

func GetFloat64(key string) float64 {
	return GetConfigServe().GetFloat64(key)
}

func GetString(key string) string {
	return GetConfigServe().GetString(key)
}

func GetStringMap(key string) map[string]interface{} {
	return GetConfigServe().GetStringMap(key)
}

func GetStringSlice(key string) []string {
	return GetConfigServe().GetStringSlice(key)
}

func SetDefault(key string, value interface{}) {
	GetConfigServe().SetDefault(key, value)
}

func GetIntWithDefault(key string, v int) int {
	GetConfigServe().SetDefault(key, v)
	return GetConfigServe().GetInt(key)
}

func GetStringWithDefault(key string, v string) string {
	GetConfigServe().SetDefault(key, v)
	return GetConfigServe().GetString(key)
}

func GetBoolWithDefault(key string, v bool) bool {
	GetConfigServe().SetDefault(key, v)
	return GetConfigServe().GetBool(key)
}
