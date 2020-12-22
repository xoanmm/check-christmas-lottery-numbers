package notifications

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

const numberToUseForCheck = "74730"

// setupEnv checks if there is an `.env' file where a series of variables
// used to perform the integration tests are defined
func setupEnv() error {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file provided")
	}
	return nil
}

func TestMain(m *testing.M) {
	// load .env file
	err := setupEnv()
	if err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	os.Exit(code)
}

func TestSendPushOverNotification(t *testing.T) {
	finalPrize := 2000.00
	number := numberToUseForCheck
	origin := "test-origin"

	expectedResultSendNotification := NotificationResult{
		Status:  1,
		Request: "",
	}

	resultSendNofitication, err := SendPushOverNotification(finalPrize, number, origin)
	if err != nil {
		log.Fatal(err)
	}

	if expectedResultSendNotification.Status != resultSendNofitication.Status {
		log.Fatal("Status Code for send notification are not expected")
	}
}
