package config

import (
	"fmt"
	"github.com/bykovme/goconfig"
	"github.com/google/uuid"
	"log"
	"strings"
)

const welcomeText = "informer.conf is not found when starting the application," +
	" trying create new one"

func startCli(usrHomePath string) (createdConfig *Config, err error) {
	fmt.Println(welcomeText)

	createdConfig = new(Config)

	for isInputEmpty(createdConfig.Ticker) {
		fmt.Println("input ticker: ")
		fmt.Scan(&createdConfig.Ticker)
	}

	fmt.Println("input location: ")
	fmt.Scan(&createdConfig.Location)

	for isInputEmpty(createdConfig.ControllerURL) {
		fmt.Println("input controller_url: ")
		fmt.Scan(&createdConfig.ControllerURL)
	}

	for isInputEmpty(createdConfig.NodeMonitoringURL) {
		fmt.Println("input node_monitoring url: ")
		fmt.Scan(&createdConfig.NodeMonitoringURL)
	}

	createdConfig.TimeForFrequentlyUpdate = 10
	createdConfig.TimeForRareUpdate = 60

	fmt.Println("input blockchain: ")
	fmt.Scan(&createdConfig.Blockchain)
	if isInputEmpty(createdConfig.Blockchain) {
		createdConfig.Blockchain = "cardano"
	}

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

func isInputEmpty(input string) bool {
	input = strings.ReplaceAll(input, " ", "")

	return len(input) == 0
}
