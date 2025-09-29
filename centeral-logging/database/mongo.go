package database

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/gunjanmistry08/centeralised-logging/centeral-logging/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
	ctx      context.Context
	cancel   context.CancelFunc
}

const MongoURI = "mongodb://localhost:27017"

func (m *MongoDB) Connect(uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return err
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
		return err
	}

	m.client = client
	m.database = client.Database("mydb") // default db, can be parameterized
	m.ctx = context.Background()
	m.cancel = cancel
	log.Println("Connected to MongoDB!")
	return nil
}

func (m *MongoDB) Disconnect() error {
	if m.cancel != nil {
		m.cancel()
	}
	if m.client != nil {
		return m.client.Disconnect(context.Background())
	}
	return nil
}

func (m *MongoDB) InsertLog(collectionName string, document interface{}) (map[string]interface{}, error) {
	collection := m.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, document)
	if err != nil {
		log.Printf("Failed to database insert log: %v", err)
		return nil, err
	}
	return map[string]interface{}{"inserted_id": res.InsertedID}, nil
}

func (m *MongoDB) GetLogs(collectionName string, filter map[string]interface{}) ([]map[string]interface{}, error) {
	collection := m.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(m.ctx, 30*time.Second)
	defer cancel()

	findOptions := options.Find()
	for key, value := range filter {
		if key == "limit" {
			limitStr, ok := value.(string)
			if ok {
				limit, err := strconv.ParseInt(limitStr, 10, 64)
				if err == nil {
					findOptions.SetLimit(limit)
				} else {
					log.Printf("Failed to parse limit: %v", err)
					return nil, err
				}
			}
		}
		if key == "sort" {
			findOptions.SetSort(bson.D{bson.E{Key: value.(string), Value: 1}})
		}
	}

	filterBson := bson.M{}
	log.Printf("Filter received: %v", filter)
	for key, value := range filter {
		if key != "limit" && key != "sort" {
			filterBson[key] = value
		}
		if key == "is.blacklisted" {
			boolVal, ok := value.(string)
			if ok {
				parsedBool, err := strconv.ParseBool(boolVal)
				if err == nil {
					filterBson["blacklisted"] = parsedBool
				} else {
					log.Printf("Failed to parse is.blacklisted: %v", err)
					filterBson["blacklisted"] = false
				}
			}
		}
	}
	log.Printf("MongoDB query filter: %v", filterBson)
	cursor, err := collection.Find(ctx, filterBson, findOptions)
	if err != nil {
		log.Printf("Failed to find logs: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []map[string]interface{}
	for cursor.Next(ctx) {
		var logEntry map[string]interface{}
		if err := cursor.Decode(&logEntry); err != nil {
			log.Printf("Failed to decode log entry: %v", err)
			continue
		}
		logs = append(logs, logEntry)
	}
	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}
	log.Println(len(logs), "logs fetched")
	return logs, nil
}

func (m *MongoDB) GetMetrics(collectionName string) (map[string]interface{}, error) {
	collection := m.database.Collection(collectionName)
	ctx, cancel := context.WithTimeout(m.ctx, 30*time.Second)
	defer cancel()
	totalCount, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("Failed to count total logs: %v", err)
		return nil, err
	}

	pipeline := mongo.Pipeline{
		bson.D{
			bson.E{Key: "$group", Value: bson.D{
				bson.E{Key: "_id", Value: bson.D{
					bson.E{Key: "eventcategory", Value: "$eventcategory"},
					bson.E{Key: "level", Value: "$level"},
				}},
				bson.E{Key: "count", Value: bson.D{
					bson.E{Key: "$sum", Value: 1},
				}},
			}},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var results []model.Result
	if err := cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	return map[string]interface{}{
		"total_count":           totalCount,
		"by_category_and_level": results,
	}, nil
}
