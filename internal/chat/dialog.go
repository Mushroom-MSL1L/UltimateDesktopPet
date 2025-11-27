package chat

import "time"

type Dialog struct {
	ID        uint
	Timestamp time.Time
	Request   string
	Response  string
}
