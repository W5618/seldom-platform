# Seldom Platform Backend (Go版本)

这是Seldom测试平台的Go语言后端实现，提供了完整的API服务来支持自动化测试管理。

## 功能特性

- 🚀 **项目管理** - 创建和管理测试项目
- 📝 **用例管理** - 测试用例的增删改查和复制
- 🌍 **环境管理** - 测试环境配置管理
- ⚡ **任务管理** - 测试任务的创建、执行和报告
- 👥 **团队管理** - 团队协作功能
- 🔐 **用户认证** - JWT认证和用户管理
- 📊 **API文档** - Swagger自动生成的API文档

## 技术栈

- **语言**: Go 1.21+
- **框架**: Gin Web Framework
- **数据库**: SQLite (支持扩展到其他数据库)
- **ORM**: GORM
- **认证**: JWT
- **文档**: Swagger
- **容器化**: Docker

## 项目结构

```
backendnew/
├── config/          # 配置管理
├── database/        # 数据库连接和迁移
├── handlers/        # HTTP处理器
├── middleware/      # 中间件
├── models/          # 数据模型
├── routes/          # 路由配置
├── services/        # 业务逻辑服务
├── utils/           # 工具函数
├── main.go          # 应用入口
├── go.mod           # Go模块文件
└── go.sum           # 依赖锁定文件
```

## 快速开始

### 本地开发

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd seldom-platform/backendnew
   ```

2. **安装依赖**
   ```bash
   go mod download
   ```

3. **构建项目**
   ```bash
   go build -o seldom-platform.exe
   ```

4. **运行服务**
   ```bash
   ./seldom-platform.exe
   ```

服务将在 `http://localhost:8080` 启动

### Docker部署

1. **使用Docker Compose**
   ```bash
   docker-compose up -d
   ```

2. **单独构建Docker镜像**
   ```bash
   docker build -t seldom-platform .
   docker run -p 8080:8080 seldom-platform
   ```

## API文档

启动服务后，访问 `http://localhost:8080/swagger/index.html` 查看完整的API文档。

### 主要API端点

- **健康检查**: `GET /health`
- **用户认证**: `POST /api/auth/login`
- **项目管理**: `GET|POST|PUT|DELETE /api/projects`
- **用例管理**: `GET|POST|PUT|DELETE /api/cases`
- **环境管理**: `GET|POST|PUT|DELETE /api/envs`
- **任务管理**: `GET|POST|PUT|DELETE /api/tasks`
- **团队管理**: `GET|POST|PUT|DELETE /api/teams`

## 配置说明

应用支持通过环境变量进行配置：

- `GIN_MODE`: Gin运行模式 (debug/release)
- `DB_TYPE`: 数据库类型 (sqlite)
- `DB_PATH`: 数据库文件路径
- `JWT_SECRET`: JWT密钥
- `SERVER_PORT`: 服务端口 (默认8080)

## 数据库

项目使用SQLite作为默认数据库，支持自动迁移。首次启动时会自动创建所需的表结构。

### 主要数据表

- `app_user_user` - 用户表
- `app_project_project` - 项目表
- `app_case_testcase` - 测试用例表
- `app_env_env` - 环境表
- `app_task_testtask` - 任务表
- `app_team_team` - 团队表

## 开发指南

### 添加新的API端点

1. 在 `models/` 中定义数据模型
2. 在 `handlers/` 中实现处理器
3. 在 `routes/` 中注册路由
4. 添加必要的中间件

### 代码规范

- 使用Go标准格式化工具: `go fmt`
- 遵循Go命名约定
- 添加适当的注释和文档
- 编写单元测试

## 部署

### 生产环境部署

1. 设置环境变量
2. 使用Docker Compose部署
3. 配置反向代理 (Nginx)
4. 设置SSL证书

### 监控和日志

- 应用日志存储在 `logs/` 目录
- 支持结构化日志记录
- 提供健康检查端点

## 贡献

欢迎提交Issue和Pull Request来改进项目。

## 许可证

本项目采用Apache 2.0许可证。

## 联系方式

如有问题或建议，请通过以下方式联系：

- 提交Issue
- 发送邮件
- 加入讨论群

---

**注意**: 这是Seldom Platform的Go语言重构版本，与原Python版本功能兼容。