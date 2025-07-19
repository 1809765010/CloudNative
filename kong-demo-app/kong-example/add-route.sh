#!/bin/bash

# Kong Admin API 地址 (根据您的实际配置修改)
KONG_ADMIN_URL="http://localhost:8001"

# 服务名称
SERVICE_NAME="test-grey"

# 首先获取服务的ID
echo "获取服务 $SERVICE_NAME 的ID..."
SERVICE_ID=$(curl -s "$KONG_ADMIN_URL/services/$SERVICE_NAME" | jq -r '.id')

if [ "$SERVICE_ID" = "null" ] || [ -z "$SERVICE_ID" ]; then
    echo "错误: 找不到服务 $SERVICE_NAME"
    exit 1
fi

echo "服务ID: $SERVICE_ID"

# 添加路由
echo "为服务 $SERVICE_NAME 添加路由..."

curl -X POST "$KONG_ADMIN_URL/services/$SERVICE_NAME/routes" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-grey-route",
    "protocols": ["http"],
    "hosts": ["localhost"],
    "paths": ["/detail"],
    "headers": {
      "x-api-version": ["v1"],
      "content-type": ["application/json"]
    },
    "strip_path": false,
    "preserve_host": false,
    "regex_priority": 0,
    "https_redirect_status_code": 426,
    "path_handling": "v1"
  }'

echo ""
echo "路由添加完成！" 