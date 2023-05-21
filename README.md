# go-firestore-layer-architecture-sample

## How to run
#### Clone the repository
```shell
git clone https://github.com/project-samples/go-firestore-layer-architecture-sample
cd go-firestore-layer-architecture-sample
```

#### To run the application
```shell
go run main.go
```

## Architecture
### Simple Layer Architecture
![Layer Architecture](https://camo.githubusercontent.com/d9b21eb50ef70dcaebf5a874559608f475e22c799bc66fcf99fb01f08576540f/68747470733a2f2f63646e2d696d616765732d312e6d656469756d2e636f6d2f6d61782f3830302f312a4a4459546c4b3030796730496c556a5a392d737037512e706e67)

### Layer Architecture with full features
![Layer Architecture with standard features: config, health check, logging, middleware log tracing](https://camo.githubusercontent.com/aa7b739a4692eaf2b363cf9caf8b021c60082c77c98d3f8c96665b5cf4640628/68747470733a2f2f63646e2d696d616765732d312e6d656469756d2e636f6d2f6d61782f3830302f312a6d79556b504969343265593477455f494446526176412e706e67)

## API Design
### Common HTTP methods
- GET: retrieve a representation of the resource
- POST: create a new resource
- PUT: update the resource
- PATCH: perform a partial update of a resource
- DELETE: delete a resource

## API design for health check
To check if the service is available.
#### *Request:* GET /health
#### *Response:*
```json
{
    "status": "UP",
    "details": {
        "firestore": {
            "status": "UP"
        }
    }
}
```

## API design for users
#### *Resource:* users

### Get all users
#### *Request:* GET /users
#### *Response:*
```json
[
    {
        "id": "spiderman",
        "username": "peter.parker",
        "email": "peter.parker@gmail.com",
        "phone": "0987654321",
        "dateOfBirth": "1962-08-25T16:59:59.999Z"
    },
    {
        "id": "wolverine",
        "username": "james.howlett",
        "email": "james.howlett@gmail.com",
        "phone": "0987654321",
        "dateOfBirth": "1974-11-16T16:59:59.999Z"
    }
]
```

### Get one user by id
#### *Request:* GET /users/:id
```shell
GET /users/wolverine
```
#### *Response:*
```json
{
    "id": "wolverine",
    "username": "james.howlett",
    "email": "james.howlett@gmail.com",
    "phone": "0987654321",
    "dateOfBirth": "1974-11-16T16:59:59.999Z"
}
```

### Create a new user
#### *Request:* POST /users 
```json
{
    "id": "wolverine",
    "username": "james.howlett",
    "email": "james.howlett@gmail.com",
    "phone": "0987654321",
    "dateOfBirth": "1974-11-16T16:59:59.999Z"
}
```
#### *Response:* 1: success, 0: duplicate key, -1: error
```json
1
```

### Update one user by id
#### *Request:* PUT /users/:id
```shell
PUT /users/wolverine
```
```json
{
    "username": "james.howlett",
    "email": "james.howlett@gmail.com",
    "phone": "0987654321",
    "dateOfBirth": "1974-11-16T16:59:59.999Z"
}
```
#### *Response:* 1: success, 0: not found, -1: error
```json
1
```

### Patch one user by id
Perform a partial update of user. For example, if you want to update 2 fields: email and phone, you can send the request body of below.
#### *Request:* PATCH /users/:id
```shell
PATCH /users/wolverine
```
```json
{
    "email": "james.howlett@gmail.com",
    "phone": "0987654321"
}
```
#### *Response:* 1: success, 0: not found, -1: error
```json
1
```

#### Problems for patch
If we pass a struct as a parameter, we cannot control what fields we need to update. So, we must pass a map as a parameter.
```go
type UserService interface {
	Update(ctx context.Context, user *User) (int64, error)
	Patch(ctx context.Context, user map[string]interface{}) (int64, error)
}
```
We must solve 2 problems:
1. At http handler layer, we must convert the user struct to map, with json format, and make sure the nested data types are passed correctly.
2. At service layer or repository layer, from json format, we must convert the json format to database format (in this case, we must convert to bson of Mongo)

#### Solutions for patch  
1. At http handler layer, we use [core-go/core](https://github.com/core-go/core), to convert the user struct to map, to make sure we just update the fields we need to update
```go
import server "github.com/core-go/core"

func (h *UserHandler) Patch(w http.ResponseWriter, r *http.Request) {
    var user User
    userType := reflect.TypeOf(user)
    _, jsonMap := sv.BuildMapField(userType)
    body, _ := sv.BuildMapAndStruct(r, &user)
    json, er1 := sv.BodyToJson(r, user, body, ids, jsonMap, nil)

    result, er2 := h.service.Patch(r.Context(), json)
    if er2 != nil {
        http.Error(w, er2.Error(), http.StatusInternalServerError)
        return
    }
    respond(w, result)
}
```

### Delete a new user by id
#### *Request:* DELETE /users/:id
```shell
DELETE /users/wolverine
```
#### *Response:* 1: success, 0: not found, -1: error
```json
1
```

## Common libraries
- [core-go/health](https://github.com/core-go/health): include Health Handler, Firestore Health Checker
- [core-go/config](https://github.com/core-go/config): to load the config file, and merge with other environments (SIT, UAT, ENV)
- [core-go/log](https://github.com/core-go/log): log and log middleware

### core-go/health
To check if the service is available, refer to [core-go/health](https://github.com/core-go/health)
#### *Request:* GET /health
#### *Response:*
```json
{
    "status": "UP",
    "details": {
        "firestore": {
            "status": "UP"
        }
    }
}
```
To create health checker, and health handler
```go
	opts := option.WithCredentialsJSON([]byte(root.Credentials))
	app, er1 := firebase.NewApp(ctx, nil, opts)
	if er1 != nil {
		return nil, er1
	}

	client, er2 := app.Firestore(ctx)
	if er2 != nil {
		return nil, er2
	}

	userService := services.NewUserService(client)
	userHandler := handlers.NewUserHandler(userService)

	firestoreChecker := firestore.NewHealthChecker(ctx, []byte(root.Credentials))
	healthHandler := health.NewHealthHandler(firestoreChecker)
```

To handler routing
```go
	r := mux.NewRouter()
	r.HandleFunc("/health", healthHandler.Check).Methods("GET")
```

### core-go/config
To load the config from "config.yml", in "configs" folder
```go
package app

import (
	"github.com/core-go/log"
	mid "github.com/core-go/log/middleware"
)

type Root struct {
	Server      ServerConfig  `mapstructure:"server"`
	Log         log.Config    `mapstructure:"log"`
	MiddleWare  mid.LogConfig `mapstructure:"middleware"`
	Credentials string        `mapstructure:"credentials"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port int64  `mapstructure:"port"`
}
```

### core-go/log *&* core-go/log/middleware
```go
import (
	"github.com/core-go/config"
	"github.com/core-go/log"
	mid "github.com/core-go/log/middleware"
	"github.com/gorilla/mux"
)

func main() {
	var conf app.Root
	config.Load(&conf, "configs/config")

	r := mux.NewRouter()

	log.Initialize(conf.Log)
	r.Use(mid.BuildContext)
	logger := mid.NewStructuredLogger()
	r.Use(mid.Logger(conf.MiddleWare, log.InfoFields, logger))
	r.Use(mid.Recover(log.ErrorMsg))
}
```
To configure to ignore the health check, use "skips":
```yaml
middleware:
  skips: /health
```