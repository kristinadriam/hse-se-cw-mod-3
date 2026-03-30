package ui

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"hse-se-cw-mod-3/internal/domain"
)

func TestPresenter_Show(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	p := Presenter{Out: &buf}
	ts := time.Unix(1700000000, 0).UTC()
	m := domain.Message{SenderName: "u", Body: "line1", SentAt: ts}
	if err := p.Show(m); err != nil {
		t.Fatal(err)
	}
	s := buf.String()
	if !strings.Contains(s, "u") || !strings.Contains(s, "line1") || !strings.Contains(s, ts.Format(time.RFC3339Nano)) {
		t.Fatalf("unexpected output: %q", s)
	}
}
