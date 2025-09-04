package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"golang.org/x/net/netutil"
)

type TemplateData struct {
	Time           string
	Hostname       string
	Host           string
	MaxConnections string
	RequestTime    string
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func main() {
	outputFormat := getEnv("OUTPUT_FORMAT", "html")
	if outputFormat != "html" && outputFormat != "text" {
		log.Fatalf("Invalid OUTPUT_FORMAT: %s\nValid OUTPUT_FORMAT values are 'html' or 'text'", outputFormat)
	}

	listenAddress := getEnv("LISTEN_ADDRESS", "0.0.0.0")
	listenPort := getEnv("LISTEN_PORT", "80")
	listen := fmt.Sprintf("%s:%s", listenAddress, listenPort)

	connectionCount, err := strconv.Atoi(getEnv("CONNECTION_COUNT", "10"))
	if err != nil {
		log.Fatalf("Invalid CONNECTION_COUNT: %v", err)
	}

	requestTime, err := strconv.Atoi(getEnv("REQUEST_TIME", "3"))
	if err != nil {
		log.Fatalf("Invalid REQUEST_TIME: %v", err)
	}

	fmt.Println("Starting Autoscaling Test HTTP Server")
	fmt.Println("Listening on:", listen)
	fmt.Println("Max Connections:", connectionCount)
	fmt.Println("Request Time:", requestTime)

	l, err := net.Listen("tcp", listen)
	if err != nil {
		log.Fatalf("Listen: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		done := make(chan int)

		for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				for {
					select {
					case <-done:
						return
					default:
					}
				}
			}()
		}

		time.Sleep(time.Duration(requestTime) * time.Second)
		close(done)

		templateData := TemplateData{
			Time:           time.Now().Format(time.RFC1123),
			Hostname:       getHostname(),
			Host:           fmt.Sprintf("%s:%s", GetOutboundIP(), listenPort),
			MaxConnections: strconv.Itoa(connectionCount),
			RequestTime:    strconv.Itoa(requestTime) + "s",
		}

		if outputFormat == "text" {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			fmt.Fprintf(w, "Host: %s\n", templateData.Host)
			fmt.Fprintf(w, "Hostname: %s\n", templateData.Hostname)
			fmt.Fprintf(w, "Date: %s\n", templateData.Time)
			fmt.Fprintf(w, "Max Connections: %s\n", templateData.MaxConnections)
			fmt.Fprintf(w, "Request Time: %s\n", templateData.RequestTime)
			return
		} else if outputFormat == "html" {
			tpl := template.Must(template.New("index").ParseFiles("index.html"))
			tpl.ExecuteTemplate(w, "index.html", templateData)
		}

	})

	defer l.Close()

	l = netutil.LimitListener(l, connectionCount)
	log.Fatal(http.Serve(l, nil))
}
