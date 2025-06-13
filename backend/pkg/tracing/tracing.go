package tracing

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// Config holds the tracing configuration
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	ExporterType   ExporterType
	OTLPEndpoint   string
	SampleRate     float64
}

// ExporterType defines the type of trace exporter to use
type ExporterType string

const (
	ExporterStdout ExporterType = "stdout"
	ExporterOTLP   ExporterType = "otlp"
)

// Initialize sets up OpenTelemetry tracing with the provided configuration
func Initialize(cfg Config) (func(context.Context) error, error) {
	// Create resource with service information
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create exporter based on configuration
	exporter, err := createExporter(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	// Create tracer provider with sampling
	samplerOption := trace.WithSampler(trace.TraceIDRatioBased(cfg.SampleRate))
	if cfg.SampleRate <= 0 {
		samplerOption = trace.WithSampler(trace.NeverSample())
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
		samplerOption,
	)

	// Set as global tracer provider
	otel.SetTracerProvider(tp)

	// Return shutdown function
	shutdown := func(ctx context.Context) error {
		return tp.Shutdown(ctx)
	}

	return shutdown, nil
}

// MustInitialize initializes tracing and panics on error
func MustInitialize(cfg Config) func(context.Context) error {
	shutdown, err := Initialize(cfg)
	if err != nil {
		panic("failed to initialize tracing: " + err.Error())
	}
	return shutdown
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig(serviceName string) Config {
	return Config{
		ServiceName:    serviceName,
		ServiceVersion: getVersion(),
		Environment:    getEnvironment(),
		ExporterType:   getExporterType(),
		OTLPEndpoint:   getOTLPEndpoint(),
		SampleRate:     getSampleRate(),
	}
}

// createExporter creates the appropriate exporter based on configuration
func createExporter(cfg Config) (trace.SpanExporter, error) {
	switch cfg.ExporterType {
	case ExporterStdout:
		return stdouttrace.New(stdouttrace.WithPrettyPrint())
	case ExporterOTLP:
		return createOTLPExporter(cfg.OTLPEndpoint)
	default:
		return stdouttrace.New(stdouttrace.WithPrettyPrint())
	}
}

// createOTLPExporter creates an OTLP HTTP exporter
func createOTLPExporter(endpoint string) (trace.SpanExporter, error) {
	if endpoint == "" {
		endpoint = "http://localhost:4318/v1/traces"
	}

	return otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(endpoint),
		),
	)
}

// Helper functions to get configuration from environment

func getVersion() string {
	if version := os.Getenv("SERVICE_VERSION"); version != "" {
		return version
	}
	return "unknown"
}

func getEnvironment() string {
	env := os.Getenv("OTEL_ENVIRONMENT")
	if env == "" {
		env = os.Getenv("GO_ENV")
	}
	if env == "" {
		env = os.Getenv("ENV")
	}
	if env == "" {
		env = "development"
	}
	return env
}

func getExporterType() ExporterType {
	exporterType := os.Getenv("OTEL_EXPORTER_TYPE")
	switch exporterType {
	case "otlp":
		return ExporterOTLP
	case "stdout":
		return ExporterStdout
	default:
		// Default to stdout for development, OTLP for production
		if getEnvironment() == "production" {
			return ExporterOTLP
		}
		return ExporterStdout
	}
}

func getOTLPEndpoint() string {
	return os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")
}

func getSampleRate() float64 {
	sampleRateStr := os.Getenv("OTEL_TRACES_SAMPLER_ARG")
	if sampleRateStr == "" {
		// Default sampling rates
		if getEnvironment() == "production" {
			return 0.1 // 10% sampling in production
		}
		return 1.0 // 100% sampling in development
	}

	// Parse the sample rate (simplified - in real code you'd want proper parsing)
	switch sampleRateStr {
	case "0":
		return 0.0
	case "1":
		return 1.0
	default:
		return 0.1
	}
}
