package main

import (
	"context"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"truetech/internal/pkg/client"
	prc "truetech/internal/pkg/processor"
)

const (
	mtsUrl = "https://mts-olimp-cloud.codenrock.com/api"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	initLogger()
	log.Info("start")

	apiClient := client.NewClient(mtsUrl)

	processor := prc.NewProcessor(apiClient)
	processor.Start(ctx)
}

func initLogger() {
	//logFile := "log.txt"
	//f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer f.Close()
	//
	//log.SetOutput(f)

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)
	log.SetLevel(log.DebugLevel)
}
