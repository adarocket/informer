package common

import (
	"log"
	"os/exec"
	"time"

	pb "github.com/adarocket/proto/proto-gen/common"
)

// GetOnlineData -
func (commonStatistic *CommonStatistic) GetOnlineData() *pb.Online {
	var onlineData pb.Online

	onlineData.SinceStart = int64(time.Since(commonStatistic.startTime).Seconds())
	onlineData.Pings = 0
	// onlineData.ServerActive = true
	onlineData.NodeActive = isNodeActive()
	onlineData.NodeActivePings = 0

	return &onlineData
}

func isNodeActive() bool {
	_, err := exec.Command("cardano-node", "version").Output()
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
