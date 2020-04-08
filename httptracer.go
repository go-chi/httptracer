package httptracer

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"
	"github.com/goware/httplog"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func Tracer(cfg Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			opts := []ddtrace.StartSpanOption{
				tracer.SpanType(ext.SpanTypeWeb),
				tracer.ServiceName(cfg.ServiceName),
				tracer.Tag(ext.HTTPMethod, r.Method),
				tracer.Tag(ext.HTTPURL, r.URL.Path),
				tracer.Measured(),
			}
			if cfg.AnalyticsRate == 0 || cfg.AnalyticsRate > 1 {
				cfg.AnalyticsRate = 1
			}
			if cfg.AnalyticsRate >= 0 { // use -1 to disable
				opts = append(opts, tracer.Tag(ext.EventSampleRate, cfg.AnalyticsRate))
			}
			if spanctx, err := tracer.Extract(tracer.HTTPHeadersCarrier(r.Header)); err == nil {
				opts = append(opts, tracer.ChildOf(spanctx))
			}
			// opts = append(opts, cfg.spanOpts...)

			span, ctx := tracer.StartSpanFromContext(r.Context(), "http.request", opts...)
			defer span.Finish()

			httplog.LogEntrySetFields(r.Context(), map[string]interface{}{
				"dd.trace_id": fmt.Sprintf("%v", span.Context().SpanID()),
				"dd.span_id":  fmt.Sprintf("%v", span.Context().TraceID()),
			})

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// pass the span through the request context and serve the request to the next middleware
			next.ServeHTTP(ww, r.WithContext(ctx))

			// set the resource name as we get it only once the handler is executed
			// resourceName := chi.RouteContext(r.Context()).RoutePattern()
			// if resourceName == "" {
			// 	resourceName = r.URL.Path
			// }
			resourceName := r.URL.Path
			resourceName = r.Method + " " + resourceName
			span.SetTag(ext.ResourceName, resourceName)

			// set the status code
			status := ww.Status()
			span.SetTag(ext.HTTPCode, strconv.Itoa(status))

			if status >= 500 && status < 600 {
				// mark 5xx server error
				span.SetTag(ext.Error, fmt.Errorf("%d: %s", status, http.StatusText(status)))
			}
		})
	}
}
