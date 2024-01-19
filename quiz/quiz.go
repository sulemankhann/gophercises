package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	ErrFileNotFound = errors.New("cannot read file, file doesn't exist")
	ErrParseError   = errors.New("unable to parse csv file")
)

type problem struct {
	question string
	answer   string
}

func ReadCSVFile(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, ErrFileNotFound
	}

	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, ErrParseError
	}

	return records, nil
}

func CreateProblems(records [][]string) []problem {
	problems := make([]problem, len(records))

	for i, record := range records {
		if len(record) == 2 {
			problems[i] = problem{
				question: record[0],
				answer:   record[1],
			}
		}
	}

	return problems
}

func AskQuestion(writer io.Writer, reader io.Reader, problems []problem) []string {
	bufReader := bufio.NewReader(reader)
	answers := make([]string, 0)

	for i, problem := range problems {
		fmt.Fprintf(writer, "Problem #%d: %s = ", i+1, problem.question)
		answer, _ := bufReader.ReadString('\n')
		answers = append(answers, strings.TrimSpace(answer))
	}

	return answers
}

func EvaluateAnswers(writer io.Writer, problems []problem, userAnswers []string) []bool {
	totalCorrectQuestions := 0
	evaluatedAnswers := []bool{}

	for i, answer := range userAnswers {
		result := problems[i].answer == answer

		if result {
			totalCorrectQuestions++
		}

		evaluatedAnswers = append(evaluatedAnswers, result)
	}

	fmt.Fprintf(
		writer,
		"Total correct: %d out of %d\n",
		totalCorrectQuestions,
		len(problems),
	)

	return evaluatedAnswers
}

func main() {
	fileFlag := flag.String("f", "problems.csv", "a csv file in the format of 'question,answer'")

	helpFlag := flag.Bool("h", false, "Show help message")

	flag.Parse()

	if *helpFlag {
		flag.PrintDefaults()
		return
	}

	filename := *fileFlag

	records, err := ReadCSVFile(filename)
	if err != nil {
		panic(err)
	}

	problems := CreateProblems(records)

	userAnswers := AskQuestion(os.Stdout, os.Stdin, problems)

	EvaluateAnswers(os.Stdout, problems, userAnswers)
}
