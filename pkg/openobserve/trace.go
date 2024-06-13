package openobserve

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	tracerKey = "tracer"
	// ScopeName is the instrumentation scope name.
	ScopeName = "github.com/tranhuyducseven/goem-template"
)

var Tracer oteltrace.Tracer

type O2Span struct {
	oteltrace.Span
}

func (s *O2Span) End() {
	if s.Span != nil {
		s.Span.End()
	}
}

func Trace(ctx context.Context, name string) (context.Context, *O2Span) {
	if Tracer == nil {
		return ctx, &O2Span{Span: nil}
	}
	ctx, span := Tracer.Start(ctx, name)
	return ctx, &O2Span{Span: span}
}

func SetupTraceHTTP() *sdktrace.TracerProvider {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	otlptracehttp.NewClient()

	otlpHTTPExporter, err := otlptracehttp.New(context.TODO(),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpointURL(config.Endpoint),
		otlptracehttp.WithURLPath(fmt.Sprintf("/api/%s/traces", config.OrgName)),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization": "Basic " + config.Credential,
			"stream-name":   config.ServiceName,
		}),
	)

	if err != nil {
		fmt.Println("Error creating HTTP OTLP exporter: ", err)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		// the service name used to display traces in backends
		semconv.ServiceNameKey.String(config.ServiceName),
		semconv.ServiceVersionKey.String("0.0.1"),
		attribute.String("environment", config.Env),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(otlpHTTPExporter),
	)
	otel.SetTracerProvider(tp)

	return tp
}

// Middleware returns echo middleware which will trace incoming requests.
func Middleware() echo.MiddlewareFunc {
	traceProvider := otel.GetTracerProvider()
	Tracer = traceProvider.Tracer(ScopeName)

	propagators := otel.GetTextMapPropagator()
	skipper := middleware.DefaultSkipper

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next(c)
			}

			c.Set(tracerKey, Tracer)
			request := c.Request()
			savedCtx := request.Context()
			defer func() {
				request = request.WithContext(savedCtx)
				c.SetRequest(request)
			}()
			ctx := propagators.Extract(savedCtx, propagation.HeaderCarrier(request.Header))
			opts := []oteltrace.SpanStartOption{
				oteltrace.WithAttributes(extRequestInfo(request)...),
				oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			}
			if path := c.Path(); path != "" {
				rAttr := semconv.HTTPRoute(path)
				opts = append(opts, oteltrace.WithAttributes(rAttr))
			}
			spanName := c.Path()
			if spanName == "" {
				spanName = fmt.Sprintf("HTTP %s route not found", request.Method)
			}

			ctx, span := Tracer.Start(ctx, spanName, opts...)
			defer span.End()

			// pass the span through the request context
			c.SetRequest(request.WithContext(ctx))

			// serve the request to the next middleware
			err := next(c)
			if err != nil {
				span.SetAttributes(attribute.String("echo.error", err.Error()))
				// invokes the registered HTTP error handler
				c.Error(err)
			}

			status := c.Response().Status
			span.SetStatus(serverStatus(status))
			if status > 0 {
				span.SetAttributes(semconv.HTTPStatusCode(status))
			}

			return err
		}
	}
}

func extRequestInfo(req *http.Request) []attribute.KeyValue {
	size := 4 // Method, scheme, proto, and host name.

	method, proto := req.Method, req.Proto

	protoName, protoVersion, _ := strings.Cut(proto, "/")
	if len(protoVersion) > 0 {
		size++
	}

	port := req.URL.Port()
	if len(port) > 0 {
		size++
	}

	remoteAddr := strings.Split(req.RemoteAddr, ":")

	var peer string
	if len(remoteAddr) > 0 {
		peer = remoteAddr[0]
		size++
	}

	var peerPort int = -1
	if len(remoteAddr) > 1 {
		p, err := strconv.Atoi(remoteAddr[1])
		if err == nil {
			size++
			peerPort = p
		}
	}

	useragent := req.UserAgent()
	if len(useragent) > 0 {
		size++
	}

	clientIP := req.Header.Get("X-Forwarded-For")
	if len(clientIP) > 0 {
		size++
	}

	attrs := make([]attribute.KeyValue, 0, size)
	attrs = append(attrs, semconv.HTTPMethod(method))
	attrs = append(attrs, semconv.NetProtocolName(protoName))

	if len(protoVersion) > 0 {
		attrs = append(attrs, semconv.NetProtocolVersion(protoVersion))
	}
	if len(peer) > 0 {
		attrs = append(attrs, semconv.NetSockPeerAddr(peer))
	}
	if peerPort > 0 {
		attrs = append(attrs, semconv.NetSockPeerPort(peerPort))
	}
	if len(useragent) > 0 {
		attrs = append(attrs, semconv.HTTPUserAgent(useragent))
	}
	if len(clientIP) > 0 {
		attrs = append(attrs, semconv.ClientAddress(clientIP))
	}
	if len(clientIP) > 0 {
		attrs = append(attrs, semconv.NetProtocolName(clientIP))
	}

	return attrs
}

func serverStatus(code int) (codes.Code, string) {
	if code < 100 || code >= 600 {
		return codes.Error, fmt.Sprintf("Invalid HTTP status code %d", code)
	}
	if code >= 500 {
		return codes.Error, ""
	}
	return codes.Unset, ""
}
