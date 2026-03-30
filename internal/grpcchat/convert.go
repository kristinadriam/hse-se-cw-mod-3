package grpcchat

import (
	"errors"
	"time"

	"hse-se-cw-mod-3/internal/domain"
	chatv1 "hse-se-cw-mod-3/proto/chat/v1"
)

func domainToProto(m domain.Message) *chatv1.ChatMessage {
	return &chatv1.ChatMessage{
		SenderName:   m.SenderName,
		SentUnixNano: m.SentAt.UnixNano(),
		Body:         m.Body,
	}
}

func fromProto(pb *chatv1.ChatMessage) (domain.Message, error) {
	if pb == nil {
		return domain.Message{}, errors.New("nil chat message")
	}
	m := domain.Message{
		SenderName: pb.SenderName,
		SentAt:     time.Unix(0, pb.SentUnixNano).UTC(),
		Body:       pb.Body,
	}
	if err := m.Validate(); err != nil {
		return domain.Message{}, err
	}
	return m, nil
}
