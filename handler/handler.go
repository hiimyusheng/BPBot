package handler

import (
	"line_bot/model"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Line struct {
}

type Gcp struct {
}

func ReceiveWebhookEvent(c *gin.Context, bot *linebot.Client) {

	userAgent := c.Request.Header.Get("User-Agent")
	switch userAgent {
	case "LineBotWebhook/2.0":
		var handler Line
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				c.Writer.WriteHeader(500)
			} else {
				c.Writer.WriteHeader(500)
			}
		}
		handler.HandleEvent(events, bot)
		// wip 回傳方式
		// wip 處理line訊息的方式

	case "Google-Alerts":
		var newGoogleAlert model.Gcp
		var handler Gcp
		if err := c.BindJSON(&newGoogleAlert); err != nil {
			log.Fatal(err)
		}
		handler.HandleEvent(newGoogleAlert, bot)

	}

}

type Handler interface {
	HandleEvent(event interface{}, bot *linebot.Client)
}
