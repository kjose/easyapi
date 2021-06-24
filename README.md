# Easy Api

Build on top of gin and gorm, Easyapi is a an extension of services to manage easily automatic logic in the creation of resources.

## Documentation

### Built on top of

- Gin router (https://github.com/gin-gonic/gin/edit/master/README.md)
- ORM Gorm (https://gorm.io/index.html)
- ODM Go bongo (https://github.com/go-bongo/bongo)

### Summary

* [Installation](#installation)
* [Getting started](#getting-started)
* [Extending easyapi](#extending-go-api)
* [Create a CRUD route](#create-a-crud-route)

## Installation

```
# 
```


## Documentation

### Router

Go api has by default 1 simple orm manager and 1 simple odm manager.

In each `CRUDL` method, you can specify the basic method you want : 
- C (Create) : Will create a `POST /{resource}/{id}` route
- R (Read) : Will create a `GET /{resource}/{id}` route
- U (Update) : Will create a `PATCH /{resource}/{id}` route
- D (Delete) : Will create a `DELETE /{resource}/{id}` route
- L (List) : Will create a `GET /{resource}` route and return a collection of resources

To enable all route just pass `""` as the fourth argument. 
To enable only some methods you can pass a parameter like `CR` to enable only Create and Read routes.

```go
...

func main() {
    r := gin.Default()

    // Init ORM
    initORM()
    initODM()

    // Resources CRUD routes
    easyapi.CRUDL(r, "/users", new(model.User), "CRUL")
    easyapi.CRUDL(r, "/banks", new(model.Bank), "")
}

func initORM() {
    err := orm.Init(mysql.Open(os.Getenv("DSN")), true, &logger.Config{
        SlowThreshold:             time.Second,
        LogLevel:                  logger.Info,
        IgnoreRecordNotFoundError: false,
        Colorful:                  true,
    })
    if err != nil {
        log.Fatal("Error init database")
    }
}

func initODM() {
    err := odm.Init(odm.Config{
        ConnectionString: os.Getenv("MONGO_SERVER"),
        Database:         os.Getenv("MONGO_DATABASE"),
    }, false)
    if err != nil {
        log.Fatal("Error init mongodb")
    }
}
```

### Security & Access management

```
# 
```

### Resource validation

```
# 
```

### Resource serializer

```
# 
```

### Filtering & Pagination

```
# 
```

### Event manager

```
# 
```




