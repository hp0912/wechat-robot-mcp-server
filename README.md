# wechat-robot-mcp-server

### 启动服务

```
go run .
```

### 打包编译

```
go build .
```

### 自定义 HTTP 接口

服务会在 `MCP_SERVER_PORT` 端口同时暴露自定义接口：

- 路径：`POST /api/v1/messages`
- 请求体(JSON)：
  ```json
  {
    "text": "hello"
  }
  ```
- 响应体(JSON)：
  ```json
  {
    "success": true,
    "echo": "hello"
  }
  ```

示例：

```bash
curl -sS -X POST "http://localhost:${MCP_SERVER_PORT}/api/v1/messages" \
  -H "Content-Type: application/json" \
  -d '{"text":"hello"}'
```