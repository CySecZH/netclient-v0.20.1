package upgrades

import (
	"github.com/gravitl/netclient/config"
	"github.com/gravitl/netmaker/logger"
)

var upgrade0146 = UpgradeInfo{
	RequiredVersions: []string{
		"v0.14.0",
		"v0.14.1",
		"v0.14.2",
		"v0.14.3",
		"v0.14.4",
		"v0.14.5",
	},
	NewVersion: "v0.14.6",
	OP:         update0146,
}

func update0146(cfg *config.ClientConfig) {
	// do stuff for 14.X -> 14.5
	// No-op
	logger.Log(0, "updating schema for 0.14.6")
}
