package grpcchat

import (
	"strings"
	"testing"
	"time"

	"hse-se-cw-mod-3/internal/domain"
	chatv1 "hse-se-cw-mod-3/proto/chat/v1"
)

func TestFromProto_roundTrip(t *testing.T) {
	t.Parallel()
	want := domain.Message{
		SenderName: "alice",
		Body:       "hello",
		SentAt:     time.Unix(1700000000, 42).UTC(),
	}
	got, err := fromProto(domainToProto(want))
	if err != nil {
		t.Fatal(err)
	}
	if got.SenderName != want.SenderName || got.Body != want.Body || !got.SentAt.Equal(want.SentAt) {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestFromProto_nil(t *testing.T) {
	t.Parallel()
	_, err := fromProto(nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFromProto_invalid_domain(t *testing.T) {
	t.Parallel()
	_, err := fromProto(domainToProto(domain.Message{SenderName: "", Body: "x"}))
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestFromProto_body_too_long(t *testing.T) {
	t.Parallel()
	body := strings.Repeat("я", domain.MaxBodyRunes+1)
	_, err := fromProto(&chatv1.ChatMessage{
		SenderName:   "a",
		SentUnixNano: 0,
		Body:         body,
	})
	if err == nil {
		t.Fatal("expected error")
	}
}
