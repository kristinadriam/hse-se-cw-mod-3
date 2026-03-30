package grpcchat

import (
	"context"

	"hse-se-cw-mod-3/internal/transport"
	chatv1 "hse-se-cw-mod-3/proto/chat/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const maxMsg = 64 << 20

var dialOpts = []grpc.DialOption{
	grpc.WithTransportCredentials(insecure.NewCredentials()),
	grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(maxMsg),
		grpc.MaxCallSendMsgSize(maxMsg),
	),
}

// Dial connects to a peer and opens a bidirectional chat stream.
func Dial(ctx context.Context, peerAddr string) (transport.Transport, error) {
	conn, err := grpc.DialContext(ctx, peerAddr, dialOpts...)
	if err != nil {
		return nil, err
	}
	client := chatv1.NewChatServiceClient(conn)
	stream, err := client.StreamChat(ctx)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	tr := NewStreamTransport(stream, func() error {
		_ = stream.CloseSend()
		return conn.Close()
	})
	return tr, nil
}
