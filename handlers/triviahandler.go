package handlers

import (
	"fmt"
	"gotriviabot/trivia"
	"math/rand"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/k3a/html2text"
)

var (
	answer = "_"
)

func RegisterTriviaHandlers(client *discordgo.Session) {
	fmt.Println("Registering Trivia Handler")
	client.AddHandler(triviaSetup)
}

func triviaSetup(s *discordgo.Session, m *discordgo.MessageCreate) {

	var activeQuestion trivia.TriviaEntry

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!getSomeTrivia" {

		rawResponse, err := trivia.GetNumOfTrivia(1)

		if err != nil {
			fmt.Println(err)
			fmt.Println("awww man :/")
			return
		}

		rawResponse.Print()
		activeQuestion = rawResponse.Results[0]
		formatedQuestion, letterAnswer := formatQuestion(activeQuestion)
		answer = letterAnswer
		s.ChannelMessageSend(m.ChannelID, formatedQuestion)

	}

	isAnswer, err := regexp.MatchString("![ABCD]", m.Content)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">"+" You Broke Me :(")
	}

	if isAnswer {

		playerAns := strings.Replace(m.Content, "!", "", 1)
		fmt.Println("Comparing : " + playerAns + " : " + answer)
		if strings.EqualFold(playerAns, answer) {
			s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">"+" Correct")
		} else {
			s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">"+" You suck")
		}

	}
}

func formatQuestion(triviaObj trivia.TriviaEntry) (string, string) {

	cleanedUpResponse := html2text.HTML2Text(triviaObj.Question)
	formatedQuestion := "Question : " + cleanedUpResponse

	//randomize array
	possibleAnswers, correctAnswer := randomizeArrEntries(triviaObj.Incorrect_answers, triviaObj.Correct_answer)

	for _, bucket := range possibleAnswers {
		formatedQuestion = formatedQuestion + "\n\t" + bucket.letterId + ": " + html2text.HTML2Text(bucket.answerString)
	}

	return formatedQuestion, correctAnswer
}

func randomizeArrEntries(incorrectAns []string, correctAns string) ([]answerBucket, string) {
	letterArray := [4]string{"A", "B", "C", "D"}

	//create answers array
	answers := make([]answerBucket, 4)

	for i, s := range incorrectAns {
		tempBucket := answerBucket{
			s,
			false,
			"-",
		}
		answers[i] = tempBucket
	}

	answers[3] = answerBucket{
		correctAns,
		true,
		"-",
	}

	//randomize the order in the array by swaping random values
	for i := 0; i < 10; i++ {
		swapIndex := rand.Intn(3) + 1
		tempBucket := answers[swapIndex]
		answers[swapIndex] = answers[0]
		answers[0] = tempBucket
	}

	var correctLetter string
	fmt.Println(answers)
	fmt.Println(letterArray)
	//find the answer
	for i, ans := range answers {
		//assign the letter to it
		ans.letterId = letterArray[i]
		//set answer letter if it is
		fmt.Print(ans)
		if ans.isAnswer {
			correctLetter = ans.letterId
		}
		answers[i] = ans
	}
	fmt.Println(answers)
	fmt.Println(correctLetter)
	return answers, correctLetter

}

type answerBucket struct {
	answerString string
	isAnswer     bool
	letterId     string
}
