package main

import (
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// func blacklistedUsername(username string) bool {
// 	blacklistedUsernames := []string{"Motabhai"}
// 	return slices.Contains(blacklistedUsernames, username)
// }

// func getSeverity(pri int) string {
// 	severityMap := map[int]string{
// 		0: "Emergency",
// 		1: "Alert",
// 		2: "Critical",
// 		3: "Error",
// 		4: "Warning",
// 		5: "Notice",
// 		6: "Informational",
// 		7: "Debug",
// 	}
// 	return severityMap[pri%8]
// }

// func getCategory(msg string) string {
// 	lower := strings.ToLower(msg)
// 	switch {
// 	case strings.Contains(lower, "logged on"):
// 		return "login.audit"
// 	case strings.Contains(lower, "logoff") || strings.Contains(lower, "logged off"):
// 		return "logout.audit"
// 	case strings.Contains(lower, "failed logon") || strings.Contains(lower, "failure"):
// 		return "login.failure"
// 	default:
// 		return "other.audit"
// 	}
// }

type LogMessage struct {
	Timestamp     time.Time `json:"timestamp"`
	Hostname      string    `json:"hostname"`
	EventSource   string    `json:"event.source.type"`
	EventCategory string    `json:"event.category"`
	Message       string    `json:"message"`
	Blacklisted   bool      `json:"is.blacklisted"`
	Level         string    `json:"severity"`
	Username      string    `json:"username"`
}

func main() {

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

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Listen on TCP port 9000
	// ln, err := net.Listen("tcp", ":9000")
	// if err != nil {
	// 	panic(err)
	// }
	// defer ln.Close()

	// fmt.Println("TCP Server listening on port 9000...")

	// for {
	// 	conn, err := ln.Accept()
	// 	if err != nil {
	// 		fmt.Println("Error accepting connection:", err)
	// 		continue
	// 	}
	// 	fmt.Println("Client connected:", conn.RemoteAddr())

	// 	// Handle each client in a goroutine
	// 	go handleConnection(conn)
	// }
	// Get routing keys from command-line arguments (after the program name)
	routingKeys := os.Args[1:]
	if len(routingKeys) == 0 {
		routingKeys = []string{"<134>", "<132>", "<131>"} // default if none provided
	}
	for _, s := range routingKeys {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name, "logs_direct", s)
		err = ch.QueueBind(
			q.Name,        // queue name
			s,             // routing key
			"logs_direct", // exchange
			false,
			nil)
		failOnError(err, "Failed to bind a queue")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever

}

// func handleConnection(conn net.Conn) {
// 	defer conn.Close()

// 	scanner := bufio.NewScanner(conn)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		fmt.Println("Received log:", line)

// 		re := regexp.MustCompile(`^<(\d+)>\s+(WIN)-([^\s]+).*?Account Name:\s*([^\r\n]+)`)
// 		matches := re.FindStringSubmatch(line)
// 		var msg LogMessage

// 		fmt.Println("pri:", matches[1], "sourceType:", matches[2], "hostname:", matches[3], "username:", strings.TrimSpace(matches[4]))

// 		if len(matches) >= 5 {
// 			pri := matches[1]
// 			sourceType := matches[2] // WIN
// 			hostname := matches[3]   // EQ5V3RA5F7H
// 			username := strings.TrimSpace(matches[4])

// 			// Convert PRI to severity
// 			var priInt int
// 			fmt.Sscanf(pri, "%d", &priInt)
// 			severity := getSeverity(priInt)

// 			// Detect category from log message
// 			category := getCategory(line)

// 			msg = LogMessage{
// 				Timestamp:     time.Now(),
// 				Hostname:      hostname,
// 				EventCategory: category,
// 				EventSource:   sourceType,
// 				Username:      username,
// 				Level:         severity,
// 				Blacklisted:   blacklistedUsername(username),
// 				Message:       line,
// 			}
// 		}

// 		fmt.Println(msg)
// 		// validate log fields
// 		requiredFields := []string{"timestamp", "event.source.type", "event.category", "hostname", "message", "username"}
// 		missingFields := []string{}

// 		for _, field := range requiredFields {
// 			var fieldValue interface{}
// 			switch field {
// 			case "timestamp":
// 				fieldValue = msg.Timestamp
// 			case "event.source.type":
// 				fieldValue = msg.EventSource
// 			case "event.category":
// 				fieldValue = msg.EventCategory
// 			case "hostname":
// 				fieldValue = msg.Hostname
// 			case "message":
// 				fieldValue = msg.Message
// 			case "username":
// 				fieldValue = msg.Username
// 			}
// 			// Check for zero value (empty string or zero time)
// 			if fieldValue == "" || (field == "timestamp" && msg.Timestamp.IsZero()) {
// 				missingFields = append(missingFields, field)
// 			}

// 			if len(missingFields) > 0 {
// 				fmt.Printf("Missing required fields: %v\n", missingFields)
// 				continue
// 			}

// 			// forward log to central logging system at /ingest
// 			fmt.Printf("Forwarding log: %+v\n", msg)
// 			data, _ := json.Marshal(msg)

// 			url := os.Getenv("centeral_logging_url")
// 			if url == "" {
// 				url = "http://localhost:8080/ingest"
// 			}

// 			resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
// 			if err != nil {
// 				fmt.Println("Error forwarding log:", err)
// 				continue
// 			}

// 			if resp.StatusCode != http.StatusOK {
// 				fmt.Println("Failed to forward log, status code:", resp.StatusCode)
// 			}
// 		}
// 	}

// 	if err := scanner.Err(); err != nil {
// 		fmt.Println("Error reading:", err)
// 	}
// }
