package testnet

import (
	"github.com/ethereum/hive/simulators/taiko/common/clients"
)

type Environment struct {
	Clients *clients.ClientDefinitionsByRole
	//Keys           []*consensus_config.ValidatorDetails
	//Secrets        *[]blsu.SecretKey
	LogEngineCalls bool
}
