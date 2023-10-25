package xlog

import (
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"time"

	rotate "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/xlib/xutil"
)

var (
	defaultConfig Config

	HTTP *logrus.Logger

	Console *logrus.Logger

	Default *logrus.Logger
)

func Init(config Config) error {
	var err error

	defaultConfig = config.WithDefault()

	HTTP, err = NewLoggerByName("http")
	if err != nil {
		return err
	}

	Console, err = NewConsoleLogger(Config{Level: "debug"})
	if err != nil {
		return err
	}

	Default, err = NewLogger(defaultConfig)
	if err != nil {
		return err
	}

	return nil
}

func NewLoggerByName(name string) (*logrus.Logger, error) {
	config := defaultConfig
	config.Name = name
	return NewLogger(config)
}

func NewLogger(config Config) (*logrus.Logger, error) {
	if config.Console {
		return NewConsoleLogger(config)
	}
	return NewFileLogger(config)
}

func NewConsoleLogger(config Config) (*logrus.Logger, error) {
	logger := logrus.New()

	// 禁止输出方法名
	logger.SetReportCaller(false)

	// 设置日志等级
	logger.SetLevel(getLevel(config.Level))

	// 设置日志格式
	logger.SetFormatter(getFormatter(config))

	return logger, nil
}

func getLevel(name string) logrus.Level {
	level, err := logrus.ParseLevel(name)
	if err != nil {
		return logrus.DebugLevel
	}
	return level
}

func getFormatter(config Config) logrus.Formatter {
	var (
		formatter        = config.Formatter
		timestampFormat  = config.TimestampFormat
		callerPrettyfier = getCallerPrettyfier(config.CallerFields)
	)

	if timestampFormat == "" {
		timestampFormat = "2006-01-02 15:04:05"
	}

	switch formatter {
	case "json":
		return &logrus.JSONFormatter{
			TimestampFormat:  timestampFormat,
			CallerPrettyfier: callerPrettyfier,
		}
	default:
		return &logrus.TextFormatter{
			DisableQuote:     true,
			PadLevelText:     true,
			TimestampFormat:  timestampFormat,
			CallerPrettyfier: callerPrettyfier,
		}
	}
}

func getCallerPrettyfier(fields string) func(frame *runtime.Frame) (function string, file string) {
	switch fields {
	case fieldsFilename:
		return func(frame *runtime.Frame) (function string, file string) {
			return "", fmt.Sprintf("%s:%d", frame.File, frame.Line)
		}
	case fieldsFunction:
		return func(frame *runtime.Frame) (function string, file string) {
			return frame.Function, ""
		}
	default:
		return func(frame *runtime.Frame) (function string, file string) {
			return frame.Function, fmt.Sprintf("%s:%d", frame.File, frame.Line)
		}
	}
}

func NewFileLogger(config Config) (*logrus.Logger, error) {
	logger := logrus.New()

	// 禁止控制台输出
	logger.SetOutput(io.Discard)

	// 禁止输出方法名
	logger.SetReportCaller(config.EnableCaller)

	// 设置日志等级
	logger.SetLevel(getLevel(config.Level))

	// 日志按天分割
	hook, err := getRotationHook(config)
	if err != nil {
		return nil, fmt.Errorf("create hook error [%v]", err)
	}
	logger.AddHook(hook)

	return logger, nil
}

func getRotationHook(config Config) (*lfshook.LfsHook, error) {
	if config.DisableLevelRotation {
		return newTimeRotationHook(config)
	}
	return newTimeLevelRotationHook(config)
}

func newTimeRotationHook(config Config) (*lfshook.LfsHook, error) {
	var (
		name         = config.Name
		path         = config.Path
		maxAge       = config.MaxAge
		rotationTime = config.RotationTime
	)

	// 确保日志目录存在
	err := xutil.MkdirIfNotExist(path)
	if err != nil {
		return nil, fmt.Errorf("create log dir error [%v]", err)
	}

	// 美化日志文件名
	if name != "" {
		name = name + "-"
	}

	// 创建分割器
	rotation, err := NewRotation(path, name+"-%Y%m%d.log", maxAge, rotationTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate error [%v]", err)
	}

	return lfshook.NewHook(rotation, getFormatter(config)), nil
}

func newTimeLevelRotationHook(config Config) (*lfshook.LfsHook, error) {
	var (
		name         = config.Name
		path         = config.Path
		maxAge       = config.MaxAge
		rotationTime = config.RotationTime
	)

	// 确保日志目录存在
	err := xutil.MkdirIfNotExist(config.Path)
	if err != nil {
		return nil, fmt.Errorf("create log dir error [%v]", err)
	}

	// 美化日志文件名
	if name != "" {
		name = name + "-"
	}

	// 创建分割器
	debugRotation, err := NewRotation(path, name+"debug-%Y%m%d.log", maxAge, rotationTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate error [%v]", err)
	}
	infoRotation, err := NewRotation(path, name+"info-%Y%m%d.log", maxAge, rotationTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate error [%v]", err)
	}
	warnRotation, err := NewRotation(path, name+"warn-%Y%m%d.log", maxAge, rotationTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate error [%v]", err)
	}
	errorRotation, err := NewRotation(path, name+"error-%Y%m%d.log", maxAge, rotationTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate error [%v]", err)
	}

	return lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: debugRotation,
		logrus.InfoLevel:  infoRotation,
		logrus.WarnLevel:  warnRotation,
		logrus.ErrorLevel: errorRotation,
		logrus.FatalLevel: errorRotation,
		logrus.PanicLevel: errorRotation,
	}, getFormatter(config)), nil
}

func NewRotation(dirname string, filename string, maxAge time.Duration, RotationTime time.Duration) (*rotate.RotateLogs, error) {
	return rotate.New(filepath.Join(dirname, filename), rotate.WithMaxAge(maxAge), rotate.WithRotationTime(RotationTime))
}

func MustNewLoggerByName(name string) *logrus.Logger {
	logger, err := NewLoggerByName(name)
	if err != nil {
		return Default
	}
	return logger
}
