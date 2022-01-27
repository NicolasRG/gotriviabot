package trivia

import "fmt"

type TriviaResponse struct {
	Response_code int8
	Results       []TriviaEntry
}

type TriviaEntry struct {
	Category          string
	QuestionType      string
	Difficulty        string
	Question          string
	Correct_answer    string
	Incorrect_answers []string
}

func (triviaResponse *TriviaResponse) Print() {
	fmt.Println("Response code", triviaResponse.Response_code)

	for index, entry := range triviaResponse.Results {
		fmt.Print(index, " ")
		fmt.Println(entry.Category, entry.Question, entry.Correct_answer, entry.Incorrect_answers, entry.Difficulty, entry.QuestionType)
	}
}

func (triviaResponse *TriviaResponse) GetLenOfResults() uint8 {
	return uint8(len(triviaResponse.Results))
}

func (triviaResponse *TriviaResponse) IsResponseValid() bool {
	if triviaResponse.Response_code == 0 {
		return true
	} else {
		return false
	}
}
