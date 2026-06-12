# WsChat-zero

基于 **go-zero** 框架重构的 WebSocket 即时通讯微服务系统。

## 架构

```
客户端 (React) ── HTTP/WS ──→ 网关服务 (gateway)
                                  │
                  ┌───────────────┼───────────────┐
                  │ gRPC         │ gRPC         │ gRPC
             ┌────▼────┐   ┌────▼─────┐  ┌─────▼──────┐
             │ 用户服务 │   │ 消息转发  │  │  好友服务  │
             │ (user)  │   │(msg-fwd) │  │ (friend)   │
             └────┬────┘   └────┬─────┘  └────────────┘
                  │              │
           ┌──────┼──────┐  ┌───┴───────┐
           │ MySQL│Redis │  │   Kafka   │
           └──────┴──────┘  └───┬───────┘
                                    │
                           ┌────────┴────────┐
                           ▼                 ▼
                    ┌─────────────┐  ┌──────────────┐
                    │ 消息存储     │  │  文件存储    │
                    │ (msg-store) │  │  (file)      │
                    └──────┬──────┘  └──────────────┘
                           │
                    ┌──────┴──────┐
                    │ MySQL│  ES  │
                    └──────┴──────┘
```

## 子服务

| 服务 | 端口 | 协议 | 说明 |
|------|------|------|------|
| gateway | 8888 | HTTP + WS | 网关入口，鉴权后转发 gRPC |
| user | 9091 | gRPC | 用户管理，JWT认证 |
| msg-forward | 9092 | gRPC | 消息转发，Kafka发布 |
| msg-store | - | Kafka consumer | 消息持久化，ES+MySQL |
| friend | 9093 | gRPC | 好友/群聊管理 |
| file | 9094 | gRPC | 文件/语音存储 |
| voice | 9095 | gRPC | 百度语音识别 |

## 技术栈

- **框架**: go-zero v1.7+ （goctl 代码生成）
- **服务发现**: etcd
- **数据库**: MySQL (GORM) + Elasticsearch
- **缓存**: Redis (go-zero redisx)
- **消息队列**: Kafka (segmentio/kafka-go)
- **通信**: gRPC + Protocol Buffers
- **WebSocket**: gorilla/websocket
- **部署**: Docker + Docker Compose

## 快速启动

```bash
# 1. 启动基础设施
docker compose -f deploy/docker-compose.yml up -d mysql redis etcd kafka

# 2. 启动子服务
make run-user     # 用户服务
make run-gateway  # 网关服务

# 3. 或一键 Docker 启动
make docker-up
```

## 目录结构

```
WsChat-zero/
├── app/
│   ├── gateway/       # 网关 (go-zero api)
│   ├── user/           # 用户管理 (go-zero rpc)
│   ├── msg-forward/    # 消息转发 (go-zero rpc)
│   ├── msg-store/      # 消息存储 (Kafka消费者)
│   ├── friend/         # 好友管理 (go-zero rpc)
│   ├── file/           # 文件存储 (go-zero rpc)
│   └── voice/          # 语音转换 (go-zero rpc)
├── proto/              # 共享 proto 定义
├── deploy/             # Docker Compose + Dockerfiles
├── pkg/                # 共享工具包
├── web/                # 前端 (React)
├── go.work             # Go workspace
└── Makefile
```
