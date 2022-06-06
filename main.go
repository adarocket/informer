package main

import (
	"log"
	"time"

	"github.com/adarocket/informer/config"
	"github.com/adarocket/informer/statistics/cardano"
	"github.com/adarocket/informer/statistics/chia"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	startTime := time.Now()
	loadedConfig, err := config.LoadConfig()
	if err != nil {
		log.Panicln(err)
		return
	}

	conn, err := grpc.Dial(loadedConfig.ControllerURL, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	switch loadedConfig.Blockchain {
	case "cardano":
		cardano.NewCardano(loadedConfig, startTime, conn)
	case "chia":
		chia.NewChia(loadedConfig, startTime, conn)
	default:
		log.Println("Blockchain \"" + loadedConfig.Blockchain + "\" not supported")
		return
	}
}
