package main

import (
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "wechat-robot-mcp-server",
		Version: "1.0.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "helloWorld",
		Description: "Say hello to the world",
	}, HelloWorld)

	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return server
	}, nil)

	if err := http.ListenAndServe("", handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
