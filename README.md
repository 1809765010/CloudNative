# Cloud-Native 微服务示例项目

本项目演示了基于 Go、gRPC、Consul 的服务注册与发现，包含服务注册端（service-register）和服务发现/网关端（service-discover）。

## 目录结构

```
Cloud-Native/
├── go.mod
├── go.sum
├── service-register/   # 服务注册与健康检查
│   ├── main.go
│   └── ...
└── service-discover/   # 服务发现与 HTTP 网关
    ├── main.go
    └── ...
```

---

## 依赖环境

- Go 1.18+
- Consul 1.9+（本地或远程均可）
- 已安装依赖：`go mod tidy`

---

## 启动 Consul

本地启动（默认端口 8500）：

```bash
consul agent -dev
```

---

## 启动服务注册端（可多实例）

分别在不同终端运行：

```bash
cd service-register
go run main.go -id=hello-service-1 -port=50051
go run main.go -id=hello-service-2 -port=50052
go run main.go -id=hello-service-3 -port=50053
go run main.go -id=hello-service-4 -port=50054
```

---

## 启动服务发现/HTTP 网关

```bash
cd service-discover
go run main.go
```

---

## 通过 HTTP 网关访问服务

多次执行，观察后端 gRPC 服务实例轮询变化：

```bash
curl "http://localhost:8080/hello?name=ConsulClient"
```

---

## 主要功能说明

- **服务注册**：服务启动时自动注册到 Consul，并支持健康检查与自动注销。
- **服务发现**：HTTP 网关通过 Consul 动态发现 gRPC 服务，采用 round_robin 负载均衡。
- **健康检查**：Consul 定期检查服务健康，异常自动剔除。

---

## 参考

- [Consul 官方文档](https://www.consul.io/docs)
- [gRPC 官方文档](https://grpc.io/docs/)
- [Gin 官方文档](https://gin-gonic.com/)

---

如需更多帮助或有定制需求，欢迎随时提问！ 