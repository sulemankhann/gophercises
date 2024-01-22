package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
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

func ShuffleProblems(problems []problem) []problem {
	fmt.Println("sss")

	shuffled := make([]problem, len(problems))
	for i := range problems {
		j := rand.Intn(len(problems) - i)
		shuffled[i], shuffled[i+j] = problems[i+j], problems[i]
	}
	return shuffled
}

func AskQuestion(
	writer io.Writer,
	reader io.Reader,
	problems []problem,
	timer *time.Timer,
) []string {
	bufReader := bufio.NewReader(reader)
	answers := make([]string, 0)

	for i, problem := range problems {
		fmt.Fprintf(writer, "Problem #%d: %s = ", i+1, problem.question)
		answerCh := make(chan string)

		go func() {
			answer, _ := bufReader.ReadString('\n')
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			return answers
		case answer := <-answerCh:
			answers = append(answers, strings.TrimSpace(answer))
		}
	}

	return answers
}

func EvaluateAnswers(writer io.Writer, problems []problem, userAnswers []string) []bool {
	totalCorrectQuestions := 0
	evaluatedAnswers := make([]bool, len(problems))

	for i, answer := range userAnswers {
		result := problems[i].answer == answer

		if result {
			totalCorrectQuestions++
		}

		evaluatedAnswers[i] = result
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
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Seed(time.Now().UnixNano())

	fileFlag := flag.String("file", "problems.csv", "a csv file in the format of 'question,answer'")

	timeFlag := flag.Int("time", 30, "the time limit for the quiz in second")

	helpFlag := flag.Bool("h", false, "Show help message")

	shuffleFlag := flag.Bool("shuffle", false, "change order of the quiz")

	flag.Parse()

	if *helpFlag {
		flag.PrintDefaults()
		return
	}

	filename := *fileFlag
	timeLimit := time.Duration(*timeFlag) * time.Second

	records, err := ReadCSVFile(filename)
	if err != nil {
		panic(err)
	}

	problems := CreateProblems(records)

	if *shuffleFlag {
		problems = ShuffleProblems(problems)
	}

	timer := time.NewTimer(timeLimit)

	userAnswers := AskQuestion(os.Stdout, os.Stdin, problems, timer)

	EvaluateAnswers(os.Stdout, problems, userAnswers)
}
