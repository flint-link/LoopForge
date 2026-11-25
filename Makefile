# LoopForge Makefile

.PHONY: build run test clean tidy vet lint fmt help

# 默认目标
all: build

# 变量
BINARY_NAME=server
BUILD_DIR=bin
CMD_DIR=cmd/server
CONFIG_FILE=configs/config.yaml

# 构建二进制文件
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# 运行服务
run: build
	@echo "Starting server..."
	@$(BUILD_DIR)/$(BINARY_NAME) -config $(CONFIG_FILE)

# 运行测试
test:
	@echo "Running tests..."
	@go test -v ./...

# 运行测试并生成覆盖率报告
test-cover:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# 代码格式化
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# 静态检查
vet:
	@echo "Running go vet..."
	@go vet ./...

# 代码检查
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run ./... || true

# 整理依赖
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# 清理构建产物
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# 开发模式运行（热重载需要安装 air）
dev:
	@echo "Starting dev server with air..."
	@air -c .air.toml || echo "Please install air: go install github.com/cosmtrek/air@latest"

# 生成 Swagger 文档（需要安装 swag）
swag:
	@echo "Generating Swagger docs..."
	@swag init -g $(CMD_DIR)/main.go -o docs || echo "Please install swag: go install github.com/swaggo/swag/cmd/swag@latest"

# 帮助信息
help:
	@echo "LoopForge - Tenant Management Service"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build       Build the binary"
	@echo "  run         Build and run the server"
	@echo "  test        Run all tests"
	@echo "  test-cover  Run tests with coverage report"
	@echo "  fmt         Format code"
	@echo "  vet         Run go vet"
	@echo "  lint        Run golangci-lint"
	@echo "  tidy        Tidy dependencies"
	@echo "  clean       Clean build artifacts"
	@echo "  dev         Run with hot reload (requires air)"
	@echo "  swag        Generate Swagger docs (requires swag)"
	@echo "  help        Show this help"