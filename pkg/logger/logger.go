package logger

import (
	"main/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

func New(module string) *Logger {
	level := zap.InfoLevel

	if config.Env.Development {
		level = zap.DebugLevel
	}

	config := zap.NewProductionConfig()

	config.Level.SetLevel(level)
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.TimeKey = "timestamp"

	config.InitialFields = map[string]interface{}{
		"module": module,
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	return &Logger{SugaredLogger: logger.Sugar()}
}
