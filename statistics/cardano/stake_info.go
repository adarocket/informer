package cardano

import pb "github.com/adarocket/proto/proto-gen/cardano"

// GetStakeInfo -
func (cardano *Cardano) GetStakeInfo() *pb.StakeInfo {
	var stakeInfo pb.StakeInfo

	stakeInfo.LiveStake = 0
	stakeInfo.ActiveStake = 0
	stakeInfo.Pledge = 0

	return &stakeInfo
}
