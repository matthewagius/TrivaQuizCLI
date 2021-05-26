/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var apiUrl = "http://localhost:8000/api/"

//single quiz question
type QuizQuestion struct {
	Id              int       `json:"id"`
	Topic           string    `json:"topic"`
	Question        string    `json:"question"`
	CorrectAnswerId string    `json:"correctAnswer"`
	Answers         [4]string `json:"answers"`
}

//get questions API response
type GetResponseData struct {
	Success bool           `json:"success"`
	Data    []QuizQuestion `json:"data"`
}

//selected answer for each question in the quiz
type SelectedAnswer struct {
	QuestionId int    `json:"questionId"`
	Answer     string `json:"answer"`
}

//all selected answers in the quiz
type SelectedAnswers []SelectedAnswer

//result to display for each question
type QuestionResult struct {
	Question       string `json:"question"`
	Correct        bool   `json:"correct"`
	SelectedAnswer string `json:"selectedAnswer"`
	CorrectAnswer  string `json:"correctAnswer"`
}

//final result including totals and percentage
type FinalResult struct {
	PersonalScore     int              `json:"personalScore"`
	QuestionResults   []QuestionResult `json:"questionResults"`
	PercentageRanking int              `json:"percentageRanking"`
}

//post questions API response
type PostResponseData struct {
	Success bool        `json:"success"`
	Data    FinalResult `json:"data"`
}

func init() {
	rootCmd.AddCommand(questionCmd)
}

// questionCmd represents the question command
var questionCmd = &cobra.Command{
	Use:   "startquiz",
	Short: "Time to test your knowledge on Sports, Technology and Science!",
	Run: func(cmd *cobra.Command, args []string) {
		//get all quiz questions
		questions, err := getQuestions()
		//check for errors
		if err != nil {
			fmt.Println(err)
		}
		if questions != nil {
			fmt.Println("Welcome to the Triva Quiz!")
			fmt.Println()
			//show questions as a multiselect
			selectedAnswers := promptQuestions(questions)

			fmt.Println("Please wait while your while your answers are verifed.")

			displayResults(selectedAnswers)
		}
	},
}

//get questions from API and unmarshal response
func getQuestions() ([]QuizQuestion, error) {

	responseBytes := getQuestionsAPIData(apiUrl + "questions")
	quizQuestion := GetResponseData{}
	if responseBytes != nil {
		if err := json.Unmarshal(responseBytes, &quizQuestion); err != nil {
			log.Printf("Could not unmarshal reponseBytes. %v", err)
			return nil, err
		}
	}
	return quizQuestion.Data, nil
}

//generic method to call api to get []byte response
func getQuestionsAPIData(baseAPI string) []byte {
	request, err := http.NewRequest(
		http.MethodGet, //method
		baseAPI,        //url
		nil,            //body
	)

	if err != nil {
		log.Printf("Could not request triva questions. %v", err)
	}

	request.Header.Add("Accept", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Could not make a request. %v", err)
	}
	if response.StatusCode == http.StatusOK {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Printf("Could not read response body. %v", err)
		}

		return responseBytes
	}
	return nil
}

//promt user with question and return selected answers
func promptQuestions(questions []QuizQuestion) (selectedAnswers SelectedAnswers) {
	for _, question := range questions {
		prompt := promptui.Select{
			Label: "[" + question.Topic + "]  " + question.Question,
			Items: question.Answers[:],
		}

		//get user result
		_, result, err := prompt.Run()
		fmt.Scan()
		//check for error
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
		}
		//save answer to selectedAnswers
		answer := SelectedAnswer{
			QuestionId: question.Id,
			Answer:     result,
		}
		selectedAnswers = append(selectedAnswers, answer)

	}
	return selectedAnswers
}

//display the quiz results to the user
func displayResults(selectedAnswers SelectedAnswers) error {
	//call api to post answers and get results
	responseBytes := postAnswersAPIData(apiUrl+"answers", selectedAnswers)

	quizScoreResult := PostResponseData{}
	if responseBytes != nil {
		if err := json.Unmarshal(responseBytes, &quizScoreResult); err != nil {
			log.Printf("Could not unmarshal reponseBytes. %v", err)
			return err
		}
	}
	r := color.New(color.FgRed, color.Bold)   // red for wrong answers
	g := color.New(color.FgGreen, color.Bold) // green for correct answers
	w := color.New(color.FgWhite)
	wb := w.Add(color.Bold)

	for _, questionResult := range quizScoreResult.Data.QuestionResults {

		w.Println(questionResult.Question)

		wb.Printf("The selected answer was %v\n", questionResult.SelectedAnswer)

		if bool(questionResult.Correct) {
			g.Println("And your answer was correct")
		} else {

			r.Printf("The correct answer was %v\n", questionResult.CorrectAnswer)
		}
		fmt.Println()
	}
	//print user score
	wb.Printf("Your Socre is of %d out of %d", quizScoreResult.Data.PersonalScore, len(quizScoreResult.Data.QuestionResults))
	fmt.Println()
	//print percentage ranking
	wb.Printf("Your ranked better than %d %% of people who took this quiz.", quizScoreResult.Data.PercentageRanking)

	return nil
}

//post answers to api
func postAnswersAPIData(baseAPI string, answers SelectedAnswers) []byte {
	answersJSON, err := json.Marshal(answers)
	if err != nil {
		log.Printf("Could not read response body. %v", err)
		os.Exit(1)
	}

	response, err := http.Post(baseAPI, "application/json", bytes.NewBuffer(answersJSON))
	if err != nil {
		log.Printf("Could not make a request. %v", err)
	}

	if response.StatusCode == http.StatusOK {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Printf("Could not read response body. %v", err)
		}

		return responseBytes
	}
	return nil
}
