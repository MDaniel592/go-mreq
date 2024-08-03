// main.go

package main

const (
	numRequests = 100000000  			     // Number of requests to perform
	serverAddr  = ":8111"  				     // Address the server listens on
	maxWorkers  = 8              			 // Maximum number of active workers
	url         = "http://192.168.1.75:8111" // URL the client sends requests to
)
func main() {
	run()
}