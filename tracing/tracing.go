package tracer

import (
	"context"
	"log"
	"os"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Tracer is an instance for tracing using Golang Opentracing and Golang Tracer client
type Tracer struct {
	ctx           context.Context
	push          func()
	TraceProvider trace.TracerProvider
	Name          string
}

//Config is configuration for Tracer
type Config struct {
	ServiceName string
	Endpoint    string
	Probability float64
}

// NewJaeger to create new Tracer instance with jaeger as destination
func NewJaeger(config *Config) func() {
	if config == nil {
		log.Fatalf("Jaeger Config Read Failed")
	}

	exporter, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint(config.Endpoint),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: config.ServiceName,
			Tags: []label.KeyValue{
				label.String("host", os.Getenv("HOSTNAME")),
			},
		}),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.TraceIDRatioBased(config.Probability)}),
	)

	if err != nil {
		log.Fatalf("Jaeger Initialization Failed")
	}

	return func() {
		exporter()
	}
}
