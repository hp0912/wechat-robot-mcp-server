package main

import (
	"context"
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

	if v, ok := meta["robotId"].(float64); ok {
		rc.RobotID = int64(v)
	}
	if v, ok := meta["robotCode"].(string); ok {
		rc.RobotCode = v
	}
	if v, ok := meta["robotRedisDb"].(float64); ok {
		rc.RobotRedisDB = uint(v)
	}
	if v, ok := meta["robotWxId"].(string); ok {
		rc.RobotWxID = v
	}
	if v, ok := meta["fromWxId"].(string); ok {
		rc.FromWxID = v
	}
	if v, ok := meta["senderWxId"].(string); ok {
		rc.SenderWxID = v
	}
	if v, ok := meta["messageId"].(float64); ok {
		rc.MessageID = int64(v)
	}
	if v, ok := meta["refMessageId"].(float64); ok {
		rc.RefMessageID = int64(v)
	}

	return rc
}
