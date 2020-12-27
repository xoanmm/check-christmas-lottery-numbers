package results

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

var christmasLotteryAPIURL = "http://api.elpais.com/ws/LoteriaNavidadPremiados"
var childDayLotteryAPIURL = "http://api.elpais.com/ws/LoteriaNinoPremiados"

// GetDrawAPIURLToCheckNumbers obtain the URL from draw parameter value
func GetDrawAPIURLToCheckNumbers(draw string) (*string, error) {
	if draw == "christmas" {
		return &christmasLotteryAPIURL, nil
	} else if draw == "childhoods" {
		return &childDayLotteryAPIURL, nil
	} else {
		err := errors.New("value provided for draw argument is not valid the only arguments are christmas or childhoods")
		return nil, err
	}
}

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
