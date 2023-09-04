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
	config Config

	std *logrus.Logger

	Http *logrus.Logger
)

func Init(c Config) error {
	var err error

	config = c.WithDefault()

	std, err = NewLogger(config.Name)
	if err != nil {
		return err
	}

	Http, err = NewLogger("http-")
	if err != nil {
		return err
	}

	return nil
}

func NewLogger(prefix string) (*logrus.Logger, error) {
	if config.Console {
		return newConsoleLogger()
	}
	return newFileLogger(prefix)
}

func newConsoleLogger() (*logrus.Logger, error) {
	logger := logrus.New()

	// 设置日志等级
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("can't parse level, error [%v]", err)
	}
	logger.SetLevel(level)

	// 禁止输出方法名
	logger.SetReportCaller(false)

	logger.SetFormatter(&logrus.TextFormatter{
		DisableQuote:    true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	return logger, nil
}

func newFileLogger(prefix string) (*logrus.Logger, error) {
	logger := logrus.New()

	// 设置日志等级
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("can't parse level, error [%v]", err)
	}
	logger.SetLevel(level)

	// 禁止控制台输出
	logger.SetOutput(io.Discard)

	// 日志按天分割
	hook, err := newRotateHook(config.Path, prefix, config.MaxAge, config.RotateTime)
	if err != nil {
		return nil, fmt.Errorf("create hook error [%v]", err)
	}
	logger.AddHook(hook)

	// 禁止输出方法名
	logger.SetReportCaller(false)

	return logger, nil
}

func newRotateHook(dirPath string, prefix string, maxAge time.Duration, rotateTime time.Duration) (*lfshook.LfsHook, error) {
	err := xutil.MkdirIfNotExist(dirPath)
	if err != nil {
		return nil, fmt.Errorf("create log dir error [%v]", err)
	}

	debugRotate, err := newRotate(dirPath, prefix+"debug-%Y%m%d.log", maxAge, rotateTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate error [%v]", err)
	}
	infoRotate, err := newRotate(dirPath, prefix+"info-%Y%m%d.log", maxAge, rotateTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate error [%v]", err)
	}
	warnRotate, err := newRotate(dirPath, prefix+"warn-%Y%m%d.log", maxAge, rotateTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate error [%v]", err)
	}
	errorRotate, err := newRotate(dirPath, prefix+"error-%Y%m%d.log", maxAge, rotateTime)
	if err != nil {
		return nil, fmt.Errorf("create rotate error [%v]", err)
	}

	return lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: debugRotate,
		logrus.InfoLevel:  infoRotate,
		logrus.WarnLevel:  warnRotate,
		logrus.ErrorLevel: errorRotate,
		logrus.FatalLevel: errorRotate,
		logrus.PanicLevel: errorRotate,
	}, &logrus.TextFormatter{
		DisableQuote:    true,
		TimestampFormat: "2006-01-02 15:04:05",
	}), nil
}

func newRotate(dirPath string, filename string, maxAge time.Duration, rotateTime time.Duration) (*rotate.RotateLogs, error) {
	return rotate.New(filepath.Join(dirPath, filename), rotate.WithMaxAge(maxAge), rotate.WithRotationTime(rotateTime))
}

func NewLoggerMust(prefix string) *logrus.Logger {
	logger, err := NewLogger(prefix)
	if err != nil {
		return std
	}
	return logger
}

func GetDefaultLogger() *logrus.Logger {
	return std
}
