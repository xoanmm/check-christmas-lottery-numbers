package results

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// OpenFilePersonsNumbersToCheck open the file where the numbers to check for the persons are stored
func OpenFilePersonsNumbersToCheck(fileNumbersToCheck string) (*PersonNumbersToCheck, error) {
	log.Printf("Opening file %s with numbers to check\n", fileNumbersToCheck)
	jsonFile, err := os.Open(fileNumbersToCheck)

	if err != nil {
		return nil, err
	}
	log.Printf("Successfully opened file %s  with numbers to check\n", fileNumbersToCheck)

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var numbers *PersonNumbersToCheck

	err = json.Unmarshal(byteValue, &numbers)
	if err != nil {
		return nil, err
	}
	return numbers, nil
}
