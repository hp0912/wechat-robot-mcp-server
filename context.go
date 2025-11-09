package main

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type RobotContext struct {
	RobotID      int64
	RobotCode    string
	RobotRedisDB uint
	RobotWxID    string
	FromWxID     string
	SenderWxID   string
	MessageID    int64
	RefMessageID int64
}

type ctxKey string

const (
	ctxKeyRobotContext ctxKey = "robot_context"
	ctxKeyDB           ctxKey = "robot_db"
)

func WithRobotContext(parent context.Context, rc RobotContext) context.Context {
	return context.WithValue(parent, ctxKeyRobotContext, rc)
}

func GetRobotContext(ctx context.Context) (RobotContext, bool) {
	val := ctx.Value(ctxKeyRobotContext)
	if val == nil {
		return RobotContext{}, false
	}
	rc, ok := val.(RobotContext)
	return rc, ok
}

func WithDB(parent context.Context, db *gorm.DB) context.Context {
	return context.WithValue(parent, ctxKeyDB, db)
}

func GetDB(ctx context.Context) (*gorm.DB, bool) {
	val := ctx.Value(ctxKeyDB)
	if val == nil {
		return nil, false
	}
	db, ok := val.(*gorm.DB)
	return db, ok
}

// GetSQLDB 尝试返回底层 *sql.DB，便于做健康检查等
func GetSQLDB(ctx context.Context) (*sql.DB, bool) {
	gdb, ok := GetDB(ctx)
	if !ok || gdb == nil {
		return nil, false
	}
	sqldb, err := gdb.DB()
	if err != nil {
		return nil, false
	}
	return sqldb, true
}
