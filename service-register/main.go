package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "cloudnative/service-register/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	consulapi "github.com/hashicorp/consul/api"
)

type server struct {
	pb.UnimplementedHelloServiceServer
	Address string
	Port    int
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in)
	// 假设你有 s.Address 和 s.Port 字段，或者用环境变量/启动参数传递
	return &pb.HelloReply{
		Message: fmt.Sprintf("Hello %s, from %s:%d", in.Name, s.Address, s.Port),
	}, nil
}

func registerServiceWithConsul(serviceID, serviceName, address string, port int) error {
	config := consulapi.DefaultConfig()
	config.Address = "127.0.0.1:8500" // 本地 Consul 地址
	client, err := consulapi.NewClient(config)
	if err != nil {
		return err
	}

	registration := &consulapi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: address,
		Port:    port,
		Check: &consulapi.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d/%s", address, port, serviceName), // Consul 通过 gRPC 方式健康检查服务，格式为 ip:端口/服务名
			Interval:                       "10s",                                               // 健康检查间隔时间，每 10 秒检查一次
			Timeout:                        "1s",                                                // 健康检查超时时间，1 秒内无响应视为失败
			DeregisterCriticalServiceAfter: "1m",                                                // 连续处于 critical 状态超过 1 分钟自动注销该服务实例
		},
	}

	return client.Agent().ServiceRegister(registration)
}

// 获取第一个非回环的IPv4地址
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", nil
}

func main() {
	var (
		serviceID string
		port      int
	)
	flag.StringVar(&serviceID, "id", "hello-service-1", "服务ID")
	flag.IntVar(&port, "port", 50051, "服务端口")
	flag.Parse()

	// 注册到 Consul 前先获取本机 IP
	ip, err := GetLocalIP()
	if err != nil {
		log.Fatalf("无法获取本机IP: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterHelloServiceServer(grpcServer, &server{
		Address: ip,
		Port:    port,
	})

	// 注册健康检查服务
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("hello-service", healthpb.HealthCheckResponse_SERVING)

	// 注册到 Consul
	err = registerServiceWithConsul(serviceID, "hello-service", ip, port)
	if err != nil {
		log.Fatalf("failed to register service with consul: %v", err)
	}
	log.Printf("服务[%s]已注册到 Consul，监听端口: %d", serviceID, port)

	// 优雅退出：监听信号，注销 Consul 服务
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("收到退出信号，注销 Consul 服务...")
		config := consulapi.DefaultConfig()
		config.Address = "127.0.0.1:8500"
		client, _ := consulapi.NewClient(config)
		client.Agent().ServiceDeregister(serviceID)
		log.Println("Consul 注销完成，优雅退出")
		os.Exit(0)
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
