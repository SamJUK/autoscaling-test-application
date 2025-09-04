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

func main() {
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

		hostname, err := os.Hostname()
		if err != nil {
			hostname = "unknown"
		}

		host := fmt.Sprintf("%s:%s", GetOutboundIP(), listenPort)

		tpl := template.Must(template.New("index").ParseFiles("index.html"))
		tpl.ExecuteTemplate(w, "index.html", map[string]string{
			"Time":     time.Now().Format(time.RFC1123),
			"Hostname": hostname,
			"Host":     host,
		})
	})

	defer l.Close()

	l = netutil.LimitListener(l, connectionCount)
	log.Fatal(http.Serve(l, nil))
}
