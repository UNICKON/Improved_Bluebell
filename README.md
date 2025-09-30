# Improved Bluebell


## 项目简介
这是一个基于 **Bluebell 社区系统** 的改进版本，针对原项目在 **性能、消息处理、排行榜** 等方面的不足进行了优化，适合学习与实践高并发场景下的社区系统设计。

Improved Bluebell 是一个基于 Go + Gin 的高性能社区/社交平台后端，支持分布式消息推送、关注关系、帖子发布、投票、热榜快照等功能，采用 Redis、MySQL、Kafka 等主流中间件，具备高可用和高扩展性。

---

## 主要功能模块

- 用户注册、登录、鉴权（JWT）
- 关注/取消关注（写入 MySQL，缓存 Redis，异步通知）
- 帖子发布、详情、列表
- 投票与热榜快照（Redis+MySQL一致性，分布式锁）
- Feed流推送（支持写扩散/读扩散混合策略，Kafka异步消息）
- 本地消息表保障消息可靠性，自动重试
- WebSocket/App推送接口预留

---

## 目录结构

```
controller/      # 路由控制器，接口层
logic/           # 业务逻辑层
models/          # 数据结构与建表SQL
dao/mysql/       # MySQL数据访问
dao/redis/       # Redis数据访问
dao/kafka/       # Kafka连接与操作
middlewares/     # Gin中间件（鉴权、限流等）
pkg/             # 雪花ID、JWT等工具包
router/          # 路由注册
settings/        # 配置加载
logger/          # 日志系统
```

---

## 运行环境

- Go 1.18+
- MySQL 5.7+
- Redis 6+
- Kafka 2.8+

---

## 快速启动

1. 初始化数据库表（见 models/*.sql）
2. 配置 config.yaml
3. 启动服务：
   ```bash
   go run main.go
   ```

---

## 典型接口说明

- 关注用户：
  ```bash
  POST /api/v1/follow
  { "follow_id": 2 }
  ```
- 取消关注：
  ```bash
  POST /api/v1/unfollow
  { "follow_id": 2 }
  ```
- 发布帖子（自动推送feed）：
  ```bash
  POST /api/v1/post
  { "title": "xxx", "content": "yyy", "community_id": 1 }
  ```
- 拉取feed：
  ```bash
  GET /api/v1/feed/get
  ```

---

## 其他说明

- 消息推送采用本地消息表+Kafka，保证高可靠性。
- 热榜快照、投票等均有分布式锁和一致性保障。
- 代码风格统一，接口均用JSON参数和自定义响应格式。

---

如需二次开发或部署，请参考各模块源码和注释。
