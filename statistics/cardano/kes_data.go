package cardano

import (
	pb "github.com/adarocket/proto/proto-gen/cardano"
	"github.com/tidwall/gjson"
)

// GetKESData -
func (cardano *Cardano) GetKESData(jsonBody string) *pb.KESData {
	var kesData pb.KESData

	kesData.KesCurrent = gjson.Get(jsonBody, "cardano.node.metrics.currentKESPeriod.int.val").Int()
	kesData.KesRemaining = gjson.Get(jsonBody, "cardano.node.metrics.remainingKESPeriods.int.val").Int()

	_, kesData.KesExpDate = cardano.getSlotTipRef(jsonBody)

	return &kesData
}
