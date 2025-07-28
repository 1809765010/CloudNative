package service

import (
	"context"
	"net"
	"os"
	"time"

	v1 "cache-service/api/cache/v1"
	"cache-service/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// CacheService is a cache service.
type CacheService struct {
	v1.UnimplementedCacheServiceServer

	uc  *biz.CacheUsecase
	log *log.Helper
}

// NewCacheService creates a new cache service.
func NewCacheService(uc *biz.CacheUsecase, logger log.Logger) *CacheService {
	return &CacheService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

// GetData implements cache.CacheServiceServer.
func (s *CacheService) GetData(ctx context.Context, req *v1.GetDataRequest) (*v1.GetDataReply, error) {
	s.log.WithContext(ctx).Infof("GetData request: key=%s", req.Key)

	data, source, err := s.uc.GetData(ctx, req.Key)
	if err != nil {
		s.log.WithContext(ctx).Errorf("GetData error: %v", err)
		return nil, err
	}

	// 添加调试日志：数据获取成功
	s.log.WithContext(ctx).Infof("GetData success: key=%s, source=%s", req.Key, source)

	// 获取 pod 信息
	podName := os.Getenv("HOSTNAME")
	podIP := os.Getenv("POD_IP")
	if podIP == "" {
		// 如果没有 POD_IP 环境变量，尝试获取本机 IP
		addrs, err := net.InterfaceAddrs()
		if err == nil {
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						podIP = ipnet.IP.String()
						break
					}
				}
			}
		}
		if podIP == "" {
			podIP = "unknown-ip"
		}
	}
	if podName == "" {
		podName = "unknown-pod"
	}

	// 添加调试日志
	s.log.WithContext(ctx).Infof("Pod info: name=%s, ip=%s", podName, podIP)

	reply := &v1.GetDataReply{
		Key:    data.Key,
		Value:  data.Value,
		Source: source,
		Pod:    podIP, // 使用 pod IP 地址
	}

	// 添加调试日志：返回结果
	s.log.WithContext(ctx).Infof("GetData returning: key=%s, pod=%s", reply.Key, reply.Pod)

	return reply, nil
}

// SetData implements cache.CacheServiceServer.
func (s *CacheService) SetData(ctx context.Context, req *v1.SetDataRequest) (*v1.SetDataReply, error) {
	s.log.WithContext(ctx).Infof("SetData request: key=%s, value=%s", req.Key, req.Value)

	err := s.uc.SetData(ctx, req.Key, req.Value)
	if err != nil {
		s.log.WithContext(ctx).Errorf("SetData error: %v", err)
		return &v1.SetDataReply{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &v1.SetDataReply{
		Success: true,
		Message: "Data set successfully",
	}, nil
}

// HealthCheck implements cache.CacheServiceServer.
func (s *CacheService) HealthCheck(ctx context.Context, req *v1.HealthCheckRequest) (*v1.HealthCheckReply, error) {
	s.log.WithContext(ctx).Info("Health check requested")

	// 检查Redis连接
	redisHealthy := s.uc.CheckRedisHealth(ctx)

	// 检查数据库连接
	dbHealthy := s.uc.CheckDBHealth(ctx)

	status := "healthy"
	if !redisHealthy || !dbHealthy {
		status = "unhealthy"
		s.log.WithContext(ctx).Warnf("Health check failed: redis=%v, db=%v", redisHealthy, dbHealthy)
	}

	return &v1.HealthCheckReply{
		Status:    status,
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   "1.0.0",
	}, nil
}
