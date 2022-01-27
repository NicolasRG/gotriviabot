package trivia

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const triviaUrl = "https://opentdb.com/api.php"

func GetNumOfTrivia(num int) (TriviaResponse, error) {
	//create url
	editedUrl := setNumOfTrivia(triviaUrl, num)

	//hit endpoint
	res, err := http.Get(editedUrl)
	if err != nil {
		fmt.Println(err)
		return TriviaResponse{}, fmt.Errorf("error in get request to trivia api")
	}

	//extract body
	var rawResponse TriviaResponse

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	parsingErr := json.NewDecoder(res.Body).Decode(&rawResponse)
	if err != nil {
		fmt.Println(parsingErr)
		return TriviaResponse{}, fmt.Errorf("error in parsing response")
	}

	return rawResponse, nil
}

func setNumOfTrivia(url string, num int) string {
	return url + "?amount=" + strconv.Itoa(num)
}
