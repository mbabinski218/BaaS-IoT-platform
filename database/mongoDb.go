package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func Connect(uri, dbName, collectionName string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("MongoDB connect error: %v", err)
	}

	collection := client.Database(dbName).Collection(collectionName)
	return &Client{
		client:     client,
		collection: collection,
	}, nil
}

func (c *Client) Add(dataId uuid.UUID, data map[string]any, deviceId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := bson.M{
		"_id":       dataId,
		"device_id": deviceId,
		"data":      data,
	}

	_, err := c.collection.InsertOne(ctx, doc)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("Insert timed out")
		} else if errors.Is(err, context.Canceled) {
			fmt.Println("Insert canceled")
		} else {
			fmt.Printf("Insert error: %v\n", err)
		}
		return err
	}

	return nil
}
