package config

import (
	"github.com/bykovme/goconfig"
	"log"
)

const cConfigPath = "/etc/ada-rocket/informer.conf"

// Config - structure of config file
type Config struct {
	Ticker                  string `json:"ticker"`
	UUID                    string `json:"uuid"`
	Location                string `json:"location"`
	ControllerURL           string `json:"controller_url"`
	TimeForFrequentlyUpdate int    `json:"time_for_frequently_update"`
	TimeForRareUpdate       int    `json:"time_for_rare_update"`

	NodeMonitoringURL             string `json:"node_monitoring_url"`
	MainnetShelleyGenesisJsonPath []byte `json:"-"`
	MainnetByronGenesisJsonPath   []byte `json:"-"`

	Blockchain           string `json:"blockchain"`
	PathToChiaBlockchain string `json:"path_to_chia_blockchain"`
}

func LoadConfig() (loadedConfig *Config, err error) {
	usrHomePath, err := goconfig.GetUserHomePath()
	if err != nil {
		return loadedConfig, err
	}

	shelley, err := getGenesisJson(shelleyName)
	if err != nil {
		log.Println(err)
		return loadedConfig, err
	}

	byron, err := getGenesisJson(byronName)
	if err != nil {
		log.Println(err)
		return loadedConfig, err
	}

	loadedConfig = new(Config)

	loadedConfig.MainnetByronGenesisJsonPath = byron
	loadedConfig.MainnetShelleyGenesisJsonPath = shelley
	err = goconfig.LoadConfig(usrHomePath+cConfigPath, loadedConfig)
	return loadedConfig, err
}
