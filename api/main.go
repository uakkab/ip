package handler

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// getClientIP extracts the client's IP address from the request headers
func getClientIP(r *http.Request) string {
	// Try to get the IP from the X-Forwarded-For header (in case the request passed through a proxy)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	// Fallback to the remote address
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// getReverseDNS performs a reverse DNS lookup on the given IP address
func getReverseDNS(ip string) (string, error) {
	names, err := net.LookupAddr(ip)
	if err != nil {
		return "", err
	}
	if len(names) > 0 {
		return names[0], nil
	}
	return "No PTR record found", nil
}

// handler handles HTTP requests and responds with the client's IP address and hostname
func handler(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientIP(r)

	hostname, err := getReverseDNS(clientIP)
	if err != nil || hostname == "No PTR record found" {
		hostname = "Unable to perform reverse DNS lookup"
	}

	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>IP Address</title>
			<style>
				body {
					font-family: Arial, sans-serif;
					display: flex;
					justify-content: center;
					align-items: center;
					height: 100vh;
					margin: 0;
					flex-direction: column;
					text-align: center;
				}
				h1, h2 {
					margin: 10px 0;
				}
			</style>
		</head>
		<body>
			<h1>Your Public IP Address is:</h1>
			<h2>%s</h2>
			<h1>Hostname:</h1>
			<h2>%s</h2>
		</body>
		</html>
	`, clientIP, hostname)
	fmt.Fprint(w, html)
}

// HandleRequest is the entry point for Vercel
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	handler(w, r)
}

func main() {
	http.HandleFunc("/", HandleRequest)
	http.ListenAndServe(":8080", nil)
}
