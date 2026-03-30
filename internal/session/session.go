package session

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"hse-se-cw-mod-3/internal/domain"
	"hse-se-cw-mod-3/internal/transport"
	"hse-se-cw-mod-3/internal/ui"
)

// Run starts stdin sender and inbound receiver until context is done or both sides finish.
func Run(ctx context.Context, localName string, t transport.Transport, in io.Reader, out io.Writer) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	present := ui.Presenter{Out: out}
	var g errgroup.Group
	g.Go(func() error {
		return sendLoop(ctx, cancel, localName, t, in)
	})
	g.Go(func() error {
		return recvLoop(ctx, cancel, t, present)
	})

	err := g.Wait()
	closeErr := t.Close()
	if err == nil {
		return closeErr
	}
	if closeErr != nil {
		return fmt.Errorf("%w; also close: %v", err, closeErr)
	}
	return err
}

func sendLoop(ctx context.Context, cancel context.CancelFunc, localName string, t transport.Transport, in io.Reader) error {
	lines := make(chan string)
	go func() {
		scanner := bufio.NewScanner(in)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			// Propagate read error by closing channel after optional log; EOF is normal.
			fmt.Fprintf(os.Stderr, "stdin: %v\n", err)
		}
		close(lines)
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case line, ok := <-lines:
			if !ok {
				cancel()
				return nil
			}
			trimmed := strings.TrimSpace(line)
			if trimmed == "/quit" || trimmed == "/exit" {
				cancel()
				return nil
			}
			if trimmed == "" {
				continue
			}
			msg := domain.Message{
				SenderName: localName,
				SentAt:     time.Now().UTC(),
				Body:       line,
			}
			if err := t.Send(ctx, msg); err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				}
				return err
			}
		}
	}
}

func recvLoop(ctx context.Context, cancel context.CancelFunc, t transport.Transport, p ui.Presenter) error {
	for {
		m, err := t.Receive(ctx)
		if err != nil {
			cancel()
			if errors.Is(err, context.Canceled) {
				return nil
			}
			if st, ok := status.FromError(err); ok && st.Code() == codes.Canceled {
				return nil
			}
			fmt.Fprintf(os.Stderr, "connection closed: %v\n", err)
			return nil
		}
		if err := p.Show(m); err != nil {
			cancel()
			return err
		}
	}
}
