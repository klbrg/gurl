package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDefaultScheme(t *testing.T) {
	req, err := buildRequest(http.MethodGet, "example.com/api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.URL.Scheme != "https" {
		t.Errorf("expected scheme 'https', got %q", req.URL.Scheme)
	}
}

func TestExplicitSchemePreserved(t *testing.T) {
	req, err := buildRequest(http.MethodGet, "http://example.com/api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.URL.Scheme != "http" {
		t.Errorf("expected scheme 'http', got %q", req.URL.Scheme)
	}
}

func TestVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"gurl", "version"}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d: %s", code, stderr.String())
	}
	expected := fmt.Sprintf("gurl version %s (%s)\n", version, commit)
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

func TestGetCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Write([]byte("hello get"))
	}))
	defer server.Close()

	var stdout, stderr bytes.Buffer
	code := run([]string{"gurl", "get", server.URL}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d: %s", code, stderr.String())
	}
	if stdout.String() != "hello get" {
		t.Errorf("expected 'hello get', got %q", stdout.String())
	}
}

func TestGetNoURL(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"gurl", "get"}, &stdout, &stderr)

	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if stderr.String() == "" {
		t.Error("expected usage message on stderr")
	}
}

func TestJSONPrettyPrint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"name":"john","age":30}`))
	}))
	defer server.Close()

	var stdout, stderr bytes.Buffer
	code := run([]string{"gurl", server.URL}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d: %s", code, stderr.String())
	}
	expected := "{\n  \"age\": 30,\n  \"name\": \"john\"\n}\n"
	if stdout.String() != expected {
		t.Errorf("expected %q, got %q", expected, stdout.String())
	}
}

func TestNonJSONNotFormatted(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("just plain text"))
	}))
	defer server.Close()

	var stdout, stderr bytes.Buffer
	code := run([]string{"gurl", server.URL}, &stdout, &stderr)

	if code != 0 {
		t.Fatalf("expected exit code 0, got %d: %s", code, stderr.String())
	}
	if stdout.String() != "just plain text" {
		t.Errorf("expected 'just plain text', got %q", stdout.String())
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
