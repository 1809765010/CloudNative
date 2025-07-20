package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func main() {
	// 创建一个控制台输出的trace exporter
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatalf("创建exporter失败: %v", err)
	}

	// 创建trace provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
	)

	// 设置全局trace provider
	otel.SetTracerProvider(tp)

	// 创建tracer
	tracer := otel.Tracer("basic-example")

	// 创建一个context
	ctx := context.Background()

	// 演示基础的span创建
	fmt.Println("=== OpenTelemetry 基础追踪示例 ===")
	fmt.Println()

	// 创建一个根span
	ctx, rootSpan := tracer.Start(ctx, "main-operation")
	rootSpan.SetAttributes(
		attribute.String("user.id", "user123"),
		attribute.String("operation.type", "demo"),
	)

	// 模拟一些工作
	doWork(ctx, tracer)

	// 结束根span
	rootSpan.End()

	// 确保所有数据都被导出
	tp.Shutdown(ctx)

	fmt.Println()
	fmt.Println("=== 追踪完成! 查看上面的JSON输出了解trace信息 ===")
}

func doWork(ctx context.Context, tracer oteltrace.Tracer) {
	// 创建子span
	ctx, span := tracer.Start(ctx, "do-work")
	defer span.End()

	span.SetAttributes(
		attribute.String("work.type", "processing"),
		attribute.Int("work.items", 100),
	)

	// 模拟处理时间
	time.Sleep(100 * time.Millisecond)

	// 调用另一个函数
	processData(ctx, tracer)

	// 添加事件
	span.AddEvent("work completed", oteltrace.WithAttributes(
		attribute.String("status", "success"),
	))
}

func processData(ctx context.Context, tracer oteltrace.Tracer) {
	// 创建另一个子span
	ctx, span := tracer.Start(ctx, "process-data")
	defer span.End()

	span.SetAttributes(
		attribute.String("data.source", "database"),
		attribute.Int("data.size", 1024),
	)

	// 模拟数据处理
	time.Sleep(50 * time.Millisecond)

	span.AddEvent("data processed")
}
