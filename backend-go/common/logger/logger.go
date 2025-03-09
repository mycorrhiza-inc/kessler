package logger

// Question I think this should probably go in the common library just because
// 1. Its probably good to standardize logs across all our go stuff
// 2. Its required by a couple of the filter objects in common, and thus cant really be ported to other libraries?
import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log  *zap.Logger
	once sync.Once
)

// Init initializes the global logger instance
func Init(env string) {
	once.Do(func() {
		var config zap.Config

		if env == "production" {
			config = zap.NewProductionConfig()
			config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		} else {
			config = zap.NewDevelopmentConfig()
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}

		var err error
		Log, err = config.Build(
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
		if err != nil {
			panic("failed to initialize logger: " + err.Error())
		}
	})
}

// GetLogger returns a named logger instance
func GetLogger(name string) *zap.Logger {
	if Log == nil {
		Init(os.Getenv("GO_ENV"))
	}
	return Log.Named(name)
}

// Sync flushes any buffered log entries
func Sync() error {
	if Log != nil {
		return Log.Sync()
	}
	return nil
}
