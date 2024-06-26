package handler

import (
	"fmt"
	"nmbot/model"
	mongodb "nmbot/mongo"
	"nmbot/utililty"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.mongodb.org/mongo-driver/mongo"
)

func (g Gcp) HandleEvent(googleAlert model.Gcp, bot *linebot.Client, db mongo.Client) {

	mongodb.InsertAlert(googleAlert, db)

	groups, err := mongodb.GetRegisteredGroup(googleAlert, db)
	if err != nil {
		utililty.Logger(3, err.Error())
		return
	}
	tm := time.Unix(googleAlert.Incident.Started, 0)
	loc, _ := time.LoadLocation("Asia/Taipei")
	triggeredTime := tm.In(loc).Format("2006-01-02 15:04:05") + " (GMT+8)"
	message := fmt.Sprintf("⚠️ *%s* 快訊觸發！ ⚠️\r\n狀態： *%s*\r\n時間： *%s*\r\n內容： *%s*\r\n請檢查系統運行狀況！", googleAlert.Incident.PolicyName, googleAlert.Incident.State, triggeredTime, googleAlert.Incident.Summary)
	for _, group := range groups {
		if _, err := bot.PushMessage(group.GroupId, linebot.NewTextMessage(message)).Do(); err != nil {
			utililty.Logger(3, err.Error())
			fmt.Println(err)
		}
	}
	fmt.Println(googleAlert)
}
