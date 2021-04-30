package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	pb "github.com/adarocket/proto"

	"github.com/tidwall/gjson"
)

func getNodeStatistic(fullStatistics bool) (*pb.SaveStatisticRequest, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", loadedConfig.NodeMonitoringURL, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		log.Println(resp.StatusCode)
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	jsonBody := string(bodyBytes)

	request := new(pb.SaveStatisticRequest)
	request.NodeAuthData = new(pb.NodeAuthData)

	request.NodeAuthData.Ticker = loadedConfig.Ticker
	request.NodeAuthData.Uuid = loadedConfig.UUID

	request.Statistic = new(pb.Statistic)

	if fullStatistics {
		// every 3600 seconds
		request.Statistic.NodeBasicData = GetNodeBasicData()
		request.Statistic.ServerBasicData = GetServerBasicData() // +
		request.Statistic.Epoch = GetEpoch(jsonBody)             // +
		request.Statistic.KesData = GetKESData(jsonBody)         // -
		request.Statistic.Blocks = GetBlocks(jsonBody)           // -
		request.Statistic.Updates = GetUpdates()
		request.Statistic.Security = GetSecurity()
		request.Statistic.StakeInfo = GetStakeInfo()
	}

	// every 10 seconds
	{
		request.Statistic.Online = GetOnlineData() // -
		if request.Statistic.MemoryState, err = GetMemoryData(); err != nil {
			log.Println(err)
			return nil, err
		}
		request.Statistic.CpuState = GetCPUState()
		request.Statistic.NodeState = GetNodeState(jsonBody)             // -
		request.Statistic.NodePerformance = GetNodePerformance(jsonBody) // -
	}

	return request, nil
}

// GetKESData -
func GetKESData(jsonBody string) *pb.KESData {
	var kesData pb.KESData

	kesData.KesCurrent = gjson.Get(jsonBody, "cardano.node.metrics.currentKESPeriod.int.val").Int()
	kesData.KesRemaining = gjson.Get(jsonBody, "cardano.node.metrics.remainingKESPeriods.int.val").Int()

	_, kesData.KesExpDate = getSlotTipRef(jsonBody)

	return &kesData
}

func getSlotTipRef(jsonBody string) (slotTipRef int64, expirationTime string) {
	currentTimeSec := time.Now().Unix()

	shelleyGenesisData, err := ioutil.ReadFile("/home/ada/cardano-my-node/mainnet-shelley-genesis.json")
	if err != nil {
		log.Println(err)
	}
	shelleyGenesisJSON := string(shelleyGenesisData)

	slotLength := gjson.Get(shelleyGenesisJSON, "slotLength").Int()

	byronGenesisData, err := ioutil.ReadFile("/home/ada/cardano-my-node/mainnet-byron-genesis.json")
	if err != nil {
		log.Println(err)
	}
	byronGenesisJSON := string(byronGenesisData)

	byronGenesisStartSec := gjson.Get(byronGenesisJSON, "startTime").Int()
	byronSlotLength := gjson.Get(byronGenesisJSON, "blockVersionData.slotDuration").String()
	byronEpochLength := gjson.Get(byronGenesisJSON, "protocolConsts.k").Int() * 10

	byronSlotLengthInt, err := strconv.ParseInt(byronSlotLength, 10, 64)

	// FIXME: Почему 208?
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

// GetBlocks -
func GetBlocks(jsonBody string) *pb.Blocks {
	var blocks pb.Blocks
	blocks.BlockLeader = gjson.Get(jsonBody, "cardano.node.metrics.Forge[\"node-is-leader\"].int.val").Int()
	blocks.BlockAdopted = gjson.Get(jsonBody, "cardano.node.metrics.Forge.adopted.int.val").Int()
	blocks.BlockInvalid = gjson.Get(jsonBody, "cardano.node.metrics.Forge[\"didnt-adopt\"].int.val").Int()
	return &blocks
}

// GetUpdates -
func GetUpdates() *pb.Updates {
	var updates pb.Updates

	updates.InformerActual = ""
	updates.InformerAvailable = ""
	updates.UpdaterActual = ""
	updates.UpdaterAvailable = ""
	updates.PackagesAvailable = 0

	return &updates
}

// GetSecurity -
func GetSecurity() *pb.Security {
	var security pb.Security

	security.SshAttackAttempts = 0
	security.SecurityPackagesAvailable = 0

	return &security
}

// GetStakeInfo -
func GetStakeInfo() *pb.StakeInfo {
	var stakeInfo pb.StakeInfo

	stakeInfo.LiveStake = 0
	stakeInfo.ActiveStake = 0
	stakeInfo.Pledge = 0

	return &stakeInfo
}

// --------------------------

// GetCPUState -
func GetCPUState() *pb.CPUState {
	var cpuState pb.CPUState

	cpuState.CpuQty = 1
	cpuState.AverageWorkload = 0.0

	return &cpuState
}

// GetNodeState -
func GetNodeState(jsonBody string) *pb.NodeState {
	var nodeState pb.NodeState

	tipNode := gjson.Get(jsonBody, "cardano.node.metrics.slotNum.int.val").Int()
	tipRef, _ := getSlotTipRef(jsonBody)
	nodeState.TipDiff = tipRef - tipNode

	nodeState.Density = float32(gjson.Get(jsonBody, "cardano.node.metrics.density.real.val").Float())

	return &nodeState
}

// GetNodePerformance -
func GetNodePerformance(jsonBody string) *pb.NodePerformance {
	var nodePerformance pb.NodePerformance

	nodePerformance.ProcessedTx = gjson.Get(jsonBody, "cardano.node.metrics.txsProcessedNum.int.val").Int()
	nodePerformance.PeersIn = getPeersIn()
	nodePerformance.PeersOut = gjson.Get(jsonBody, "cardano.node.metrics.connectedPeers.int.val").Int()

	return &nodePerformance
}

func getPeersIn() int64 {
	// $(netstat -an|awk "\$4 ~ /${cardanoport}/"|grep -c ESTABLISHED)
	// netstat -an|awk "\$4 ~ /6000/"|grep -c ESTABLISHED
	out, err := exec.Command("bash", "-c", `netstat -an|awk "\$4 ~ 6000"`).Output()
	if err != nil {
		log.Println(err.Error())
		return 0
	}

	return count("ESTABLISHED", string(out))
}

func count(str, value string) int64 {
	re := regexp.MustCompile(str)
	// Find all matches and return count.
	results := re.FindAllString(value, -1)
	return int64(len(results))
}
