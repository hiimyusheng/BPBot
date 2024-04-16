package main

import (
	"bpbot/handler"
	mongodb "bpbot/mongo"
	"bpbot/utililty"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/spf13/viper"
)

type Token struct {
	Secret string `mapstructure:"channel_secret"`
	Token  string `mapstructure:"channel_token"`
}

func main() {
	router := gin.Default()

	router.POST("/", HandleWebhookEvent)
	router.Run(":80")
}

func HandleWebhookEvent(c *gin.Context) {
	conf := readTokenConfig()

	bot, err := linebot.New(conf.Secret, conf.Token)
	if err != nil {
		utililty.Logger(3, err.Error())
		log.Fatal(err)
	}

	db, DBerr := mongodb.ConnectDB()
	if DBerr != nil {
		utililty.Logger(3, DBerr.Error())
		log.Fatal(DBerr)
	}
	handler.ReceiveWebhookEvent(c, bot, db)

	c.JSON(200, gin.H{})
}

func readTokenConfig() *Token {
	var Token = new(Token)
	viper.SetConfigFile("./config/token.json")
	if err := viper.ReadInConfig(); err != nil {
		utililty.Logger(3, err.Error())
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	if err := viper.Unmarshal(Token); err != nil {
		utililty.Logger(3, err.Error())
		panic(fmt.Errorf("unmarshal conf fail: %s", err))
	}
	return Token
}
