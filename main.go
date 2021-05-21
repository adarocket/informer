package main

import (
	"adarocket/informer/config"
	"adarocket/informer/statistics/cardano"
	"adarocket/informer/statistics/chia"
	"context"
	"log"
	"time"

	pb "github.com/adarocket/proto/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	var node NodeInterface

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

	client := pb.NewInformerClient(conn)

	durationForFrequentlyUpdate := time.Second * time.Duration(loadedConfig.TimeForFrequentlyUpdate)
	durationForRareUpdate := time.Second * time.Duration(loadedConfig.TimeForRareUpdate)

	switch loadedConfig.Blockchain {
	case "cardano":
		node = cardano.NewCardano(loadedConfig, startTime)
	case "chia":
		node = chia.NewChia(loadedConfig, startTime)
	default:
		log.Println("Blockchain \"" + loadedConfig.Blockchain + "\" not supported")
		return
	}

	sendStatistic(client, true, node)
	for {
		timer := time.NewTimer(durationForFrequentlyUpdate)
		<-timer.C
		durationForRareUpdate -= durationForFrequentlyUpdate
		if durationForRareUpdate <= 0 {
			sendStatistic(client, true, node)
		} else {
			sendStatistic(client, false, node)
		}
	}
}

func sendStatistic(client pb.InformerClient, fullStatistics bool, node NodeInterface) {
	request, err := node.GetNodeStatistic(fullStatistics)
	if err != nil {
		grpclog.Infoln(err)
		return
	}

	response, err := client.SaveStatistic(context.Background(), request)
	if err != nil {
		grpclog.Infoln(err)
		return
	}

	log.Println(response.Status)
}

type NodeInterface interface {
	GetNodeStatistic(fullStatistics bool) (*pb.SaveStatisticRequest, error)
}
