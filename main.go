package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, getUsage())
		os.Exit(1)
	}
	host := os.Args[1]

	ports := []int{22, 80, 443, 8080, 12345}

	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			if isOpen := scanPort(host, port); isOpen {
				fmt.Printf("%-5d : Open\n", port)
			} else {
				fmt.Printf("%-5d : Closed\n", port)
			}
		}(port)
	}

	wg.Wait()
}

func getUsage() string {
	msg := fmt.Sprintf("Usage: %s <host>\n", filepath.Base(os.Args[0]))
	msg += fmt.Sprintln("\nArguments:")
	msg += fmt.Sprintln("  <host>\t\t호스트 주소(도메인 또는 IP 주소)")
	return msg
}

func scanPort(host string, port int) bool {
	address := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
