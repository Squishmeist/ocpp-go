package core

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

type Telemeter struct {
	serviceName string
	namespace   string
	endpoint    string
}

func (t *Telemeter) NewTracerProvider() *trace.TracerProvider {
	ctx := context.Background()
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(t.endpoint),
		otlptracehttp.WithInsecure(),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		panic(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return tp
}

func NewTelemeter(serviceName, endpoint, namespace string) *Telemeter {
	if endpoint == "" || serviceName == "" || namespace == "" {
		panic(fmt.Sprintf("telemeter is not configured: %s, %s, %s", serviceName, endpoint, namespace))
	}

	return &Telemeter{
		serviceName: serviceName,
		endpoint:    endpoint,
		namespace:   namespace,
	}
}

