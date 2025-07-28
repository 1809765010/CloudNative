.PHONY: build push deploy test

# 构建镜像
build:
	cd cache-service && make build
	cd cache-service && docker build -t cache-service:latest .

# 推送镜像
push: build
	docker tag cache-service:latest 1809765010/cache-service:latest
	docker push 1809765010/cache-service:latest

# 部署到集群
deploy: push
	kubectl rollout restart deployment cache-service -n cache-service
	kubectl rollout status deployment cache-service -n cache-service

# 测试
test:
	curl -s http://localhost:30080/api/v1/data/test | jq .

# 一键部署
all: deploy test 