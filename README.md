# wechat-robot-mcp-server

### 启动服务

```
go run .
```

### 打包编译

```
go build .
```

向微信客户端发送消息示例:

Post(fmt.Sprintf("http://client_%s:%s/api/v1/robot/message/send/longtext", rc.RobotCode, rc.WeChatClientPort))

注意：前面的`client_`前缀