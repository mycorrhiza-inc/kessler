package logger

// Question I think this should probably go in the common library just because
// 1. Its probably good to standardize logs across all our go stuff
// 2. Its required by a couple of the filter objects in common, and thus cant really be ported to other libraries?
import (
	"context"
	"net/http"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	_ "github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/trace"
)

const LoggerKey = "__internal__logger"

// Init initializes the global logger instance
func WithLogger(ctx context.Context) context.Context {

	var config zap.Config
	env := ctx.Value("system_env")
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
	BaseLog := *otelzap.New(bareLog)

	return context.WithValue(ctx, LoggerKey, BaseLog)
}

func LoggingMiddleware(baseLogger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = WithLogger(ctx)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// GetLogger returns a named logger instance
func GetLoggerFromContext(ctx context.Context) *otelzap.Logger {
	if logger, ok := ctx.Value(LoggerKey).(*otelzap.Logger); ok {
		return logger
	}
	// Return default/global logger as fallback
	return otelzap.L() // or your global logger instance
}
func GetLogger(name string) *otelzap.Logger {
	ctx := WithLogger(context.Background())
	log := GetLoggerFromContext(ctx)
	return otelzap.New(log.Named(name))
}

func Info(ctx context.Context, msg string, fields ...zapcore.Field) {
	log := GetLoggerFromContext(ctx)
	log.InfoContext(ctx, msg, fields...)
}
func Error(ctx context.Context, msg string, fields ...zapcore.Field) {
	log := GetLoggerFromContext(ctx)
	log.ErrorContext(ctx, msg, fields...)
}
func Debug(ctx context.Context, msg string, fields ...zapcore.Field) {
	log := GetLoggerFromContext(ctx)
	log.DebugContext(ctx, msg, fields...)
}
func Warn(ctx context.Context, msg string, fields ...zapcore.Field) {
	log := GetLoggerFromContext(ctx)
	log.WarnContext(ctx, msg, fields...)
}

func GetTrace(ctx context.Context) trace.Tracer {
	return
}

// Sync flushes any buffered log entries
func Sync(ctx context.Context) error {
	log := GetLoggerFromContext(ctx)
	return log.Sync()
}
