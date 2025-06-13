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
	"github.com/mbabinski218/BaaS-IoT-platform/utils"
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

func (c *Client) Add(dataId uuid.UUID, data map[string]any, deviceId uuid.UUID) (uuid.UUID, time.Duration, error) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := bson.M{
		"_id":       utils.ToBinaryUUID(dataId),
		"device_id": utils.ToBinaryUUID(deviceId),
		"data":      data,
	}

	res, err := c.collection.InsertOne(ctx, doc)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Insert timed out")
		} else if errors.Is(err, context.Canceled) {
			log.Println("Insert canceled")
		} else {
			fmt.Printf("Insert error: %v\n", err)
		}
		return uuid.Nil, 0, err
	}

	createdId, err := utils.ToUUID(res.InsertedID)
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("failed to convert inserted ID to UUID: %v", err)
	}

	duration := time.Since(start)
	return createdId, duration, nil
}

func (c *Client) Get(dataId uuid.UUID) (map[string]any, time.Duration, error) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := utils.ToBinaryUUID(dataId)

	var result struct {
		Data map[string]any `bson:"data"`
	}

	err := c.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, 0, fmt.Errorf("data with ID %s not found", dataId)
		}
		return nil, 0, fmt.Errorf("failed to get data: %v", err)
	}

	duration := time.Since(start)

	return result.Data, duration, nil
}

func (c *Client) GetAuditData(n int64) ([]types.DocData, error) {
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

	var results []types.DocData
	for _, r := range rawData {
		uuid, err := uuid.FromBytes(r["_id"].(primitive.Binary).Data)
		if err != nil {
			return nil, fmt.Errorf("failed to convert record Id to UUID: %v", err)
		}

		results = append(results, types.DocData{
			Id:   uuid,
			Data: r["data"].(map[string]any),
		})
	}

	return results, nil
}

func (c *Client) GetFromTo(from, to time.Time) ([]types.DocData, time.Duration, error) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$addFields", Value: bson.M{
			"timestampDate": bson.M{"$toDate": "$data.timestamp"},
		}}},
		{{Key: "$match", Value: bson.M{
			"timestampDate": bson.M{
				"$gte": primitive.NewDateTimeFromTime(from),
				"$lte": primitive.NewDateTimeFromTime(to),
			},
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":  1,
			"data": 1,
		}}},
	}

	cursor, err := c.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get data: %v", err)
	}
	defer cursor.Close(ctx)

	var rawData []map[string]any
	if err = cursor.All(ctx, &rawData); err != nil {
		return nil, 0, err
	}

	var results []types.DocData
	for _, r := range rawData {
		uuid, err := uuid.FromBytes(r["_id"].(primitive.Binary).Data)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to convert record Id to UUID: %v", err)
		}

		results = append(results, types.DocData{
			Id:   uuid,
			Data: r["data"].(map[string]any),
		})
	}

	duration := time.Since(start)

	return results, duration, nil
}
