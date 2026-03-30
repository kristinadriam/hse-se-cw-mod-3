package grpcchat

import (
	"context"
	"net"
	"testing"
	"time"

	"hse-se-cw-mod-3/internal/domain"
	"hse-se-cw-mod-3/internal/transport"
)

func TestListenOn_Dial_roundTrip(t *testing.T) {
	t.Parallel()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := lis.Addr().String()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	type listenResult struct {
		tr  transport.Transport
		err error
	}
	resCh := make(chan listenResult, 1)
	go func() {
		tr, err := listenOn(ctx, lis)
		resCh <- listenResult{tr, err}
	}()

	clientTr, err := Dial(ctx, addr)
	if err != nil {
		t.Fatal(err)
	}
	defer clientTr.Close()

	lr := <-resCh
	if lr.err != nil {
		t.Fatal(lr.err)
	}
	serverTr := lr.tr
	defer serverTr.Close()

	msg := domain.Message{SenderName: "a", Body: "ping", SentAt: time.Unix(1, 0).UTC()}
	if err := clientTr.Send(ctx, msg); err != nil {
		t.Fatal(err)
	}
	rm, err := serverTr.Receive(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if rm.Body != "ping" || rm.SenderName != "a" {
		t.Fatalf("got %+v", rm)
	}

	reply := domain.Message{SenderName: "b", Body: "pong", SentAt: time.Unix(2, 0).UTC()}
	if err := serverTr.Send(ctx, reply); err != nil {
		t.Fatal(err)
	}
	cm, err := clientTr.Receive(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if cm.Body != "pong" || cm.SenderName != "b" {
		t.Fatalf("client got %+v", cm)
	}
}
