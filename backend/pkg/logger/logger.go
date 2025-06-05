package logger

import (
	"context"
	"log"
	"net/http"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	_ "github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const LoggerKey = "__internal__logger"

func Init() {
	// Create a stdout exporter
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}

	// Create a tracer provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
	)

	// Set as global tracer provider
	otel.SetTracerProvider(tp)
}

// Init initializes the global logger instance
func WithLogger(ctx context.Context) context.Context {
	var config zap.Config
	env := "dev"
	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	var err error
	bareLog, err := config.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	BaseLog := otelzap.New(bareLog)

	return context.WithValue(ctx, LoggerKey, BaseLog)
}

func LoggingMiddleware(baseLogger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := WithLogger(r.Context()) // Add logger to context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// GetLogger returns a named logger instance

func GetLoggerFromContext(ctx context.Context) *otelzap.Logger {
	if logger, ok := ctx.Value(LoggerKey).(*otelzap.Logger); ok && logger != nil {
		return logger
	}

	// Create a safe fallback logger instead of relying on global
	config := zap.NewProductionConfig()
	bareLog, err := config.Build()
	if err != nil {
		// Last resort - use standard log
		log.Printf("Failed to create fallback logger: %v", err)
		return nil
	}
	return otelzap.New(bareLog)
}

func GetLogger(name string) *otelzap.Logger {
	ctx := WithLogger(context.Background())
	log := GetLoggerFromContext(ctx)
	return log
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

// Sync flushes any buffered log entries
func Sync(ctx context.Context) error {
	log := GetLoggerFromContext(ctx)
	return log.Sync()
}
