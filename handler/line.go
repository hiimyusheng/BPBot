package handler

import (
	"bpbot/handler/command"
	"bpbot/http_response"
	"bpbot/model"
	mongodb "bpbot/mongo"
	"bpbot/utililty"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.mongodb.org/mongo-driver/mongo"
)

func (l Line) HandleEvent(events []*linebot.Event, bot *linebot.Client, db mongo.Client) {

	for _, event := range events {
		mongodb.InsertEvent(*event, db)
		switch event.Type {
		case linebot.EventTypeMessage:
			var newMessage model.Message
			var replyMessage string
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if b := isCommand(message.Text); b && event.Source.Type == "group" {
					if isAuthed := checkAuth(event, db); !isAuthed {
						replyMessage = "您沒有權限使用此指令！"
					} else {
						replyMessage = handleCommand(event, bot, db)
					}
					if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						utililty.Logger(3, err.Error())
						log.Print(err)
					}
				}
				newMessage.Message = message.Text
				if event.Source.Type == "group" {
					newMessage.SourceType = "group"
					newMessage.Id = string(event.Source.GroupID)
					newMessage.UserId = string(event.Source.UserID)
				}
				if event.Source.Type == "room" {
					newMessage.SourceType = "room"
					newMessage.Id = string(event.Source.RoomID)
					newMessage.UserId = string(event.Source.UserID)
				}
				if event.Source.Type == "user" {
					newMessage.SourceType = "user"
					newMessage.Id = string(event.Source.UserID)
					newMessage.UserId = string(event.Source.UserID)
				}
				newMessage.ReplyToken = event.ReplyToken
				newMessage.Time = event.Timestamp
				mongodb.RecieveMessage(newMessage, db)
			}
		default:
		}
	}
}

func isCommand(text string) bool {
	if m, _ := regexp.MatchString(`^(\/[^ ]+)(?:\s+([^ ]+))?`, text); m {
		return true
	}
	return false
}

func handleCommand(event *linebot.Event, bot *linebot.Client, db mongo.Client) string {
	r := regexp.MustCompile(`^(\/[^ ]+)(?:\s+([^ ]+))?`)
	match := r.FindStringSubmatch(event.Message.(*linebot.TextMessage).Text)
	if match == nil {
		utililty.Logger(3, "Nil match")
		return ""
	}
	fmt.Println(match)
	switch match[1] {
	case "/add":
		if match[2] == "" {
			utililty.Logger(3, "Missing project id")
		}
		message := command.Add(match[2], event, bot, db)
		return message
	default:
		utililty.Logger(3, "Unrecognized command")
		return "無法辨認的指令"
	}
}

func checkAuth(event *linebot.Event, db mongo.Client) bool {
	result, err := mongodb.GetAuthedUsers(event.Source.UserID, db)
	if err != nil {
		utililty.Logger(3, err.Error())
		return false
	}
	fmt.Println("result is", result)
	return result != nil
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
		message, err := mongodb.QueryMessage(id, db)
		if err != nil {
			utililty.Logger(3, err.Error())
			return
		}
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
		groups, err := mongodb.GetAllJoinedGroupSummary(db)
		if err != nil {
			utililty.Logger(3, err.Error())
			return
		}
		c.JSON(http.StatusOK, groups)
	}
	return gin.HandlerFunc(fn)
}
