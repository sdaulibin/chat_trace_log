# 聊天日志API服务

这是一个基于Golang的Web服务，提供符合OpenAPI 3.1.0规范的接口，用于记录聊天信息。该服务将输入的信息记录在本地文件中，每条记录包含文本、验证结果和输入时间三个字段。

## 功能特点

- 符合OpenAPI 3.1.0规范的RESTful API
- 提供Swagger UI界面进行API文档浏览和测试
- 按日期组织存储聊天日志
- 支持日期范围查询
- 并发安全的文件操作

## 安装与运行

### 前置条件

- Go 1.20或更高版本
- Git

### 安装步骤

1. 克隆仓库
```bash
git clone https://github.com/binginx/bqd_chat_log.git
cd bqd_chat_log
```

2. 安装依赖
```bash
go mod download
```

3. 生成Swagger文档
```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

4. 构建并运行
```bash
go build -o chat_log_server
./chat_log_server
```

服务将在 http://localhost:8080 上启动

## API文档

启动服务后，可以通过访问 http://localhost:8080/swagger/index.html 查看完整的API文档。

### 主要API端点

#### 创建日志

- **URL**: `/api/v1/logs`
- **方法**: `POST`
- **请求体**:
  ```json
  {
    "text": "这是一条聊天记录",
    "validation_result": true
  }
  ```
- **成功响应** (200 OK):
  ```json
  {
    "code": 200,
    "message": "日志保存成功",
    "data": {
      "id": "1672567200000000000",
      "text": "这是一条聊天记录",
      "validation_result": true,
      "input_time": "2023-01-01 12:00:00.000"
    }
  }
  ```

#### 获取日志

- **URL**: `/api/v1/logs`
- **方法**: `GET`
- **查询参数**:
  - `start_date`: 开始日期 (YYYY-MM-DD) (可选，默认为结束日期前7天)
  - `end_date`: 结束日期 (YYYY-MM-DD) (可选，默认为当前日期)
- **成功响应** (200 OK):
  ```json
  {
    "code": 200,
    "message": "获取日志成功",
    "data": [
      {
        "id": "1672567200000000000",
        "text": "这是一条聊天记录",
        "validation_result": true,
        "input_time": "2023-01-01T12:00:00Z"
      }
    ]
  }
  ```

## 数据存储

日志数据按日期存储在 `./storage` 目录下，每个日期对应一个JSON文件。

## 许可证

本项目采用MIT许可证。详情请参阅LICENSE文件。