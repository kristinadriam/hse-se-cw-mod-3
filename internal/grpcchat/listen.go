package grpcchat

import (
	"context"
	"errors"
	"net"

	"hse-se-cw-mod-3/internal/transport"
	chatv1 "hse-se-cw-mod-3/proto/chat/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type chatServer struct {
	chatv1.UnimplementedChatServiceServer
	streamCh chan chatv1.ChatService_StreamChatServer
}

func (s *chatServer) StreamChat(stream chatv1.ChatService_StreamChatServer) error {
	select {
	case s.streamCh <- stream:
	default:
		return status.Error(codes.ResourceExhausted, "only one peer is allowed")
	}
	<-stream.Context().Done()
	return nil
}

var serverOpts = []grpc.ServerOption{
	grpc.MaxRecvMsgSize(maxMsg),
	grpc.MaxSendMsgSize(maxMsg),
}

// Listen waits for the first incoming peer and returns a transport for that stream.
func Listen(ctx context.Context, listenAddr string) (transport.Transport, error) {
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	return listenOn(ctx, lis)
}

// listenOn is the core listener used by Listen and tests that supply a net.Listener.
func listenOn(ctx context.Context, lis net.Listener) (transport.Transport, error) {
	streamCh := make(chan chatv1.ChatService_StreamChatServer, 1)
	srv := grpc.NewServer(serverOpts...)
	chatv1.RegisterChatServiceServer(srv, &chatServer{streamCh: streamCh})

	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve(lis) }()

	select {
	case stream := <-streamCh:
		tr := NewStreamTransport(stream, func() error {
			srv.GracefulStop()
			return nil
		})
		return tr, nil
	case err := <-serveErr:
		if err != nil {
			return nil, err
		}
		return nil, errors.New("grpc server stopped before a peer connected")
	case <-ctx.Done():
		srv.Stop()
		return nil, ctx.Err()
	}
}
