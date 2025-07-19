# Cloud-Native 学习项目

这是一个用于学习云原生技术的项目集合，包含多个微服务和网关技术的演示。

## 项目结构

```
Cloud-Native/
├── kong-demo-app/          # Kong 网关演示项目
│   ├── main.go            # Go HTTP 服务
│   ├── go.mod             # Go 模块文件
│   ├── Dockerfile         # Docker 构建文件
│   └── README.md          # Kong 项目说明
├── kong-example/           # Kong 网关示例
├── consul-example/         # Consul 服务发现示例
│   ├── service-discover/   # 服务发现
│   ├── service-register/   # 服务注册
│   └── README.md          # Consul 项目说明
└── README.md              # 项目总览
```

## 技术栈

- **API 网关**: Kong
- **服务发现**: Consul
- **编程语言**: Go
- **容器化**: Docker
- **微服务架构**: 分布式系统

## 快速开始

### 1. Kong 网关演示

```bash
cd kong-demo-app
go run main.go
```

访问：
- 服务：http://localhost:3000
- 网关：http://localhost:8000/demo

### 2. Consul 服务发现

```bash
cd consul-example
# 查看具体说明
cat README.md
```

## 学习目标

1. **API 网关技术**
   - Kong 网关的安装和配置
   - 服务注册和路由管理
   - 插件系统使用

2. **服务发现**
   - Consul 集群搭建
   - 服务注册和发现
   - 健康检查和负载均衡

3. **微服务架构**
   - 服务间通信
   - 分布式系统设计
   - 容器化部署

## 开发环境

- macOS 15.5.0
- Go 1.21+
- Docker Desktop
- Kong Gateway
- Consul

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个学习项目。

## 许可证

MIT License 