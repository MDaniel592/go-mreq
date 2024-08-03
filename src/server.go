// +build server

package main

import (
	"log"
	"runtime"
	"time"
	"sync"

	"github.com/valyala/fasthttp"
)

var requestQueue chan *fasthttp.RequestCtx
var wg sync.WaitGroup

// Handler function for incoming HTTP requests
func handler(ctx *fasthttp.RequestCtx) {
	requestQueue <- ctx
}

// Worker function to process requests
func worker(id int) {
	defer wg.Done()
	for ctx := range requestQueue {
		start := time.Now()
		// Simulate some processing time
		time.Sleep(10 * time.Millisecond)
		ctx.SetBodyString("Hello, World!")
		ctx.Response.Header.Set("Connection", "close") // Ensure the connection is closed after the response
		duration := time.Since(start)
		log.Printf("Worker %d handled request in %v\n", id, duration)
	}
}

func run() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	requestQueue = make(chan *fasthttp.RequestCtx, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker(i)
	}

	server := &fasthttp.Server{
		Handler: handler,
		// Optionally, configure other server settings if needed
	}

	log.Printf("Server is running on port %s\n", serverAddr)
	if err := server.ListenAndServe(serverAddr); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}

	close(requestQueue)
	wg.Wait()
}
