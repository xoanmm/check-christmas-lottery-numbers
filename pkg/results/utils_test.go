package results

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestOpenFileNumbersToCheckFile(t *testing.T) {
	numbersToCheck := []Number{{
		Number:    numberToUseForCheck,
		BetAmount: 20,
		Origin:    "test-origin",
	}}
	personNumbersToCheck := PersonNumbersToCheck{[]PersonNumbers{{
		Numbers: numbersToCheck,
		Owner:   "test-owner",
	}}}
	personNumbersToCheckJSON, err := json.Marshal(personNumbersToCheck)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("/tmp/results_to_check.json", personNumbersToCheckJSON, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	_, err = OpenFilePersonsNumbersToCheck("/tmp/results_to_check.json")
	if err != nil {
		log.Fatal(err)
	}
}
