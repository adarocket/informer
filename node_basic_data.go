package main

import (
	"log"
	"os/exec"
	"regexp"

	pb "github.com/adarocket/proto"
)

// GetNodeBasicData -
func GetNodeBasicData() *pb.NodeBasicData {
	var nodeBasicData pb.NodeBasicData

	nodeBasicData.Ticker = loadedConfig.Ticker
	nodeBasicData.Type = "" // from node info
	nodeBasicData.Location = loadedConfig.Location
	nodeBasicData.NodeVersion = getNodeVersion()

	return &nodeBasicData
}

func getNodeVersion() string {
	out, err := exec.Command("cardano-node", "version").Output()
	if err != nil {
		log.Println(err)
		return ""
	}

	var validNodeVersion = regexp.MustCompile(`\d{1,2}\.\d{1,3}\.\d{1,3}`)
	return validNodeVersion.FindString(string(out))
}
