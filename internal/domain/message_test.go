package domain

import (
	"strings"
	"testing"
)

func TestMessage_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		m       Message
		wantErr bool
	}{
		{name: "ok", m: Message{SenderName: "a", Body: "hi"}, wantErr: false},
		{name: "empty sender", m: Message{Body: "x"}, wantErr: true},
		{name: "body too long", m: Message{SenderName: "a", Body: strings.Repeat("я", MaxBodyRunes+1)}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.m.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
