package mongo

import (
	"context"
	"fmt"
	"line_bot/model"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const uri = "mongodb://mongodb:27017"

func ConnectDB() (mongo.Client, error) {
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		// if err = client.Disconnect(context.TODO()); err != nil {
		// 	panic(err)
		// }
	}()
	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected and pinged.")
	return *client, nil
}

func RecieveMessage(message model.Message, client mongo.Client) {
	coll := client.Database("line").Collection("message")
	_, err := coll.InsertOne(context.TODO(), message)
	if err != nil {
		panic(err)
	}
	fmt.Println("Insert Successfully")
}

func QueryMessage(user string, client mongo.Client) []model.Message {
	var result []model.Message
	filter := bson.D{{"id", user}}
	coll := client.Database("line").Collection("message")
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	for cursor.Next(context.TODO()) {
		var message model.Message
		err := cursor.Decode(&message)
		if err != nil {
			panic(err)
		}
		result = append(result, message)
	}
	return result
}

func InsertGroup(group model.Group, client mongo.Client) {
	coll := client.Database("line").Collection("group")
	filter := bson.D{{"groupid", group.GroupId}}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	isEmpty := !cursor.Next(context.TODO())
	if isEmpty {
		_, err := coll.InsertOne(context.TODO(), group)
		if err != nil {
			fmt.Println("Cursor is empty, no documents found")
			panic(err)
		}
	}
	for cursor.Next(context.TODO()) {
		var document bson.M
		err := cursor.Decode(&document)
		if err != nil {
			log.Fatal(err)
		}
		if document["groupid"] == group.GroupId {
			continue
		}
		_, insertErr := coll.InsertOne(context.TODO(), group)
		if insertErr != nil {
			panic(insertErr)
		}
	}
	fmt.Println("Insert Group Successfully")
}

func GetAllJoinedGroupSummary(client mongo.Client) []model.Group {
	var result []model.Group
	coll := client.Database("line").Collection("group")
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	if err != nil {
		panic(err)
	}
	for cursor.Next(context.TODO()) {
		var group model.Group
		err := cursor.Decode(&group)
		if err != nil {
			panic(err)
		}
		result = append(result, group)
	}
	return result
}
