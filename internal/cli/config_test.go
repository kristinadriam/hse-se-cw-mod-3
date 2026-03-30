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
