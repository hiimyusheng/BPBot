package handler

import (
	"line_bot/model"
	"line_bot/utililty"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.mongodb.org/mongo-driver/mongo"
)

type Line struct {
}

type Gcp struct {
}

func ReceiveWebhookEvent(c *gin.Context, bot *linebot.Client, db mongo.Client) {

	userAgent := c.Request.Header.Get("User-Agent")
	switch userAgent {
	case "LineBotWebhook/2.0":
		var handler Line
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				utililty.Logger(3, err.Error())
				c.Writer.WriteHeader(500)
			} else {
				utililty.Logger(3, err.Error())
				c.Writer.WriteHeader(500)
			}
		}
		handler.HandleEvent(events, bot, db)
		// wip 回傳方式
		// wip 處理line訊息的方式

	case "Google-Alerts":
		var newGoogleAlert model.Gcp
		var handler Gcp
		if err := c.BindJSON(&newGoogleAlert); err != nil {
			utililty.Logger(3, err.Error())
			log.Fatal(err)
		}
		handler.HandleEvent(newGoogleAlert, bot, db)

	}

}

type Handler interface {
	HandleEvent(event interface{}, bot *linebot.Client)
}
