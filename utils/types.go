package utils

import "time"

// JSONResponse represents the JSON event data. Currently we only need the message.
type JSONResponse struct {
	Message string
}

// EventData represents the websocket event to be retained.
type EventData struct {
	Response JSONResponse
	RealTime time.Time
}

// WordPlace represents the location of the world in relation to others.
type WordPlace struct {
	RealWord string
	Word     string
	Count    int
}
