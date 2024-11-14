package config

import (
	"os"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	m.Run()
}

func createTempConfigFile(content string) (string, error) {
	tmpfile, err := os.CreateTemp("", "config_test_*.yaml")
	if err != nil {
		return "", err
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		tmpfile.Close()
		return "", err
	}

	if err := tmpfile.Close(); err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}

func TestViperConfig(t *testing.T) {
	Convey("Given a ViperConfig with a test config file", t, func() {
		configContent := `
some_string_key: "expected_string_value"
some_int_key: 42
some_float_key: 3.14
some_bool_key: true
some_map_key:
  key1: "value1"
  key2: "value2"
some_slice_key:
  - "item1"
  - "item2"
  - "item3"
`
		configFile, err := createTempConfigFile(configContent)
		So(err, ShouldBeNil)
		defer os.Remove(configFile) // 确保在测试结束时删除临时文件

		viperConfig := NewViperConfig(configFile)

		Convey("When getting a string value", func() {
			value := viperConfig.GetString("some_string_key")
			So(value, ShouldEqual, "expected_string_value")
		})

		Convey("When getting an int value", func() {
			value := viperConfig.GetInt("some_int_key")
			So(value, ShouldEqual, 42)
		})

		Convey("When getting a float64 value", func() {
			value := viperConfig.GetFloat64("some_float_key")
			So(value, ShouldEqual, 3.14)
		})

		Convey("When getting a bool value", func() {
			value := viperConfig.GetBool("some_bool_key")
			So(value, ShouldBeTrue)
		})

		Convey("When getting a string map", func() {
			value := viperConfig.GetStringMap("some_map_key")
			So(value, ShouldResemble, map[string]interface{}{"key1": "value1", "key2": "value2"})
		})

		Convey("When getting a string slice", func() {
			value := viperConfig.GetStringSlice("some_slice_key")
			So(value, ShouldResemble, []string{"item1", "item2", "item3"})
		})

		Convey("When setting a string value", func() {
			viperConfig.SetString("new_string_key", "new_value")
			value := viperConfig.GetString("new_string_key")
			So(value, ShouldEqual, "new_value")
		})

		Convey("When setting a default value", func() {
			viperConfig.SetDefault("default_key", "default_value")
			value := viperConfig.GetString("default_key")
			So(value, ShouldEqual, "default_value")
		})
	})
}

func TestConfigConcurrent(t *testing.T) {
	Convey("Given a ViperConfig with a test config file", t, func() {
		configContent := `
some_string_key: "expected_string_value"
`
		configFile, err := createTempConfigFile(configContent)
		So(err, ShouldBeNil)
		defer os.Remove(configFile) // 确保在测试结束时删除临时文件

		viperConfig := NewViperConfig(configFile)

		var w sync.WaitGroup
		for i := 0; i < 100; i++ {
			w.Add(1)
			go func() {
				defer w.Done()
				value := viperConfig.GetString("some_string_key")
				assert.Equal(t, value, "expected_string_value")
			}()
		}
		w.Wait()
	})
}

func TestInitConfig(t *testing.T) {
	Convey("direct get config instance", t, func() {
		config, err := GetConfig()
		So(err, ShouldBeNil)
		So(config, ShouldNotBeNil)
	})

	Convey("Given a config file path", t, func() {
		err := InitConfig("../conf/config.yaml")
		So(err, ShouldBeNil)
	})

}
