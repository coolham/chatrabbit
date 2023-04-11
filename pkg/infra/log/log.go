package log

import (
	"chatrabbit/config"
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

//type Fields map[string]interface{}

type Logger struct {
	Log *logrus.Logger
}

var logger = logrus.New()

func Init(logLevel string) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic("can not parse log level str")
	}

	logger.SetLevel(level)
	logger.Formatter = &logrus.JSONFormatter{
		//TimestampFormat: time.RFC3339Nano,
	}
	// logger.Formatter = &logrus.TextFormatter{
	// 	TimestampFormat: "2006-01-02 15:04:05",
	// }
}

func GetLogger() *Logger {
	return &Logger{
		Log: logger,
	}
}

// CreateLogger creates a logger instance for all components
func CreateLogger(logLevel string) (*Logger, error) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic("can not parse log level str")
	}
	logger.SetLevel(level)
	logger.Formatter = &logrus.JSONFormatter{
		//TimestampFormat: time.RFC3339Nano,
	}
	// logger.Formatter = &logrus.TextFormatter{
	// 	TimestampFormat: "2006-01-02 15:04:05",
	// }

	return &Logger{
		Log: logger,
	}, nil
}

func SetLogLevel(levelStr string) (err error) {
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		return
	}
	logger.SetLevel(level)
	return
}

func SetLogFormatter(formatter logrus.Formatter) {
	logger.Formatter = formatter
}

func ConfigLogger() {
	logPath := config.GetString("debug.log.filepath")
	fileName := config.GetString("debug.log.filename")
	configLocalFilesystemLogger(logPath, fileName, 30*24*time.Hour, 7*24*time.Hour)
}

// config logrus log to local filesystem, with file rotation
func configLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	baseLogPaht := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPaht+".%Y%m%d%H%M.log",
		rotatelogs.WithLinkName(baseLogPaht),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		logger.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.TraceLevel: writer,
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, logger.Formatter)
	logger.AddHook(lfHook)
}

func fileInfo() string {
	_, file, line, ok := runtime.Caller(2) // skip 2 step to get where error happened
	if !ok {
		file = "<???>"
		line = -1
	} else {
		// TODO: 取最后两级文件结构
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			slash = strings.LastIndex(file[:slash], "/")
			if slash >= 0 {
				file = file[slash+1:]
			}
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func WithFields(fields map[string]interface{}) *logrus.Entry {
	var fs logrus.Fields = fields
	entry := logger.WithFields(fs)
	entry.Data["zfile"] = fileInfo()
	return entry
}

func Trace(args ...interface{}) {
	if logger.Level >= logrus.TraceLevel {
		logger.WithFields(logrus.Fields{"zfile": fileInfo()}).Trace(args...)
	}
}

func Tracef(format string, args ...interface{}) {
	if logger.Level >= logrus.TraceLevel {
		logger.WithFields(logrus.Fields{"zfile": fileInfo()}).Tracef(format, args...)
	}
}

func Debug(args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		logger.WithFields(logrus.Fields{"zfile": fileInfo()}).Debug(args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if logger.Level >= logrus.DebugLevel {
		logger.WithFields(logrus.Fields{"zfile": fileInfo()}).Debugf(format, args...)
	}
}

func Info(args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		logger.WithFields(logrus.Fields{"zfile": fileInfo()}).Info(args...)
	}
}

func Infof(format string, args ...interface{}) {
	if logger.Level >= logrus.InfoLevel {
		logger.WithFields(logrus.Fields{"zfile": fileInfo()}).Infof(format, args...)
	}
}

func Warn(args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		logger.WithFields(logrus.Fields{"zfile": fileInfo()}).Warn(args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if logger.Level >= logrus.WarnLevel {
		logger.WithFields(logrus.Fields{"zfile": fileInfo()}).Warnf(format, args...)
	}
}

func Error(args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		logger.WithFields(logrus.Fields{"zfile": fileInfo()}).Error(args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if logger.Level >= logrus.ErrorLevel {
		logger.WithFields(logrus.Fields{"zfile": fileInfo()}).Errorf(format, args...)
	}
}

func Fatal(args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		logger.WithFields(logrus.Fields{"zfile": fileInfo()}).Error(args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	if logger.Level >= logrus.FatalLevel {
		logger.WithFields(logrus.Fields{"zfile": fileInfo()}).Logger.Fatalf(format, args...)
	}
}

func Panic(args ...interface{}) {
	if logger.Level >= logrus.PanicLevel {
		logger.WithFields(logrus.Fields{"zfile": fileInfo()}).Panic(args...)
	}
}

func Panicf(format string, args ...interface{}) {
	if logger.Level >= logrus.PanicLevel {
		logger.WithFields(logrus.Fields{"zfile": fileInfo()}).Logger.Panicf(format, args...)
	}
}
