package common

import pb "github.com/adarocket/proto/proto-gen/common"

// GetUpdates -
func (commonStatistic *CommonStatistic) GetUpdates() *pb.Updates {
	var updates pb.Updates

	updates.InformerActual = ""
	updates.InformerAvailable = ""
	updates.UpdaterActual = ""
	updates.UpdaterAvailable = ""
	updates.PackagesAvailable = 0

	return &updates
}
