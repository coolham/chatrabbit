package config

import (
	"os"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfigInstance(t *testing.T) {

	// Convey("Test env not set", t, func() {
	// 	So(func() { GetConfigServe() }, ShouldPanic)
	// })

	Convey("Test sprcify config file", t, func() {
		configFile := "../conf/config.yaml"
		_, err := os.Stat(configFile)
		So(err, ShouldBeNil)
		LoadConfig(configFile)
		c := GetConfigServe()
		So(c, ShouldNotBeNil)
	})
}

func correntGetConfig(t *testing.T, w *sync.WaitGroup) {
	defer w.Done()
	Convey("Test sprcify config file", t, func() {
		c := GetConfigServe()
		So(c, ShouldNotBeNil)
	})
}
func TestConfigConcurrent(t *testing.T) {
	configFile := "../conf/config.yaml"
	LoadConfig(configFile)

	var w sync.WaitGroup
	for i := 0; i < 100; i++ {
		w.Add(1)
		go correntGetConfig(t, &w)
	}
	w.Wait()
}
