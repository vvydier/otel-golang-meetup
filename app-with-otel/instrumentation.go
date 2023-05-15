package main

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func newExporter(ctx context.Context) (*otlptrace.Exporter, error) {

	// "OTEL_EXPORTER_OTLP_ENDPOINT" should be specified without the schema part. ex: "127.0.0.1:4317"
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(collectorURL),
	)

	// We create an exporter using otlptracegrpc.New(). An exporter creates trace data in the OTLP wire format.
	// The one we are using here is backed by GRPC, you can also use other exporters backed with other transport mechanisms like http etc.
	return otlptrace.New(
		context.Background(),
		client,
	)
}

func newTraceProvider(exporter sdktrace.SpanExporter) (*sdktrace.TracerProvider, error) {
	log.Println("initializing resource")
	/*
		resources, err := resource.New(
			context.Background(),
			resource.WithSchemaURL(
				semconv.SchemaURL,
			),
			resource.WithAttributes(
				semconv.ServiceNameKey.String(serviceName),
				semconv.ServiceVersionKey.String(serviceVersion),
				attribute.String("library.language", "go"),
			),
		)
	*/
	resources, err := resource.New(context.Background(),
		resource.WithFromEnv(),   // pull attributes from OTEL_RESOURCE_ATTRIBUTES and OTEL_SERVICE_NAME environment variables
		resource.WithProcess(),   // This option configures a set of Detectors that discover process information
		resource.WithOS(),        // This option configures a set of Detectors that discover OS information
		resource.WithContainer(), // This option configures a set of Detectors that discover container information
		resource.WithHost(),      // This option configures a set of Detectors that discover host information
		// resource.WithDetectors(thirdparty.Detector{}),           // Bring your own external Detector implementation
		resource.WithAttributes(attribute.String("serviceVersion", serviceVersion)), // Or specify resource attributes directly
	)

	if err != nil {
		log.Printf("Could not set resources")
		return nil, err
	}

	log.Println("initializing tracer")
	// batchSpanProcessor := trace.NewBatchSpanProcessor(otlpTraceExporter)

	tracerProvider := sdktrace.NewTracerProvider(
		// sdktrace.WithSampler(sdktrace.NeverSample()),
		sdktrace.WithBatcher(exporter),
		// trace.WithSpanProcessor(batchSpanProcessor),
		sdktrace.WithResource(resources),
	)
	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, // W3C Trace Context format; https://www.w3.org/TR/trace-context/
			propagation.Baggage{},
		),
	)

	return tracerProvider, nil
}
