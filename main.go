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

var bot *linebot.Client

func main() {

	router := gin.Default()

	router.POST("/", HandleWebhookEvent)
	// router.POST("/api/pushMessage", handler.PushMessageHandler(bot))
	// router.GET("api/getUserInfo/:user_id", handler.GetUserInfoHandler(bot))
	// router.GET("/api/queryMessage/:user_id", handler.QueryMessageHandler(db))
	// router.GET("api/getJoinedGroup", handler.GetAllJoinedGroupSummary(db))
	// router.POST("api/setNotifyGroup", setNotifyGroupHandler(client))

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
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	if err := viper.Unmarshal(Token); err != nil {
		panic(fmt.Errorf("unmarshal conf fail: %s \n", err))
	}
	return Token
}

// func setNotifyGroupHandler(client mongo.Client) gin.HandlerFunc {
// 	fn := func(c *gin.Context) {
// 		var NotifyGroup struct {
// 			RequestId string `form:"request_id", json:"request_id"`
// 			GroupName string `form:"group_name", json:"group_name"`
// 		}
// 		if err := c.BindJSON(&NotifyGroup); err != nil {
// 			c.JSON(http.StatusBadRequest, http_response.NewErrorResp(1, "Invalid parameter format or missing necessary parameter."))
// 			return
// 		}
// 		if _, err := bot.PushMessage(pushMessage.User, linebot.NewTextMessage(pushMessage.Text)).Do(); err != nil {
// 			c.JSON(http.StatusBadRequest, http_response.NewErrorResp(1, err.Error()))
// 		}
// 	}

// 	return gin.HandlerFunc(fn)
// }
