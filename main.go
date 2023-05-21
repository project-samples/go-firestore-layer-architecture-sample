package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/core-go/config"
	"github.com/core-go/log"
	mid "github.com/core-go/log/middleware"
	"github.com/gorilla/mux"

	"go-service/internal/app"
)

func main() {
	var cfg app.Config
	err := config.Load(&cfg, "configs/config")
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	log.Initialize(cfg.Log)
	r.Use(mid.BuildContext)
	logger := mid.NewLogger()
	if log.IsInfoEnable() {
		r.Use(mid.Logger(cfg.MiddleWare, log.InfoFields, logger))
	}
	r.Use(mid.Recover(log.ErrorMsg))

	err = app.Route(context.Background(), r, cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println("Start server")
	server := ""
	if cfg.Server.Port > 0 {
		server = ":" + strconv.FormatInt(cfg.Server.Port, 10)
	}
	if err = http.ListenAndServe(server, r); err != nil {
		fmt.Println(err.Error())
	}
}
