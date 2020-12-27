package notifications

import "go.mongodb.org/mongo-driver/bson/primitive"

// Notification struct which contain information
// about notification for PushOver app
type Notification struct {
	Token   string `json:"token"`
	User    string `json:"user"`
	Message string `json:"message"`
}

// NotificationResult struct which contain information
// about notification result for PushOver app
type NotificationResult struct {
	Status  int    `json:"status"`
	Request string `json:"request"`
}

// NotificationMongo struct which contain information
// about notification send with a number result with prize
type NotificationMongo struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Draw     string             `bson:"draw"`
	Year     int                `bson:"year"`
	Owner    string             `bson:"owner"`
	Number   string             `bson:"number"`
	Prize    float64            `bson:"prize"`
	Origin   string             `bson:"origin"`
	Notified bool               `bson:"notified"`
}
