package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	loggerOptions := options.
		Logger().
		SetMaxDocumentLength(25).
		SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		os.Getenv("DATABASE_URL")).SetRetryWrites(true).SetServerMonitor(&event.ServerMonitor{
		ServerHeartbeatSucceeded: func(event *event.ServerHeartbeatSucceededEvent) {
			// fmt.Println(fmt.Sprintf("ConnectionID: %v, Duration: %v", event.ConnectionID, event.Duration))
		},
	}).SetServerSelectionTimeout(5*time.Second).SetLoggerOptions(loggerOptions))

	if err != nil {
		cancel()
		return nil, err
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		cancel()
		return nil, err
	}

	db := client.Database(os.Getenv("MONGO_DB"))
	DB = db

	fmt.Println("Database connected successfully.")

	return cancel, nil
}
