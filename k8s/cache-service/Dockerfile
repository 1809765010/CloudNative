FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# 复制构建好的Linux二进制文件
COPY bin/cache-service-linux ./server

# 创建配置目录
RUN mkdir -p /data/conf

EXPOSE 8000 9000

CMD ["./server", "-conf", "/data/conf/config.yaml"]
