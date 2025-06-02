package logger

// Question I think this should probably go in the common library just because
// 1. Its probably good to standardize logs across all our go stuff
// 2. Its required by a couple of the filter objects in common, and thus cant really be ported to other libraries?
import (
	"context"
	"os"
	"sync"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	_ "github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log  *otelzap.Logger
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
		bareLog, err := config.Build(
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
		if err != nil {
			panic("failed to initialize logger: " + err.Error())
		}
		Log = otelzap.New(bareLog)
	})
}

// GetLogger returns a named logger instance
func GetLogger(name string) *otelzap.Logger {
	if Log == nil {
		Init(os.Getenv("GO_ENV"))
	}
	return otelzap.New(Log.Named(name))
}

// GetLogger returns a named logger instance
func GetLoggerCtx(ctx context.Context, name string) otelzap.LoggerWithCtx {
	if Log == nil {
		Init(os.Getenv("GO_ENV"))
	}
	tmp_logger := otelzap.New(Log.Named(name))
	return tmp_logger.Ctx(ctx)
}

// Sync flushes any buffered log entries
func Sync() error {
	if Log != nil {
		return Log.Sync()
	}
	return nil
}
