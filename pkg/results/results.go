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
func GetAPILotteryDrawStatus(draw string, lotteryDrawResultsAPIURL string) (*LotteryDrawStatus, error) {
	log.Printf("Checking the status of the %s lottery draw\n", draw)

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
func CheckNumber(lotteryDrawResultsAPIURL string, number string, bet int, origin string) (*Result, error) {

	log.Printf("Checking prize obtained for number %s with %d€ bet and origin %s\n", number, bet, origin)

	res, err := requests.DoGetRequest(fmt.Sprintf("%s?s=1&n=%s", lotteryDrawResultsAPIURL, number))
	if err != nil {
		return nil, err
	}
	numberCheckResult, err := converStrResponseToResult(*res)
	if err != nil {
		return nil, err
	}
	return numberCheckResult, nil
}

func isAlreadyNotified(mongoHostURL string, mongoRootUsername string, mongoRootPassword string, draw string, year int, owner string, number string, prize float64, origin string) (*bool, error) {
	ctx, mongoClient, _, err := notifications.ConnectToMongo(mongoHostURL, mongoRootUsername, mongoRootPassword)
	if err != nil {
		return nil, err
	}
	notificationMongo := notifications.NotificationMongo{
		Draw:     draw,
		Year:     year,
		Owner:    owner,
		Number:   number,
		Prize:    prize,
		Origin:   origin,
		Notified: true,
	}
	christmasLotteryCollection := notifications.GetMongoCollection(mongoClient, "lottery", draw)
	exists, err := notifications.CheckNotificationExistsInMongoCollection(ctx, notificationMongo, christmasLotteryCollection)

	if err != nil {
		return nil, err
	}

	return exists, nil
}

// isNecessaryNotify check if is necessary notify for each number result
func isNecessaryNotify(finalPrize float64, notify bool, storeNotifications bool, mongoHostURL string, mongoRootUsername string, mongoRootPassword string, draw string, year int, owner string, number string, prize float64, origin string) (*bool, error) {
	isNecessaryNotify := false
	if finalPrize > 0 && notify {
		if storeNotifications {
			isAlreadyNotified, err := isAlreadyNotified(mongoHostURL, mongoRootUsername, mongoRootPassword, draw, year, owner, number, prize, origin)
			if err != nil {
				return nil, err
			}
			if !*isAlreadyNotified {
				isNecessaryNotify = true
				return &isNecessaryNotify, nil
			}
		} else {
			isNecessaryNotify := true
			return &isNecessaryNotify, nil
		}
	}
	return &isNecessaryNotify, nil
}

// GetProbabilityOfWin calculate the probability to win a prize
// with the number of numbers provided
func getProbabilityOfWin(number int) float64 {
	return 0.001 * float64(number)
}

// GetFinalPrizeFromBet calculate the prize to obtain using the prize
// and the bet
func getFinalPrizeFromBet(bet int, premio int) float64 {
	porcentage := float64(premio / 20.0)
	return float64(bet) * porcentage
}

// CheckPersonsNumbers checks the prize for the numbers of a list of persons
func CheckPersonsNumbers(lotteryDrawResultsAPIURL string, personsNumbersToCheck *PersonNumbersToCheck,
	draw string, notify bool, storeNotifications bool, mongoHostURL string,
	mongoRootUsername string, mongoRootPassword string, year int) error {

	personsNumbersToCheckNum := len(personsNumbersToCheck.PersonsNumbers)
	log.Printf("Numbers are going to be check from %d different owners\n", personsNumbersToCheckNum)
	for i := 0; i < personsNumbersToCheckNum; i++ {
		personNumbersToCheck := personsNumbersToCheck.PersonsNumbers[i]
		personNumbersToCheckNum := len(personNumbersToCheck.Numbers)

		err := checkPersonNumbers(lotteryDrawResultsAPIURL, personNumbersToCheck, personNumbersToCheckNum, draw, notify, storeNotifications, mongoHostURL, mongoRootUsername, mongoRootPassword, year)

		if err != nil {
			return err
		}
	}
	return nil
}

// checkPersonNumbers check the number for a specific person
func checkPersonNumbers(lotteryDrawResultsAPIURL string, personNumbersToCheck PersonNumbers, personNumbersToCheckNum int,
	draw string, notify bool, storeNotifications bool, mongoHostURL string, mongoRootUsername string, mongoRootPassword string, year int) error {

	owner := personNumbersToCheck.Owner
	probabilityOfWin := getProbabilityOfWin(personNumbersToCheckNum)

	log.Printf("%d numbers are going to be check from owner %s. The probability to win a prize is %g%%\n", personNumbersToCheckNum, owner, probabilityOfWin)
	for j := 0; j < personNumbersToCheckNum; j++ {
		number := personNumbersToCheck.Numbers[j].Number
		bet := personNumbersToCheck.Numbers[j].BetAmount
		origin := personNumbersToCheck.Numbers[j].Origin
		numberCheckResult, err := CheckNumber(lotteryDrawResultsAPIURL, number, bet, origin)
		if err != nil {
			return err
		}
		finalPrize := getFinalPrizeFromBet(bet, numberCheckResult.Premio)
		log.Printf("%s won %.2f€ in number %s with %d € bet and origin %s\n", owner, finalPrize, number, bet, origin)
		isNecessaryNotify, err := isNecessaryNotify(finalPrize, notify, storeNotifications, mongoHostURL, mongoRootUsername, mongoRootPassword, draw, year, owner, number, finalPrize, origin)
		if err != nil {
			return err
		}
		if *isNecessaryNotify {
			notificationResult, err := notifications.SendPushOverNotification(owner, finalPrize, number, origin)
			if err != nil {
				return err
			}
			log.Printf("Notification send correctly with id %s to PushOver App\n", notificationResult.Request)
			if storeNotifications {
				err = notifications.AddNotificationToMongo(mongoHostURL, mongoRootUsername, mongoRootPassword, draw, year, owner, number, finalPrize, origin)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
