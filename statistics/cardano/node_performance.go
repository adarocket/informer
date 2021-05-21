package cardano

import (
	"log"
	"os/exec"
	"regexp"

	"github.com/tidwall/gjson"

	pb "github.com/adarocket/proto/proto"
)

// GetNodePerformance -
func (cardano *Cardano) GetNodePerformance(jsonBody string) *pb.NodePerformance {
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
