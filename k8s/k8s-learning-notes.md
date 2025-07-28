# Kubernetes 学习笔记

## 日期：2025-07-28

### 今天学到的核心知识

#### 1. 负载均衡机制
- **kube-proxy默认使用iptables模式**
- **负载均衡算法**：基于源IP的哈希算法，不是真正的轮询
- **会话亲和性**：相同客户端IP的请求总是路由到同一个Pod
- **测试方法**：使用`curl --no-keepalive -H "Connection: close"`禁用HTTP长连接

#### 2. Pod反亲和性配置
```yaml
affinity:
  podAntiAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
    - weight: 100
      podAffinityTerm:
        labelSelector:
          matchExpressions:
          - key: app
            operator: In
            values:
            - cache-service
        topologyKey: kubernetes.io/hostname
```

#### 3. Service配置
- **sessionAffinity: None** - 禁用会话亲和性
- **NodePort映射** - 在kind-config.yaml中配置
- **Endpoints** - 查看Service的后端Pod

#### 4. 网络知识
- **iptables**：Linux网络包过滤和NAT工具
- **kube-proxy**：使用iptables实现Service负载均衡
- **哈希算法**：基于源IP、目标IP、端口计算哈希值

#### 5. 测试负载均衡的命令
```bash
# 禁用HTTP长连接测试负载均衡
for i in {1..10}; do 
  echo "请求 $i:"; 
  curl -s --no-keepalive -H "Connection: close" \
    http://localhost:30080/api/v1/data/loadtest_$i | jq -r '.pod'; 
done
```

#### 6. 重要发现
- **iptables模式不是真正的轮询**，而是基于哈希的随机分发
- **相同源IP总是路由到同一个Pod**，这是默认行为
- **Pod反亲和性**可以确保Pod分布在不同节点上

### 下一步学习计划
1. 搭建博客系统
2. 深入学习Kubernetes网络
3. 实践CI/CD流程
4. 学习Service Mesh（如Istio）

### 有用的命令
```bash
# 查看Pod分布
kubectl get pods -n cache-service -o wide

# 查看Service配置
kubectl get svc -n cache-service -o yaml

# 查看Endpoints
kubectl get endpoints -n cache-service

# 测试负载均衡
curl -s --no-keepalive -H "Connection: close" http://localhost:30080/api/v1/data/test
``` 