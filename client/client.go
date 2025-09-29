package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

func getRandomElement(arr []string) string {
	return arr[rand.Intn(len(arr))]
}

func main() {
	rand.Seed(time.Now().UnixNano())

	usernames := []string{"Motadata", "Motabhai", "Motakaka"}
	hostnames := []string{"EQ5V3RA5F7H", "EQ5V3RA5F7Q", "EQ5V3RA5F7T", "EQ5V3RA5F7A", "EQ5V3RA5F7Q"}
	levels := []string{"<134>", "<132>", "<131>"}
	activity := []string{"logged on", "logged off", "failed to log in"}

	// Keep one persistent connection to the server
	url := os.Getenv("server_url")
	if url == "" {
		url = "localhost:9000"
	}

	conn, err := net.Dial("tcp", url)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		// Pick random elements
		pri := getRandomElement(levels)
		hostname := getRandomElement(hostnames)
		user := getRandomElement(usernames)
		action := getRandomElement(activity)

		// Build log line
		data := fmt.Sprintf("%s WIN-%s Microsoft-Windows-Security-Auditing: A user account was successfully %s. Account Name: %s\r\n",
			pri, hostname, action, user)

		// Send log
		_, err := conn.Write([]byte(data))
		if err != nil {
			fmt.Println("Error sending log:", err)
			break
		}

		fmt.Println("Sent log:", data)

		// Wait 1-2 seconds
		time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second)
	}
}
