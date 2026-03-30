package domain

import (
	"errors"
	"fmt"
	"time"
	"unicode/utf8"
)

const MaxBodyRunes = 65536

// Message is the domain model for a single chat line.
type Message struct {
	SenderName string
	SentAt     time.Time
	Body       string
}

func (m Message) Validate() error {
	if m.SenderName == "" {
		return errors.New("sender name is empty")
	}
	if utf8.RuneCountInString(m.Body) > MaxBodyRunes {
		return fmt.Errorf("body exceeds %d runes", MaxBodyRunes)
	}
	return nil
}
