package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"

	// "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	// serviceName    = os.Getenv("SERVICE_NAME")
	serviceVersion = "1.0.0"
	collectorURL   = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	// insecure       = os.Getenv("INSECURE_MODE")
)

var (
	// fooKey     = attribute.Key("ex.com/foo")
	// barKey     = attribute.Key("ex.com/bar")
	orderIDKey = attribute.Key("order_id")
)

func handler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	order_id := params["orderID"]
	log.Printf("received order for id: %s", order_id)
	ctx := r.Context()

	// simulate work
	doWork(ctx)

	json := simplejson.New()
	json.Set("order_id", order_id)
	json.Set("status", "received")

	payload, err := json.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func doWork(ctx context.Context) {
	/*
		m0, _ := baggage.NewMember(string(fooKey), "foo1")
		m1, _ := baggage.NewMember(string(barKey), "bar1")
		b, _ := baggage.New(m0, m1)
		newctx := baggage.ContextWithBaggage(ctx, b)
	*/
	r := rand.Intn(1729)
	time.Sleep(time.Duration(r) * time.Microsecond)

	_, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("exampleTracer").Start(ctx, "doWork")
	defer span.End()
	// span.AddEvent("Nice operation!", trace.WithAttributes(attribute.Int("bogons", 100)))
	span.SetAttributes(orderIDKey.String("yes"))
	doMoreWork(ctx)
	doMoreWork(ctx)
	doMoreWork(ctx)

}

func doMoreWork(ctx context.Context) {
	r := rand.Intn(100000)
	time.Sleep(time.Duration(r) * time.Microsecond)

	_, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("exampleTracer").Start(ctx, "doMoreWork")
	defer span.End()
	span.SetAttributes(orderIDKey.String("yes"))
}

func catchAllHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Invalid URI path. Use /shipping/{orderID} to create shipping for order"))
}

// Route declaration
func router() *mux.Router {
	r := mux.NewRouter()

	// handlerFunc := http.HandlerFunc(handler)
	// wrappedHandler := otelhttp.NewHandler(handlerFunc, "handle_request")
	// r.Handle("/shipping/{orderID}", wrappedHandler)
	r.HandleFunc("/shipping/{orderID}", handler)

	// r.Handle("/", (otelhttp.NewHandler(http.HandlerFunc(catchAllHandler), "handle_catchall_request")))
	r.HandleFunc("/", catchAllHandler)

	return r
}

// otel error handler
type OtelSpanErrorHandler struct{}

func (OtelSpanErrorHandler) Handle(err error) {
	fmt.Println(err)
	switch err.(type) {
	// case *SpanExporterError:
	default:
		fmt.Println(err)
	}
}

// Initiate web server
func main() {
	log.Println("Collector URL is ", collectorURL)
	ctx := context.Background()

	exp, err := newExporter(ctx)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}
	// exp.Start(ctx)

	// Create a new tracer provider with a batch span processor and the given exporter.
	tp, err := newTraceProvider(exp)
	if err != nil {
		log.Fatalf("failed to initialize tracer provider: %v", err)
	}

	// Handle shutdown properly so nothing leaks.
	defer func() {
		_ = tp.Shutdown(ctx)
		log.Println("otel tracer shutting down")
	}()

	otel.SetTracerProvider(tp)

	otel.SetErrorHandler(OtelSpanErrorHandler{})

	// Finally, set the tracer that can be used for this package.
	// tracer = tp.Tracer("shipping-service")

	router := router()
	router.Use(otelmux.Middleware("shipping-service"))

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:9100",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
