package model

import "time"

type LogMessage struct {
	Timestamp     time.Time `json:"timestamp"`
	Hostname      string    `json:"hostname"`
	EventSource   string    `json:"event.source.type"`
	EventCategory string    `json:"event.category"`
	Message       string    `json:"message"`
	Blacklisted   bool      `json:"is.blacklisted"`
	Level         string    `json:"severity"`
	Username      string    `json:"username"`
}

type Result struct {
	ID struct {
		EventCategory string `bson:"eventcategory" json:"eventcategory"`
		Level         string `bson:"level" json:"level"`
	} `bson:"_id" json:"result"`
	Count int `bson:"count" json:"count"`
}
