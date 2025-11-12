package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MysqlSettingS struct {
	Host     string
	Port     string
	User     string
	Password string
}

var (
	MCPServerPort int
	MysqlSettings = &MysqlSettingS{}
)

// TenantDBManager 负责基于 RobotCode 缓存和创建不同的 *gorm.DB 连接
type TenantDBManager struct {
	mu      sync.RWMutex
	tenants map[string]*gorm.DB
}

var tenantDB = &TenantDBManager{
	tenants: make(map[string]*gorm.DB),
}

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

	port, err := strconv.Atoi(os.Getenv("MCP_SERVER_PORT"))
	if err != nil {
		log.Fatal("环境变量 [MCP_SERVER_PORT] 未配置")
	}
	if port == 0 {
		port = 9000
	}
	if port < 1 || port > 65535 {
		log.Fatal("MCPServerPort 必须在 1 到 65535 之间")
	}
	MCPServerPort = port

	MysqlSettings.Host = os.Getenv("MYSQL_HOST")
	MysqlSettings.Port = os.Getenv("MYSQL_PORT")
	MysqlSettings.User = os.Getenv("MYSQL_USER")
	MysqlSettings.Password = os.Getenv("MYSQL_PASSWORD")
}

// GetDBByRobotCode 获取指定 RobotCode 对应的 *gorm.DB（带缓存）
func GetDBByRobotCode(robotCode string) (*gorm.DB, error) {
	if robotCode == "" {
		return nil, fmt.Errorf("robotCode 为空")
	}
	// 读缓存
	tenantDB.mu.RLock()
	db, ok := tenantDB.tenants[robotCode]
	tenantDB.mu.RUnlock()
	if ok && db != nil {
		return db, nil
	}

	// 双重检查加锁创建
	tenantDB.mu.Lock()
	defer tenantDB.mu.Unlock()
	if db, ok := tenantDB.tenants[robotCode]; ok && db != nil {
		return db, nil
	}

	dsn := buildDSNForRobot(robotCode)
	newDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败(%s): %w", robotCode, err)
	}
	// 基础连接池设置
	sqlDB, err := newDB.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层连接失败(%s): %w", robotCode, err)
	}
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("数据库不可用(%s): %w", robotCode, err)
	}

	tenantDB.tenants[robotCode] = newDB
	return newDB, nil
}

func buildDSNForRobot(robotCode string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		MysqlSettings.User,
		MysqlSettings.Password,
		MysqlSettings.Host,
		MysqlSettings.Port,
		robotCode,
	)
}
