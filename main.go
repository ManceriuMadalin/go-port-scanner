package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ScanResult struct {
	Port    int
	Open    bool
	Service string
	Banner  string
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run main.go <host> <startPort> <endPort>")
		fmt.Println("Example: go run main.go localhost 20 1024")
		return
	}

	host := os.Args[1]

	startPort, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid startPort:", os.Args[2])
		return
	}

	endPort, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Invalid endPort:", os.Args[3])
		return
	}

	fmt.Printf("Scanning host %s ports %d-%d...\n", host, startPort, endPort)
	fmt.Printf("Detailed scan with service detection enabled\n")
	fmt.Printf("=" + strings.Repeat("=", 60) + "\n\n")

	results := make(chan ScanResult)
	var wg sync.WaitGroup

	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			result := scanPort(host, p)
			results <- result
		}(port)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var allResults []ScanResult
	openPorts := []int{}
	closedCount := 0

	for result := range results {
		allResults = append(allResults, result)
		if result.Open {
			openPorts = append(openPorts, result.Port)
		} else {
			closedCount++
		}
	}

	for i := 0; i < len(allResults); i++ {
		for j := i + 1; j < len(allResults); j++ {
			if allResults[i].Port > allResults[j].Port {
				allResults[i], allResults[j] = allResults[j], allResults[i]
			}
		}
	}

	for _, result := range allResults {
		if result.Open {
			fmt.Printf("Port %d is OPEN", result.Port)
			if result.Service != "" {
				fmt.Printf(" - Service: %s", result.Service)
			}
			if result.Banner != "" {
				fmt.Printf(" - Banner: %s", result.Banner)
			}
			fmt.Println()
		} else {
			fmt.Printf("Port %d is CLOSED\n", result.Port)
		}
	}

	fmt.Printf("\n" + strings.Repeat("=", 60) + "\n")
	fmt.Printf("SCAN SUMMARY:\n")
	fmt.Printf("Total ports scanned: %d\n", len(allResults))
	fmt.Printf("Open ports: %d\n", len(openPorts))
	fmt.Printf("Closed ports: %d\n", closedCount)

	if len(openPorts) > 0 {
		fmt.Printf("Open ports list: %v\n", openPorts)
	}
}

func scanPort(host string, port int) ScanResult {
	result := ScanResult{
		Port:    port,
		Open:    false,
		Service: "",
		Banner:  "",
	}

	address := fmt.Sprintf("%s:%d", host, port)
	timeout := time.Second * 2

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return result
	}
	defer conn.Close()

	result.Open = true
	result.Service = identifyService(port)
	result.Banner = grabBanner(conn, port)

	return result
}

func identifyService(port int) string {
	services := map[int]string{
		20:    "FTP Data",
		21:    "FTP Control",
		22:    "SSH",
		23:    "Telnet",
		25:    "SMTP",
		53:    "DNS",
		80:    "HTTP",
		110:   "POP3",
		143:   "IMAP",
		443:   "HTTPS",
		993:   "IMAPS",
		995:   "POP3S",
		1433:  "MS SQL Server",
		1521:  "Oracle DB",
		3306:  "MySQL",
		3389:  "RDP",
		5432:  "PostgreSQL",
		5000:  "Flask/Python Dev Server",
		5001:  "Flask/Python Dev Server",
		6379:  "Redis",
		7000:  "Cassandra/Custom",
		8000:  "HTTP Alt/Django",
		8080:  "HTTP Alt/Tomcat",
		8443:  "HTTPS Alt",
		9200:  "Elasticsearch",
		27017: "MongoDB",
	}

	if service, exists := services[port]; exists {
		return service
	}
	return "Unknown"
}

func grabBanner(conn net.Conn, port int) string {
	conn.SetReadDeadline(time.Now().Add(time.Second * 3))

	switch port {
	case 80, 8000, 8080:
		return grabHTTPBanner(conn)
	case 21:
		return grabFTPBanner(conn)
	case 22:
		return grabSSHBanner(conn)
	case 25:
		return grabSMTPBanner(conn)
	case 5432:
		return "PostgreSQL Database"
	case 3306:
		return grabMySQLBanner(conn)
	default:
		return grabGenericBanner(conn)
	}
}

func grabHTTPBanner(conn net.Conn) string {
	request := "GET / HTTP/1.1\r\nHost: localhost\r\nUser-Agent: PortScanner\r\n\r\n"
	conn.Write([]byte(request))

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "HTTP Service"
	}

	response := string(buffer[:n])
	lines := strings.Split(response, "\r\n")
	if len(lines) > 0 {
		return "HTTP - " + lines[0]
	}
	return "HTTP Service"
}

func grabFTPBanner(conn net.Conn) string {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "FTP Service"
	}
	banner := strings.TrimSpace(string(buffer[:n]))
	return "FTP - " + banner
}

func grabSSHBanner(conn net.Conn) string {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "SSH Service"
	}
	banner := strings.TrimSpace(string(buffer[:n]))
	return "SSH - " + banner
}

func grabSMTPBanner(conn net.Conn) string {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "SMTP Service"
	}
	banner := strings.TrimSpace(string(buffer[:n]))
	return "SMTP - " + banner
}

func grabMySQLBanner(conn net.Conn) string {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "MySQL Service"
	}
	if n > 0 {
		return "MySQL Database Server"
	}
	return "MySQL Service"
}

func grabGenericBanner(conn net.Conn) string {
	scanner := bufio.NewScanner(conn)
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500))

	if scanner.Scan() {
		banner := strings.TrimSpace(scanner.Text())
		if banner != "" {
			return banner
		}
	}
	return ""
}
