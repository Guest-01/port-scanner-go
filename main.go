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

	ports := []int{20, 21, 22, 23, 25, 53, 80, 110, 111, 135, 139, 443, 445, 993, 995, 1723, 3306, 3389, 5000, 5432, 8080, 8443}

	var wg sync.WaitGroup

	for _, port := range ports {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			scanPort(host, port)
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

func scanPort(host string, port int) {
	address := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		fmt.Printf("%d: Closed\n", port)
		return
	}
	defer conn.Close()
	fmt.Printf("%d: Open\n", port)
}
