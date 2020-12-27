package notifications

import (
	"errors"
	"log"
	"os"
	"testing"
	"time"
)

const mongoTestsHostURL = "localhost:27017"
const mongoDatabaseForTests = "lottery"
const mongoCollectionForTests = "christmas"
const prizeToUseForTests = 200.0
const ownerToUseForTests = "test-owner"
const originToUserForTests = "test-origin"

var year = time.Now().Year()

func TestConnectToMongo(t *testing.T) {
	ctx, client, _, err := ConnectToMongo(mongoTestsHostURL, "", "")
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal(err)
	}
}

func TestConnectToMongoErrorNoUsername(t *testing.T) {
	originalMongoRootUsername := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	os.Setenv("MONGO_INITDB_ROOT_USERNAME", "")
	_, _, _, err := ConnectToMongo(mongoTestsHostURL, "", "")

	expectedError := errors.New("is not possible store notifications results in mongodb no value provided for MongoDB username")

	if err.Error() != expectedError.Error() {
		log.Fatal("Error expected is not equal to error obtained")
	}
	defer os.Setenv("MONGO_INITDB_ROOT_USERNAME", originalMongoRootUsername)
}

func TestConnectToMongoErrorNoPassword(t *testing.T) {
	originalMongoRootPassword := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	os.Setenv("MONGO_INITDB_ROOT_PASSWORD", "")
	_, _, _, err := ConnectToMongo(mongoTestsHostURL, "", "")

	expectedError := errors.New("is not possible store notifications results in mongodb no value provided for MongoDB password")

	if err.Error() != expectedError.Error() {
		log.Println(err)
		log.Fatal("Error expected is not equal to error obtained")
	}
	os.Setenv("MONGO_INITDB_ROOT_PASSWORD", originalMongoRootPassword)
}

func TestGetMongoCollection(t *testing.T) {
	_, client, _, err := ConnectToMongo(mongoTestsHostURL, "", "")
	if err != nil {
		log.Fatal(err)
	}
	mongoCollection := GetMongoCollection(client, mongoDatabaseForTests, mongoCollectionForTests)
	if mongoCollection.Name() != mongoCollectionForTests {
		log.Fatal(errors.New("mongoCollections obtained are not equal"))
	}
}

func TestAddNotificationToMongo(t *testing.T) {
	notificationAddResultErr := AddNotificationToMongo(mongoTestsHostURL, "", "", mongoCollectionForTests, year, ownerToUseForTests, numberToUseForCheck, prizeToUseForTests, originToUserForTests)
	if notificationAddResultErr != nil {
		log.Fatal(notificationAddResultErr)
	}
	defer testFindAndDeleteNotificationsFromMongo(mongoCollectionForTests, year, ownerToUseForTests, numberToUseForCheck, prizeToUseForTests, originToUserForTests)
}

func TestCheckNotificationAlreadyExistsInMongoCollection(t *testing.T) {
	expectedResult := true
	notificationAddResultErr := AddNotificationToMongo(mongoTestsHostURL, "", "", mongoCollectionForTests, year, ownerToUseForTests, numberToUseForCheck, prizeToUseForTests, originToUserForTests)
	if notificationAddResultErr != nil {
		log.Fatal(notificationAddResultErr)
	}
	defer testFindAndDeleteNotificationsFromMongo(mongoCollectionForTests, year, ownerToUseForTests, numberToUseForCheck, prizeToUseForTests, originToUserForTests)

	ctx, client, _, err := ConnectToMongo(mongoTestsHostURL, "", "")
	if err != nil {
		log.Fatal(err)
	}
	notificationToCheck := NotificationMongo{
		Draw:     mongoCollectionForTests,
		Year:     year,
		Owner:    ownerToUseForTests,
		Number:   numberToUseForCheck,
		Prize:    prizeToUseForTests,
		Origin:   originToUserForTests,
		Notified: true,
	}
	mongoCollection := GetMongoCollection(client, mongoDatabaseForTests, mongoCollectionForTests)

	notificationExists, err := CheckNotificationExistsInMongoCollection(ctx, notificationToCheck, mongoCollection)
	if err != nil {
		log.Fatal(err)
	}
	if *notificationExists != expectedResult {
		log.Fatal("notifactionExists hasn't expected value")
	}
}

func TestCheckNotificationNoExistsInMongoCollection(t *testing.T) {
	expectedResult := false
	ctx, client, _, err := ConnectToMongo(mongoTestsHostURL, "", "")
	if err != nil {
		log.Fatal(err)
	}
	notificationToCheck := NotificationMongo{
		Draw:     mongoCollectionForTests,
		Year:     year,
		Owner:    ownerToUseForTests,
		Number:   numberToUseForCheck,
		Prize:    prizeToUseForTests,
		Origin:   originToUserForTests,
		Notified: true,
	}
	mongoCollection := GetMongoCollection(client, mongoDatabaseForTests, mongoCollectionForTests)

	notificationExists, err := CheckNotificationExistsInMongoCollection(ctx, notificationToCheck, mongoCollection)
	if err != nil {
		log.Fatal(err)
	}
	if *notificationExists != expectedResult {
		log.Fatal("notifactionExists hasn't expected value")
	}
}

func testFindAndDeleteNotificationsFromMongo(draw string, year int, owner string, number string, prize float64, origin string) {
	ctx, client, _, err := ConnectToMongo(mongoTestsHostURL, "", "")
	mongoCollection := GetMongoCollection(client, mongoDatabaseForTests, mongoCollectionForTests)
	if err != nil {
		log.Fatal(err)
	}
	notificationToDelete := NotificationMongo{
		Draw:     draw,
		Year:     year,
		Owner:    owner,
		Number:   number,
		Prize:    prize,
		Origin:   origin,
		Notified: true,
	}
	var notificationDeleted NotificationMongo
	deleteResult := mongoCollection.FindOneAndDelete(ctx, notificationToDelete)
	err = deleteResult.Decode(&notificationDeleted)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Deleted successfully notification from collection %s with id %v\n", draw, notificationDeleted.ID)
}
