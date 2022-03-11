package results

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const numberToUseForCheck = "74730"
const correctLotteryDrawResultsAPIURL = "http://api.elpais.com/ws/LoteriaNavidadPremiados"
const inCorrectLotteryDrawResultsAPIURLForJSONUnmarshall = "http://api.elpais.com/ws/LoteriaNavidadPremiadoss"
const actualLotteryDrawStatusCode = 4
const mongoCollectionForTests = "christmas"

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestGetAPILotteryDrawStatus(t *testing.T) {
	drawStatus, err := GetAPILotteryDrawStatus(mongoCollectionForTests, correctLotteryDrawResultsAPIURL)

	if err != nil {
		log.Fatal(err)
	}

	expectedResultStatus := actualLotteryDrawStatusCode
	expectedResultError := 0
	expectedDrawStatus := NewLotteryDrawStatus(expectedResultStatus, expectedResultError)
	if !(cmp.Equal(drawStatus, expectedDrawStatus)) {
		log.Fatal("drawStatus and expectedResultStatus for GetAPILotteryDrawStatus are not equal")
	}
}

func TestGetAPILotteryDrawStatusIncorrectUrlForJSONUnmarshall(t *testing.T) {
	_, err := GetAPILotteryDrawStatus(mongoCollectionForTests, inCorrectLotteryDrawResultsAPIURLForJSONUnmarshall)

	expectedError := errors.New("invalid character '>' after top-level value")
	if err.Error() != expectedError.Error() {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestGetAPILotteryDrawStatusIncorrectUrlForRequest(t *testing.T) {
	inCorrectLotteryDrawResultsAPIURLForRequest := "http:/api.elpais.com/ws/LoteriaNavidadPremiadoss"
	_, err := GetAPILotteryDrawStatus(mongoCollectionForTests, inCorrectLotteryDrawResultsAPIURLForRequest)
	inCorrectLotteryDrawResultsAPIURLForRequestInRequest := "http:///api.elpais.com/ws/LoteriaNavidadPremiadoss"
	expectedError := fmt.Errorf(
		"Get \"%s?s=1\": http: no Host in request URL",
		inCorrectLotteryDrawResultsAPIURLForRequestInRequest,
	)
	if err.Error() != expectedError.Error() {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestCheckAPILOtteryDrawStatusDrawNotStarted(t *testing.T) {
	_, err := CheckAPILOtteryDrawStatus(0)

	expectedError := errors.New("the draw has not started yet, it is not possible to get the results for the numbers")
	if err.Error() != expectedError.Error() {
		t.Fatal("The error obtained is not the expected")
	}
}

func TestCheckAPILOtteryDrawStatusDrawStartedWithoutAllNumbers(t *testing.T) {
	message, _ := CheckAPILOtteryDrawStatus(1)

	expectedMessage := "The draw has started but not all the numbers are in the list"

	if *message != expectedMessage {
		log.Fatal("The messages obtained for draw status are not equal")
	}
}

func TestCheckAPILOtteryDrawStatusDrawFinishedWithProvisionalNumbers(t *testing.T) {
	message, _ := CheckAPILOtteryDrawStatus(2)

	expectedMessage := "The draw has ended but the list of numbers is provisional"

	if *message != expectedMessage {
		log.Fatal("The messages obtained for draw status are not equal")
	}
}

func TestCheckAPILOtteryDrawStatusDrawFinishedWithSemiofficialNumbers(t *testing.T) {
	message, _ := CheckAPILOtteryDrawStatus(3)

	expectedMessage := "The draw is over but the list of numbers is semi-official"

	if *message != expectedMessage {
		log.Fatal("The messages obtained for draw status are not equal")
	}
}

func TestCheckAPILOtteryDrawStatusDrawFinishedWithOfficialNumbers(t *testing.T) {
	message, _ := CheckAPILOtteryDrawStatus(4)

	expectedMessage := "The draw is over and the list of numbers is official"

	if *message != expectedMessage {
		log.Fatal("The messages obtained for draw status are not equal")
	}
}

func TestCheckAPILOtteryDrawStatusDrawFinishedWithUnknownCode(t *testing.T) {
	_, err := CheckAPILOtteryDrawStatus(5)

	expectedError := errors.New("code obtained for draw status is unknown")
	if err.Error() != expectedError.Error() {
		t.Fatal("The error obtained is not the expected")
	}
}

func resultEqualWithoutTimestamp(expected Result, obtained Result) bool {
	if expected.Error == obtained.Error && expected.Premio == obtained.Premio &&
		expected.Status == obtained.Status && expected.Numero == obtained.Numero {
		return true
	}
	return false
}

func TestCheckNumberCorrect(t *testing.T) {
	expectedResult, err := NewResult(numberToUseForCheck, 0, 51515151, actualLotteryDrawStatusCode, 0)
	if err != nil {
		log.Fatal(err)
	}
	result, err := CheckNumber(correctLotteryDrawResultsAPIURL, numberToUseForCheck, 20, "origin_test")

	if err != nil {
		log.Fatal(err)
	}

	if !resultEqualWithoutTimestamp(*expectedResult, *result) {
		log.Fatal("Result and expected result for number are not equal")
	}
}

func TestCheckNumberIncorrect(t *testing.T) {
	expectedResult, err := NewResult("0", 0, 0, 0, 1)
	if err != nil {
		log.Fatal(err)
	}
	result, err := CheckNumber(correctLotteryDrawResultsAPIURL, "7455730", 20, "origin_test")

	if err != nil {
		log.Fatal(err)
	}

	if !resultEqualWithoutTimestamp(*expectedResult, *result) {
		fmt.Println(result)
		log.Fatal("Result and expected result for number are not equal")
	}
}

func TestCheckNumbersWithoutNotify(t *testing.T) {
	numbersToCheck := []Number{{
		Number:    numberToUseForCheck,
		BetAmount: 20,
		Origin:    "test-origin",
	}}
	personNumbersToCheck := PersonNumbersToCheck{[]PersonNumbers{{
		Numbers: numbersToCheck,
		Owner:   "test-owner",
	}}}

	err := CheckPersonsNumbers(correctLotteryDrawResultsAPIURL, &personNumbersToCheck, mongoCollectionForTests, false, false, "", "", "", 0)

	if err != nil {
		log.Fatal(err)
	}
}
