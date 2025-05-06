package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	// Parse CLI argument
	dir := flag.String("dir", ".", "Directory to serve")
	port := flag.Int("port", 8088, "Port to listen on")
	flag.Parse()

	// Check if directory exists
	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		log.Fatalf("Directory %s does not exist", *dir)
	}

	// File server handler
	fs := http.FileServer(http.Dir(*dir))

	http.Handle("/", fs)

	// Print local network IPs
	fmt.Println("Serving files from:", *dir)
	fmt.Println("Access the server at:")
	for _, ip := range getLocalIPs() {
		fmt.Printf("  http://%s:%d\n", ip, *port)
	}

	// Start server
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), nil))
}

func getLocalIPs() []string {
	var ips []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return ips
	}

	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil && ip.To4() != nil && !ip.IsLoopback() {
				ips = append(ips, ip.String())
			}
		}
	}
	return ips
}
