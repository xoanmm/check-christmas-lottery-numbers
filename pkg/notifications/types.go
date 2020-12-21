package notifications

// Notification struct which contain information
// about notification for PushOver app
type Notification struct {
	Token string `json:"token"`
	User string `json:"user"`
	Message string `json:"message"`
}

// NotificationResult struct which contain information
// about notification result for PushOver app
type NotificationResult struct {
	Status int `json:"status"`
	Request string `json:"request"`
}
