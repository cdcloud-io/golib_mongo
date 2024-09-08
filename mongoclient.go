package mongoclient

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client wraps the MongoDB client and provides additional functionality
type Client struct {
	*mongo.Client
}

// ClientOptions represents options for creating a new Client
type ClientOptions struct {
	URI                    string
	ConnectTimeout         time.Duration
	ServerSelectionTimeout time.Duration
}

// QueryParams abstracts the MongoDB query parameters
type QueryParams struct {
	Database   string
	Collection string
	Filter     bson.M
}

// NewClient creates and returns a new Client with the given options
func NewClient(opts ClientOptions) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), opts.ConnectTimeout)
	defer cancel()

	clientOpts := options.Client().ApplyURI(opts.URI).
		SetServerSelectionTimeout(opts.ServerSelectionTimeout)

	mongoClient, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		mongoClient.Disconnect(ctx)
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return &Client{Client: mongoClient}, nil
}

// Close disconnects the client from MongoDB
func (c *Client) Close(ctx context.Context) error {
	return c.Disconnect(ctx)
}

// QueryOne executes a query to find a single document using QueryParams
func (c *Client) QueryOne(ctx context.Context, params QueryParams, result interface{}) error {
	collection := c.Database(params.Database).Collection(params.Collection)
	err := collection.FindOne(ctx, params.Filter).Decode(result)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to execute FindOne query: %w", err)
	}
	return nil
}

// QueryMany executes a query to find multiple documents using QueryParams
func (c *Client) QueryMany(ctx context.Context, params QueryParams) ([]interface{}, error) {
	collection := c.Database(params.Database).Collection(params.Collection)
	cursor, err := collection.Find(ctx, params.Filter)
	if err != nil {
		return nil, fmt.Errorf("failed to execute Find query: %w", err)
	}
	defer cursor.Close(ctx)

	var results []interface{}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode query results: %w", err)
	}
	return results, nil
}

// InsertOne inserts a single document using QueryParams
func (c *Client) InsertOne(ctx context.Context, params QueryParams, document interface{}) (*mongo.InsertOneResult, error) {
	result, err := c.Database(params.Database).Collection(params.Collection).InsertOne(ctx, document)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %w", err)
	}
	return result, nil
}

// UpdateOne updates a single document using QueryParams
func (c *Client) UpdateOne(ctx context.Context, params QueryParams, update interface{}) (*mongo.UpdateResult, error) {
	result, err := c.Database(params.Database).Collection(params.Collection).UpdateOne(ctx, params.Filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}
	return result, nil
}

// DeleteOne deletes a single document using QueryParams
func (c *Client) DeleteOne(ctx context.Context, params QueryParams) (*mongo.DeleteResult, error) {
	result, err := c.Database(params.Database).Collection(params.Collection).DeleteOne(ctx, params.Filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete document: %w", err)
	}
	return result, nil
}

// QueryMongoDB executes a MongoDB query with abstracted parameters
func (c *Client) QueryMongoDB(ctx context.Context, params QueryParams) (map[string]interface{}, error) {
	collection := c.Database(params.Database).Collection(params.Collection)

	// Store the result in a map[string]interface{} since the structure is unknown
	var result map[string]interface{}

	err := collection.FindOne(ctx, params.Filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("no documents found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query MongoDB: %w", err)
	}

	return result, nil
}
