package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func main() {
	// 创建Jaeger导出器
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		log.Fatalf("创建Jaeger导出器失败: %v", err)
	}

	// 创建资源
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName("timing-demo"),
		),
	)
	if err != nil {
		log.Fatalf("创建资源失败: %v", err)
	}

	// 创建trace provider，配置更短的发送间隔
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp,
			trace.WithBatchTimeout(2*time.Second), // 2秒发送一次
			trace.WithMaxExportBatchSize(3),       // 3个spans就发送
		),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	tracer := otel.Tracer("timing-demo")
	ctx := context.Background()

	fmt.Println("=== 数据发送时机演示 ===")
	fmt.Println("配置：每2秒或累积3个spans时发送")
	fmt.Println()

	// 演示1: 累积发送
	fmt.Println("1️⃣ 创建3个spans (应该触发自动发送)")
	for i := 1; i <= 3; i++ {
		_, span := tracer.Start(ctx, fmt.Sprintf("操作-%d", i))
		span.SetAttributes(attribute.Int("序号", i))
		span.End()
		fmt.Printf("   ✓ span-%d 已结束，进入缓冲区\n", i)
		time.Sleep(500 * time.Millisecond) // 给时间让数据发送
	}

	fmt.Println("   🚀 3个spans应该已经自动发送到Jaeger!")
	fmt.Println()

	// 演示2: 时间间隔发送
	fmt.Println("2️⃣ 创建2个spans，等待时间触发发送")
	for i := 4; i <= 5; i++ {
		_, span := tracer.Start(ctx, fmt.Sprintf("操作-%d", i))
		span.SetAttributes(attribute.Int("序号", i))
		span.End()
		fmt.Printf("   ✓ span-%d 已结束，进入缓冲区\n", i)
	}

	fmt.Println("   ⏰ 等待3秒让批处理器自动发送...")
	time.Sleep(3 * time.Second)
	fmt.Println("   🚀 剩余spans应该已经自动发送到Jaeger!")
	fmt.Println()

	// 演示3: 手动强制发送
	fmt.Println("3️⃣ 创建1个span，使用Shutdown强制发送")
	_, span := tracer.Start(ctx, "最后的操作")
	span.SetAttributes(attribute.String("类型", "强制发送演示"))
	span.End()
	fmt.Println("   ✓ span 已结束，进入缓冲区")

	fmt.Println("   💪 调用Shutdown强制发送...")
	if err := tp.Shutdown(ctx); err != nil {
		log.Printf("Shutdown出错: %v", err)
	}
	fmt.Println("   🚀 所有数据已强制发送到Jaeger!")
	fmt.Println()

	fmt.Println("✅ 演示完成! 请到Jaeger UI查看:")
	fmt.Println("   🌐 http://localhost:16686")
	fmt.Println("   📊 服务名: timing-demo")
	fmt.Println("   🔍 你应该能看到6个不同的trace记录")
}
