package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var (
	isSinkRegistered bool
)

type apiLogger struct {
	base  *zap.SugaredLogger
	level string
}

// NewLogger инициализирует логгер с уровнем логгирования level, в папке path по относительному пути, с добавлением
// postfix к имени файла (postfix будет перед временем)
func NewLogger(level string, path string, postfix string) (Interface, error) {
	if path == "" {
		path = "logs"
	}
	parsedLevel, parseErr := zapcore.ParseLevel(level)
	if parseErr != nil {
		return nil, parseErr
	}

	levelSelected := normalizeLevel(parsedLevel)

	atom := zap.NewAtomicLevelAt(levelSelected)

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	filename := getOSFilePath(filepath.Join(path,
		fmt.Sprintf("log_%s%s.log", postfix,
			time.Now().Format("02012006"))))

	config := zap.Config{
		Level:    atom,
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:     "time",
			MessageKey:  "message",
			LevelKey:    "level",
			NameKey:     "namekey",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
			EncodeTime:  zapcore.ISO8601TimeEncoder},
		OutputPaths:      []string{filename},
		ErrorOutputPaths: []string{filename},
	}
	newWinFileSink := func(u *url.URL) (zap.Sink, error) {
		// Remove leading slash left by url.Parse()
		return os.OpenFile(u.Path[1:], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	}

	if !isSinkRegistered {
		err = zap.RegisterSink("winfile", newWinFileSink)
		if err != nil {
			return nil, err
		}
		isSinkRegistered = true
	}

	logger := zap.Must(config.Build())

	logger.Info("logger construction succeeded")

	errSync := logger.Sync()
	if errSync != nil {
		fmt.Println("logger sync failed: ", errSync)
		return nil, errSync
	}
	al := apiLogger{base: logger.Sugar(), level: level}
	return &al, nil
}

func (al *apiLogger) Debugf(template string, args ...interface{}) {
	al.base.Debugf(template, args)
}
func (al *apiLogger) Debugw(msg string, keysAndValues ...interface{}) {
	al.base.Debugw(msg, keysAndValues)
}
func (al *apiLogger) Error(args ...interface{}) {
	al.base.Error(args)
}
func (al *apiLogger) Errorf(template string, args ...interface{}) {
	al.base.Errorf(template, args)
}
func (al *apiLogger) Errorw(msg string, keysAndValues ...interface{}) {
	al.base.Errorw(msg, keysAndValues)
}
func (al *apiLogger) Info(args ...interface{}) {
	al.base.Info(args)
}
func (al *apiLogger) Infof(template string, args ...interface{}) {
	al.base.Infof(template, args)
}
func (al *apiLogger) Infow(msg string, keysAndValues ...interface{}) {
	al.base.Infow(msg, keysAndValues)
}
func (al *apiLogger) Warn(args ...interface{}) {
	al.base.Warn(args)
}
func (al *apiLogger) Warnf(template string, args ...interface{}) {
	al.base.Warnf(template, args)
}
func (al *apiLogger) Warnw(msg string, keysAndValues ...interface{}) {
	al.base.Warnw(msg, keysAndValues)
}

func (al *apiLogger) Level() string {
	return al.level
}

// normalizeLevel выравнивает уровень для Dev и Prod, возвращая DebugLevel или ErrorLevel соответственно.
func normalizeLevel(level zapcore.Level) zapcore.Level {
	if level >= zapcore.ErrorLevel {
		return zapcore.DebugLevel
	}
	return zapcore.ErrorLevel
}
