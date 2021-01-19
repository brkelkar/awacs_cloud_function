package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//Log  is globle varibale used for logging
var (
	log *zap.Logger
)

func init() {

	logConfig := zap.Config{
		OutputPaths: []string{"stdout"},
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "level",
			TimeKey:      "time",
			MessageKey:   "msg",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	var err error

	if log, err = logConfig.Build(); err != nil {
		panic(err)
	}
}

// GetLogger returns logger object
func GetLogger() *zap.Logger {
	return log
}

// Info log level debug message
func Info(msg string, tags ...zap.Field) {
	log.Info(msg, tags...)
	log.Sync()

}

// Debug log level debug message
func Debug(msg string, tags ...zap.Field) {
	log.Debug(msg, tags...)
	log.Sync()

}

// Error level debug message
func Error(msg string, err error, tags ...zap.Field) {
	if err != nil {
		tags = append(tags, zap.NamedError("Error", err))
	}
	log.Error(msg, tags...)
	log.Sync()

}
