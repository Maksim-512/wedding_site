package models

import "time"

type RSVP struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Attendance bool      `json:"attendance"`
	Companion  string    `json:"companion,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type RSVPRequest struct {
	Name       string `json:"name"`
	Attendance string `json:"attendance"`
	Companion  string `json:"companion"`
}
