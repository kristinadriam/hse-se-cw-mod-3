package session

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"hse-se-cw-mod-3/internal/domain"
	"hse-se-cw-mod-3/internal/transport"
)

type pipeTransport struct {
	mu     sync.Mutex
	recv   chan domain.Message
	sendFn func(context.Context, domain.Message) error
	closed bool
}

func (p *pipeTransport) Send(ctx context.Context, m domain.Message) error {
	if p.sendFn != nil {
		return p.sendFn(ctx, m)
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.recv <- m:
		return nil
	}
}

func (p *pipeTransport) Receive(ctx context.Context) (domain.Message, error) {
	select {
	case <-ctx.Done():
		return domain.Message{}, ctx.Err()
	case m, ok := <-p.recv:
		if !ok {
			return domain.Message{}, errors.New("closed")
		}
		return m, nil
	}
}

func (p *pipeTransport) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil
	}
	p.closed = true
	close(p.recv)
	return nil
}

func newPipeTransport() *pipeTransport {
	// Unbuffered: Send blocks until Receive runs, so a line cannot be "queued"
	// ahead of /quit while recvLoop is still idle (avoids flaky empty output).
	return &pipeTransport{recv: make(chan domain.Message)}
}

func TestRun_quit(t *testing.T) {
	t.Parallel()
	tr := newPipeTransport()
	in := strings.NewReader("/quit\n")
	var out bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := Run(ctx, "me", tr, in, &out); err != nil {
		t.Fatal(err)
	}
}

func TestRun_exit_alias(t *testing.T) {
	t.Parallel()
	tr := newPipeTransport()
	in := strings.NewReader("/exit\n")
	var out bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := Run(ctx, "me", tr, in, &out); err != nil {
		t.Fatal(err)
	}
}

// Pipe transport feeds sent messages back into Receive; recv loop prints them (self-echo scenario).
func TestRun_line_echo_to_output(t *testing.T) {
	t.Parallel()
	tr := newPipeTransport()
	in := strings.NewReader("hello\n/quit\n")
	var out bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := Run(ctx, "me", tr, in, &out); err != nil {
		t.Fatal(err)
	}
	s := out.String()
	if !strings.Contains(s, "hello") || !strings.Contains(s, "me") {
		t.Fatalf("expected sender and body in output, got %q", s)
	}
}

var _ transport.Transport = (*pipeTransport)(nil)
