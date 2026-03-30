package transport

import (
	"context"

	"hse-se-cw-mod-3/internal/domain"
)

// Transport is a duplex message channel (gRPC or test double).
type Transport interface {
	Send(ctx context.Context, m domain.Message) error
	Receive(ctx context.Context) (domain.Message, error)
	Close() error
}
