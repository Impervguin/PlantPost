package opentelemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type TracingConfig struct {
	Enabled        bool
	JaegerEndpoint string
	ServiceName    string
	ServiceVersion string
	Environment    string
}

func InitTracer(cfg *TracingConfig) (func(), error) {
	if !cfg.Enabled {
		return func() {}, nil
	}

	// Создаем Jaeger exporter
	exp, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(cfg.JaegerEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	// Создаем ресурс с атрибутами сервиса
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			attribute.String("environment", cfg.Environment),
		)),
	)

	// Устанавливаем глобальный TracerProvider
	otel.SetTracerProvider(tp)

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			fmt.Printf("Error shutting down tracer provider: %v\n", err)
		}
	}, nil
}
