package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const log_filename = "exc-1-10.log"

func main() {
	start := time.Now()
	ch := make(chan string)
	log, err := os.OpenFile(log_filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		fmt.Printf("exc-1-10: can't open logfile: %v\n", err)
		os.Exit(1)
	}
	for _, url := range os.Args[1:] {
		go fetch(url, ch) // start a goroutine
	}
	for range os.Args[1:] {
		msg := <- ch // receive from channel ch
		fmt.Println(msg)
		fmt.Fprintln(log, msg)
	}
	elapsed := time.Since(start).Seconds()
	fmt.Printf("%.2fs elapsed\n", elapsed)
	fmt.Fprintf(log, "%.2fs elapsed\n", elapsed)
	log.Close()
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // don't leak resources

	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)
}
