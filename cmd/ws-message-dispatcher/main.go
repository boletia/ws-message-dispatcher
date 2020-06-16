package main

import (
	"os"

	"github.com/boletia/ws-message-dispatcher/config"
	"github.com/boletia/ws-message-dispatcher/pkg/sender"
	"github.com/boletia/ws-message-dispatcher/pkg/service"
	"github.com/boletia/ws-message-dispatcher/pkg/store/dynamodb"
	"github.com/labstack/echo"
	"github.com/sevenNt/echo-pprof"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})

	cnf, err := config.Read()
	if err != nil {
		os.Exit(1)
	}

	srv := service.New(
		dynamodb.New(cnf.Dynamo.Region, cnf.Dynamo.Table),
		sender.New(cnf.Lambda.Region, cnf.Lambda.Function),
	)

	e := echo.New()
	e.POST("/", srv.TakeIn)
	echopprof.Wrap(e)

	e.Logger.Fatal(e.Start(cnf.Service.Host))
}
