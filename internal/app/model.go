package app

type Notification struct {
	Message       string   `json:"message"`
	AlertChannels []string `json:"alert_channels"`
	UserId        string   `json:"user_id"`
	GUID          string   `json:"guid"`
}
