package main

import (
	"fmt"
	"line_bot/http_response"
	"line_bot/model"
	mongodb "line_bot/mongo"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

type Token struct {
	Secret string `mapstructure:"channel_secret"`
	Token  string `mapstructure:"channel_token"`
}

var bot *linebot.Client

func main() {

	conf := readTokenConfig()
	client, DBerr := mongodb.ConnectDB()
	if DBerr != nil {
		log.Fatal(DBerr)
	}
	bot, err := linebot.New(conf.Secret, conf.Token)
	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()

	router.POST("/", receiveMessageHandler(bot, client))
	router.POST("/api/pushMessage", pushMessageHandler(bot))
	router.GET("/api/queryMessage/:user_id", queryMessageHandler(client))
	router.GET("api/getUserInfo/:user_id", getUserInfoHandler(bot))
	router.GET("api/getJoinedGroup", getAllJoinedGroupSummary(client))
	router.POST("api/setNotifyGroup", setNotifyGroupHandler(client))

	router.Run(":80")

}

func readTokenConfig() *Token {
	var Token = new(Token)
	viper.SetConfigFile("./config/token.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	if err := viper.Unmarshal(Token); err != nil {
		panic(fmt.Errorf("unmarshal conf fail: %s \n", err))
	}
	return Token
}

func receiveMessageHandler(bot *linebot.Client, client mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				c.Writer.WriteHeader(400)
			} else {
				c.Writer.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				var newMessage model.Message
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					newMessage.MessageType = "text"
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
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
							mongodb.InsertGroup(group, client)
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
					mongodb.RecieveMessage(newMessage, client)
					// if message.Text
				}
			}
		}
	}
	return gin.HandlerFunc(fn)
}

func pushMessageHandler(bot *linebot.Client) gin.HandlerFunc {
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

func queryMessageHandler(client mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("user_id")
		var message []model.Message
		message = mongodb.QueryMessage(id, client)
		c.JSON(http.StatusOK, message)
	}
	return gin.HandlerFunc(fn)
}

func getUserInfoHandler(bot *linebot.Client) gin.HandlerFunc {
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

func getAllJoinedGroupSummary(client mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		groups := mongodb.GetAllJoinedGroupSummary(client)
		c.JSON(http.StatusOK, groups)
	}
	return gin.HandlerFunc(fn)
}

func setNotifyGroupHandler(client mongo.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var NotifyGroup struct {
			RequestId string `form:"request_id", json:"request_id"`
			GroupName string `form:"group_name", json:"group_name"`
		}
		if err := c.BindJSON(&NotifyGroup); err != nil {
			c.JSON(http.StatusBadRequest, http_response.NewErrorResp(1, "Invalid parameter format or missing necessary parameter."))
			return
		}
		if _, err := bot.PushMessage(pushMessage.User, linebot.NewTextMessage(pushMessage.Text)).Do(); err != nil {
			c.JSON(http.StatusBadRequest, http_response.NewErrorResp(1, err.Error()))
		}
	}

	return gin.HandlerFunc(fn)
}
