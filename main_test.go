package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"gurl", "version"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d: %s", code, stderr.String())
	}
	expected := "gurl version " + version + "\n"
	if stdout.String() != expected {
		t.Errorf("expected %q, got %q", expected, stdout.String())
	}
}

func TestUserAgent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		expected := "gurl/" + version
		if ua != expected {
			t.Errorf("expected User-Agent %q, got %q", expected, ua)
		}
		w.Write([]byte("ok"))
	}))
	defer server.Close()

	var stdout, stderr bytes.Buffer
	code := run([]string{"gurl", server.URL}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d: %s", code, stderr.String())
	}
}

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
