package chia

import (
	"adarocket/informer/config"
	"adarocket/informer/statistics/common"
	"context"
	"log"
	"time"

	pbChia "github.com/adarocket/proto/proto-gen/chia"
	pbCommon "github.com/adarocket/proto/proto-gen/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type Chia struct {
	loadedConfig    *config.Config
	commonStatistic *common.CommonStatistic
}

func NewChia(config *config.Config, startTime time.Time, conn *grpc.ClientConn) {
	chia := new(Chia)
	chia.loadedConfig = config
	chia.commonStatistic = common.NewCommonStatistic(config, startTime)

	// ----------------------------------------------------------------------

	client := pbChia.NewChiaClient(conn)

	durationForFrequentlyUpdate := time.Second * time.Duration(chia.loadedConfig.TimeForFrequentlyUpdate)
	durationForRareUpdate := time.Second * time.Duration(chia.loadedConfig.TimeForRareUpdate)

	sendStatistic(client, true, chia)
	for {
		timer := time.NewTimer(durationForFrequentlyUpdate)
		<-timer.C
		durationForRareUpdate -= durationForFrequentlyUpdate
		if durationForRareUpdate <= 0 {
			sendStatistic(client, true, chia)
		} else {
			sendStatistic(client, false, chia)
		}
	}

}

func sendStatistic(client pbChia.ChiaClient, fullStatistics bool, node *Chia) {
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

func (chia *Chia) GetNodeStatistic(fullStatistics bool) (request *pbChia.SaveStatisticRequest, err error) {
	request = new(pbChia.SaveStatisticRequest)

	request.NodeAuthData = new(pbCommon.NodeAuthData)
	request.NodeAuthData.Ticker = chia.loadedConfig.Ticker
	request.NodeAuthData.Uuid = chia.loadedConfig.UUID

	request.Statistic = new(pbChia.Statistic)

	if fullStatistics {
		// every 3600 seconds
		// Common
		request.Statistic.NodeBasicData = chia.commonStatistic.GetNodeBasicData()
		request.Statistic.ServerBasicData = chia.commonStatistic.GetServerBasicData()
		request.Statistic.Updates = chia.commonStatistic.GetUpdates()
		request.Statistic.Security = chia.commonStatistic.GetSecurity()
	}

	// every 10 seconds
	{
		// Common
		request.Statistic.Online = chia.commonStatistic.GetOnlineData()
		request.Statistic.MemoryState, err = chia.commonStatistic.GetMemoryData()
		if err != nil {
			return nil, err
		}
		request.Statistic.CpuState = chia.commonStatistic.GetCPUState()

		// Chia
		request.Statistic.ChiaNodeFarming = chia.GetChiaNodeFarming()
	}

	return request, nil
}
