package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/mbabinski218/BaaS-IoT-platform/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping: %v", err)
	}

	if err := ensureCollectionExists(client.Database(dbName), collectionName); err != nil {
		return nil, fmt.Errorf("ensure collection exists error: %v", err)
	}

	collection := client.Database(dbName).Collection(collectionName)

	log.Println("MongoDB connected successfully")
	return &Client{
		client:     client,
		collection: collection,
	}, nil
}

func ensureCollectionExists(db *mongo.Database, collectionName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	if slices.Contains(collections, collectionName) {
		return nil
	}

	if err = db.CreateCollection(ctx, collectionName); err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	return nil
}

func (c *Client) Add(dataId uuid.UUID, data map[string]any, deviceId uuid.UUID) (time.Duration, error) {
	start := time.Now()

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
			log.Println("Insert timed out")
		} else if errors.Is(err, context.Canceled) {
			log.Println("Insert canceled")
		} else {
			fmt.Printf("Insert error: %v\n", err)
		}
		return 0, err
	}

	duration := time.Since(start)
	return duration, nil
}

func (c *Client) Get(dataId uuid.UUID) (map[string]any, time.Duration, error) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result struct {
		Data map[string]any `bson:"data"`
	}

	err := c.collection.FindOne(ctx, bson.M{"_id": dataId}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, 0, fmt.Errorf("data with ID %s not found", dataId)
		}
		return nil, 0, fmt.Errorf("failed to get data: %v", err)
	}

	duration := time.Since(start)

	return result.Data, duration, nil
}

func (c *Client) GetAuditData(n int64) ([]types.AuditData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$sample", Value: bson.D{{Key: "size", Value: n}}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "data", Value: 1},
		}}},
	}

	cursor, err := c.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rawData []map[string]any
	if err = cursor.All(ctx, &rawData); err != nil {
		return nil, err
	}

	var results []types.AuditData
	for _, r := range rawData {
		uuid, err := uuid.FromBytes(r["_id"].(primitive.Binary).Data)
		if err != nil {
			return nil, fmt.Errorf("failed to convert record Id to UUID: %v", err)
		}

		results = append(results, types.AuditData{
			Id:   uuid,
			Data: r["data"].(map[string]any),
		})
	}

	return results, nil
}
