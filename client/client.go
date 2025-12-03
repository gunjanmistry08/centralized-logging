package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func getRandomElement(arr []string) string {
	return arr[rand.Intn(len(arr))]
}

func main() {
	rand.Seed(time.Now().UnixNano())

	usernames := []string{"Motadata", "Motabhai", "Motakaka"}
	hostnames := []string{"EQ5V3RA5F7H", "EQ5V3RA5F7Q", "EQ5V3RA5F7T", "EQ5V3RA5F7A", "EQ5V3RA5F7Q"}
	levels := []string{"<134>", "<132>", "<131>"}
	activity := []string{"logged on", "logged off", "failed to log in"}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_direct", // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Keep one persistent connection to the server
	// url := os.Getenv("server_url")
	// if url == "" {
	// 	url = "localhost:9000"
	// }

	// conn, err := net.Dial("tcp", url)
	// if err != nil {
	// 	panic(err)
	// }
	// defer conn.Close()

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
		// _, err := conn.Write([]byte(data))
		// if err != nil {
		// 	fmt.Println("Error sending log:", err)
		// 	break
		// }

		// fmt.Println("Sent log:", data)

		// Publish log to RabbitMQ
		err = ch.PublishWithContext(ctx,
			"logs_direct", // exchange
			pri,           // routing key
			false,         // mandatory
			false,         // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(data),
			})
		failOnError(err, "Failed to publish a message")

		log.Printf(" [x] Sent %s", data)
		// Wait 1-2 seconds
		time.Sleep(time.Duration(rand.Intn(2)+1) * time.Second)
	}
}
