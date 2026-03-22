package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Write([]byte("hello world"))
	}))
	defer server.Close()

	var stdout, stderr bytes.Buffer
	code := run([]string{"gurl", server.URL}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d: %s", code, stderr.String())
	}
	if stdout.String() != "hello world" {
		t.Errorf("expected 'hello world', got %q", stdout.String())
	}
}

func TestMissingURL(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"gurl"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if stderr.String() == "" {
		t.Error("expected usage message on stderr")
	}
}

func TestInvalidURL(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"gurl", "not-a-url"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if stderr.String() == "" {
		t.Error("expected error message on stderr")
	}
}
