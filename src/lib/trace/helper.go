// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trace

import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// GetGlobalTracer returns the global tracer.
func GetGlobalTracer(instrumentationName string, opts ...oteltrace.TracerOption) oteltrace.Tracer {
	return otel.GetTracerProvider().Tracer(instrumentationName, opts...)
}

// StartSpan starts a span with the given name.
func StartSpan(ctx context.Context, name string) (context.Context, oteltrace.Span) {
	return otel.Tracer("goharbor/harbor/src/lib/trace").Start(ctx, name)
}

// SpanFromContext returns the span from the context.
func SpanFromHTTPRequest(req *http.Request) oteltrace.Span {
	ctx := req.Context()
	return oteltrace.SpanFromContext(ctx)
}

// RecordError records the error in the span from context.
func RecordError(span oteltrace.Span, err error, description string) {
	if span == nil {
		return
	}
	span.RecordError(err)
	span.SetStatus(codes.Error, description)
}

// NewHandler returns a handler that wraps the given handler with tracing.
func NewHandler(h http.Handler, operation string) http.Handler {
	httpOptions := []otelhttp.Option{
		otelhttp.WithTracerProvider(otel.GetTracerProvider()),
		otelhttp.WithPropagators(otel.GetTextMapPropagator()),
	}
	return otelhttp.NewHandler(h, operation, httpOptions...)
}

// StarTrace returns a new span with the given name.
func StartTrace(ctx context.Context, tracerName string, spanName string, opts ...oteltrace.SpanStartOption) (context.Context, oteltrace.Span) {
	return otel.Tracer(tracerName).Start(ctx, spanName, opts...)
}
