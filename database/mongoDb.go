package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/mbabinski218/BaaS-IoT-platform/configs"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.MongoContextTimeout)*2*time.Second)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.MongoContextTimeout)*time.Second)
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

func (c *Client) Add(document bson.M) (uuid.UUID, time.Duration, error) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.MongoContextTimeout)*time.Second)
	defer cancel()

	res, err := c.collection.InsertOne(ctx, document)
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

func (c *Client) Delete(dataId uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.MongoContextTimeout)*time.Second)
	defer cancel()

	id := utils.ToBinaryUUID(dataId)
	_, err := c.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("data with id: %s not found", dataId)
		} else if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("delete operation timed out for id: %s", dataId)
		} else if errors.Is(err, context.Canceled) {
			return fmt.Errorf("delete operation canceled for id: %s", dataId)
		}
		return fmt.Errorf("failed to delete data: %v", err)
	}

	return nil
}

func (c *Client) Get(dataId uuid.UUID) (map[string]any, [][32]byte, time.Duration, error) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.MongoContextTimeout)*time.Second)
	defer cancel()

	id := utils.ToBinaryUUID(dataId)

	var result struct {
		Data  map[string]any `bson:"data"`
		Proof [][32]byte     `bson:"proof"`
	}

	err := c.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil, 0, fmt.Errorf("data with ID %s not found", dataId)
		}
		return nil, nil, 0, fmt.Errorf("failed to get data: %v", err)
	}

	duration := time.Since(start)

	return result.Data, result.Proof, duration, nil
}

func (c *Client) GetAuditData(n int64) ([]types.DocData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.MongoContextTimeout)*time.Second)
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$sample", Value: bson.D{{Key: "size", Value: n}}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "data", Value: 1},
			{Key: "proof", Value: 1},
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.MongoContextTimeout)*time.Second)
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
			"_id":   1,
			"data":  1,
			"proof": 1,
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

func (c *Client) UpdateProof(dataId uuid.UUID, proof [][]byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.MongoContextTimeout)*time.Second)
	defer cancel()

	id := utils.ToBinaryUUID(dataId)

	update := bson.M{
		"$set": bson.M{"proof": proof},
	}

	_, err := c.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("data with id: %s not found", dataId)
		} else if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("update operation timed out for id: %s", dataId)
		} else if errors.Is(err, context.Canceled) {
			return fmt.Errorf("update operation canceled for id: %s", dataId)
		}

		return fmt.Errorf("failed to update proof: %v", err)
	}

	return nil
}

func (c *Client) GetFirstDocumentId() (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.MongoContextTimeout)*time.Second)
	defer cancel()

	count, err := c.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to count documents: %v", err)
	}

	if count == 0 {
		return uuid.Nil, fmt.Errorf("no documents found")
	}

	opts := options.FindOne().SetSkip(10)

	var result map[string]any
	err = c.collection.FindOne(ctx, bson.M{}, opts).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return uuid.Nil, fmt.Errorf("no documents found at center")
		}
		return uuid.Nil, fmt.Errorf("failed to get center document: %v", err)
	}

	id, err := utils.ToUUID(result["_id"])
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to convert first document ID to UUID: %v", err)
	}

	return id, nil
}

func (c *Client) GetCenterDocumentId() (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.MongoContextTimeout)*time.Second)
	defer cancel()

	count, err := c.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to count documents: %v", err)
	}

	if count == 0 {
		return uuid.Nil, fmt.Errorf("no documents found")
	}

	skip := count / 2
	opts := options.FindOne().SetSkip(skip)

	var result map[string]any
	err = c.collection.FindOne(ctx, bson.M{}, opts).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return uuid.Nil, fmt.Errorf("no documents found at center")
		}
		return uuid.Nil, fmt.Errorf("failed to get center document: %v", err)
	}

	id, err := utils.ToUUID(result["_id"])
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to convert first document ID to UUID: %v", err)
	}

	return id, nil
}

func (c *Client) GetLastDocumentId() (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.Envs.MongoContextTimeout)*time.Second)
	defer cancel()

	count, err := c.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to count documents: %v", err)
	}

	if count == 0 {
		return uuid.Nil, fmt.Errorf("no documents found")
	}

	skip := count - 10
	opts := options.FindOne().SetSkip(skip)

	var result map[string]any
	err = c.collection.FindOne(ctx, bson.M{}, opts).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return uuid.Nil, fmt.Errorf("no documents found at center")
		}
		return uuid.Nil, fmt.Errorf("failed to get center document: %v", err)
	}

	id, err := utils.ToUUID(result["_id"])
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to convert first document ID to UUID: %v", err)
	}

	return id, nil
}
