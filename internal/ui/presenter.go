package ui

import (
	"fmt"
	"io"
	"time"

	"hse-se-cw-mod-3/internal/domain"
)

// Presenter formats inbound messages for the console.
type Presenter struct {
	Out io.Writer
}

// Show writes sender, timestamp (UTC RFC3339), and body on separate lines for readability.
func (p Presenter) Show(m domain.Message) error {
	ts := m.SentAt.UTC().Format(time.RFC3339Nano)
	_, err := fmt.Fprintf(p.Out, "[%s] %s\n%s\n", ts, m.SenderName, m.Body)
	return err
}
