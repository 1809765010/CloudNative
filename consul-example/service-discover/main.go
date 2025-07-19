package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	pb "cloudnative/service-register/pb"

	"github.com/gin-gonic/gin"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"google.golang.org/grpc"
)

func main() {
	// 建立全局 gRPC 连接
	conn, err := grpc.Dial(
		"consul://127.0.0.1:8500/hello-service?healthy=true",
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		log.Fatalf("无法连接服务: %v", err)
	}
	defer conn.Close()
	client := pb.NewHelloServiceClient(conn)

	// 启动 HTTP 服务
	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		name := c.DefaultQuery("name", "ConsulClient")
		resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: name})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.String(http.StatusOK, resp.Message)
	})

	fmt.Println("HTTP 网关已启动，访问: http://localhost:8080/hello?name=xxx")
	r.Run(":8080")
}
