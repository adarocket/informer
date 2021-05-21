package cardano

import (
	pb "github.com/adarocket/proto/proto"
	"github.com/tidwall/gjson"
)

// GetEpoch -
func (cardano *Cardano) GetEpoch(jsonBody string) *pb.Epoch {
	var epoch pb.Epoch
	epoch.EpochNumber = gjson.Get(jsonBody, "cardano.node.metrics.epoch.int.val").Int()
	return &epoch
}
