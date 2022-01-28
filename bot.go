package main

import (
	"fmt"
	"gotriviabot/discordclient"
	"gotriviabot/triviaHelper"
)

func main() {
	fmt.Println("Welcome to the the trivia bot!!")

	err, client := discordclient.InitConnection()

	if err != nil {
		fmt.Println(err)
		return
	}
	//setup for trivia flow
	triviaStore := &(triviaHelper.ReplyStore{})
	triviaHelper.RegisterTriviaHandler(client, triviaStore)

	discordclient.CloseConnection(client)
}

type QuestionBucket struct {
}
