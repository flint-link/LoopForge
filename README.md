# LoopForge - Tenant Management Service

基于 Go + Gin + PostgreSQL 的租户管理后端服务，提供完整的 CRUD API 接口。

## 功能特性

- ✅ 租户创建、查询、更新、删除 (CRUD)
- ✅ 分页查询，支持按名称、编码、状态筛选
- ✅ 租户状态切换 (启用/禁用)
- ✅ 统一响应格式和错误码
- ✅ 参数校验和 SQL 注入防护
- ✅ 结构化日志和请求追踪
- ✅ 配置文件和环境变量支持
- ✅ 优雅启动和关闭

## 技术栈

| 组件 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | 1.21+ |
| Web框架 | Gin | v1.9+ |
| 数据库驱动 | pgx | v5+ |
| 配置管理 | Viper | v1.18+ |
| 日志库 | Zap | v1.26+ |
| 参数校验 | go-playground/validator | v10+ |

## 项目结构

```
LoopForge/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口
├── internal/
│   ├── config/
│   │   └── config.go            # 配置加载
│   ├── database/
│   │   └── postgres.go          # 数据库连接
│   ├── models/
│   │   └── tenant.go            # 数据模型
│   ├── repository/
│   │   └── tenant_repository.go # 数据访问层
│   ├── service/
│   │   └── tenant_service.go    # 业务逻辑层
│   ├── handler/
│   │   └── tenant_handler.go    # HTTP处理器
│   ├── middleware/
│   │   ├── logger.go            # 日志中间件
│   │   └── recovery.go          # 异常恢复
│   └── response/
│       └── response.go          # 统一响应
├── pkg/
│   └── logger/
│       └── logger.go            # 日志工具
├── configs/
│   └── config.yaml.example      # 配置模板
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 快速开始

### 1. 环境要求

- Go 1.21+
- PostgreSQL 12+ (运行在 127.0.0.1:5432)
- 数据库 `ruoyi-vue-pro` 已创建
- 表 `system_tenant` 已存在
- 账号 `postgres` / 密码 `Flint@123456` 有读写权限

### 2. 克隆项目

```bash
git clone https://github.com/flint-link/LoopForge.git
cd LoopForge
```

### 3. 配置文件

复制配置模板并修改：

```bash
cp configs/config.yaml.example configs/config.yaml
# 编辑 configs/config.yaml，根据需要修改数据库连接等配置
```

### 4. 安装依赖

```bash
make tidy
```

### 5. 构建和运行

```bash
# 构建
make build

# 运行
make run

# 或者直接运行
./bin/server -config configs/config.yaml
```

服务启动后监听 `http://0.0.0.0:8080`

### 6. 环境变量覆盖

支持通过环境变量覆盖配置（前缀 `LOOPFORGE_`）：

```bash
export LOOPFORGE_DATABASE_HOST=127.0.0.1
export LOOPFORGE_DATABASE_PASSWORD=Flint@123456
export LOOPFORGE_SERVER_PORT=8080
./bin/server
```

## API 文档

### 统一响应格式

```json
{
    "code": 0,
    "message": "success",
    "data": {}
}
```

### 错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 参数错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 409 | 资源冲突 (如编码重复) |
| 500 | 服务器内部错误 |

### 租户管理接口

#### 创建租户
```http
POST /api/v1/tenants
Content-Type: application/json

{
    "name": "测试租户",
    "code": "TEST001",
    "domain": "test.example.com",
    "package_id": 1,
    "remark": "测试用租户"
}
```

#### 查询租户列表
```http
GET /api/v1/tenants?page=1&page_size=10&name=测试&status=0
```

#### 查询租户详情
```http
GET /api/v1/tenants/1
```

#### 更新租户
```http
PUT /api/v1/tenants/1
Content-Type: application/json

{
    "name": "更新后的租户名称",
    "domain": "new.example.com",
    "remark": "更新备注"
}
```

#### 删除租户
```http
DELETE /api/v1/tenants/1
```

#### 切换租户状态
```http
PATCH /api/v1/tenants/1/status
Content-Type: application/json

{
    "status": 1,
    "update_by": "admin"
}
```

### 健康检查
```http
GET /health
```

## 开发指南

### 代码规范

```bash
# 格式化代码
make fmt

# 静态检查
make vet

# 代码检查 (需要安装 golangci-lint)
make lint
```

### 运行测试

```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-cover
```

### 热重载开发

安装 air：
```bash
go install github.com/cosmtrek/air@latest
```

启动开发模式：
```bash
make dev
```

### 生成 Swagger 文档

安装 swag：
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

生成文档：
```bash
make swag
```

访问 `http://localhost:8080/swagger/index.html` 查看 API 文档。

## 配置说明

### Server 配置
| 字段 | 说明 | 默认值 |
|------|------|--------|
| host | 监听地址 | 0.0.0.0 |
| port | 监听端口 | 8080 |
| mode | 运行模式 | debug |

### Database 配置
| 字段 | 说明 | 默认值 |
|------|------|--------|
| host | 数据库地址 | 127.0.0.1 |
| port | 数据库端口 | 5432 |
| user | 用户名 | postgres |
| password | 密码 | - |
| dbname | 数据库名 | ruoyi-vue-pro |
| max_open_conns | 最大打开连接数 | 20 |
| max_idle_conns | 最大空闲连接数 | 10 |
| conn_max_lifetime | 连接最大存活时间 | 1h |

### Log 配置
| 字段 | 说明 | 默认值 |
|------|------|--------|
| level | 日志级别 | info |
| format | 输出格式 | json |
| output | 输出目标 | stdout |
| file_path | 文件路径 | ./logs/app.log |

## 部署

### 使用 systemd

创建 `/etc/systemd/system/loopforge.service`：

```ini
[Unit]
Description=LoopForge Tenant Service
After=network.target postgresql.service

[Service]
Type=simple
User=loopforge
WorkingDirectory=/opt/loopforge
ExecStart=/opt/loopforge/bin/server -config /opt/loopforge/configs/config.yaml
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启用并启动：
```bash
sudo systemctl daemon-reload
sudo systemctl enable loopforge
sudo systemctl start loopforge
```

### 使用 Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/server ./cmd/server

FROM alpine:3.18
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/bin/server .
COPY --from=builder /app/configs/config.yaml.example ./configs/config.yaml
EXPOSE 8080
CMD ["./server", "-config", "configs/config.yaml"]
```

构建镜像：
```bash
docker build -t loopforge:latest .
```

运行容器：
```bash
docker run -d \
  --name loopforge \
  -p 8080:8080 \
  -e LOOPFORGE_DATABASE_HOST=host.docker.internal \
  -e LOOPFORGE_DATABASE_PASSWORD=Flint@123456 \
  loopforge:latest
```

## 许可证

MIT License