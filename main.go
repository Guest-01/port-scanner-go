package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ScanResult struct {
	Port   int
	IsOpen bool
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, getUsage())
		os.Exit(1)
	}

	host := os.Args[1]

	ports, err := parseCommaSeparatedPorts(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid port(s): %s. Details: %v\n", os.Args[2], err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	results := make([]ScanResult, len(ports))

	for i, port := range ports {
		wg.Add(1)
		go func(i, port int) {
			defer wg.Done()
			isOpen := scanPort(host, port)
			results[i] = ScanResult{Port: port, IsOpen: isOpen}
		}(i, port)
	}

	wg.Wait()

	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	for _, result := range results {
		if result.IsOpen {
			fmt.Printf("%-5d : Open\n", result.Port)
		} else {
			fmt.Printf("%-5d : Closed\n", result.Port)
		}
	}
}

func getUsage() string {
	msg := fmt.Sprintf("Usage: %s <host> <port(s)>\n", filepath.Base(os.Args[0]))
	msg += fmt.Sprintln("\nArguments:")
	msg += fmt.Sprintln("  <host>\t\t호스트 주소(도메인 또는 IP 주소)")
	msg += fmt.Sprintln("  <port(s)>\t\t포트(단일 또는 여러 개의 포트를 콤마(,)로 구분)")
	return msg
}

// "80,443,8080-8082" -> []int{80, 443, 8080, 8081, 8082}
func parseCommaSeparatedPorts(portsStr string) ([]int, error) {
	ports := make([]int, 0)
	for _, portStr := range strings.Split(portsStr, ",") {
		if strings.Contains(portStr, "-") {
			parsedRanges, err := parseRangePorts(portStr)
			if err != nil {
				return nil, err
			}
			ports = append(ports, parsedRanges...)
			continue
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}
		ports = append(ports, port)
	}
	for _, port := range ports {
		if port > 65535 || port < 1 {
			return nil, fmt.Errorf("port out of range: %v", port)
		}
	}
	return ports, nil
}

// "8080-8082" -> []int{8080, 8081, 8082}
func parseRangePorts(portsStr string) ([]int, error) {
	ports := make([]int, 0)
	split := strings.Split(portsStr, "-")
	if len(split) != 2 {
		return nil, fmt.Errorf("len(split) != 2: %v", split)
	}

	start, err := strconv.Atoi(split[0])
	if err != nil {
		return nil, fmt.Errorf("start is not a number: %v", split[0])
	}

	end, err := strconv.Atoi(split[1])
	if err != nil {
		return nil, fmt.Errorf("end is not a number: %v", split[1])
	}

	if start >= end {
		return nil, fmt.Errorf("start is bigger than end: %v", split)
	}

	for i := start; i <= end; i++ {
		ports = append(ports, i)
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
