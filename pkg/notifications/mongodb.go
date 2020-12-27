package notifications

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

// getMongoCredentials create a mongoDB credentials using username and
// password provided stored in environment variables or provided by arguments
func getMongoCredentials(mongoDBRootUsernameFromArg string, mongoDBRootPasswordFromArg string) (*options.Credential, error) {
	mongoDBRootUsername := ""
	mongoDBRootPassword := ""
	mongoDBRootUsernameFromEnv := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongoDBRootPasswordFromEnv := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	if mongoDBRootUsernameFromEnv != "" {
		mongoDBRootUsername = mongoDBRootUsernameFromEnv
	} else if mongoDBRootUsernameFromArg != "" {
		mongoDBRootUsername = mongoDBRootUsernameFromArg
	} else {
		err := errors.New("is not possible store notifications results in mongodb no value provided for MongoDB username")
		return nil, err
	}
	if mongoDBRootPasswordFromEnv != "" {
		mongoDBRootPassword = mongoDBRootPasswordFromEnv
	} else if mongoDBRootPasswordFromArg != "" {
		mongoDBRootPassword = mongoDBRootPasswordFromArg
	} else {
		err := errors.New("is not possible store notifications results in mongodb no value provided for MongoDB password")
		return nil, err
	}
	credentials := options.Credential{
		Username: mongoDBRootUsername,
		Password: mongoDBRootPassword,
	}
	return &credentials, nil
}

// ConnectToMongo check if is possible connect to mongo host provided and returns
// a mongoClient for interact with these host
func ConnectToMongo(mongoURL string, mongoRootUsername string, mongoRootPassword string) (context.Context, *mongo.Client, context.CancelFunc, error) {
	credentials, err := getMongoCredentials(mongoRootUsername, mongoRootPassword)
	if err != nil {
		return nil, nil, nil, err
	}
	clientOpts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", mongoURL)).SetAuth(*credentials)

	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		return nil, nil, nil, err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, nil, cancelFunc, err
	}

	return ctx, client, cancelFunc, nil
}

// GetMongoCollection obtain a collection into a specific database in mongo
func GetMongoCollection(mongoClient *mongo.Client, databaseName string, collectionName string) *mongo.Collection {
	notificationsDatabase := mongoClient.Database(databaseName)
	mongoCollection := notificationsDatabase.Collection(collectionName)

	return mongoCollection
}

// CheckNotificationExistsInMongoCollection checks is a Document that uses
// NotificationMongo struct exists in a specific collection
func CheckNotificationExistsInMongoCollection(ctx context.Context, notification NotificationMongo, mongoCollection *mongo.Collection) (*bool, error) {
	count, err := mongoCollection.CountDocuments(ctx, notification)
	exists := false
	if err != nil {
		return nil, err
	}
	if count >= 1 {
		log.Printf("Notification already inserted in collection %s\n", mongoCollection.Name())
		exists := true
		return &exists, nil
	}
	log.Printf("Notification is not inserted in collection %s yet\n", mongoCollection.Name())
	return &exists, nil
}

// insertNotificationInMongoCollection inserts a Document with struct NotificationMongo in a mongo host
func insertNotificationInMongoCollection(ctx context.Context, notification NotificationMongo, mongoCollection *mongo.Collection) (interface{}, error) {
	insertResult, err := mongoCollection.InsertOne(ctx, notification)
	if err != nil {
		return nil, err
	}
	return insertResult.InsertedID, nil
}

// AddNotificationToMongo add a Document with Notification struct to mongoHost
func AddNotificationToMongo(mongoHostURL string, mongoRootUsername string, mongoRootPassword string, draw string, year int, owner string, number string, finalPrize float64, origin string) error {
	ctx, client, _, err := ConnectToMongo(mongoHostURL, mongoRootUsername, mongoRootPassword)
	if err != nil {
		return err
	}
	notificationMongo := NotificationMongo{
		Draw:     draw,
		Year:     year,
		Owner:    owner,
		Number:   number,
		Prize:    finalPrize,
		Origin:   origin,
		Notified: true,
	}
	christmasLotteryCollection := GetMongoCollection(client, "lottery", draw)

	id, err := insertNotificationInMongoCollection(ctx, notificationMongo, christmasLotteryCollection)
	if err != nil {
		return err
	}
	log.Printf("Notification inserted sucessfully in mongodb collection with id %v\n", id)
	return nil
}
