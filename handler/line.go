package handler

import (
	"fmt"
	"line_bot/http_response"
	"line_bot/model"
	mongodb "line_bot/mongo"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.mongodb.org/mongo-driver/mongo"
)

func ReceiveMessageHandler(events []*linebot.Event, bot *linebot.Client, db mongo.Client) {
	for _, event := range events {
		mongodb.InsertEvent(*event, db)
		switch event.Type {
		case linebot.EventTypeMessage:
			var newMessage model.Message
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				newMessage.MessageType = "text"
				if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Print(err)
				}
				newMessage.UserId = event.Source.UserID
				newMessage.Message = message.Text
				if event.Source.Type == "group" {
					newMessage.SourceType = "group"
					newMessage.Id = string(event.Source.GroupID)
					if summary, err := bot.GetGroupSummary(event.Source.GroupID).Do(); err == nil {
						var group model.Group
						group.GroupId = summary.GroupID
						group.GroupName = summary.GroupName
						mongodb.InsertGroup(group, db)
					} else {
						log.Print(err)
						fmt.Println(err)
					}
				}
				if event.Source.Type == "room" {
					newMessage.SourceType = "room"
					newMessage.Id = string(event.Source.RoomID)
				}
				if event.Source.Type == "user" {
					newMessage.SourceType = "user"
					newMessage.Id = string(event.Source.UserID)
				}
				newMessage.ReplyToken = event.ReplyToken
				newMessage.Time = event.Timestamp
				mongodb.RecieveMessage(newMessage, db)
			}
		default:
			break
		}
	}
}

func PushMessageHandler(bot *linebot.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var pushMessage struct {
			User string `form:"user", json:"user"`
			Type string `form:"type", json:"type"`
			Text string `form:"text", json:"text"`
		}
		if err := c.BindJSON(&pushMessage); err != nil {
			c.JSON(http.StatusBadRequest, http_response.NewErrorResp(1, "Invalid parameter format or missing necessary parameter."))
			return
		}
		switch pushMessage.Type {
		case "text":
			if _, err := bot.PushMessage(pushMessage.User, linebot.NewTextMessage(pushMessage.Text)).Do(); err != nil {
				c.JSON(http.StatusBadRequest, http_response.NewErrorResp(1, err.Error()))
			}
		}
	}
	return gin.HandlerFunc(fn)
}

func QueryMessageHandler(db mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("user_id")
		var message []model.Message
		message = mongodb.QueryMessage(id, db)
		c.JSON(http.StatusOK, message)
	}
	return gin.HandlerFunc(fn)
}

func GetUserInfoHandler(bot *linebot.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("user_id")
		if info, err := bot.GetProfile(id).Do(); err != nil {
			c.JSON(http.StatusBadRequest, http_response.NewErrorResp(1, err.Error()))
		} else {
			c.JSON(http.StatusOK, info)
		}

	}
	return gin.HandlerFunc(fn)
}

func GetAllJoinedGroupSummary(db mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		groups := mongodb.GetAllJoinedGroupSummary(db)
		c.JSON(http.StatusOK, groups)
	}
	return gin.HandlerFunc(fn)
}
