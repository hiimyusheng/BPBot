package command

import (
	"line_bot/model"
	mongodb "line_bot/mongo"
	"line_bot/utililty"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.mongodb.org/mongo-driver/mongo"
)

func Add(projectId string, event *linebot.Event, bot *linebot.Client, db mongo.Client) string {

	if summary, err := bot.GetGroupSummary(event.Source.GroupID).Do(); err == nil {
		var group model.Group
		group.GroupId = summary.GroupID
		group.GroupName = summary.GroupName
		group.ProjectId = projectId
		dbErr := mongodb.InsertProject(group, db)
		if dbErr != nil {
			utililty.Logger(3, dbErr.Error())
			return "新增失敗"
		}
		return "新增成功"
	} else {
		utililty.Logger(3, err.Error())
		return "新增失敗"
	}
}
