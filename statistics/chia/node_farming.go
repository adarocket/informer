package chia

import (
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sergey-shpilevskiy/go-bytesize"

	pbChia "github.com/adarocket/proto/proto-gen/chia"
)

// GetBlocks -
func (chia *Chia) GetChiaNodeFarming() *pbChia.ChiaNodeFarming {
	out, err := exec.Command("bash", "-c", `cd `+chia.loadedConfig.PathToChiaBlockchain+`; . ./activate; chia farm summary`).Output()
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return parse(string(out))
}

func parse(text string) *pbChia.ChiaNodeFarming {
	var err error
	var chiaNodeFarming pbChia.ChiaNodeFarming
	textMap := make(map[string]string)

	rows := strings.Split(text, "\n")

	for _, row := range rows {
		if row != "" {
			s := strings.Split(row, ": ")
			if len(s) != 2 {
				log.Println("invalid input")
				continue
			}

			textMap[s[0]] = s[1]
		}
	}
	chiaNodeFarming.FarmingStatus = textMap["Farming status"]

	value, ok := textMap["Total chia farmed"]
	if ok {
		chiaNodeFarming.TotalChiaFarmed = stringToFloat32(value)
	}

	value, ok = textMap["User transaction fees"]
	if ok {
		chiaNodeFarming.UserTransactionFees = stringToFloat32(value)
	}

	value, ok = textMap["Block rewards"]
	if ok {
		chiaNodeFarming.BlockRewards = stringToFloat32(value)
	}

	value, ok = textMap["Last height farmed"]
	if ok {
		if chiaNodeFarming.LastHeightFarmed, err = strconv.ParseUint(value, 10, 64); err != nil {
			log.Println(err)
		}
	}

	value, ok = textMap["Plot count"]
	if ok {
		if chiaNodeFarming.PlotCount, err = strconv.ParseUint(value, 10, 64); err != nil {
			log.Println(err)
		}
	}

	value, ok = textMap["Total size of plots"]
	if ok {
		chiaNodeFarming.TotalSizeOfPlots = stringToByte(value)
	}

	value, ok = textMap["Estimated network space"]
	if ok {
		chiaNodeFarming.EstimatedNetworkSpace = stringToByte(value)
	}

	chiaNodeFarming.ExpectedTimeToWin = textMap["Expected time to win"]

	return &chiaNodeFarming
}

func stringToFloat32(s string) float32 {
	value, err := strconv.ParseFloat(s, 32)
	if err != nil {
		log.Println(err)
		return 0
	}

	return float32(value)
}

func stringToByte(s string) uint64 {
	b, err := bytesize.Parse(strings.ReplaceAll(s, "i", ""))
	if err != nil {
		log.Println(err)
		return 0
	}

	return uint64(b)
}
