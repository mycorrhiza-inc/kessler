package logger

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"go.opentelemetry.io/otel"
// 	"go.opentelemetry.io/otel/attribute"
// 	"go.opentelemetry.io/otel/trace"
// 	"go.uber.org/zap"
// 	"go.uber.org/zap/zapcore"
// )

// // ResponseWriter wrapper to capture status code
// type responseWriter struct {
// 	http.ResponseWriter
// 	statusCode int
// 	written    int64
// }

// func (rw *responseWriter) WriteHeader(code int) {
// 	rw.statusCode = code
// 	rw.ResponseWriter.WriteHeader(code)
// }

// func (rw *responseWriter) Write(b []byte) (int, error) {
// 	if rw.statusCode == 0 {
// 		rw.statusCode = http.StatusOK
// 	}
// 	n, err := rw.ResponseWriter.Write(b)
// 	rw.written += int64(n)
// 	return n, err
// }

// // TracingMiddleware creates a middleware that traces and logs every HTTP request
// func TracingMiddleware() func(http.Handler) http.Handler {
// 	tracer := otel.Tracer("http-server")

// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			start := time.Now()

// 			// Create a new span for this request
// 			ctx, span := tracer.Start(r.Context(), r.Method+" "+r.URL.Path,
// 				trace.WithAttributes(
// 					attribute.String("http.method", r.Method),
// 					attribute.String("http.url", r.URL.String()),
// 					attribute.String("http.route", r.URL.Path),
// 					attribute.String("http.scheme", r.URL.Scheme),
// 					attribute.String("http.host", r.Host),
// 					attribute.String("http.user_agent", r.UserAgent()),
// 					attribute.String("http.remote_addr", r.RemoteAddr),
// 				),
// 			)
// 			defer span.End()

// 			// Wrap the response writer to capture status code and response size
// 			wrapped := &responseWriter{
// 				ResponseWriter: w,
// 				statusCode:     0,
// 			}

// 			// Add request ID for correlation (you can generate this however you prefer)
// 			requestID := span.SpanContext().SpanID().String()
// 			type contextKey string
// 			const requestIDKey contextKey = "request_id"
// 			ctx = context.WithValue(ctx, requestIDKey, requestID)

// 			// Ensure logger is available in context
// 			if GetLoggerFromContext(ctx) == nil {
// 				ctx = WithLogger(ctx)
// 			}

// 			// Log request start
// 			Info(ctx, "HTTP request started",
// 				zap.String("method", r.Method),
// 				zap.String("path", r.URL.Path),
// 				zap.String("query", r.URL.RawQuery),
// 				zap.String("remote_addr", r.RemoteAddr),
// 				zap.String("user_agent", r.UserAgent()),
// 				zap.String("request_id", requestID),
// 				zap.String("host", r.Host),
// 			)

// 			// Process the request
// 			next.ServeHTTP(wrapped, r.WithContext(ctx))

// 			// Calculate duration
// 			duration := time.Since(start)

// 			// Add response attributes to span
// 			span.SetAttributes(
// 				attribute.Int("http.status_code", wrapped.statusCode),
// 				attribute.Int64("http.response_size", wrapped.written),
// 				attribute.Int64("http.duration_ms", duration.Milliseconds()),
// 			)

// 			// Determine log level based on status code
// 			logLevel := zapcore.InfoLevel
// 			if wrapped.statusCode >= 400 && wrapped.statusCode < 500 {
// 				logLevel = zapcore.WarnLevel
// 			} else if wrapped.statusCode >= 500 {
// 				logLevel = zapcore.ErrorLevel
// 			}

// 			// Log request completion
// 			logFields := []zapcore.Field{
// 				zap.String("method", r.Method),
// 				zap.String("path", r.URL.Path),
// 				zap.String("query", r.URL.RawQuery),
// 				zap.Int("status_code", wrapped.statusCode),
// 				zap.Duration("duration", duration),
// 				zap.Int64("response_size", wrapped.written),
// 				zap.String("request_id", requestID),
// 				zap.String("remote_addr", r.RemoteAddr),
// 			}

// 			switch logLevel {
// 			case zapcore.ErrorLevel:
// 				Error(ctx, "HTTP request completed with error", logFields...)
// 			case zapcore.WarnLevel:
// 				Warn(ctx, "HTTP request completed with warning", logFields...)
// 			default:
// 				Info(ctx, "HTTP request completed", logFields...)
// 			}
// 		})
// 	}
// }

// // Enhanced version with more detailed tracing and error handling
// func DetailedTracingMiddleware() func(http.Handler) http.Handler {
// 	tracer := otel.Tracer("http-server-detailed")

// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			start := time.Now()

// 			// Create span with more detailed attributes
// 			ctx, span := tracer.Start(r.Context(), r.Method+" "+r.URL.Path,
// 				trace.WithAttributes(
// 					attribute.String("http.method", r.Method),
// 					attribute.String("http.url", r.URL.String()),
// 					attribute.String("http.route", r.URL.Path),
// 					attribute.String("http.scheme", r.URL.Scheme),
// 					attribute.String("http.host", r.Host),
// 					attribute.String("http.user_agent", r.UserAgent()),
// 					attribute.String("http.remote_addr", r.RemoteAddr),
// 					attribute.String("http.x_forwarded_for", r.Header.Get("X-Forwarded-For")),
// 					attribute.String("http.x_real_ip", r.Header.Get("X-Real-IP")),
// 					attribute.String("http.content_type", r.Header.Get("Content-Type")),
// 					attribute.String("http.accept", r.Header.Get("Accept")),
// 					attribute.Int64("http.content_length", r.ContentLength),
// 				),
// 			)
// 			defer func() {
// 				// Handle any panics and add to span
// 				if err := recover(); err != nil {
// 					// Convert panic to error for recording
// 					var panicErr error
// 					if e, ok := err.(error); ok {
// 						panicErr = e
// 					} else {
// 						panicErr = fmt.Errorf("panic: %v", err)
// 					}

// 					span.RecordError(panicErr)
// 					span.SetAttributes(
// 						attribute.String("error.panic", "true"),
// 						attribute.String("error.message", panicErr.Error()),
// 					)
// 					Error(ctx, "HTTP request panicked",
// 						zap.Any("panic", err),
// 						zap.String("method", r.Method),
// 						zap.String("path", r.URL.Path),
// 					)
// 					panic(err) // Re-throw the panic
// 				}
// 				span.End()
// 			}()

// 			wrapped := &responseWriter{
// 				ResponseWriter: w,
// 				statusCode:     0,
// 			}

// 			requestID := span.SpanContext().SpanID().String()
// 			ctx = context.WithValue(ctx, "request_id", requestID)

// 			// Ensure logger context
// 			if GetLoggerFromContext(ctx) == nil {
// 				ctx = WithLogger(ctx)
// 			}

// 			// Log detailed request start
// 			Info(ctx, "HTTP request started",
// 				zap.String("method", r.Method),
// 				zap.String("path", r.URL.Path),
// 				zap.String("query", r.URL.RawQuery),
// 				zap.String("remote_addr", r.RemoteAddr),
// 				zap.String("user_agent", r.UserAgent()),
// 				zap.String("request_id", requestID),
// 				zap.String("host", r.Host),
// 				zap.String("x_forwarded_for", r.Header.Get("X-Forwarded-For")),
// 				zap.String("x_real_ip", r.Header.Get("X-Real-IP")),
// 				zap.String("content_type", r.Header.Get("Content-Type")),
// 				zap.Int64("content_length", r.ContentLength),
// 				zap.String("referer", r.Referer()),
// 			)

// 			// Process request
// 			next.ServeHTTP(wrapped, r.WithContext(ctx))

// 			duration := time.Since(start)

// 			// Add comprehensive response attributes
// 			span.SetAttributes(
// 				attribute.Int("http.status_code", wrapped.statusCode),
// 				attribute.Int64("http.response_size", wrapped.written),
// 				attribute.Int64("http.duration_ms", duration.Milliseconds()),
// 				attribute.Float64("http.duration_seconds", duration.Seconds()),
// 			)

// 			// Determine log level and message
// 			var logMsg string
// 			var logLevel zapcore.Level

// 			switch {
// 			case wrapped.statusCode >= 500:
// 				logLevel = zapcore.ErrorLevel
// 				logMsg = "HTTP request failed with server error"
// 			case wrapped.statusCode >= 400:
// 				logLevel = zapcore.WarnLevel
// 				logMsg = "HTTP request failed with client error"
// 			case wrapped.statusCode >= 300:
// 				logLevel = zapcore.InfoLevel
// 				logMsg = "HTTP request redirected"
// 			default:
// 				logLevel = zapcore.InfoLevel
// 				logMsg = "HTTP request completed successfully"
// 			}

// 			// Comprehensive logging
// 			logFields := []zapcore.Field{
// 				zap.String("method", r.Method),
// 				zap.String("path", r.URL.Path),
// 				zap.String("query", r.URL.RawQuery),
// 				zap.Int("status_code", wrapped.statusCode),
// 				zap.String("status_text", http.StatusText(wrapped.statusCode)),
// 				zap.Duration("duration", duration),
// 				zap.Float64("duration_seconds", duration.Seconds()),
// 				zap.Int64("response_size", wrapped.written),
// 				zap.String("request_id", requestID),
// 				zap.String("remote_addr", r.RemoteAddr),
// 				zap.String("user_agent", r.UserAgent()),
// 				zap.String("trace_id", span.SpanContext().TraceID().String()),
// 				zap.String("span_id", span.SpanContext().SpanID().String()),
// 			}

// 			switch logLevel {
// 			case zapcore.ErrorLevel:
// 				Error(ctx, logMsg, logFields...)
// 			case zapcore.WarnLevel:
// 				Warn(ctx, logMsg, logFields...)
// 			default:
// 				Info(ctx, logMsg, logFields...)
// 			}
// 		})
// 	}
// }
