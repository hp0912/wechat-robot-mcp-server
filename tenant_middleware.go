package main

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TenantMiddleware(next mcp.MethodHandler) mcp.MethodHandler {
	return func(
		ctx context.Context,
		method string,
		req mcp.Request,
	) (mcp.Result, error) {
		if ctr, ok := req.(*mcp.CallToolRequest); ok {
			if ctr.Params.Meta != nil {
				// 1. 将 Meta map 结构转换成 RobotContext
				// 2. 将 RobotContext 添加到上下文
				// 3. 根据 RobotContext 获取对应的 *gorm.DB
				// 4. 将 *gorm.DB 添加到上下文
			}
		}
		return next(ctx, method, req)
	}
}
