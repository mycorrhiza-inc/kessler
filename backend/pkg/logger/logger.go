package logger

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// contextKey is a custom type to avoid context key collisions
type contextKey string

const loggerKey contextKey = "logger"

var (
	// Global logger instance for non-context scenarios
	globalLogger *otelzap.Logger
)

// Config holds logger configuration
type Config struct {
	Level       zapcore.Level
	Environment string
	ServiceName string
}

// Init initializes both OpenTelemetry and the global logger
func Init(cfg Config) error {
	// Initialize OpenTelemetry
	if err := initTracing(cfg.ServiceName); err != nil {
		return err
	}

	// Initialize global logger
	var err error
	globalLogger, err = createLogger(cfg)
	if err != nil {
		return err
	}

	return nil
}

// initTracing sets up OpenTelemetry tracing
func initTracing(serviceName string) error {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)
	return nil
}

// createLogger creates a new logger instance
func createLogger(cfg Config) (*otelzap.Logger, error) {
	var config zap.Config

	if cfg.Environment == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	config.Level = zap.NewAtomicLevelAt(cfg.Level)

	zapLogger, err := config.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}

	return otelzap.New(zapLogger), nil
}

// WithLogger adds a logger to the context
func WithLogger(ctx context.Context) context.Context {
	if ctx.Value(loggerKey) != nil {
		return ctx // Logger already exists
	}
	return context.WithValue(ctx, loggerKey, globalLogger)
}

// FromContext retrieves a logger from context, or returns the global logger
func FromContext(ctx context.Context) *otelzap.Logger {
	if logger, ok := ctx.Value(loggerKey).(*otelzap.Logger); ok && logger != nil {
		return logger
	}
	return globalLogger
}

// Global returns the global logger instance
func Global() *otelzap.Logger {
	return globalLogger
}

// Named returns a named logger
func Named(name string) *otelzap.Logger {
	if globalLogger == nil {
		return nil
	}
	// Get the underlying zap logger, create named version, then wrap with otelzap
	namedZapLogger := globalLogger.Logger.Named(name)
	return otelzap.New(namedZapLogger)
}

// Middleware returns an HTTP middleware that adds logger to request context
func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := WithLogger(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ResponseWriter wrapper to capture status code and response size
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

// TracingMiddleware creates HTTP request tracing and logging middleware
func TracingMiddleware() func(http.Handler) http.Handler {
	tracer := otel.Tracer("http-server")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create span
			ctx, span := tracer.Start(r.Context(), r.Method+" "+r.URL.Path,
				oteltrace.WithAttributes(
					attribute.String("http.method", r.Method),
					attribute.String("http.url", r.URL.String()),
					attribute.String("http.route", r.URL.Path),
					attribute.String("http.host", r.Host),
					attribute.String("http.user_agent", r.UserAgent()),
					attribute.String("http.remote_addr", r.RemoteAddr),
				),
			)
			defer span.End()

			// Ensure logger is in context
			ctx = WithLogger(ctx)
			logger := FromContext(ctx)

			// Wrap response writer
			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     0,
			}

			requestID := span.SpanContext().SpanID().String()

			// Log request start
			logger.InfoContext(ctx, "HTTP request started",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("query", r.URL.RawQuery),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
				zap.String("request_id", requestID),
				zap.String("host", r.Host),
			)

			// Process request
			next.ServeHTTP(wrapped, r.WithContext(ctx))

			// Calculate duration
			duration := time.Since(start)

			// Add response attributes to span
			span.SetAttributes(
				attribute.Int("http.status_code", wrapped.statusCode),
				attribute.Int64("http.response_size", wrapped.written),
				attribute.Int64("http.duration_ms", duration.Milliseconds()),
			)

			// Determine log level based on status code
			logFields := []zap.Field{
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("query", r.URL.RawQuery),
				zap.Int("status_code", wrapped.statusCode),
				zap.Duration("duration", duration),
				zap.Int64("response_size", wrapped.written),
				zap.String("request_id", requestID),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("trace_id", span.SpanContext().TraceID().String()),
			}

			switch {
			case wrapped.statusCode >= 500:
				logger.ErrorContext(ctx, "HTTP request completed with server error", logFields...)
			case wrapped.statusCode >= 400:
				logger.WarnContext(ctx, "HTTP request completed with client error", logFields...)
			default:
				logger.InfoContext(ctx, "HTTP request completed", logFields...)
			}
		})
	}
}

// Context-aware logging functions
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).InfoContext(ctx, msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).ErrorContext(ctx, msg, fields...)
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).DebugContext(ctx, msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).WarnContext(ctx, msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).FatalContext(ctx, msg, fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() Config {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	level := zapcore.InfoLevel
	if env == "development" {
		level = zapcore.DebugLevel
	}

	return Config{
		Level:       level,
		Environment: env,
		ServiceName: "kessler",
	}
}
