package handler

import (
	"fmt"
	"line_bot/model"
	mongodb "line_bot/mongo"
	"line_bot/utililty"
	"log"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func (g Gcp) HandleEvent(googleAlert model.Gcp, bot *linebot.Client) {
	db, DBerr := mongodb.ConnectDB()
	if DBerr != nil {
		utililty.Logger(3, DBerr.Error())
		log.Fatal(DBerr)
	}
	mongodb.InsertAlert(googleAlert, db)

	groups := mongodb.GetAllJoinedGroupSummary(db)
	tm := time.Unix(googleAlert.Incident.Started, 0)
	loc, _ := time.LoadLocation("Asia/Taipei")
	triggeredTime := tm.In(loc).Format("2006-01-02 15:04:05") + " (GMT+8)"
	message := fmt.Sprintf(`⚠️ *%s* alert triggered！ ⚠️:

	State： *%s*
	Time： *%s*
	Summary： *%s*

	Please check it out！`, googleAlert.Incident.PolicyName, googleAlert.Incident.State, triggeredTime, googleAlert.Incident.Summary)

	for _, group := range groups {
		if _, err := bot.PushMessage(group.GroupId, linebot.NewTextMessage(message)).Do(); err != nil {
			utililty.Logger(3, err.Error())
			fmt.Println(err)
		}
	}
	fmt.Println(googleAlert)
}
