package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func sendRequest(url string, headers map[string]string, rest int, i int, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()
	sem <- struct{}{}
	defer func() { <-sem }()

	start := time.Now()

	// print the infomation of the redirect
	client := &http.Client{
		Timeout: 10 * time.Second, // timeout for 10 seconds
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 0 {
				fmt.Printf("\nrequest %d has been redirect to: %s\n", i+1, req.URL.String())
			}
			return nil
		},
	}

	// print the create requests' infomation
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("create request %d failed: %v\n\n", i+1, err)
		return
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// print the failed requests' infomation
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("\nrequest %d failed, error: %v\n", i+1, err)
		time.Sleep(time.Duration(rest) * time.Second)
		return
	}
	defer resp.Body.Close()

	// print the succeed requests' infomation
	if resp.StatusCode == http.StatusOK {
		fmt.Printf("\n\nrequest %d succeed, final address: %s\n\n", i+1, resp.Request.URL.String())
		elapsed := time.Since(start).Seconds()

		fileSizeStr := resp.Header.Get("Content-Length")
		var fileSize float64
		if fileSizeStr != "" {
			sizeInt, err := strconv.Atoi(fileSizeStr)
			if err == nil {
				fileSize = float64(sizeInt) / 1024 / 1024
			}
		}
		speed := fileSize / elapsed
		fmt.Printf("usage %.2f seconds, file size: %.2f MB, speed: %.2f MB/s\n\n%s\n", elapsed, fileSize, speed, "----------------------------------------")
	} else {
		fmt.Printf("request %d failed, status code: %d\n\n%s\n", i+1, resp.StatusCode, "----------------------------------------")
	}

	time.Sleep(time.Duration(rest) * time.Second)
}

func main() {
	// get the flags
	times := flag.Int("t", 1, "request times")
	workers := flag.Int("w", 1, "request threads")
	ua := flag.String("ua", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 Edg/128.0.0.0", "user agent")
	rest := flag.Int("r", 1, "request rest time")
	flag.Parse()

	// the final string without the flag is the request address
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("request address")
		os.Exit(1)
	}
	url := args[0]

	// set the user agent
	headers := map[string]string{
		"User-Agent": *ua,
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, *workers)

	// send the requests
	for i := 0; i < *times; i++ {
		wg.Add(1)
		go sendRequest(url, headers, *rest, i, &wg, sem)
	}
	wg.Wait()

	fmt.Println("Done, please check the results.")
}
