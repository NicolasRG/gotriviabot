package triviaHelper

import (
	"fmt"
	"gotriviabot/triviaApi"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/k3a/html2text"
)

func RegisterTriviaHandler(client *discordgo.Session, store *ReplyStore) {
	fmt.Println("Registering Trivia Handler")
	store.clearStore()
	client.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		triviaSetup(s, m, store)
	})
}

func triviaSetup(s *discordgo.Session, m *discordgo.MessageCreate, store *ReplyStore) {

	var activeQuestion triviaApi.TriviaEntry

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Println(m.Content)
	if m.Content == "!getSomeTrivia" {

		if store.active {
			s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">"+"Question already in progress")
			return
		}

		rawResponse, err := triviaApi.GetNumOfTrivia(1)

		if err != nil {
			fmt.Println(err)
			fmt.Println("awww man :/")
			return
		}

		rawResponse.Print()
		activeQuestion = rawResponse.Results[0]
		typeOfQuestion := activeQuestion.QuestionType
		formatedQuestion, letterAnswer := formatQuestion(activeQuestion, typeOfQuestion)
		store.answerLetter = letterAnswer
		store.question = activeQuestion.Question
		store.active = true
		store.channelID = m.ChannelID
		s.ChannelMessageSend(m.ChannelID, formatedQuestion)

		//create timeout
		go func() {
			time.Sleep(30 * time.Second)
			resolveQuestion(s, store)
		}()
	}

	isAnswer, err := regexp.MatchString("![ABCD]", m.Content)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">"+" You Broke Me :(")
		return
	}

	if isAnswer && store.active {

		playerAns := strings.Replace(m.Content, "!", "", 1)

		store.replys[m.Author.ID] = playerAns

		s.MessageReactionAdd(m.ChannelID, m.ID, "✔️")
	}

	//some meme stuff
	//TODO :  MOVE TO ITS OWN HANLDER LATER
	//create a probability to react randomly to messages
	shouldReply := rand.Intn(101)

	if shouldReply > 75 {
		fmt.Println("reply :)")
		//get guild emojis
		emojis, err := s.GuildEmojis(m.GuildID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">"+" You Broke Me :(")
			return
		}

		lenOfEmojis := len(emojis)
		indOfEmoji := rand.Intn(lenOfEmojis)

		fmt.Println(emojis[indOfEmoji])

		reactErr := s.MessageReactionAdd(m.ChannelID, m.ID, ":"+emojis[indOfEmoji].Name+":"+emojis[indOfEmoji].ID)
		if reactErr != nil {
			fmt.Println(reactErr)
		}
	}
}

func resolveQuestion(s *discordgo.Session, store *ReplyStore) {

	resolvedAnswer := "Correct Answer : " + store.answerLetter
	for id, playerAns := range store.replys {

		fmt.Println("Comparing : " + playerAns + " : " + store.answerLetter)

		if strings.EqualFold(playerAns, store.answerLetter) {
			//s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">"+" Correct")
			resolvedAnswer = resolvedAnswer + "\n" + "<@" + id + ">" + " Correct"
		} else {
			//s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">"+" You suck")
			resolvedAnswer = resolvedAnswer + "\n" + "<@" + id + ">" + " You suck"
		}
	}

	s.ChannelMessageSend(store.channelID, resolvedAnswer)
	store.clearStore()
}

func formatQuestion(triviaObj triviaApi.TriviaEntry, typeOfQuestion string) (string, string) {

	cleanedUpResponse := html2text.HTML2Text(triviaObj.Question)
	formatedQuestion := "Question : " + cleanedUpResponse

	//randomize array
	possibleAnswers, correctAnswer := randomizeArrEntries(triviaObj.Incorrect_answers, triviaObj.Correct_answer, typeOfQuestion)

	for _, bucket := range possibleAnswers {
		formatedQuestion = formatedQuestion + "\n\t" + bucket.letterId + ": " + html2text.HTML2Text(bucket.answerString)
	}

	return formatedQuestion, correctAnswer
}

func randomizeArrEntries(incorrectAns []string, correctAns string, typeOfQuestion string) ([]answerBucket, string) {
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
	if typeOfQuestion == "boolean" {
		answers[1] = answerBucket{
			correctAns,
			true,
			"-",
		}

		swapTimes := rand.Intn(10) + 1
		for i := 0; i < swapTimes; i++ {
			tempBucket := answers[1]
			answers[1] = answers[0]
			answers[0] = tempBucket
		}

	} else {

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

type ReplyStore struct {
	replys       map[string]string
	channelID    string
	question     string
	answerLetter string
	active       bool
}

func (store *ReplyStore) clearStore() {
	store.replys = make(map[string]string)
	store.active = false
	store.answerLetter = "_"
	store.question = "_"
}
