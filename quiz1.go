package main

import (
	"os"
	"encoding/csv"
	"fmt"
	"flag"
	"log"
	"bufio"
	"strconv"
	"strings"
	"github.com/Knetic/govaluate"
	"time"
)

func main() {
	// get input from input args
	fileName, timerSec := setupArgParsing()

	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(fmt.Errorf("could not open file: %s", fileName))
	}
	defer f.Close()

	// open file and read contents
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(fmt.Errorf("Something went wrong while reading the csv file: %s!", fileName))
	}

	processProblems(records, timerSec)
}

func setupArgParsing() (string, int) {
	csvPtr := flag.String("csv", "problems.csv", "path to the quiz csv file")
  	timePtr := flag.Int("time", 30, "seconds to allot for the quiz")

	flag.Parse()

	if len(flag.Args()) > 0 {
		flag.PrintDefaults();
		os.Exit(1)
	}

	return *csvPtr, *timePtr;
}

// processProblems takes a list of quiz problems and asks the user to solve them,
// keeping track of the number of correct/incorrect answers
func processProblems(problems [][]string, timerSec int) {

	totalNumProblems := len(problems)
	numCorrect := 0
	numWrong := 0

	timer1 := time.NewTimer(time.Second * time.Duration(timerSec))

	go func() {

		reader := bufio.NewReader(os.Stdin)

		for _, problemArg := range problems {

			problem := problemArg[0]
			//answer, _ := strconv.ParseFloat(problemArg[1], 64)

			fmt.Println(fmt.Sprintf("Solve: %s", problem))
			userAnswerBytes, _ := reader.ReadString('\n')
			userAnswer, err := strconv.ParseFloat(strings.TrimSpace(string(userAnswerBytes)), 32)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(userAnswer)
			solution := eval(problem)

			if solution == userAnswer {
				numCorrect += 1
			} else {
				numWrong += 1
			}
		}
	}()

	<-timer1.C
	fmt.Printf("Correct: %d\n", numCorrect)
	fmt.Printf("Wrong: %d\n", numWrong)
	fmt.Printf("Total: %d\n\n", totalNumProblems)
}

func eval(problem string) float64 {
	expression, err := govaluate.NewEvaluableExpression(problem)
	if err != nil {
		log.Fatal(err)
	}
	result, err := expression.Evaluate(nil);
	if err != nil {
		log.Fatal(err);
	}

	return result.(float64)
}
