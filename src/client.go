// +build client

package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

var requestQueue chan int
var completedRequests int
var mu sync.Mutex

// Worker function to perform a single HTTP request
func worker(id int, wg *sync.WaitGroup, client *fasthttp.Client, url string) {
	defer wg.Done()
	for range requestQueue {
		_, _, err := client.Get(nil, url)
		if err != nil {
			// Log the error to stderr but don't interrupt worker
			log.Printf("Worker %d: Error: %v\n", id, err)
			continue
		}

		mu.Lock()
		completedRequests++
		mu.Unlock()
	}
}

// Log requests per second
func logRate(rateLogger *log.Logger) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		mu.Lock()
		reqs := completedRequests
		completedRequests = 0
		mu.Unlock()
		rateLogger.Printf("Requests per second: %d\n", reqs)
	}
}

func run() {
	var wg sync.WaitGroup

	client := &fasthttp.Client{
		MaxConnsPerHost: 10000, // Increase the number of connections per host
	}

	rateFile, err := os.Create("client_rate.log")
	if err != nil {
		log.Fatalf("Error creating rate log file: %v\n", err)
	}
	defer rateFile.Close()

	rateLogger := log.New(rateFile, "", log.LstdFlags)

	requestQueue = make(chan int, 100000)
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker(i, &wg, client, url)
	}

	go logRate(rateLogger)

	start := time.Now()
	for i := 0; i < numRequests; i++ {
		requestQueue <- i
	}
	close(requestQueue)
	wg.Wait()

	duration := time.Since(start)
	log.Printf("Completed %d requests in %v\n", numRequests, duration)
}
