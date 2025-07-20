#!/bin/bash

echo "🚀 OpenTelemetry HTTP API 测试脚本"
echo "====================================="
echo

# 测试函数
test_api() {
    local method=$1
    local endpoint=$2
    local description=$3
    
    echo "📍 测试: $description"
    echo "🔗 请求: $method $endpoint"
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -i "http://localhost:8080$endpoint")
    else
        response=$(curl -s -i -X "$method" "http://localhost:8080$endpoint")
    fi
    
    # 提取响应头中的Trace ID
    trace_id=$(echo "$response" | grep -i "X-Trace-ID" | cut -d' ' -f2 | tr -d '\r\n')
    
    # 提取响应体
    body=$(echo "$response" | sed -n '/^{/,$p')
    
    echo "🏷️  响应头中的Trace ID: $trace_id"
    echo "📦 响应内容:"
    echo "$body" | jq '.' 2>/dev/null || echo "$body"
    echo "----------------------------------------"
    echo
}

echo "正在测试各个API端点的Trace ID功能..."
echo

# 测试各个端点
test_api "GET" "/" "首页接口"
test_api "GET" "/users" "获取用户列表"
test_api "GET" "/products" "获取商品列表"
test_api "POST" "/orders" "创建订单"
test_api "GET" "/health" "健康检查"

echo "✅ 测试完成!"
echo
echo "🔍 观察要点:"
echo "1. 每个请求都有唯一的Trace ID"
echo "2. 响应头中包含 'X-Trace-ID'"
echo "3. 响应体JSON中也包含相同的 'trace_id'"
echo "4. 可以用这个Trace ID在Jaeger UI中查找对应的追踪"
echo
echo "📊 打开Jaeger UI查看追踪:"
echo "   🌐 http://localhost:16686"
echo "   📋 服务名: http-api-server"
echo "   🔍 复制上面的任一Trace ID进行搜索" 