package app

import (
	"context"
	"firebase.google.com/go"
	v "github.com/core-go/core/v10"
	f "github.com/core-go/firestore"
	"github.com/core-go/firestore/query"
	"github.com/core-go/health"
	"github.com/core-go/health/firestore"
	"github.com/core-go/log"
	"github.com/core-go/search"
	"google.golang.org/api/option"
	"reflect"

	"go-service/internal/handler"
	"go-service/internal/model"
	"go-service/internal/service"
)

type ApplicationContext struct {
	Health *health.Handler
	User   handler.UserPort
}

func NewApp(ctx context.Context, cfg Config) (*ApplicationContext, error) {
	opts := option.WithCredentialsJSON([]byte(cfg.Credentials))
	app, err := firebase.NewApp(ctx, nil, opts)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	logError := log.LogError
	validator := v.NewValidator()

	userType := reflect.TypeOf(model.User{})
	userQuery := query.NewBuilder(userType)
	userSearchBuilder := f.NewSearchBuilder(client, "users", userType, userQuery.BuildQuery, search.GetSort, "createTime", "updateTime")
	userRepository := f.NewRepository(client, "users", userType, "CreateTime", "UpdateTime")
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userSearchBuilder.Search, userService, validator.Validate, logError)

	firestoreChecker := firestore.NewHealthChecker(ctx, []byte(cfg.Credentials))
	healthHandler := health.NewHandler(firestoreChecker)

	return &ApplicationContext{
		Health: healthHandler,
		User:   userHandler,
	}, nil
}
