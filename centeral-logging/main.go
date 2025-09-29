package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gunjanmistry08/centeralised-logging/centeral-logging/database"
	"github.com/gunjanmistry08/centeralised-logging/centeral-logging/handler"
)

func main() {

	database.Init()
	defer database.DB.Disconnect()

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"message": "pong"}`)
	})

	http.HandleFunc("/ingest", handler.Ingest)
	http.HandleFunc("/logs", handler.GetLogs)
	http.HandleFunc("/metrics", handler.GetMetrics)
	log.Println("HTTP Server listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
