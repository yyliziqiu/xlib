package xlog

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	rotate "github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"

	"github.com/rifflock/lfshook"

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

	Console, err = NewConsoleLogger("debug")
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
		return NewConsoleLogger(config.Level)
	}
	return NewFileLogger(config)
}

func NewConsoleLogger(level string) (*logrus.Logger, error) {
	logger := logrus.New()

	// 禁止输出方法名
	logger.SetReportCaller(false)

	// 设置日志等级
	logger.SetLevel(getLevel(level))

	// 设置日志格式
	logger.SetFormatter(getFormatter(textFormatter))

	return logger, nil
}

func getLevel(name string) logrus.Level {
	level, err := logrus.ParseLevel(name)
	if err != nil {
		return logrus.DebugLevel
	}
	return level
}

func getFormatter(name string) logrus.Formatter {
	timestampFormat := "2006-01-02 15:04:05"

	switch name {
	case "json":
		return &logrus.JSONFormatter{TimestampFormat: timestampFormat}
	}

	return &logrus.TextFormatter{TimestampFormat: timestampFormat, DisableQuote: true}
}

func NewFileLogger(config Config) (*logrus.Logger, error) {
	logger := logrus.New()

	// 禁止控制台输出
	logger.SetOutput(io.Discard)

	// 禁止输出方法名
	logger.SetReportCaller(false)

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
		formatter    = config.Formatter
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

	return lfshook.NewHook(rotation, getFormatter(formatter)), nil
}

func newTimeLevelRotationHook(config Config) (*lfshook.LfsHook, error) {
	var (
		name         = config.Name
		path         = config.Path
		formatter    = config.Formatter
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
	}, getFormatter(formatter)), nil
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
