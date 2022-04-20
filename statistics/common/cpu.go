package common

import (
	pb "github.com/adarocket/proto/proto-gen/common"
	"log"
	"runtime"
)

// GetCPUState -
func (commonStatistic *CommonStatistic) GetCPUState() *pb.CPUState {
	var cpuState pb.CPUState

	cpuState.CpuQty = int64(runtime.NumCPU())
	load, err := avgLoad()
	if err != nil {
		log.Println(err)
	}

	cpuState.AverageWorkload = float32(load.Load15)

	return &cpuState
}

type avgStat struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

func avgLoad() (avgStat, error) {
	/*
	type loadavg struct {
		load  [3]uint32
		scale int
	}

	b, err := unix.SysctlRaw("vm.loadavg")
	if err != nil {
		log.Println(err)
		return avgStat{}, err
	}

	load := *(*loadavg)(unsafe.Pointer((&b[0])))
	scale := float64(load.scale)
	ret := avgStat{
		Load1:  float64(load.load[0]) / scale,
		Load5:  float64(load.load[1]) / scale,
		Load15: float64(load.load[2]) / scale,
	}
	*/
	return ret, nil
}
