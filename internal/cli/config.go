package cli

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"strings"
)

// Config holds parsed CLI flags.
type Config struct {
	Name       string
	ListenAddr string // host:port for listen mode
	PeerAddr   string // host:port for dial mode; empty => listen
}

// Parse reads argv (typically os.Args[1:]).
func Parse(args []string) (Config, error) {
	fs := flag.NewFlagSet("chat", flag.ContinueOnError)
	out := flag.CommandLine.Output()
	fs.SetOutput(out)
	fs.Usage = func() {
		fmt.Fprintf(out, "Usage: chat -name <name> [-listen host:port] [-connect host:port]\n")
		fmt.Fprintln(out, "  Omit -connect to listen for one peer (default -listen :50051).")
		fs.PrintDefaults()
	}

	name := fs.String("name", "", "display name (required)")
	listen := fs.String("listen", ":50051", "listen address in server mode (host:port)")
	peer := fs.String("connect", "", "peer address host:port; if omitted, listen mode")

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}
	if strings.TrimSpace(*name) == "" {
		return Config{}, errors.New("flag -name is required")
	}
	cfg := Config{
		Name:       strings.TrimSpace(*name),
		ListenAddr: *listen,
		PeerAddr:   strings.TrimSpace(*peer),
	}
	if cfg.PeerAddr != "" {
		if _, err := net.ResolveTCPAddr("tcp", cfg.PeerAddr); err != nil {
			return Config{}, fmt.Errorf("invalid -connect address: %w", err)
		}
	} else {
		if _, err := net.ResolveTCPAddr("tcp", cfg.ListenAddr); err != nil {
			return Config{}, fmt.Errorf("invalid -listen address: %w", err)
		}
	}
	return cfg, nil
}
