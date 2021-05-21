package cardano

import (
	pb "github.com/adarocket/proto/proto"
	"github.com/tidwall/gjson"
)

// GetNodeState -
func (cardano *Cardano) GetNodeState(jsonBody string) *pb.NodeState {
	var nodeState pb.NodeState

	tipNode := gjson.Get(jsonBody, "cardano.node.metrics.slotNum.int.val").Int()
	tipRef, _ := cardano.getSlotTipRef(jsonBody)
	nodeState.TipDiff = tipRef - tipNode

	nodeState.Density = float32(gjson.Get(jsonBody, "cardano.node.metrics.density.real.val").Float())

	return &nodeState
}
