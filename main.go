package main

import (
	"fmt"
	"line_bot/model"
	"line_bot/mongo"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/viper"
)

type Token struct {
	Secret string `mapstructure:"channel_secret"`
	Token  string `mapstructure:"channel_token"`
}

var bot *linebot.Client

func main() {

	conf := readConfig()
	client, DBerr := mongo.ConnectDB()
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
					mongo.RecieveMessage(newMessage, client)
				}
			}
		}
	})
	router.Run(":80")

}

func readConfig() *Token {
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

func getMessage(c *gin.Context) {
	fmt.Println("testttt")
	c.JSON(http.StatusOK, gin.H{})
}
