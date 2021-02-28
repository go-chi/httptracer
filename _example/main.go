package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httptracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/opentracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	// Tracer setup

	// Start the regular tracer and return it as an opentracing.Tracer interface. You
	// may use the same set of options as you normally would with the Datadog tracer.
	tr := opentracer.New(
		tracer.WithService("htdemo"),
		tracer.WithAgentAddr("0.0.0.0:8126"),
	)

	// Stop it using the regular Stop call for the tracer package.
	// defer tracer.Stop()

	// Set the global OpenTracing tracer.
	// opentracing.SetGlobalTracer(t)

	// HTTP service
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(httptracer.Tracer(tr, httptracer.Config{
		ServiceName:    "htdemo",
		ServiceVersion: "v0.1.0",
		SampleRate:     1,
		SkipFunc: func(r *http.Request) bool {
			return r.URL.Path == "/health"
		},
		Tags: map[string]interface{}{
			"_dd.measured": 1, // datadog, turn on metrics for http.request stats
			// "_dd1.sr.eausr": 1, // datadog, event sample rate
		},
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("index"))
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("good"))
	})

	r.Get("/fail", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte("ah, internal server error"))
	})

	r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("oh no")
	})

	http.ListenAndServe(":3333", r)
}
