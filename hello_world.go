package main

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type HelloWorldInput struct {
	Name string `json:"name" jsonschema:"Name to say hello to"`
}

func HelloWorld(ctx context.Context, req *mcp.CallToolRequest, params *HelloWorldInput) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		IsError: false,
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Hello, %s!", params.Name),
			},
		},
	}, nil, nil
}
