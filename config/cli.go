package config

import (
	"fmt"
	"github.com/bykovme/goconfig"
	"github.com/google/uuid"
	"log"
)

const welcomeText = "informer.conf is not found when starting the application," +
	" trying create new one"

func startCli(usrHomePath string) (createdConfig *Config, err error) {
	fmt.Println(welcomeText)

	createdConfig = new(Config)

	fmt.Println("input ticker: ")
	fmt.Scan(&createdConfig.Ticker)
	fmt.Println("input location: ")
	fmt.Scan(&createdConfig.Location)
	fmt.Println("input controller_url: ")
	fmt.Scan(&createdConfig.ControllerURL)
	fmt.Println("input node_monitoring url: ")
	fmt.Scan(&createdConfig.NodeMonitoringURL)
	fmt.Println("input time_for_frequently_update: ")
	fmt.Scan(&createdConfig.TimeForFrequentlyUpdate)
	fmt.Println("input time_for_rare_update: ")
	fmt.Scan(&createdConfig.TimeForRareUpdate)
	fmt.Println("input blockchain: ")
	fmt.Scan(&createdConfig.Blockchain)

	createdConfig.UUID = uuid.New().String()
	fmt.Println("Your uuid:", createdConfig.UUID)

	err = goconfig.SaveConfig(usrHomePath+cConfigPath, createdConfig)
	if err != nil {
		log.Println(err)
		return createdConfig, err
	}

	fmt.Println("config path:", usrHomePath+cConfigPath)
	return createdConfig, nil
}
