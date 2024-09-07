# golib_mongo_c

## usage

```go
Example usage:
func main() {
    opts := mongoclient.ClientOptions{
        URI:               "mongodb://localhost:27017",
        ConnectTimeout:    10 * time.Second,
        ServerSelectionTimeout: 5 * time.Second,
    }
    client, err := mongoclient.NewClient(opts)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close(context.Background())

    var user User
    query := mongoclient.Query{
        Collection: "users",
        Filter:     bson.M{"iboNumber": 123456},
        Result:     &user,
    }
    if err := client.FindOne(context.Background(), "mydb", query); err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found user: %+v\n", user)
}
```

