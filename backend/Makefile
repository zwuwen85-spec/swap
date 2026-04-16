.PHONY: all dev build test clean run

# 变量定义
APP_NAME=campus-swap-shop
BUILD_DIR=bin
CMD_DIR=cmd/server

# Go参数
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# 所有任务
all: dev

# 开发
dev:
	@echo "Starting development server..."
	$(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/main.go
	./$(BUILD_DIR)/$(APP_NAME)

# 运行
run:
	@echo "Running..."
	$(GOCMD) run $(CMD_DIR)/main.go

# 构建
build:
	@echo "Building..."
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

# 测试
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# 测试覆盖率
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# 清理
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# 下载依赖
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# 格式化代码
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# 代码检查
vet:
	@echo "Vetting code..."
	$(GOCMD) vet ./...

# 安装开发工具
install-tools:
	@echo "Installing development tools..."
	$(GOGET) -u github.com/gin-gonic/gin
	$(GOGET) -u gorm.io/gorm
	$(GOGET) -u gorm.io/driver/mysql
	$(GOGET) -u github.com/redis/go-redis/v9

# Docker构建
docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):latest .

# Docker运行
docker-run:
	@echo "Running Docker container..."
	docker-compose up -d
