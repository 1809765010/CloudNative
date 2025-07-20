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
	oteltrace "go.opentelemetry.io/otel/trace"
)

func main() {
	// 创建Jaeger导出器
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		log.Fatalf("创建Jaeger导出器失败: %v", err)
	}

	// 创建资源信息
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName("go-opentelemetry-demo"),
			semconv.ServiceVersion("v1.0.0"),
			semconv.DeploymentEnvironment("development"),
		),
	)
	if err != nil {
		log.Fatalf("创建资源失败: %v", err)
	}

	// 创建trace provider，使用Jaeger导出器
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)

	// 设置全局trace provider
	otel.SetTracerProvider(tp)

	// 创建tracer
	tracer := otel.Tracer("jaeger-example")

	// 创建context
	ctx := context.Background()

	fmt.Println("=== OpenTelemetry + Jaeger 追踪示例 ===")
	fmt.Println("正在发送追踪数据到Jaeger...")
	fmt.Println()

	// 创建根span
	ctx, rootSpan := tracer.Start(ctx, "购物流程")
	rootSpan.SetAttributes(
		attribute.String("用户ID", "user-12345"),
		attribute.String("会话ID", "session-abcde"),
		attribute.String("用户类型", "VIP"),
	)

	// 模拟购物流程
	simulateShoppingFlow(ctx, tracer)

	// 结束根span
	rootSpan.End()

	// 确保所有数据都被发送到Jaeger
	if err := tp.Shutdown(ctx); err != nil {
		log.Printf("关闭trace provider时出错: %v", err)
	}

	fmt.Println("✅ 追踪数据已发送到Jaeger!")
	fmt.Println("🌐 打开浏览器访问: http://localhost:16686")
	fmt.Println("📊 在Jaeger UI中搜索服务: 'go-opentelemetry-demo'")
	fmt.Println()
}

func simulateShoppingFlow(ctx context.Context, tracer oteltrace.Tracer) {
	// 1. 用户登录
	login(ctx, tracer)

	// 2. 浏览商品
	browsProducts(ctx, tracer)

	// 3. 添加到购物车
	addToCart(ctx, tracer)

	// 4. 结算支付
	checkout(ctx, tracer)
}

func login(ctx context.Context, tracer oteltrace.Tracer) {
	ctx, span := tracer.Start(ctx, "用户登录")
	defer span.End()

	span.SetAttributes(
		attribute.String("登录方式", "用户名密码"),
		attribute.Bool("记住我", true),
	)

	// 模拟登录验证时间
	time.Sleep(50 * time.Millisecond)

	// 验证用户
	authenticateUser(ctx, tracer)

	span.AddEvent("登录成功", oteltrace.WithAttributes(
		attribute.String("结果", "成功"),
		attribute.String("用户角色", "VIP"),
	))
}

func authenticateUser(ctx context.Context, tracer oteltrace.Tracer) {
	ctx, span := tracer.Start(ctx, "用户认证")
	defer span.End()

	span.SetAttributes(
		attribute.String("认证方式", "数据库验证"),
		attribute.String("数据库", "MySQL"),
	)

	time.Sleep(30 * time.Millisecond)
	span.AddEvent("密码验证通过")
}

func browsProducts(ctx context.Context, tracer oteltrace.Tracer) {
	ctx, span := tracer.Start(ctx, "浏览商品")
	defer span.End()

	span.SetAttributes(
		attribute.String("商品类别", "电子产品"),
		attribute.Int("浏览页数", 3),
	)

	time.Sleep(200 * time.Millisecond)

	// 模拟加载商品数据
	loadProductData(ctx, tracer)

	span.AddEvent("商品浏览完成", oteltrace.WithAttributes(
		attribute.Int("查看商品数量", 15),
	))
}

func loadProductData(ctx context.Context, tracer oteltrace.Tracer) {
	ctx, span := tracer.Start(ctx, "加载商品数据")
	defer span.End()

	span.SetAttributes(
		attribute.String("数据源", "Redis缓存"),
		attribute.Int("缓存命中率", 85),
	)

	time.Sleep(80 * time.Millisecond)
	span.AddEvent("数据加载完成")
}

func addToCart(ctx context.Context, tracer oteltrace.Tracer) {
	ctx, span := tracer.Start(ctx, "添加到购物车")
	defer span.End()

	span.SetAttributes(
		attribute.String("商品ID", "product-001"),
		attribute.String("商品名", "iPhone 15 Pro"),
		attribute.Int("数量", 1),
		attribute.Float64("价格", 8999.00),
	)

	time.Sleep(100 * time.Millisecond)

	// 检查库存
	checkInventory(ctx, tracer)

	span.AddEvent("商品已添加到购物车", oteltrace.WithAttributes(
		attribute.String("状态", "成功"),
	))
}

func checkInventory(ctx context.Context, tracer oteltrace.Tracer) {
	ctx, span := tracer.Start(ctx, "检查库存")
	defer span.End()

	span.SetAttributes(
		attribute.String("仓库", "华东仓"),
		attribute.Int("可用库存", 50),
	)

	time.Sleep(60 * time.Millisecond)
	span.AddEvent("库存检查完成")
}

func checkout(ctx context.Context, tracer oteltrace.Tracer) {
	ctx, span := tracer.Start(ctx, "结算支付")
	defer span.End()

	span.SetAttributes(
		attribute.String("支付方式", "支付宝"),
		attribute.Float64("订单金额", 8999.00),
		attribute.String("优惠券", "VIP95折"),
		attribute.Float64("实付金额", 8549.05),
	)

	time.Sleep(300 * time.Millisecond)

	// 处理支付
	processPayment(ctx, tracer)

	span.AddEvent("订单支付成功", oteltrace.WithAttributes(
		attribute.String("订单号", "ORDER-20240720-001"),
		attribute.String("支付状态", "已支付"),
	))
}

func processPayment(ctx context.Context, tracer oteltrace.Tracer) {
	ctx, span := tracer.Start(ctx, "处理支付")
	defer span.End()

	span.SetAttributes(
		attribute.String("支付网关", "支付宝网关"),
		attribute.String("交易ID", "ALIPAY-123456789"),
	)

	time.Sleep(150 * time.Millisecond)
	span.AddEvent("支付处理完成")
}
