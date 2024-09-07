package mongoclient

import (
	"context"
	"fmt"
	"time"

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

// Query represents a MongoDB query
type Query struct {
	Collection string
	Filter     interface{}
	Result     interface{}
}

// ## Public functions ##

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

// ## Public methods ##

// Close disconnects the client from MongoDB
func (c *Client) Close(ctx context.Context) error {
	return c.Disconnect(ctx)
}

func (c *Client) QueryIBOtoStruct(ctx context.Context, )

// FindOne executes a query to find a single document
func (c *Client) QueryOne(ctx context.Context, db string, query Query) error {
	collection := c.Database(db).Collection(query.Collection)
	err := collection.FindOne(ctx, query.Filter).Decode(query.Result)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to execute FindOne query: %w", err)
	}
	return nil
}

// FindMany executes a query to find multiple documents
func (c *Client) QueryMany(ctx context.Context, db string, query Query) ([]interface{}, error) {
	collection := c.Database(db).Collection(query.Collection)
	cursor, err := collection.Find(ctx, query.Filter)
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

// InsertOne inserts a single document into the specified collection
func (c *Client) InsertOne(ctx context.Context, db, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	result, err := c.Database(db).Collection(collection).InsertOne(ctx, document)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %w", err)
	}
	return result, nil
}

// UpdateOne updates a single document in the specified collection
func (c *Client) UpdateOne(ctx context.Context, db, collection string, filter, update interface{}) (*mongo.UpdateResult, error) {
	result, err := c.Database(db).Collection(collection).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}
	return result, nil
}

// DeleteOne deletes a single document from the specified collection
func (c *Client) DeleteOne(ctx context.Context, db, collection string, filter interface{}) (*mongo.DeleteResult, error) {
	result, err := c.Database(db).Collection(collection).DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete document: %w", err)
	}
	return result, nil
}


// ## LEGACY ##
/*

package mongoclient

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateClient initializes and returns a MongoDB client
func CreateClient(uri string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(uri)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		client.Disconnect(context.TODO())
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client, nil
}

// QueryMongoDB queries the MongoDB collection for a user with the specified IBO number
func QueryMongoDB(db *mongo.Database, collectionName string, iboNumber uint64) (*data.MongoUser, error) {
	collection := db.Collection(collectionName)
	filter := bson.M{"iboNumber": iboNumber}
	var user data.MongoUser
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query MongoDB: %w", err)
	}
	return &user, nil
}

*/
