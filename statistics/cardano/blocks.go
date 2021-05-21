package cardano

import (
	pb "github.com/adarocket/proto/proto"
	"github.com/tidwall/gjson"
)

// GetBlocks -
func (cardano *Cardano) GetBlocks(jsonBody string) *pb.Blocks {
	var blocks pb.Blocks
	blocks.BlockLeader = gjson.Get(jsonBody, "cardano.node.metrics.Forge[\"node-is-leader\"].int.val").Int()
	blocks.BlockAdopted = gjson.Get(jsonBody, "cardano.node.metrics.Forge.adopted.int.val").Int()
	blocks.BlockInvalid = gjson.Get(jsonBody, "cardano.node.metrics.Forge[\"didnt-adopt\"].int.val").Int()
	return &blocks
}
