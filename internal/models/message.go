package models

import "time"

type Message struct {
	ID   int
	From *User
	Date time.Time
	Text string
}

func NewMessage(id int, from *User, date int, text string) *Message {
	return &Message{
		ID:   id,
		From: from,
		Date: time.Unix(int64(date), 0),
		Text: text,
	}
}
