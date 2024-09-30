package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		path := r.URL.Path
		ip := r.RemoteAddr

		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("Welcome to %s", path)))

		timeSpent := time.Since(startTime).Milliseconds()
		timestamp := startTime.Unix()

		logEntry := fmt.Sprintf("%s,%s,%d,%d\n", path, ip, timestamp, timeSpent)
		if _, err := logFile.WriteString(logEntry); err != nil {
			fmt.Println("Error writing to log file:", err)
		}
	})

	http.ListenAndServe(":4664", nil)
}
