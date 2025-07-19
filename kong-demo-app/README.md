# Kong Demo App

这是一个用于学习 Kong 网关的演示项目，包含一个简单的 Go HTTP 服务。

## 项目结构

```
.
├── main.go          # Go HTTP 服务
├── go.mod           # Go 模块文件
├── Dockerfile       # Docker 构建文件
└── README.md        # 项目说明
```

## 服务接口

### 直接访问（端口 3000）

- `GET /` - 根路径，返回服务信息
- `GET /health` - 健康检查
- `GET /api/users` - 获取用户列表
- `GET /api/users/{id}` - 获取指定用户

### 通过 Kong 网关访问（端口 8000）

- `GET /demo/` - 根路径
- `GET /demo/health` - 健康检查
- `GET /demo/api/users` - 获取用户列表
- `GET /demo/api/users/{id}` - 获取指定用户

## 启动服务

### 1. 启动 Go 服务

```bash
go run main.go
```

服务将在 `http://localhost:3000` 启动。

### 2. 启动 Kong 网关

确保 Kong 和 PostgreSQL 已经启动：

```bash
# 检查 Kong 状态
curl http://localhost:8001/

# 检查 Kong 代理
curl http://localhost:8000/
```

### 3. 配置 Kong 服务

```bash
# 添加服务
curl -X POST http://localhost:8001/services \
  -H "Content-Type: application/json" \
  -d '{
    "name": "demo-app",
    "url": "http://host.docker.internal:3000"
  }'

# 添加路由
curl -X POST http://localhost:8001/services/demo-app/routes \
  -H "Content-Type: application/json" \
  -d '{
    "name": "demo-app-route",
    "paths": ["/demo"],
    "strip_path": true
  }'
```

## 测试

### 直接访问服务

```bash
curl http://localhost:3000/
curl http://localhost:3000/health
curl http://localhost:3000/api/users
```

### 通过 Kong 网关访问

```bash
curl http://localhost:8000/demo/health
curl http://localhost:8000/demo/api/users
```

## Kong 管理

### 查看服务列表

```bash
curl http://localhost:8001/services
```

### 查看路由列表

```bash
curl http://localhost:8001/routes
```

### 查看插件列表

```bash
curl http://localhost:8001/plugins
```

## 学习要点

1. **服务注册**：将后端服务注册到 Kong
2. **路由配置**：配置访问路径和转发规则
3. **路径处理**：`strip_path` 参数控制路径前缀处理
4. **网关代理**：通过 Kong 统一对外提供服务

## 下一步

- 添加认证插件
- 配置限流插件
- 添加日志插件
- 配置 SSL/TLS
- 学习更多 Kong 插件功能 