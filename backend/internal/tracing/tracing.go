package tracing

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// TracerProvider 全局追踪器提供者
var TracerProvider *sdktrace.TracerProvider

// InitTracing 初始化 OpenTelemetry 链路追踪
// endpoint: Jaeger OTLP gRPC 地址，例如 "localhost:4317"
func InitTracing(serviceName, endpoint string) (func(context.Context) error, error) {
	// 创建 OTLP gRPC 导出器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(), // 本地开发使用 HTTP，线上请使用 TLS
	)
	if err != nil {
		return nil, fmt.Errorf("创建 OTLP 导出器失败: %w", err)
	}

	// 创建资源信息
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("1.0.0"),
			attribute.String("environment", "development"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("创建资源失败: %w", err)
	}

	// 创建 TracerProvider
	TracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// 设置全局 TracerProvider
	otel.SetTracerProvider(TracerProvider)

	// 设置全局 propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return TracerProvider.Shutdown, nil
}

// StartSpan 开始一个新的追踪 span
func StartSpan(ctx context.Context, name string, opts ...TraceOption) (context.Context, trace.Span) {
	tracer := TracerProvider.Tracer("gomall")

	options := &traceOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var attrs []attribute.KeyValue
	if len(options.attrs) > 0 {
		attrs = options.attrs
	}

	ctx, span := tracer.Start(ctx, name, trace.WithAttributes(attrs...))

	return ctx, span
}

// TraceOption Trace 选项
type TraceOption func(*traceOptions)

type traceOptions struct {
	attrs []attribute.KeyValue
}

// WithAttributes 设置 span 属性
func WithAttributes(attrs ...attribute.KeyValue) TraceOption {
	return func(o *traceOptions) {
		o.attrs = attrs
	}
}

// RecordError 记录错误
func RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
}

// GetTracer 获取 Tracer
func GetTracer(name string) trace.Tracer {
	return TracerProvider.Tracer(name)
}

// Shutdown 关闭追踪器
func Shutdown(ctx context.Context) error {
	if TracerProvider != nil {
		return TracerProvider.Shutdown(ctx)
	}
	return nil
}

// 添加自定义属性
func WithUserID(userID uint) TraceOption {
	return WithAttributes(attribute.Int64("user_id", int64(userID)))
}

// WithProductID 添加商品ID属性
func WithProductID(productID uint) TraceOption {
	return WithAttributes(attribute.Int64("product_id", int64(productID)))
}

// WithOrderNo 添加订单号属性
func WithOrderNo(orderNo string) TraceOption {
	return WithAttributes(attribute.String("order_no", orderNo))
}
