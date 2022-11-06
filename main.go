package main

import (
	"fmt"
	"line_bot/http_response"
	"line_bot/model"
	mongodb "line_bot/mongo"
	"log"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

type Token struct {
	Secret string `mapstructure:"channel_secret"`
	Token  string `mapstructure:"channel_token"`
}
type User struct {
	Id string `mapstructure:"user_id"`
}

var bot *linebot.Client

func main() {

	conf := readTokenConfig()
	user := readUserConfig()
	cmd := exec.Command("bash", "./init.sh")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to start mongodb: %v", err)
	}
	client, DBerr := mongodb.ConnectDB()
	if DBerr != nil {
		log.Fatal(DBerr)
	}
	bot, err := linebot.New(conf.Secret, conf.Token)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("bot:", bot)
	router := gin.Default()
	router.POST("/callback", func(c *gin.Context) {
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
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
						log.Print(err)
					}
					newMessage.Id = event.Source.UserID
					newMessage.Message = message.Text
					mongodb.RecieveMessage(newMessage, client)
				}
			}
		}
	})
	router.POST("/api/pushmessage", pushMessageHandler(bot, user.Id))
	router.GET("/api/querymessage", queryMessageHandler(client, user.Id))
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

func readUserConfig() *User {
	var User = new(User)
	viper.SetConfigFile("./config/user.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	if err := viper.Unmarshal(User); err != nil {
		panic(fmt.Errorf("unmarshal conf fail: %s \n", err))
	}
	return User
}

func pushMessageHandler(bot *linebot.Client, user string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var pushMessage struct {
			Type string `form:"type", json:"type"`
			Text string `form:"text", json:"text"`
		}
		if err := c.BindJSON(&pushMessage); err != nil {
			c.JSON(http.StatusBadRequest, http_response.NewErrorResp(1, "Invalid parameter format or missing necessary parameter."))
			// return
		}
		switch pushMessage.Type {
		case "text":
			if _, err := bot.PushMessage(user, linebot.NewTextMessage(pushMessage.Text)).Do(); err != nil {
				c.JSON(http.StatusBadRequest, http_response.NewErrorResp(1, "Push Failed"))
			}
		}
	}
	return gin.HandlerFunc(fn)
}

func queryMessageHandler(client mongo.Client, user string) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var message []model.Message
		message = mongodb.QueryMessage(user, client)
		c.JSON(http.StatusOK, message)
	}
	return gin.HandlerFunc(fn)
}
