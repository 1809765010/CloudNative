package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "Cloud-Native/service-register/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	consulapi "github.com/hashicorp/consul/api"
)

type server struct {
	pb.UnimplementedHelloServiceServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
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
			GRPC:                           fmt.Sprintf("%s:%d/%s", address, port, serviceName),
			Interval:                       "10s",
			Timeout:                        "1s",
			DeregisterCriticalServiceAfter: "1m",
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
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterHelloServiceServer(grpcServer, &server{})

	// 注册健康检查服务
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("hello-service", healthpb.HealthCheckResponse_SERVING)

	// 注册到 Consul
	ip, err := GetLocalIP()
	if err != nil {
		log.Fatalf("无法获取本机IP: %v", err)
	}
	err = registerServiceWithConsul("hello-service-1", "hello-service", ip, 50051)
	if err != nil {
		log.Fatalf("failed to register service with consul: %v", err)
	}
	log.Println("服务已注册到 Consul")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
