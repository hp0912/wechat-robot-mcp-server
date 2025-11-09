package main

import (
	"context"
	"encoding/json"
	"log"

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
				rc := parseRobotContext(ctr.Params.Meta)
				ctx = WithRobotContext(ctx, rc)
				if rc.RobotCode != "" {
					db, err := GetDBByRobotCode(rc.RobotCode)
					if err != nil {
						log.Printf("获取数据库连接失败(RobotCode:%s): %v", rc.RobotCode, err)
					} else {
						ctx = WithDB(ctx, db)
					}
				}
			}
		}
		return next(ctx, method, req)
	}
}

func parseRobotContext(meta map[string]any) RobotContext {
	rc := RobotContext{}

	data, err := json.Marshal(meta)
	if err != nil {
		log.Printf("序列化 meta 失败: %v", err)
		return rc
	}

	if err := json.Unmarshal(data, &rc); err != nil {
		log.Printf("反序列化到 RobotContext 失败: %v", err)
		return rc
	}

	return rc
}
