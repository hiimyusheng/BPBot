package model

import "github.com/line/line-bot-sdk-go/v7/linebot"

type Event struct {
	EventType   linebot.EventType
	EventSource struct {
		Type    string
		GroupId string
		UserId  string
	}
	WebhookEvent string
}
