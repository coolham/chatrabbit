package log

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateLogger(t *testing.T) {
	level := "info"
	Convey("create logger", t, func() {
		logger, err := CreateLogger(level)
		So(err, ShouldBeNil)
		So(logger, ShouldNotBeNil)
	})
}

func TestSetLogLevel(t *testing.T) {
	level := "info"
	Convey("create logger", t, func() {
		logger, err := CreateLogger(level)
		So(err, ShouldBeNil)
		So(logger, ShouldNotBeNil)
		err = SetLogLevel("error")
		So(err, ShouldBeNil)
	})
}

// func TestLogTrace(t *testing.T) {
// 	SetLogLevel("trace")
// 	Trace("hello trace")
// 	Tracef("jerry:%d", 123)
// 	WithFields(Fields{"username": "tom"}).Tracef("hello:%s", "jack")
// }

// func TestLogsInfo(t *testing.T) {
// 	SetLogLevel("info")
// 	Info("hello trace")
// 	Infof("jerry:%d", 123)
// 	WithFields(Fields{"username": "tom"}).Infof("hello:%s", "jack")
// }

// func TestLogsLevel(t *testing.T) {
// 	SetLogLevel("trace")
// 	WithFields(Fields{"username": "tom"}).Tracef("hello:%s", "jack")
// 	WithFields(Fields{"username": "tom"}).Debugf("hello:%s", "jack")
// 	WithFields(Fields{"username": "tom"}).Infof("hello:%s", "jack")
// 	WithFields(Fields{"username": "tom"}).Warnf("hello:%s", "jack")
// 	WithFields(Fields{"username": "tom"}).Errorf("hello:%s", "jack")
// 	WithFields(Fields{"username": "tom"}).Fatalf("hello:%s", "jack")
// }
