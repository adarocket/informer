package common

import pb "github.com/adarocket/proto/proto-gen/common"

// GetSecurity -
func (commonStatistic *CommonStatistic) GetSecurity() *pb.Security {
	var security pb.Security

	security.SshAttackAttempts = 0
	security.SecurityPackagesAvailable = 0

	return &security
}
