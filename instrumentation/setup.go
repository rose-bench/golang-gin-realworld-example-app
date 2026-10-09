package instrumentation

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func Setup(ctx context.Context, serviceName string) (func(context.Context) error, error) {
	noopShutdown := func(context.Context) error { return nil }
	res, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("service.version", "unknown"),
			attribute.String("deployment.environment.name", "development"),
			attribute.String("service.owner.url", "https://github.com/rose-bench/golang-gin-realworld-example-app"),
		),
		resource.WithFromEnv(),
		resource.WithAttributes(attribute.String("service.name", serviceName)),
	)
	if err != nil {
		return noopShutdown, err
	}

	exporter, err := otlptracehttp.New(ctx)
	if err != nil {
		return noopShutdown, err
	}
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{},
	))
	return provider.Shutdown, nil
}
