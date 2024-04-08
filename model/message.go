package model

import "time"

type Message struct {
	Id          string
	SourceType  string
	UserId      string
	Message     string
	MessageType string
	ReplyToken  string
	Time        time.Time
}
