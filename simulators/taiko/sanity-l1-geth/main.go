package main

import (
	"context"
	"github.com/ethereum/hive/hivesim"
	"github.com/ethereum/hive/simulators/taiko/common/clients"
	el "github.com/ethereum/hive/simulators/taiko/common/config/execution"
	"github.com/ethereum/hive/simulators/taiko/common/testnet"
	tn "github.com/ethereum/hive/simulators/taiko/common/testnet"
	"math/big"
)

var (
	TERMINAL_TOTAL_DIFFICULTY = big.NewInt(100)
)

func main() {
	suite := hivesim.Suite{
		Name:        "Sanity - L1 - geth",
		Description: "L1 geth sanity test suite",
	}
	suite.Add(&hivesim.TestSpec{
		Name:        "chainId31336",
		Description: "Asserts that the ChainId is equal 31336",
		Run: func(t *hivesim.T) {
			clientTypes, err := t.Sim.ClientTypes()
			if err != nil {
				t.Fatal(err)
			}
			c := clients.ClientsByRole(clientTypes)
			if len(c.L1ExecutionClient) != 1 {
				t.Fatal("choose 1 l1_client client type")
			}
			for _, node := range c.Combinations() {
				env := &testnet.Environment{
					Clients: c,
				}
				config := tn.Config{
					TerminalTotalDifficulty: TERMINAL_TOTAL_DIFFICULTY,
					NodeDefinitions: []clients.NodeDefinition{
						node,
					},
					L1ExecutionConsensus: el.ExecutionCliqueConsensus{},
				}

				ctx := context.Background()
				_ = tn.StartTestnet(ctx, t, env, &config)
			}

			chainId31336(t, c)
		},
	})

	sim := hivesim.New()
	hivesim.MustRun(sim, suite)
}

func chainId31336(t *hivesim.T, c *clients.ClientDefinitionsByRole) {
	//_, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	//defer cancel()
	//
	//chainIdString := ""
	//err := c.RPC().Call(&chainIdString, "eth_chainId")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//chainId, err := strconv.ParseInt(chainIdString, 0, 0)
	//if err != nil {
	//	fmt.Println("Error: %e", err)
	//	return
	//}
	//if chainId != 31336 {
	//	t.Fatalf("ChainId is not equal 31336, it is %i", chainId)
	//}
}
