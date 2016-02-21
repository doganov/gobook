package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	prefix_http  = "http://"
	prefix_https = "https://"
)

func main() {
	for _, url := range os.Args[1:] {
		if !(strings.HasPrefix(url, prefix_http) || strings.HasPrefix(url, prefix_https)) {
			url = prefix_http + url
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "exc-1-7: %v\n", err)
			os.Exit(1)
		}
		_, err = io.Copy(os.Stdout, resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "exc-1-7: reading %s: %v\n", url, err)
			os.Exit(1)
		}
	}
}
