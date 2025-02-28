package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, getUsage())
		os.Exit(1)
	}
	host := os.Args[1]
	// TODO: Do multi port scan
	scanPort(host, 8443)
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
