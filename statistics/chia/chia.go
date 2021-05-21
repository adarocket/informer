package chia

import (
	"adarocket/informer/config"
	"adarocket/informer/statistics/common"
	"time"

	pb "github.com/adarocket/proto/proto"
)

type Chia struct {
	loadedConfig    *config.Config
	commonStatistic *common.CommonStatistic
}

func NewChia(config *config.Config, startTime time.Time) (chia *Chia) {
	chia = new(Chia)
	chia.loadedConfig = config
	chia.commonStatistic = common.NewCommonStatistic(config, startTime)
	return chia
}

func (chia *Chia) GetNodeStatistic(fullStatistics bool) (request *pb.SaveStatisticRequest, err error) {
	request = new(pb.SaveStatisticRequest)

	request.NodeAuthData = new(pb.NodeAuthData)
	request.NodeAuthData.Ticker = chia.loadedConfig.Ticker
	request.NodeAuthData.Uuid = chia.loadedConfig.UUID

	request.Statistic = new(pb.Statistic)

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
