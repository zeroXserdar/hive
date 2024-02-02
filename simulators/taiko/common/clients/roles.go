package clients

import "github.com/taikoxyz/hive/hivesim"

type ClientDefinitionsByRole struct {
	L1ExecutionClient          []*hivesim.ClientDefinition `json:"l1_client"`
	L2ExecutionClient          []*hivesim.ClientDefinition `json:"l2_client"`
	L1L2ProtocolDeployerClient []*hivesim.ClientDefinition `json:"l1l2_protocol_deployer"`
	L2DriverClient             []*hivesim.ClientDefinition `json:"l2_driver"`
	L2ProposerClient           []*hivesim.ClientDefinition `json:"l2_proposer"`
	L2ProverClient             []*hivesim.ClientDefinition `json:"l2_prover"`
}

func ClientsByRole(
	available []*hivesim.ClientDefinition,
) *ClientDefinitionsByRole {
	var out ClientDefinitionsByRole
	for _, client := range available {
		if client.HasRole("l1_client") {
			out.L1ExecutionClient = append(out.L1ExecutionClient, client)
		}
		if client.HasRole("l2_client") {
			out.L2ExecutionClient = append(out.L2ExecutionClient, client)
		}
		if client.HasRole("l1l2_protocol_deployer") {
			out.L1L2ProtocolDeployerClient = append(out.L1L2ProtocolDeployerClient, client)
		}
		if client.HasRole("l2_driver") {
			out.L2DriverClient = append(out.L2DriverClient, client)
		}
		if client.HasRole("l2_proposer") {
			out.L2DriverClient = append(out.L2DriverClient, client)
		}
		if client.HasRole("l2_prover") {
			out.L2ProverClient = append(out.L2ProverClient, client)
		}
	}
	return &out
}

func (c *ClientDefinitionsByRole) ClientByNameAndRole(
	name, role string,
) *hivesim.ClientDefinition {
	switch role {
	case "l1_client":
		return byName(c.L1ExecutionClient, name)
	case "l2_client":
		return byName(c.L2ExecutionClient, name)
	case "l1l2_protocol_deployer":
		return byName(c.L1L2ProtocolDeployerClient, name)
	case "l2_driver":
		return byName(c.L2DriverClient, name)
	case "l2_proposer":
		return byName(c.L2ProposerClient, name)
	case "l2_prover":
		return byName(c.L2ProverClient, name)
	}
	return nil
}

func byName(
	clients []*hivesim.ClientDefinition,
	name string,
) *hivesim.ClientDefinition {
	for _, client := range clients {
		if client.Name == name {
			return client
		}
	}
	return nil
}

func (c *ClientDefinitionsByRole) Combinations() NodeDefinitions {
	var nodes NodeDefinitions
	//for _, beacon := range c.Beacon {
	//	for _, eth1 := range c.Eth1 {
	//		nodes = append(
	//			nodes,
	//			NodeDefinition{
	//				ExecutionClient: eth1.Name,
	//				ConsensusClient: beacon.Name,
	//			},
	//		)
	//	}
	//}
	nodes = append(
		nodes,
		NodeDefinition{
			L1ExecutionClient: c.L1ExecutionClient[0].Name,
		})
	if len(c.L1L2ProtocolDeployerClient) == 1 {
		nodes[0].L1L2ProtocolDeployerClient = c.L1L2ProtocolDeployerClient[0].Name
	}
	return nodes
}
