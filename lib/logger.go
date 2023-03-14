package lib

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func init() {
	Logger = newMultiOutputLogger()
}
func newMultiOutputLogger() *zap.Logger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	stdoutSink := zapcore.Lock(os.Stdout)
	stderrSink := zapcore.Lock(os.Stderr)

	enc := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	core := zapcore.NewTee(
		zapcore.NewCore(enc, stderrSink, highPriority),
		zapcore.NewCore(enc, stdoutSink, lowPriority),
	)

	logger := zap.New(core)
	return logger
}
