package clients

import (
	"github.com/taikoxyz/hive-taiko-clients/clients"
	"github.com/taikoxyz/hive/simulators/taiko/common/utils"
)

type ProtocolDeployerClientConfig struct {
	Layer  string
	Subnet string
}

type ProtocolDeployerClient struct {
	clients.Client
	Logger  utils.Logging
	Config  ProtocolDeployerClientConfig
	Builder interface{}
}
