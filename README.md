# mongoclient Library

A simple and flexible MongoDB client library written in Go. This library abstracts MongoDB operations such as querying, inserting, updating, and deleting documents. It is designed to support **Hexagonal Architecture (Ports and Adapters)** by providing a clean separation between the application's core logic and the database layer.

## Features

- MongoDB connection management
- Query single and multiple documents
- Insert, update, and delete documents
- Abstracted query parameters for flexibility
- Facilitates **Hexagonal Architecture**

## Installation

To install the MongoDB client library, use the following command:

```sh
go get github.com/cdcloud-io/golib_mongo_c
```

## Usage

### 1. Connecting to MongoDB

To create a new MongoDB client, initialize it with connection options:

```go
package main

import (
    "context"
    "log"
    "time"
    "github.com/yourusername/mongoclient"
)

func main() {
    // MongoDB client options
    clientOptions := mongoclient.ClientOptions{
        URI:                    "mongodb://localhost:27017",
        ConnectTimeout:         10 * time.Second,
        ServerSelectionTimeout: 5 * time.Second,
    }

    // Create the MongoDB client
    client, err := mongoclient.NewClient(clientOptions)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    defer client.Close(context.Background())
}
```

### 2. Querying Documents

You can query MongoDB for single or multiple documents using the `QueryParams` struct to abstract the parameters.

#### Query One Document

```go
var result map[string]interface{}

params := mongoclient.QueryParams{
    Database:   "mydb",
    Collection: "users",
    Filter:     bson.M{"username": "johndoe"},
}

err := client.QueryOne(context.Background(), params, &result)
if err != nil {
    log.Fatalf("QueryOne failed: %v", err)
}

fmt.Printf("User: %+v\n", result)
```

#### Query Multiple Documents

```go
params := mongoclient.QueryParams{
    Database:   "mydb",
    Collection: "users",
    Filter:     bson.M{"age": bson.M{"$gt": 25}},
}

results, err := client.QueryMany(context.Background(), params)
if err != nil {
    log.Fatalf("QueryMany failed: %v", err)
}

fmt.Printf("Users: %+v\n", results)
```

### 3. Inserting Documents

You can insert a document into MongoDB using the `InsertOne` method:

```go
doc := bson.M{"username": "newuser", "age": 30}

params := mongoclient.QueryParams{
    Database:   "mydb",
    Collection: "users",
}

insertResult, err := client.InsertOne(context.Background(), params, doc)
if err != nil {
    log.Fatalf("InsertOne failed: %v", err)
}

fmt.Printf("Inserted ID: %v\n", insertResult.InsertedID)
```

### 4. Updating Documents

To update an existing document, use the `UpdateOne` method:

```go
update := bson.M{"$set": bson.M{"age": 31}}

params := mongoclient.QueryParams{
    Database:   "mydb",
    Collection: "users",
    Filter:     bson.M{"username": "newuser"},
}

updateResult, err := client.UpdateOne(context.Background(), params, update)
if err != nil {
    log.Fatalf("UpdateOne failed: %v", err)
}

fmt.Printf("Matched %v document(s) and updated %v document(s)\n", updateResult.MatchedCount, updateResult.ModifiedCount)
```

### 5. Deleting Documents

To delete a document, use the `DeleteOne` method:

```go
params := mongoclient.QueryParams{
    Database:   "mydb",
    Collection: "users",
    Filter:     bson.M{"username": "newuser"},
}

deleteResult, err := client.DeleteOne(context.Background(), params)
if err != nil {
    log.Fatalf("DeleteOne failed: %v", err)
}

fmt.Printf("Deleted %v document(s)\n", deleteResult.DeletedCount)
```

## Hexagonal Architecture

This library is designed to support **Hexagonal Architecture (Ports and Adapters Architecture)** by abstracting the MongoDB interaction behind interfaces. The core application logic communicates with the MongoDB adapter through **ports** like the `QueryParams` struct, ensuring a clean separation between business logic and infrastructure.

- **Ports**: The application core uses `QueryParams` and other abstracted data types to communicate with MongoDB.
- **Adapters**: The `Client` is an adapter that handles MongoDB-specific operations.

### Key Sections

- **Installation**: Provides instructions to install the library.
- **Usage**: Includes code examples for connecting to MongoDB, querying documents, inserting, updating, and deleting documents.
- **Hexagonal Architecture**: Describes how the library supports Hexagonal Architecture by abstracting MongoDB operations behind ports and adapters.
