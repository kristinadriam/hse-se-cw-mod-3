package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"hse-se-cw-mod-3/internal/cli"
	"hse-se-cw-mod-3/internal/grpcchat"
	"hse-se-cw-mod-3/internal/session"
	"hse-se-cw-mod-3/internal/transport"
)

func main() {
	if err := run(); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := cli.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var tr transport.Transport
	if cfg.PeerAddr != "" {
		tr, err = grpcchat.Dial(ctx, cfg.PeerAddr)
	} else {
		fmt.Fprintf(os.Stderr, "listening on %s ...\n", cfg.ListenAddr)
		tr, err = grpcchat.Listen(ctx, cfg.ListenAddr)
	}
	if err != nil {
		return err
	}

	return session.Run(ctx, cfg.Name, tr, os.Stdin, os.Stdout)
}
