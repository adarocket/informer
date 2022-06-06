package config

import (
	"log"

	"github.com/bykovme/goconfig"
)

const cConfigPath = "/etc/ada-rocket/informer.conf"

// Config - structure of config file
type Config struct {
	Ticker string `json:"ticker"`
	Type   string `json:"type"`
	Name   string `json:"name"`

	UUID                    string `json:"uuid"`
	Location                string `json:"location"`
	ControllerURL           string `json:"controller_url"`
	TimeForFrequentlyUpdate int    `json:"time_for_frequently_update"`
	TimeForRareUpdate       int    `json:"time_for_rare_update"`

	NodeMonitoringURL string `json:"node_monitoring_url"`

	Blockchain           string `json:"blockchain"`
	PathToChiaBlockchain string `json:"path_to_chia_blockchain"`
}

func LoadConfig() (loadedConfig *Config, err error) {
	usrHomePath, err := goconfig.GetUserHomePath()
	if err != nil {
		return loadedConfig, err
	}

	loadedConfig = new(Config)
	err = goconfig.LoadConfig(usrHomePath+cConfigPath, loadedConfig)
	if err != nil {
		log.Println(err)
		loadedConfig, err = startCli(usrHomePath)
	}

	return loadedConfig, err
}
