package xlog

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	rotate "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"

	"github.com/yyliziqiu/xlib/xfile"
)

var (
	_config Config

	Default *logrus.Logger

	Console *logrus.Logger
)

func Init(config Config) (err error) {
	_config = config.Default()

	Default, err = New(_config)
	if err != nil {
		return err
	}

	Console, err = NewConsoleLogger(_config)
	if err != nil {
		return err
	}

	return nil
}

func New(config Config) (*logrus.Logger, error) {
	if config.Console {
		return NewConsoleLogger(config)
	}
	return NewFileLogger(config)
}

func NewConsoleLogger(config Config) (*logrus.Logger, error) {
	logger := logrus.New()

	// 禁止输出方法名
	logger.SetReportCaller(config.EnableCaller)

	// 设置日志等级
	logger.SetLevel(level(config.Level))

	// 设置日志格式
	logger.SetFormatter(formatter(config))

	return logger, nil
}

func level(name string) logrus.Level {
	lvl, err := logrus.ParseLevel(name)
	if err != nil {
		return logrus.DebugLevel
	}
	return lvl
}

func formatter(config Config) logrus.Formatter {
	var (
		formatterName   = config.Formatter
		timestampFormat = config.TimestampFormat
	)

	if timestampFormat == "" {
		timestampFormat = "2006-01-02 15:04:05"
	}

	switch formatterName {
	case JSONFormatterName:
		return &logrus.JSONFormatter{
			TimestampFormat: timestampFormat,
		}
	default:
		return &logrus.TextFormatter{
			DisableQuote:    true,
			TimestampFormat: timestampFormat,
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
	logger.SetLevel(level(config.Level))

	// 日志按天分割
	hook, err := getRotationHook(config)
	if err != nil {
		return nil, fmt.Errorf("create hook failed [%v]", err)
	}
	logger.AddHook(hook)

	return logger, nil
}

func getRotationHook(config Config) (*lfshook.LfsHook, error) {
	switch config.RotationLevel {
	case 0:
		return newRotationHook0(config)
	case 1:
		return newRotationHook1(config)
	default:
		return newRotationHook2(config)
	}
}

func newRotationHook0(config Config) (*lfshook.LfsHook, error) {
	var (
		name         = config.Name
		path         = config.Path
		maxAge       = config.MaxAge
		rotationTime = config.RotationTime
	)

	// 确保日志目录存在
	err := xfile.MkdirIfNotExist(config.Path)
	if err != nil {
		return nil, fmt.Errorf("create logs dir failed [%v]", err)
	}

	// 美化日志文件名
	if !strings.HasSuffix(name, "-") {
		name = name + "-"
	}

	// 创建分割器
	rotation, err := NewRotation(path, name+"%Y%m%d.log", maxAge, rotationTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate failed [%v]", err)
	}
	errorRotation, err := NewRotation(path, name+"error-%Y%m%d.log", maxAge, rotationTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate failed [%v]", err)
	}

	return lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: rotation,
		logrus.InfoLevel:  rotation,
		logrus.WarnLevel:  rotation,
		logrus.ErrorLevel: errorRotation,
		logrus.FatalLevel: errorRotation,
		logrus.PanicLevel: errorRotation,
	}, formatter(config)), nil
}

func newRotationHook1(config Config) (*lfshook.LfsHook, error) {
	var (
		name         = config.Name
		path         = config.Path
		maxAge       = config.MaxAge
		rotationTime = config.RotationTime
	)

	// 确保日志目录存在
	err := xfile.MkdirIfNotExist(path)
	if err != nil {
		return nil, fmt.Errorf("create logs dir failed [%v]", err)
	}

	// 美化日志文件名
	if !strings.HasSuffix(name, "-") {
		name = name + "-"
	}

	// 创建分割器
	rotation, err := NewRotation(path, name+"%Y%m%d.log", maxAge, rotationTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate failed [%v]", err)
	}

	return lfshook.NewHook(rotation, formatter(config)), nil
}

func newRotationHook2(config Config) (*lfshook.LfsHook, error) {
	var (
		name         = config.Name
		path         = config.Path
		maxAge       = config.MaxAge
		rotationTime = config.RotationTime
	)

	// 确保日志目录存在
	err := xfile.MkdirIfNotExist(config.Path)
	if err != nil {
		return nil, fmt.Errorf("create logs dir failed [%v]", err)
	}

	// 美化日志文件名
	if !strings.HasSuffix(name, "-") {
		name = name + "-"
	}

	// 创建分割器
	debugRotation, err := NewRotation(path, name+"debug-%Y%m%d.log", maxAge, rotationTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate failed [%v]", err)
	}
	infoRotation, err := NewRotation(path, name+"info-%Y%m%d.log", maxAge, rotationTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate failed [%v]", err)
	}
	warnRotation, err := NewRotation(path, name+"warn-%Y%m%d.log", maxAge, rotationTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate failed [%v]", err)
	}
	errorRotation, err := NewRotation(path, name+"error-%Y%m%d.log", maxAge, rotationTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate failed [%v]", err)
	}

	return lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: debugRotation,
		logrus.InfoLevel:  infoRotation,
		logrus.WarnLevel:  warnRotation,
		logrus.ErrorLevel: errorRotation,
		logrus.FatalLevel: errorRotation,
		logrus.PanicLevel: errorRotation,
	}, formatter(config)), nil
}

func NewRotation(dirname string, filename string, maxAge time.Duration, RotationTime time.Duration) (*rotate.RotateLogs, error) {
	return rotate.New(filepath.Join(dirname, filename), rotate.WithMaxAge(maxAge), rotate.WithRotationTime(RotationTime))
}

func NewWithName(name string) (*logrus.Logger, error) {
	config := _config
	config.Name = name
	return New(config)
}

func NewWithNameMust(name string) *logrus.Logger {
	logger, err := NewWithName(name)
	if err != nil {
		return Default
	}
	return logger
}
