package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func worker(host string, ports, results chan int) {
	for p := range ports {
		if checkTcpPort(host, p) {
			results <- p
		} else {
			results <- 0
		}
	}
}

// checkTcpPort performs a TCP handshake with the specified host:port. Returns
// true if connection succeeds, false otherwise.
func checkTcpPort(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	timeout := time.Duration(5) * time.Second
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err == nil {
		conn.Close()
		return true
	}
	return false
}

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [-top N] HOST\n", os.Args[0])
		flag.PrintDefaults()
	}
	usageError := func() {
		flag.Usage()
		os.Exit(1)
	}

	// get command line args
	var top int
	flag.IntVar(&top, "top", 65535, "Specify the top N most common ports to scan (1 <= N <= 65535).")
	flag.Parse()

	// check command line args
	if flag.NArg() != 1 {
		fmt.Println("Exactly one host required for scanning!")
		usageError()
	}
	if top < 1 || top > 65535 {
		fmt.Println("Value for 'top' must be between 1 and 65535!")
		usageError()
	}

	host := flag.Arg(0)
	todo := make(chan int, 1024)
	done := make(chan int)

	fmt.Printf("Scanning top %d ports of '%s'...\n", top, host)

	// start worker pool
	for i := 0; i < cap(todo); i++ {
		go worker(host, todo, done)
	}

	// assign ports to workers
	go func() {
		for _, p := range sortedPorts[:top] {
			todo <- p
		}
	}()

	// collect results
	for i := 0; i < top; i++ {
		port := <-done
		if port != 0 {
			fmt.Printf("%d open\n", port)
		}
	}
	close(todo)
	close(done)
}
