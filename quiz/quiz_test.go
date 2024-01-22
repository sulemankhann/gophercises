package main

import (
	"bytes"
	"encoding/csv"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestEvaluateAnswers(t *testing.T) {
	problems := []problem{
		{
			question: "5+5",
			answer:   "10",
		},
		{
			question: "7+3",
			answer:   "10",
		},
		{
			question: "1+1",
			answer:   "2",
		},
	}
	buffer := bytes.Buffer{}
	userResponses := []string{"10", "9", "2"}

	result := EvaluateAnswers(&buffer, problems, userResponses)

	got := buffer.String()
	want := "Total correct: 2 out of 3\n"

	if got != want {
		t.Errorf("Output mismatch:\n got: %q\nwant: %q", got, want)
	}

	// Check if the evaluation results match the expected results
	expectedResult := []bool{true, false, true}

	for i, expected := range expectedResult {
		if result[i] != expected {
			t.Errorf(
				"Unexpected evaluation result for question #%d. Got %t, expected %t",
				i+1,
				result[i],
				expected,
			)
		}
	}
}

func TestAskQuestion(t *testing.T) {
	problems := []problem{
		{
			question: "5+5",
			answer:   "10",
		},
		{
			question: "7+3",
			answer:   "10",
		},
		{
			question: "1+1",
			answer:   "2",
		},
	}
	buffer := bytes.Buffer{}
	answers := strings.NewReader("10\n10\n2\n") // Simulate user input

	timer := time.NewTimer(5 * time.Second)
	capturedAnswers := AskQuestion(&buffer, answers, problems, timer)
	expectedAnswers := []string{"10", "10", "2"}

	got := buffer.String()
	want := `Problem #1: 5+5 = Problem #2: 7+3 = Problem #3: 1+1 = `

	if got != want {
		t.Errorf("Output mismatch:\n got: %q\nwant: %q", got, want)
	}

	if !reflect.DeepEqual(capturedAnswers, expectedAnswers) {
		t.Errorf("Answers mismatch:\n got: %q\nwant: %q", capturedAnswers, expectedAnswers)
	}
}

func TestCreateProblems(t *testing.T) {
	records := [][]string{
		{"5+5", "10"},
		{"1+1", "2"},
		{"8+3", "11"},
		{"1+2", "3"},
	}

	problems := CreateProblems(records)

	// Check if the Questions and Answers slices are correctly populated
	expectedQuestions := []string{"5+5", "1+1", "8+3", "1+2"}
	expectedAnswers := []string{"10", "2", "11", "3"}

	for i, question := range expectedQuestions {
		if problems[i].question != question {
			t.Errorf(
				"Unexpected question at index %d. Got %s, expected %s",
				i,
				problems[i].question,
				question,
			)
		}
	}

	for i, answer := range expectedAnswers {
		if problems[i].answer != answer {
			t.Errorf(
				"Unexpected answer at index %d. Got %s, expected %s",
				i,
				problems[i].answer,
				answer,
			)
		}
	}
}

func TestReadCSVFile(t *testing.T) {
	t.Run("Read CSV file", func(t *testing.T) {
		// Create a CSV file for testing
		testData := [][]string{
			{"5+5", "10"},
			{"1+1", "2"},
			{"8+3", "11"},
			{"1+2", "3"},
		}
		filename := "test_data.csv"
		defer os.Remove(filename)

		err := createTestCSVFile(filename, testData)
		if err != nil {
			t.Fatalf("Error creating test CSV file: %v", err)
		}

		// Test the ReadCSVFile function
		records, err := ReadCSVFile(filename)
		if err != nil {
			t.Fatalf("Error reading CSV file: %v", err)
		}

		// Compare the expected and actual records
		if len(records) != len(testData) {
			t.Fatalf("Expected %d records, got %d", len(testData), len(records))
		}

		for row := 0; row < len(testData); row++ {
			if len(records[row]) != len(testData[row]) {
				t.Fatalf("Mismatch in record length at row %d", row)
			}

			for column := 0; column < len(testData[row]); column++ {
				if records[row][column] != testData[row][column] {
					t.Fatalf("Mismatch in record data at row %d, column %d", row, column)
				}
			}
		}
	})

	t.Run("Read CSV file, which does not exit", func(t *testing.T) {
		_, err := ReadCSVFile("unkown_file.csv")

		assertError(t, err, ErrFileNotFound)
	})

	t.Run("Read CSV file, which does not have valid data", func(t *testing.T) {
		// Create a CSV file for testing
		testData := [][]string{
			{"5+5", "10"},
			{"1+1", "2"},
			{"8+3", "11"},
			{"1+2"},
		}
		filename := "test_data.csv"
		defer os.Remove(filename)

		error := createTestCSVFile(filename, testData)
		if error != nil {
			t.Fatalf("Error creating test CSV file: %v", error)
		}

		// Test the ReadCSVFile function
		_, err := ReadCSVFile(filename)
		assertError(t, err, ErrParseError)
	})
}

func createTestCSVFile(filename string, data [][]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.WriteAll(data)
	writer.Flush()

	return nil
}

func assertError(t testing.TB, got error, want error) {
	t.Helper()

	if got == nil {
		t.Fatal("wanted an error but didn't get one")
	}

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
