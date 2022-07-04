package utils

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"time"
)

const LOCAL_HOST = "127.0.0.1"

var IP_PATTERN = regexp.MustCompile(`\b(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`)

func IsFoundHost(host string, port uint16) bool {
	target := fmt.Sprintf("%s:%d", host, port)

	_, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil {
		fmt.Printf("%s %v\n", target, err)
		return false
	}

	return true
}

func FindNeighbors(
	myHost string,
	myPort uint16,
	startIp uint8,
	endIp uint8,
	startPort uint16,
	endPort uint16,
) []string {
	address := fmt.Sprintf("%s:%d", myHost, myPort)

	m := IP_PATTERN.FindStringSubmatch(myHost)
	if m == nil {
		return nil
	}

	prefixHost := fmt.Sprintf("%s.%s.%s", m[1], m[2], m[3])
	lastIp := 1
	neighbors := make([]string, 0)

	for port := startPort; port < endPort; port++ {
		for ip := startIp; ip < endIp; ip++ {
			guessHost := fmt.Sprintf("%s.%d", prefixHost, lastIp+int(ip))
			guessTarget := fmt.Sprintf("%s:%d", guessHost, port)

			if guessTarget != address && IsFoundHost(guessHost, port) {
				neighbors = append(neighbors, guessTarget)
			}
		}
	}

	return neighbors
}

func GetHost() string {
	hostname, err := os.Hostname()
	if err != nil {
		return LOCAL_HOST
	}

	address, err := net.LookupHost(hostname)
	if err != nil {
		return LOCAL_HOST
	}

	return address[0]
}
