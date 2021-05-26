## About the Project

A simple quiz with 10 questions, having one correct answer per question, composed of REST API and CLI that talks with the API

Just in memery, so no database.

## Requirements

- user should be presented with questions having multiple answers
- user should be able to select one answer per question
- user should be able to answer all the questions and then post his answers to get back how many correct answers they had and be displayed to the user
- user should see how well he rated compared to others that have taken the quiz "You scored higher than 60% of all quizzers

## Built With

- [Go](https://golang.org/)
- [Cobra](https://github.com/spf13/cobra)

## Prerequisites

- go installed on your machine
- TriviaQuestionsApi running locally on your machine on port 8000(see readme included with that project)

## Running this CLI

- download TriviaQuizClient
- run $ go run main.go startquiz


## References

- https://dev.to/divrhino/series/12228 
- https://github.com/manifoldco/promptui
- https://github.com/fatih/color
 