package common

import pb "github.com/adarocket/proto/proto"

// GetCPUState -
func (commonStatistic *CommonStatistic) GetCPUState() *pb.CPUState {
	var cpuState pb.CPUState

	cpuState.CpuQty = 1
	cpuState.AverageWorkload = 0.0

	return &cpuState
}
