package main

import (
	"time"

	pb "github.com/adarocket/proto"
)

// GetOnlineData -
func GetOnlineData() *pb.Online {
	var onlineData pb.Online

	onlineData.SinceStart = int64(time.Now().Sub(startTime).Seconds())
	onlineData.Pings = 0
	// onlineData.ServerActive = true
	onlineData.NodeActive = true
	onlineData.NodeActivePings = 0

	return &onlineData
}
