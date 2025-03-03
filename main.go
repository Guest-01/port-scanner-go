package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, getUsage())
		os.Exit(1)
	}

	host := os.Args[1]

	ports, err := parsePorts(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid port(s):", os.Args[2])
		os.Exit(1)
	}

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
	msg := fmt.Sprintf("Usage: %s <host> <port(s)>\n", filepath.Base(os.Args[0]))
	msg += fmt.Sprintln("\nArguments:")
	msg += fmt.Sprintln("  <host>\t\t호스트 주소(도메인 또는 IP 주소)")
	msg += fmt.Sprintln("  <port(s)>\t\t포트(단일 또는 여러 개의 포트를 콤마(,)로 구분)")
	return msg
}

func parsePorts(portsStr string) ([]int, error) {
	if !strings.Contains(portsStr, ",") {
		port, err := strconv.Atoi(portsStr)
		if err != nil {
			return nil, err
		}
		if port > 65535 || port < 1 {
			return nil, fmt.Errorf("port out of range: %v", port)
		}
		return []int{port}, nil
	}
	ports := make([]int, 0)
	for _, portStr := range strings.Split(portsStr, ",") {
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}
		if port > 65535 || port < 1 {
			return nil, fmt.Errorf("port out of range: %v", port)
		}
		ports = append(ports, port)
	}
	return ports, nil
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
