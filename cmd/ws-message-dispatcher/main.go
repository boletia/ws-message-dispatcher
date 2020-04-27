package main

import (
	"github.com/boletia/container-events/pkg/sender"
	"github.com/boletia/container-events/pkg/service"
	"github.com/boletia/container-events/pkg/store/dynamodb"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})

	srv := service.New(
		dynamodb.New(),
		sender.New(),
	)

	e := echo.New()
	e.POST("/", srv.TakeIn)

	e.Logger.Fatal(e.Start(":8080"))
}
