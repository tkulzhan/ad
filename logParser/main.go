package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type LogEntry struct {
	URL       string `json:"url"`
	IP        string `json:"ip"`
	Timestamp int64  `json:"timestamp"`
	TimeSpent int64  `json:"timeSpent"`
}

func main() {
	tail := flag.Int("tail", 0, "Number of lines to read from the end of the file")
	flag.Parse()

	logFile, err := os.Open("server.log.csv")
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer logFile.Close()

	var lines []string
	scanner := bufio.NewScanner(logFile)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if *tail > 0 && *tail < len(lines) {
		lines = lines[len(lines)-*tail:]
	}

	jsonFile, err := os.Create("logs.json")
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
		return
	}
	defer jsonFile.Close()

	var logs []LogEntry
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) != 4 {
			fmt.Println("Invalid log format:", line)
			continue
		}

		var log LogEntry
		log.URL = parts[0]
		log.IP = parts[1]
		fmt.Sscanf(parts[2], "%d", &log.Timestamp)
		fmt.Sscanf(parts[3], "%d", &log.TimeSpent)

		logs = append(logs, log)
	}

	jsonData, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		fmt.Println("Error converting logs to JSON:", err)
		return
	}

	if _, err := jsonFile.Write(jsonData); err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

	fmt.Printf("Parsed the last %d log entries and written them to logs.json\n", len(logs))
}
