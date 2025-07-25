package core

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

type Telemeter struct {
	serviceName string
	namespace   string
	endpoint    string
}

func (t *Telemeter) NewTracerProvider() *trace.TracerProvider {
	ctx := context.Background()
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(t.endpoint),
		otlptracegrpc.WithInsecure(),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		panic(err)
	}

	resource := getResource(t.serviceName, t.namespace)

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource),
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

func getResource(serviceName, namespace string) *resource.Resource {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceNamespaceKey.String(namespace),
			semconv.ServiceInstanceID(getInstanceID()),
			semconv.TelemetrySDKLanguageGo,
			semconv.TelemetrySDKName("otel"),
		),
	)
	if err != nil {
		panic(err)
	}
	return r
}

func getInstanceID() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}
