package main

import (
	"context"
	"log"
	"time"

	pb "github.com/adarocket/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var startTime time.Time

func main() {
	startTime = time.Now()
	loadConfig()

	conn, err := grpc.Dial(loadedConfig.ControllerURL, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewInformerClient(conn)

	durationForFrequentlyUpdate := time.Second * time.Duration(loadedConfig.TimeForFrequentlyUpdate)
	durationForRareUpdate := time.Second * time.Duration(loadedConfig.TimeForRareUpdate)

	sendStatistic(client, true)
	for {
		timer := time.NewTimer(durationForFrequentlyUpdate)
		<-timer.C
		durationForRareUpdate -= durationForFrequentlyUpdate
		if durationForRareUpdate <= 0 {
			sendStatistic(client, true)
		} else {
			sendStatistic(client, false)
		}
	}
}

func sendStatistic(client pb.InformerClient, fullStatistics bool) {
	request, err := getNodeStatistic(fullStatistics)
	if err != nil {
		grpclog.Println(err)
		return
	}

	response, err := client.SaveStatistic(context.Background(), request)
	if err != nil {
		grpclog.Println(err)
		return
	}

	log.Println(response.Status)
}
