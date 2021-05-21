package common

import (
	"time"

	pb "github.com/adarocket/proto/proto"
)

// GetOnlineData -
func (commonStatistic *CommonStatistic) GetOnlineData() *pb.Online {
	var onlineData pb.Online

	onlineData.SinceStart = int64(time.Since(commonStatistic.startTime).Seconds())
	onlineData.Pings = 0
	// onlineData.ServerActive = true
	onlineData.NodeActive = true
	onlineData.NodeActivePings = 0

	return &onlineData
}
