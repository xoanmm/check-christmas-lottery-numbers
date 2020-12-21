package results

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xoanmm/check-christmas-lottery-numbers/pkg/notifications"
	"github.com/xoanmm/check-christmas-lottery-numbers/pkg/requests"
	"log"
	"strings"
)

// codeDrawNotStartedYet represents the code returned by the Christmas lottery
// API when the draw has not started
const codeDrawNotStartedYet = 0
// codeDrawStartedWithoutAllNumbers represents the code returned by the Christmas lottery
// API when the draw has started but all the numbers have not been published yet
const codeDrawStartedWithoutAllNumbers = 1
// codeDrawEndWithProvisionalNumbers represents the code returned by the Christmas lottery
// API when the draw has started but the list of numbers is provisional
const codeDrawEndWithProvisionalNumbers = 2
// codeDrawEndWithSemiofficialNumbers represents the code returned by the Christmas lottery
// API when the draw has started but the list of numbers is semi-official
const codeDrawEndWithSemiofficialNumbers = 3
// codeDrawEndWithOfficialNumbers represents the code returned by the Christmas lottery
// API when the draw has started but the list of numbers is official
const codeDrawEndWithOfficialNumbers = 4

// convertStrResponseToLotteryDrawStatus converts the result obtained
// from checking the status of the lotterydraw as a string to the LotteryDrawStatus struct
func convertStrResponseToLotteryDrawStatus(response string) (*LotteryDrawStatus, error) {
	strResultSplit := strings.Split(response, "=")
	resultData := strResultSplit[1]
	var lotteryDrawStatus *LotteryDrawStatus
	err := json.Unmarshal([]byte(resultData), &lotteryDrawStatus)
	if err != nil {
		return nil, err
	}
	return lotteryDrawStatus, nil
}

// GetAPILotteryDrawStatus check the actual status of the lottery draw
func GetAPILotteryDrawStatus(lotteryDrawResultsAPIURL string) (*LotteryDrawStatus, error) {
	log.Println("Checking the status of the Christmas lottery draw")

	res, err := requests.DoGetRequest(fmt.Sprintf("%s?s=1", lotteryDrawResultsAPIURL))
	if err != nil {
		return nil, err
	}
	lotteryDrawStatus, err := convertStrResponseToLotteryDrawStatus(*res)
	if err != nil {
		return nil, err
	}
	return lotteryDrawStatus, nil
}

// CheckAPILOtteryDrawStatus check the status of the christmas lottery draw based on the status code provided
func CheckAPILOtteryDrawStatus(statusCode int) (*string, error) {
	if statusCode == codeDrawNotStartedYet {
		var err = errors.New("the draw has not started yet, it is not possible to get the results for the numbers")
		return nil, err
	} else if statusCode == codeDrawStartedWithoutAllNumbers {
		message := "The draw has started but not all the numbers are in the list"
		return &message, nil
	} else if statusCode == codeDrawEndWithProvisionalNumbers {
		message := "The draw has ended but the list of numbers is provisional"
		return &message, nil
	} else if statusCode == codeDrawEndWithSemiofficialNumbers {
		message := "The draw is over but the list of numbers is semi-official"
		return &message, nil
	} else if statusCode == codeDrawEndWithOfficialNumbers {
		message := "The draw is over and the list of numbers is official"
		return &message, nil
	} else {
		var err = errors.New("code obtained for draw status is unknown")
		return nil, err
	}
}

// converStrResponseToResult converts the result obtained
// from checking a number as a string to the Result struct
func converStrResponseToResult(response string) (*Result, error) {
	strResultSplit := strings.Split(response, "=")
	resultData := strResultSplit[1]

	var numberCheckResult *Result
	err := json.Unmarshal([]byte(resultData), &numberCheckResult)
	if err != nil {
		return nil, err
	}
	return numberCheckResult, nil
}

// CheckNumber check if a specific lottery number has prize
func CheckNumber(lotteryDrawResultsAPIURL string, number int, bet int, origin string) (*Result, error) {

	log.Printf("Checking prize obtained for number %d with %d€ bet and origin %s\n", number, bet, origin)

	res, err := requests.DoGetRequest(fmt.Sprintf("%s?s=1&n=%d", lotteryDrawResultsAPIURL, number))
	if err != nil {
		return nil, err
	}
	numberCheckResult, err := converStrResponseToResult(*res)
	if err != nil {
		return nil, err
	}
	return numberCheckResult, nil
}

// isNecessaryNotify check if is necessary notify for each number result
func isNecessaryNotify(finalPrize int, notify bool) bool{
	if finalPrize > 0 && notify {
		return true
	}
	return false
}

// GetProbabilityOfWin calculate the probability to win a prize
// with the number of numbers provided
func getProbabilityOfWin(number int) float64 {
	return 0.001 * float64(number)
}

// GetFinalPrizeFromBet calculate the prize to obtain using the prize
// and the bet
func getFinalPrizeFromBet(bet int, premio int) int {
	return (20 / bet) * (premio)
}

// CheckPersonsNumbers checks the prize for the numbers of a list of persons
func CheckPersonsNumbers(lotteryDrawResultsAPIURL string, personsNumbersToCheck *PersonNumbersToCheck, notify bool) error {

	personsNumbersToCheckNum := len(personsNumbersToCheck.PersonsNumbers)
	log.Printf("Numbers are going to be check from %d different owners\n", personsNumbersToCheckNum)
	for i:= 0; i < personsNumbersToCheckNum; i++ {
		personNumbersToCheck := personsNumbersToCheck.PersonsNumbers[i]
		personNumbersToCheckNum := len(personNumbersToCheck.Numbers)

		err := checkPersonNumbers(lotteryDrawResultsAPIURL, personNumbersToCheck, personNumbersToCheckNum, notify)

		if err != nil {
			return err
		}
	}
	return nil
}

// checkPersonNumbers check the number for a specific person
func checkPersonNumbers(lotteryDrawResultsAPIURL string, personNumbersToCheck PersonNumbers, personNumbersToCheckNum int,
	notify bool) error {
	owner := personNumbersToCheck.Owner
	probabilityOfWin := getProbabilityOfWin(personNumbersToCheckNum)

	log.Printf("%d numbers are going to be check from owner %s. The probabily to win a prize is %g%%\n", personNumbersToCheckNum, owner, probabilityOfWin)
	for j := 0; j < personNumbersToCheckNum; j++ {
		number := personNumbersToCheck.Numbers[j].Number
		bet := personNumbersToCheck.Numbers[j].BetAmount
		origin := personNumbersToCheck.Numbers[j].Origin
		numberCheckResult, err := CheckNumber(lotteryDrawResultsAPIURL, number, bet, origin)
		if err != nil {
			return err
		}
		finalPrize := getFinalPrizeFromBet(bet, numberCheckResult.Premio)
		log.Printf("Prize obtained for number %d with %d € bet and origin is %d€\n", number, bet, finalPrize)
		if isNecessaryNotify(finalPrize, notify) {
			notificationResult, err := notifications.SendPushOverNotification(finalPrize, number, origin)
			if err != nil {
				return err
			}
			log.Printf("Notification send correctly with id %s to PushOver App\n", notificationResult.Request)
		}
	}
	return nil
}