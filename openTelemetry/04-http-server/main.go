package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// 响应结构体
type APIResponse struct {
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	TraceID   string      `json:"trace_id"`
	Timestamp string      `json:"timestamp"`
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func main() {
	// 初始化OpenTelemetry
	cleanup := initTracing()
	defer cleanup()

	// 创建HTTP服务器
	mux := http.NewServeMux()

	// 注册路由
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/users", usersHandler)
	mux.HandleFunc("/products", productsHandler)
	mux.HandleFunc("/orders", ordersHandler)
	mux.HandleFunc("/health", healthHandler)

	// 包装所有路由使用tracing中间件
	handler := tracingMiddleware(mux)

	fmt.Println("🚀 HTTP服务启动成功!")
	fmt.Println("📍 服务地址: http://localhost:8080")
	fmt.Println("🔍 可用端点:")
	fmt.Println("   GET  /              - 首页")
	fmt.Println("   GET  /users         - 获取用户列表")
	fmt.Println("   GET  /products      - 获取商品列表")
	fmt.Println("   POST /orders        - 创建订单")
	fmt.Println("   GET  /health        - 健康检查")
	fmt.Println()
	fmt.Println("📊 Jaeger UI: http://localhost:16686")
	fmt.Println("🏷️  每个响应都会在header中包含 X-Trace-ID")
	fmt.Println()

	// 启动服务器
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// 初始化OpenTelemetry
func initTracing() func() {
	// 创建Jaeger导出器
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		log.Fatalf("创建Jaeger导出器失败: %v", err)
	}

	// 创建资源
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName("http-api-server"),
			semconv.ServiceVersion("v1.0.0"),
			semconv.DeploymentEnvironment("development"),
		),
	)
	if err != nil {
		log.Fatalf("创建资源失败: %v", err)
	}

	// 创建trace provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)

	// 设置全局provider和propagator
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("关闭trace provider出错: %v", err)
		}
	}
}

// Tracing中间件 - 关键部分！
func tracingMiddleware(next http.Handler) http.Handler {
	tracer := otel.Tracer("http-server")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从请求中提取trace context (支持分布式追踪)
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		// 创建span
		ctx, span := tracer.Start(ctx, fmt.Sprintf("%s %s", r.Method, r.URL.Path))
		defer span.End()

		// 添加请求信息到span
		span.SetAttributes(
			semconv.HTTPMethod(r.Method),
			semconv.HTTPRoute(r.URL.Path),
			semconv.HTTPUserAgent(r.UserAgent()),
			semconv.HTTPClientIP(r.RemoteAddr),
		)

		// ⭐ 关键：获取trace ID并添加到响应头
		traceID := span.SpanContext().TraceID().String()
		w.Header().Set("X-Trace-ID", traceID)
		w.Header().Set("Access-Control-Expose-Headers", "X-Trace-ID") // 让前端可以读取

		// 将context传递给下一个handler
		r = r.WithContext(ctx)

		// 调用实际的handler
		next.ServeHTTP(w, r)

		// 记录响应状态
		span.AddEvent("请求处理完成")
	})
}

// 首页处理器
func homeHandler(w http.ResponseWriter, r *http.Request) {
	span := oteltrace.SpanFromContext(r.Context())
	span.AddEvent("处理首页请求")

	response := APIResponse{
		Message:   "欢迎使用OpenTelemetry演示API!",
		TraceID:   span.SpanContext().TraceID().String(),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// 用户列表处理器
func usersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := oteltrace.SpanFromContext(ctx)

	span.SetAttributes(attribute.String("handler", "users"))
	span.AddEvent("开始获取用户列表")

	// 模拟数据库查询
	users := fetchUsers(ctx)

	response := APIResponse{
		Message:   "用户列表获取成功",
		Data:      users,
		TraceID:   span.SpanContext().TraceID().String(),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	span.AddEvent("用户列表获取完成", oteltrace.WithAttributes(
		attribute.Int("用户数量", len(users)),
	))

	writeJSONResponse(w, http.StatusOK, response)
}

// 商品列表处理器
func productsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	span := oteltrace.SpanFromContext(ctx)

	span.SetAttributes(attribute.String("handler", "products"))
	span.AddEvent("开始获取商品列表")

	// 模拟数据库查询
	products := fetchProducts(ctx)

	response := APIResponse{
		Message:   "商品列表获取成功",
		Data:      products,
		TraceID:   span.SpanContext().TraceID().String(),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	span.AddEvent("商品列表获取完成", oteltrace.WithAttributes(
		attribute.Int("商品数量", len(products)),
	))

	writeJSONResponse(w, http.StatusOK, response)
}

// 创建订单处理器
func ordersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	span := oteltrace.SpanFromContext(ctx)

	span.SetAttributes(attribute.String("handler", "create_order"))
	span.AddEvent("开始处理订单创建")

	// 模拟订单处理流程
	orderID := processOrder(ctx)

	response := APIResponse{
		Message:   "订单创建成功",
		Data:      map[string]interface{}{"order_id": orderID},
		TraceID:   span.SpanContext().TraceID().String(),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	span.AddEvent("订单创建完成", oteltrace.WithAttributes(
		attribute.String("order_id", orderID),
	))

	writeJSONResponse(w, http.StatusCreated, response)
}

// 健康检查处理器
func healthHandler(w http.ResponseWriter, r *http.Request) {
	span := oteltrace.SpanFromContext(r.Context())
	span.AddEvent("健康检查")

	response := APIResponse{
		Message:   "服务运行正常",
		Data:      map[string]string{"status": "healthy"},
		TraceID:   span.SpanContext().TraceID().String(),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	writeJSONResponse(w, http.StatusOK, response)
}

// 模拟获取用户数据
func fetchUsers(ctx context.Context) []User {
	tracer := otel.Tracer("database")
	ctx, span := tracer.Start(ctx, "fetch_users_from_db")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.operation", "SELECT"),
		attribute.String("db.table", "users"),
	)

	// 模拟数据库查询延迟
	time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)

	users := []User{
		{ID: 1, Name: "张三", Email: "zhangsan@example.com"},
		{ID: 2, Name: "李四", Email: "lisi@example.com"},
		{ID: 3, Name: "王五", Email: "wangwu@example.com"},
	}

	span.AddEvent("数据库查询完成", oteltrace.WithAttributes(
		attribute.Int("返回行数", len(users)),
	))

	return users
}

// 模拟获取商品数据
func fetchProducts(ctx context.Context) []Product {
	tracer := otel.Tracer("database")
	ctx, span := tracer.Start(ctx, "fetch_products_from_db")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.operation", "SELECT"),
		attribute.String("db.table", "products"),
	)

	// 模拟数据库查询延迟
	time.Sleep(time.Duration(30+rand.Intn(70)) * time.Millisecond)

	products := []Product{
		{ID: 1, Name: "iPhone 15 Pro", Price: 8999.00},
		{ID: 2, Name: "MacBook Pro", Price: 18999.00},
		{ID: 3, Name: "AirPods Pro", Price: 1999.00},
	}

	span.AddEvent("商品数据查询完成", oteltrace.WithAttributes(
		attribute.Int("返回商品数", len(products)),
	))

	return products
}

// 模拟订单处理
func processOrder(ctx context.Context) string {
	tracer := otel.Tracer("business")
	ctx, span := tracer.Start(ctx, "process_order")
	defer span.End()

	orderID := fmt.Sprintf("ORDER-%d", time.Now().Unix())

	span.SetAttributes(
		attribute.String("order.id", orderID),
		attribute.String("order.status", "processing"),
	)

	// 模拟订单处理步骤
	validateOrder(ctx)
	calculatePrice(ctx)
	saveOrder(ctx, orderID)

	span.AddEvent("订单处理流程完成")
	return orderID
}

func validateOrder(ctx context.Context) {
	tracer := otel.Tracer("business")
	_, span := tracer.Start(ctx, "validate_order")
	defer span.End()

	time.Sleep(20 * time.Millisecond)
	span.AddEvent("订单验证通过")
}

func calculatePrice(ctx context.Context) {
	tracer := otel.Tracer("business")
	_, span := tracer.Start(ctx, "calculate_price")
	defer span.End()

	time.Sleep(15 * time.Millisecond)
	span.AddEvent("价格计算完成")
}

func saveOrder(ctx context.Context, orderID string) {
	tracer := otel.Tracer("database")
	_, span := tracer.Start(ctx, "save_order_to_db")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.operation", "INSERT"),
		attribute.String("db.table", "orders"),
		attribute.String("order.id", orderID),
	)

	time.Sleep(30 * time.Millisecond)
	span.AddEvent("订单保存成功")
}

// 写入JSON响应
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON编码错误: %v", err)
		http.Error(w, "内部服务器错误", http.StatusInternalServerError)
	}
}
