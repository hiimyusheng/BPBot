package utililty

import (
	"context"
	"fmt"
	mongodb "line_bot/mongo"
	"log"
)

type Log struct {
	Level   string
	Message string
}

func Logger(level int, message string) {
	db, DBerr := mongodb.ConnectDB()
	if DBerr != nil {
		log.Fatal(DBerr)
	}
	levelString := map[int]string{
		0: "Debug",
		1: "Info",
		2: "Warn",
		3: "Error",
	}
	var log Log
	log.Level = levelString[level]
	log.Message = message

	coll := db.Database("system").Collection("log")
	_, err := coll.InsertOne(context.TODO(), log)
	if err != nil {
		panic(err)
	}
	fmt.Println("Log Successfully")
}
