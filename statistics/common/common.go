package common

import (
	"adarocket/informer/config"
	"time"
)

type CommonStatistic struct {
	loadedConfig *config.Config
	startTime    time.Time
}

func NewCommonStatistic(config *config.Config, startTime time.Time) (commonStatistic *CommonStatistic) {
	commonStatistic = new(CommonStatistic)
	commonStatistic.loadedConfig = config
	commonStatistic.startTime = startTime

	return commonStatistic
}
