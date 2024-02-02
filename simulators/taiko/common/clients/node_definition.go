package clients

import (
	"fmt"
	"golang.org/x/exp/slices"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	cg "github.com/taikoxyz/hive/simulators/taiko/common/chain_generators"
)

// Describe a node setup, which consists of:
// - L1 Execution Client
// - L2 Execution Client
// - L1L2 Protocol Deployer Client
// - L2 Driver Client
// - L2 Proposer Client
// - L2 Prover Client
type NodeDefinition struct {
	// Client Types
	L1ExecutionClient          string `json:"l1_execution_client"`
	L2ExecutionClient          string `json:"l2_execution_client"`
	L1L2ProtocolDeployerClient string `json:"l1l2_protocol_deployer_client"`
	L2DriverClient             string `json:"l2_driver_client"`
	L2ProposerClient           string `json:"l2_proposer_client"`
	L2ProverClient             string `json:"l2_prover_client"`

	// L1 Execution Config
	L1ExecutionClientTTD *big.Int          `json:"l1_execution_client_ttd,omitempty"`
	L1ChainGenerator     cg.ChainGenerator `json:"-"`
	L1Chain              []*types.Block    `json:"l1_chain,omitempty"`

	// L2 Execution Config
	L2ExecutionClientTTD *big.Int          `json:"l2_execution_client_ttd,omitempty"`
	L2ChainGenerator     cg.ChainGenerator `json:"-"`
	L2Chain              []*types.Block    `json:"l2_chain,omitempty"`

	// Node Config
	//TestVerificationNode bool `json:"test_verification_node"`
	//DisableStartup       bool `json:"disable_startup"`

	// Subnet Configuration
	L1Subnet string `json:"l1_subnet"`
	L2Subnet string `json:"l2_subnet"`
	Subnet   string `json:"subnet"`
}

func (n *NodeDefinition) String() string {
	return fmt.Sprintf("%s-%s", n.L1ExecutionClient, n.L2ExecutionClient)
}

func (n *NodeDefinition) L1ExecutionClientName() string {
	return n.L1ExecutionClient
}

func (n *NodeDefinition) L2ExecutionClientName() string {
	return n.L2ExecutionClient
}

func (n *NodeDefinition) L1L2ProtocolDeployerClientName() string {
	return n.L1L2ProtocolDeployerClient
}

func (n *NodeDefinition) L2DriverClientName() string {
	return n.L2DriverClient
}

func (n *NodeDefinition) L2ProposerClientName() string {
	return n.L2ProposerClient
}

func (n *NodeDefinition) L2ProverClientName() string {
	return n.L2ProverClient
}
func (n *NodeDefinition) GetL1Subnet() string {
	if n.L1Subnet != "" {
		return n.L1Subnet
	}
	if n.Subnet != "" {
		return n.Subnet
	}
	return ""
}

func (n *NodeDefinition) GetL2Subnet() string {
	if n.L2Subnet != "" {
		return n.L2Subnet
	}
	if n.Subnet != "" {
		return n.Subnet
	}
	return ""
}

//func beaconNodeToValidator(name string) string {
//	name, branch, hasBranch := strings.Cut(name, "_")
//	name = strings.TrimSuffix(name, "-bn")
//	validator := name + "-vc"
//	if hasBranch {
//		validator += "_" + branch
//	}
//	return validator
//}

type NodeDefinitions []NodeDefinition

func (nodes NodeDefinitions) ClientTypes() []string {
	types := make([]string, 0)
	for _, n := range nodes {
		if !slices.Contains(types, n.L1ExecutionClient) {
			types = append(types, n.L1ExecutionClient)
		}
		if !slices.Contains(types, n.L2ExecutionClient) {
			types = append(types, n.L2ExecutionClient)
		}
		if !slices.Contains(types, n.L1L2ProtocolDeployerClient) {
			types = append(types, n.L1L2ProtocolDeployerClient)
		}
		if !slices.Contains(types, n.L2DriverClient) {
			types = append(types, n.L2DriverClient)
		}
		if !slices.Contains(types, n.L2ProposerClient) {
			types = append(types, n.L2ProposerClient)
		}
		if !slices.Contains(types, n.L2ProverClient) {
			types = append(types, n.L2ProverClient)
		}
	}
	return types
}

//TODO: decide if it makes sense to implement
//func (all NodeDefinitions) FilterByCL(filters []string) NodeDefinitions {
//	ret := make(NodeDefinitions, 0)
//	for _, n := range all {
//		for _, filter := range filters {
//			if strings.Contains(n.ConsensusClient, filter) {
//				ret = append(ret, n)
//				break
//			}
//		}
//	}
//	return ret
//}

//TODO: decide if it makes sense to implement
//func (all NodeDefinitions) FilterByEL(filters []string) NodeDefinitions {
//	ret := make(NodeDefinitions, 0)
//	for _, n := range all {
//		for _, filter := range filters {
//			if strings.Contains(n.ExecutionClient, filter) {
//				ret = append(ret, n)
//				break
//			}
//		}
//	}
//	return ret
//}
