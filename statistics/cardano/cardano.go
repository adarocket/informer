package cardano

import (
	"adarocket/informer/config"
	"adarocket/informer/statistics/common"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	pb "github.com/adarocket/proto/proto"
	"github.com/tidwall/gjson"
)

type Cardano struct {
	loadedConfig    *config.Config
	commonStatistic *common.CommonStatistic
}

func NewCardano(config *config.Config, startTime time.Time) (cardano *Cardano) {
	cardano = new(Cardano)
	cardano.loadedConfig = config

	cardano.commonStatistic = common.NewCommonStatistic(config, startTime)

	return cardano
}

func (cardano *Cardano) GetNodeStatistic(fullStatistics bool) (*pb.SaveStatisticRequest, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", cardano.loadedConfig.NodeMonitoringURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("Status Code not 200")
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	jsonBody := string(bodyBytes)

	request := new(pb.SaveStatisticRequest)

	request.NodeAuthData = new(pb.NodeAuthData)
	request.NodeAuthData.Ticker = cardano.loadedConfig.Ticker
	request.NodeAuthData.Uuid = cardano.loadedConfig.UUID

	request.Statistic = new(pb.Statistic)

	if fullStatistics {
		// every 3600 seconds

		// Common
		request.Statistic.NodeBasicData = cardano.commonStatistic.GetNodeBasicData()
		request.Statistic.ServerBasicData = cardano.commonStatistic.GetServerBasicData()
		request.Statistic.Updates = cardano.commonStatistic.GetUpdates()
		request.Statistic.Security = cardano.commonStatistic.GetSecurity()

		request.Statistic.Epoch = cardano.GetEpoch(jsonBody)
		request.Statistic.KesData = cardano.GetKESData(jsonBody)
		request.Statistic.Blocks = cardano.GetBlocks(jsonBody)
		request.Statistic.StakeInfo = cardano.GetStakeInfo()
	}

	// every 10 seconds
	{
		// Common
		request.Statistic.Online = cardano.commonStatistic.GetOnlineData() // -
		if request.Statistic.MemoryState, err = cardano.commonStatistic.GetMemoryData(); err != nil {
			log.Println(err)
			return nil, err
		}
		request.Statistic.CpuState = cardano.commonStatistic.GetCPUState()

		request.Statistic.NodeState = cardano.GetNodeState(jsonBody)
		request.Statistic.NodePerformance = cardano.GetNodePerformance(jsonBody)
	}

	return request, nil
}

func (cardano *Cardano) getSlotTipRef(jsonBody string) (slotTipRef int64, expirationTime string) {
	currentTimeSec := time.Now().Unix()

	shelleyGenesisData, err := ioutil.ReadFile(cardano.loadedConfig.MainnetShelleyGenesisJsonPath)
	if err != nil {
		log.Println(err)
	}
	shelleyGenesisJSON := string(shelleyGenesisData)

	slotLength := gjson.Get(shelleyGenesisJSON, "slotLength").Int()

	byronGenesisData, err := ioutil.ReadFile(cardano.loadedConfig.MainnetByronGenesisJsonPath)
	if err != nil {
		log.Println(err)
	}
	byronGenesisJSON := string(byronGenesisData)

	byronGenesisStartSec := gjson.Get(byronGenesisJSON, "startTime").Int()
	byronSlotLength := gjson.Get(byronGenesisJSON, "blockVersionData.slotDuration").String()
	byronEpochLength := gjson.Get(byronGenesisJSON, "protocolConsts.k").Int() * 10

	byronSlotLengthInt, err := strconv.ParseInt(byronSlotLength, 10, 64)
	if err != nil {
		log.Println(err)
	}

	shelleyTransEpoch := int64(208)

	byronSlots := shelleyTransEpoch * byronEpochLength

	byronEndTime := byronGenesisStartSec + ((shelleyTransEpoch * byronEpochLength * byronSlotLengthInt) / 1000)

	if currentTimeSec < byronEndTime {
		slotTipRef = ((currentTimeSec - byronGenesisStartSec) * 1000) / byronSlotLengthInt
	} else {
		slotTipRef = byronSlots + ((currentTimeSec - byronEndTime) / slotLength)
	}

	slotsPerKESPeriod := gjson.Get(shelleyGenesisJSON, "slotsPerKESPeriod").Int()

	remainingKesPeriods := gjson.Get(jsonBody, "cardano.node.metrics.remainingKESPeriods.int.val").Int()

	expirationTimeSec := currentTimeSec - (slotLength * (slotTipRef % slotsPerKESPeriod)) + (slotLength * slotsPerKESPeriod * remainingKesPeriods)
	expirationTime = time.Unix(expirationTimeSec, 0).String()

	return slotTipRef, expirationTime
}