package notifications

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xoanmm/check-christmas-lottery-numbers/pkg/requests"
	"log"
	"os"
)

const pushOverURL = "https://api.pushover.net/1/messages.json"


// GetNecessaryNotificationParams checks if all the necessary
// params are set in the system for make calls to pushOver
func getNecessaryNotificationParams() (*string, *string, error){
	pushOverNotificationToken := os.Getenv("PUSH_OVER_NOTIFICATION_TOKEN")
	pushOverNotificationUser := os.Getenv("PUSH_OVER_NOTIFICATION_USER")

	if pushOverNotificationToken != "" && pushOverNotificationUser != "" {
		return &pushOverNotificationToken, &pushOverNotificationUser, nil
	}

	var err = errors.New("environment variables 'PUSH_OVER_NOTIFICATION_TOKEN' and " +
		"'PUSH_OVER_NOTIFICATION_USER' need to be set to send notifications with PushOver")
	return nil, nil, err
}

// sendPushOverNotificationMessage is in charge of make calls to pushOver with notifications
func sendPushOverNotificationMessage(notificationToken string, notificationUser string,
	message string) (*NotificationResult, error) {
	var notification Notification
	notification.Message = message
	notification.Token = notificationToken
	notification.User = notificationUser

	b, err := json.Marshal(notification)
	if err != nil {
		return nil, err
	}
	notificationRequestResultBody, err := requests.DoPostRequestWithBody(pushOverURL, b)
	if err != nil {
		return nil, err
	}
	var notificationResult *NotificationResult
	err = json.Unmarshal(notificationRequestResultBody, &notificationResult)
	if err != nil {
		return nil, err
	}
	return notificationResult, nil
}

// SendPushOverNotification makes a notification with pushOver to warn that a prize has been won
func SendPushOverNotification(finalPrize int, number int, origin string) (*NotificationResult, error) {
	log.Println("A notification is going to be send with pushOver")
	notificationToken, notificationUser, err := getNecessaryNotificationParams()
	if err != nil {
		return nil, err
	}
	message := fmt.Sprintf("You won %d € with number %d obtained from %s", finalPrize, number, origin)
	notificationResult, err := sendPushOverNotificationMessage(*notificationToken, *notificationUser, message)
	if err != nil {
		return nil, err
	}
	return notificationResult, err
}
