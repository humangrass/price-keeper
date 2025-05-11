package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zap *zap.SugaredLogger
}

func (l *Logger) Sugar() *zap.SugaredLogger {
	return l.zap
}

func (l *Logger) Drop() error {
	return l.zap.Sync()
}

func (l *Logger) DropMsg() string {
	return "close zap logger"
}

func New(isProduction bool) (*Logger, error) {
	var zapConfig zap.Config
	if isProduction {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
	}

	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	zapConfig.EncoderConfig.TimeKey = "time"

	zapLogger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{
		zap: zapLogger.Sugar(),
	}, nil
}
