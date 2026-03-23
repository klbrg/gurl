package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const version = "0.1.0"

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
		req, err = http.NewRequest(http.MethodGet, args[2], nil)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	default:
		req, err = http.NewRequest(http.MethodGet, args[1], nil)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}

	req.Header.Set("User-Agent", "gurl/"+version)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	defer resp.Body.Close()

	io.Copy(stdout, resp.Body)
	return 0
}

func main() {
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}
