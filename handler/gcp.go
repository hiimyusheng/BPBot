package handler

import (
	"fmt"
	"line_bot/model"
	mongodb "line_bot/mongo"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.mongodb.org/mongo-driver/mongo"
)

func ReceiveWebhookEvent(googleAlert model.Gcp, bot *linebot.Client, db mongo.Client) {
	mongodb.InsertAlert(googleAlert, db)

	groups := mongodb.GetAllJoinedGroupSummary(db)
	tm := time.Unix(googleAlert.Incident.Started, 0)
	triggeredTime := tm.Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(`⚠️ %s Triggered! ⚠️: 

	State: %s
	Time: %s
	Summary: %s
	
	Please check it out!`, googleAlert.Incident.PolicyName, googleAlert.Incident.State, triggeredTime, googleAlert.Incident.Summary)

	for _, group := range groups {
		if _, err := bot.PushMessage(group.GroupId, linebot.NewTextMessage(message)).Do(); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println(googleAlert)
}
