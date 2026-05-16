package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTruncateEnd(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		maxLen int
		want   string
	}{
		{name: "short", value: "mk_short", maxLen: 20, want: "mk_short"},
		{name: "long", value: "mk_live_abcdefghijklmnopqrstuvwxyz", maxLen: 20, want: "mk_live_abcdefghi..."},
		{name: "zero", value: "mk_live_abcdefghijklmnopqrstuvwxyz", maxLen: 0, want: "mk_live_abcdefghijklmnopqrstuvwxyz"},
		{name: "tiny", value: "abcdef", maxLen: 3, want: "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncateEnd(tt.value, tt.maxLen); got != tt.want {
				t.Fatalf("truncateEnd() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDateOnly(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "rfc3339", value: "2026-05-16T10:20:30Z", want: "2026-05-16"},
		{name: "short", value: "today", want: "today"},
		{name: "empty", value: "", want: "-"},
		{name: "spaces", value: "   ", want: "-"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dateOnly(tt.value); got != tt.want {
				t.Fatalf("dateOnly() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolveBodyReadsAtFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "body.txt")
	want := "hello from a file"

	if err := os.WriteFile(path, []byte(want), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := resolveBody("@" + path)
	if err != nil {
		t.Fatalf("resolveBody() error = %v", err)
	}
	if got != want {
		t.Fatalf("resolveBody() = %q, want %q", got, want)
	}
}

func TestResolveBodyReturnsInlineBody(t *testing.T) {
	want := "inline body"

	got, err := resolveBody(want)
	if err != nil {
		t.Fatalf("resolveBody() error = %v", err)
	}
	if got != want {
		t.Fatalf("resolveBody() = %q, want %q", got, want)
	}
}
