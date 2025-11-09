package main

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
