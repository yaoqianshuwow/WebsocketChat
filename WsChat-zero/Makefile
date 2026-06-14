.PHONY: all gen build build-all run-user run-gateway run-friend run-file run-voice run-msg-forward run-all clean proto docker

# 默认目标：生成代码 + 编译所有服务
all: gen build-all

# ── 代码生成 ──
gen:
	@echo ">>> 生成 proto 代码..."
	goctl rpc protoc proto/user.proto --proto_path=proto --go_out=app/user/pb --go-grpc_out=app/user/pb --zrpc_out=app/user --client
	goctl rpc protoc proto/message.proto --proto_path=proto --go_out=app/msg-forward/pb --go-grpc_out=app/msg-forward/pb --zrpc_out=app/msg-forward --client
	goctl rpc protoc proto/friend.proto --proto_path=proto --go_out=app/friend/pb --go-grpc_out=app/friend/pb --zrpc_out=app/friend --client
	goctl rpc protoc proto/file.proto --proto_path=proto --go_out=app/file/pb --go-grpc_out=app/file/pb --zrpc_out=app/file --client
	goctl rpc protoc proto/voice.proto --proto_path=proto --go_out=app/voice/pb --go-grpc_out=app/voice/pb --zrpc_out=app/voice --client
	@echo ">>> 生成 API 代码..."
	goctl api go --api app/gateway/gateway.api --dir app/gateway
	@echo ">>> 生成完成"

# ── 编译（工作空间模式）──
build-all:
	@echo ">>> 编译所有服务..."
	go build -o bin/user-rpc ./app/user
	go build -o bin/gateway ./app/gateway
	go build -o bin/friend-rpc ./app/friend
	go build -o bin/file-rpc ./app/file
	go build -o bin/voice-rpc ./app/voice
	go build -o bin/msg-forward ./app/msg-forward
	go build -o bin/msg-store ./app/msg-store/cmd
	@echo ">>> 全部编译完成，输出在 bin/ 目录"

build: build-all

# ── 单服务编译 ──
build-user:
	go build -o bin/user-rpc ./app/user
build-gateway:
	go build -o bin/gateway ./app/gateway
build-friend:
	go build -o bin/friend-rpc ./app/friend
build-file:
	go build -o bin/file-rpc ./app/file
build-voice:
	go build -o bin/voice-rpc ./app/voice
build-msg-forward:
	go build -o bin/msg-forward ./app/msg-forward
build-msg-store:
	go build -o bin/msg-store ./app/msg-store/cmd

# ── 本地运行 ──
run-user:
	cd app/user && go run user.go -f etc/user.yaml

run-gateway:
	cd app/gateway && go run gateway.go -f etc/gateway-api.yaml

run-friend:
	cd app/friend && go run friend.go -f etc/friend.yaml

run-file:
	cd app/file && go run file.go -f etc/file.yaml

run-voice:
	cd app/voice && go run voice.go -f etc/voice.yaml

run-msg-forward:
	cd app/msg-forward && go run *.go -f etc/msg.yaml

run-msg-store:
	cd app/msg-store && go run cmd/main.go -f etc/msg-store.yaml

run-all:
	@echo ">>> 启动所有服务（需要先启动 etcd、mysql、redis、kafka）..."
	cd app/user && go run user.go -f etc/user.yaml &
	cd app/friend && go run friend.go -f etc/friend.yaml &
	cd app/file && go run file.go -f etc/file.yaml &
	cd app/voice && go run voice.go -f etc/voice.yaml &
	cd app/msg-forward && go run *.go -f etc/msg.yaml &
	cd app/msg-store && go run cmd/main.go -f etc/msg-store.yaml &
	sleep 3
	cd app/gateway && go run gateway.go -f etc/gateway-api.yaml &

# ── Docker ──
docker-build-all:
	docker build -t wschat-gateway -f deploy/Dockerfile.gateway .
	docker build -t wschat-user -f deploy/Dockerfile.user .
	docker build -t wschat-friend -f deploy/Dockerfile.friend .
	docker build -t wschat-file -f deploy/Dockerfile.file .
	docker build -t wschat-voice -f deploy/Dockerfile.voice .
	docker build -t wschat-msg-forward -f deploy/Dockerfile.msg-forward .
	docker build -t wschat-msg-store -f deploy/Dockerfile.msg-store .
	docker build -t wschat-web -f deploy/Dockerfile.web .

docker-up:
	docker compose -f deploy/docker-compose.yml up -d --build

docker-down:
	docker compose -f deploy/docker-compose.yml down

docker-logs:
	docker compose -f deploy/docker-compose.yml logs -f

# ── 清理 ──
clean:
	rm -rf bin/
	go clean -cache
