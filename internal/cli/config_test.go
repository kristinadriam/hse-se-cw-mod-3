package cli

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	t.Parallel()
	cfg, err := Parse([]string{"-name", "Alice", "-listen", ":0", "-connect", "127.0.0.1:9"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Name != "Alice" || cfg.PeerAddr != "127.0.0.1:9" {
		t.Fatalf("got %+v", cfg)
	}

	_, err = Parse([]string{})
	if err == nil {
		t.Fatal("expected error without name")
	}
	if !strings.Contains(err.Error(), "name") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestParse_invalid_listen(t *testing.T) {
	t.Parallel()
	_, err := Parse([]string{"-name", "x", "-listen", "not-a-host:port"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "listen") {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestParse_defaults(t *testing.T) {
	t.Parallel()
	cfg, err := Parse([]string{"-name", "Bob"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.ListenAddr != ":50051" || cfg.PeerAddr != "" {
		t.Fatalf("got %+v", cfg)
	}
}
