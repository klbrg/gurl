package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) < 2 {
		fmt.Fprintln(stderr, "usage: gurl <url>")
		return 1
	}

	client := &http.Client{}

	resp, err := client.Get(args[1])
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
