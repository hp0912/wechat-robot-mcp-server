package main

import (
	"regexp"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// NormalizeAIBaseURL 规范化AI BaseURL，确保以/v+数字结尾，如果没有则添加/v1
func NormalizeAIBaseURL(baseURL string) string {
	baseURL = strings.TrimRight(baseURL, "/")
	versionRegex := regexp.MustCompile(`/v\d+$`)
	if !versionRegex.MatchString(baseURL) {
		baseURL += "/v1"
	}
	return baseURL
}

func CallToolResultError(errmsg string) (*mcp.CallToolResult, any, error) {
	result := &mcp.CallToolResult{}
	result.IsError = true
	result.Content = []mcp.Content{
		&mcp.TextContent{
			Text: errmsg,
		},
	}
	return result, nil, nil
}
