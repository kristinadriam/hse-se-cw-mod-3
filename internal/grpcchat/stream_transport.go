package grpcchat

import (
	"context"

	"hse-se-cw-mod-3/internal/domain"
	"hse-se-cw-mod-3/internal/transport"
	chatv1 "hse-se-cw-mod-3/proto/chat/v1"
)

// bidi is implemented by both client and server gRPC stream types.
type bidi interface {
	Send(*chatv1.ChatMessage) error
	Recv() (*chatv1.ChatMessage, error)
	Context() context.Context
}

// StreamTransport adapts a gRPC bidi stream to transport.Transport.
type StreamTransport struct {
	s bidi
	// close releases the stream / connection (CloseSend + conn.Close, or GracefulStop server).
	closeFn func() error
}

var _ transport.Transport = (*StreamTransport)(nil)

// NewStreamTransport wraps a bidi stream. closeFn may be nil (caller manages lifecycle).
func NewStreamTransport(s bidi, closeFn func() error) *StreamTransport {
	return &StreamTransport{s: s, closeFn: closeFn}
}

func (t *StreamTransport) Send(ctx context.Context, m domain.Message) error {
	if err := m.Validate(); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return t.s.Send(domainToProto(m))
}

func (t *StreamTransport) Receive(ctx context.Context) (domain.Message, error) {
	type res struct {
		msg *chatv1.ChatMessage
		err error
	}
	ch := make(chan res, 1)
	go func() {
		msg, err := t.s.Recv()
		ch <- res{msg, err}
	}()
	select {
	case <-ctx.Done():
		return domain.Message{}, ctx.Err()
	case r := <-ch:
		if r.err != nil {
			return domain.Message{}, r.err
		}
		return fromProto(r.msg)
	}
}

func (t *StreamTransport) Close() error {
	if t.closeFn == nil {
		return nil
	}
	return t.closeFn()
}
