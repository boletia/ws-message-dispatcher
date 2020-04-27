package main

import (
	"os"

	"github.com/boletia/container-events/config"
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

	cnf, err := config.ReadConfig()
	if err != nil {
		os.Exit(1)
	}

	srv := service.New(
		dynamodb.New(cnf.Dynamo.Region, cnf.Dynamo.Table),
		sender.New(cnf.Lambda.Region, cnf.Lambda.Function),
	)

	e := echo.New()
	e.POST("/", srv.TakeIn)

	e.Logger.Fatal(e.Start(cnf.Service.Host))
}
