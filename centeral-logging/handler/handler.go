package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gunjanmistry08/centeralised-logging/centeral-logging/database"
	"github.com/gunjanmistry08/centeralised-logging/centeral-logging/model"
)

func Ingest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {

		log.Printf("[ERROR] %v", err)
		http.Error(w, `{"error": "Failed to read request body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Printf("Received log: %s\n", string(body))

	var logE model.LogMessage
	err = json.Unmarshal(body, &logE)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		http.Error(w, `{"error": "Invalid log format"}`, http.StatusBadRequest)
		return
	}
	// write in database
	_, err = database.DB.InsertLog("logs", logE)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		http.Error(w, `{"error": "Failed to ingest log"}`, http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, `{"status": "log ingested"}`)
}

func GetLogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	query := r.URL.Query()
	filter := make(map[string]interface{})
	for key, value := range query {
		filter[key] = value[0]
	}
	logs, err := database.DB.GetLogs("logs", filter)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		http.Error(w, `{"error": "Failed to fetch logs"}`, http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(logs)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		http.Error(w, `{"error": "Failed to marshal logs"}`, http.StatusInternalServerError)
		return
	}
	w.Write(response)

}

func GetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	metrics, err := database.DB.GetMetrics("logs")
	if err != nil {
		log.Printf("[ERROR] %v", err)
		http.Error(w, `{"error": "Failed to fetch metrics"}`, http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(metrics)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		http.Error(w, `{"error": "Failed to marshal metrics"}`, http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
