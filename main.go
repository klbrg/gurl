package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const version = "0.1.0"

func buildRequest(method string, rawURL string) (*http.Request, error) {
	if !strings.Contains(rawURL, "://") {
		rawURL = "https://" + rawURL
	}
	req, err := http.NewRequest(method, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gurl/"+version)
	return req, nil
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) < 2 {
		fmt.Fprintln(stderr, "usage: gurl <url>")
		return 1
	}

	var req *http.Request
	var err error

	switch args[1] {
	case "version":
		fmt.Fprintln(stdout, "gurl version "+version)
		return 0
	case "get":
		if len(args) < 3 {
			fmt.Fprintln(stderr, "usage: gurl get <url>")
			return 1
		}
		req, err = buildRequest(http.MethodGet, args[2])
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	default:
		req, err = buildRequest(http.MethodGet, args[1])
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	defer resp.Body.Close()
	if strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		var parsed any
		if err := json.Unmarshal(body, &parsed); err == nil {
			pretty, err := json.MarshalIndent(parsed, "", "  ")
			if err == nil {
				fmt.Fprintln(stdout, string(pretty))
				return 0
			}
		}
		stdout.Write(body)
	} else {
		io.Copy(stdout, resp.Body)
	}
	return 0
}

func main() {
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}
