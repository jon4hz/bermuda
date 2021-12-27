package logger

import (
	"log"
	"os"

	"github.com/jon4hz/bermuda/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log       *zap.Logger
	LogCloser func()
)

func New(cfg *config.LoggingConfig) {
	// info level enabler
	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.InfoLevel
	})

	// error and fatal level enabler
	errorFatalLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zapcore.ErrorLevel || level == zapcore.FatalLevel
	})

	anyLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return true
	})

	// write syncers
	stdoutSyncer := zapcore.Lock(os.Stdout)
	stderrSyncer := zapcore.Lock(os.Stderr)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.CallerKey = ""
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	cores := []zapcore.Core{
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			stdoutSyncer,
			infoLevel,
		),
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			stderrSyncer,
			errorFatalLevel,
		),
	}
	if cfg.LogFile != "" {
		fileSyncer, closer, err := zap.Open(cfg.LogFile)
		if err != nil {
			log.Fatal(err)
		}
		LogCloser = closer
		cores = append(cores, zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			fileSyncer,
			anyLevel,
		))
	}
	core := zapcore.NewTee(
		cores...,
	)
	Log = zap.New(core)
}
