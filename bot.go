package main

import (
	"fmt"
	"gotriviabot/discordclient"
	"gotriviabot/handlers"
)

func main() {
	fmt.Println("Welcome to the the trivia bot!!")

	err, client := discordclient.InitConnection()

	if err != nil {
		fmt.Println(err)
		return
	}

	handlers.RegisterTriviaHandlers(client)

	discordclient.CloseConnection(client)
}
