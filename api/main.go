package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"
)

func getPublicIP() (string, error) {
	cmd := exec.Command("curl", "ifconfig.me")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

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

func handler(w http.ResponseWriter, r *http.Request) {
	ip, err := getPublicIP()
	if err != nil {
		http.Error(w, "Unable to get public IP", http.StatusInternalServerError)
		return
	}

	hostname, err := getReverseDNS(ip)
	if err != nil {
		http.Error(w, "Unable to perform reverse DNS lookup", http.StatusInternalServerError)
		return
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
    `, ip, hostname)
	fmt.Fprint(w, html)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server is running")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
