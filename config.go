package main

import (
	"log"

	"github.com/bykovme/goconfig"
)

// Config - structure of config file
type Config struct {
	Ticker                  string `json:"ticker"`
	UUID                    string `json:"uuid"`
	Location                string `json:"location"`
	NodeMonitoringURL       string `json:"node_monitoring_url"`
	ControllerURL           string `json:"controller_url"`
	TimeForFrequentlyUpdate int    `json:"time_for_frequently_update"`
	TimeForRareUpdate       int    `json:"time_for_rare_update"`
}

const cConfigPath = "/etc/ada-rocket/informer.conf"

var loadedConfig Config

func loadConfig() {
	usrHomePath, err := goconfig.GetUserHomePath()
	if err == nil {
		err = goconfig.LoadConfig(usrHomePath+cConfigPath, &loadedConfig)
		if err != nil {
			log.Println(err.Error())
			panic(err.Error())
		}
	} else {
		log.Println(err.Error())
		panic(err.Error())
	}
}
