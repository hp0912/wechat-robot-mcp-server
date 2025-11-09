package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var (
	MCPServerPort int
)

func LoadConfig() error {
	loadEnvConfig()
	return nil
}

func loadEnvConfig() {
	// 本地开发模式
	isDevMode := strings.ToLower(os.Getenv("GO_ENV")) == "dev"
	if isDevMode {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("加载本地环境变量失败，请检查是否存在 .env 文件")
		}
	}

	MCPServerPort, err := strconv.Atoi(os.Getenv("MCP_SERVER_PORT"))
	if err != nil {
		log.Fatal("加载本地环境变量失败，请检查是否存在 .env 文件")
	}
	if MCPServerPort == 0 {
		MCPServerPort = 9000
	}
	if MCPServerPort < 1 || MCPServerPort > 65535 {
		log.Fatal("MCPServerPort 必须在 1 到 65535 之间")
	}
}
