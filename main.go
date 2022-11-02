package main

import (
	"fmt"
	"line_bot/mongo"
)

func main() {

	defer func() {
		/* Catch fatal error from panic */
		if panicMsg := recover(); panicMsg != nil {
			fmt.Println("XInsight catch error.")
			fmt.Println(panicMsg.(error).Error())
		}

	}()

	mongo.ConnectDB()

}
