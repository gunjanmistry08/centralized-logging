package database

import (
	"log"
	"os"
)

type Database interface {
	Connect(uri string) error
	InsertLog(collectionName string, document interface{}) (map[string]interface{}, error)
	GetLogs(collectionName string, filter map[string]interface{}) ([]map[string]interface{}, error)
	GetMetrics(collectionName string) (map[string]interface{}, error)
	Disconnect() error
}

var DB Database

func Init() {
	mongoDB := &MongoDB{}
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = MongoURI
	}
	err := mongoDB.Connect(uri)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	DB = mongoDB
}
