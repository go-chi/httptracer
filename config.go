package httptracer

import "net/http"

type Config struct {
	// Service details to record
	ServiceName    string
	ServiceVersion string

	// Operation name
	//
	// The span operation name record for the http request trace.
	// Default value if empty is set to "http.request"
	OperationName string

	// SampleRate sample rate, value between 0 to 1.0
	//
	// Only span a percentage of the spans. Default value is
	// set to 1.0
	SampleRate float64

	// Skip particular requests from the tracer
	SkipFunc func(r *http.Request) bool

	// SetTagFunc
	//
	// ...
	// SetTagFunc func(r *http.Request) map[string]interface{}

	// Tags
	//
	// Extra tags to include with a span
	Tags map[string]interface{}
}
