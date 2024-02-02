package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/hive/hivesim"
	"github.com/taikoxyz/hive/simulators/taiko/common/clients"
	el "github.com/taikoxyz/hive/simulators/taiko/common/config/execution"
	"github.com/taikoxyz/hive/simulators/taiko/common/testnet"
	tn "github.com/taikoxyz/hive/simulators/taiko/common/testnet"
	"math/big"
	"net/http"
	"strconv"
	"time"
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
				testnetInstance := tn.StartTestnet(ctx, t, env, &config)
				testnetInstance.Logf("testnet started")
				//time.Sleep(10 * time.Second)
				chainId31336(t, testnetInstance)
			}

		},
	})

	sim := hivesim.New()
	hivesim.MustRun(sim, suite)
}

func chainId31336(t *hivesim.T, testnetInstance *testnet.Testnet) {
	_, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	l1_client := testnetInstance.L1ExecutionClients().Running()[0]

	client := &http.Client{}

	userRPCAddress, err := l1_client.UserRPCAddress()
	if err != nil {
		t.Fatal(err)
	}
	ethRpcClient, err := rpc.DialHTTPWithClient(userRPCAddress, client)
	if err != nil {
		t.Fatal(err)
	}

	chainIdString := ""
	err = ethRpcClient.Call(&chainIdString, "eth_chainId")
	if err != nil {
		t.Fatal(err)
	}

	chainId, err := strconv.ParseInt(chainIdString, 0, 0)
	if err != nil {
		fmt.Println("Error: %e", err)
		return
	}
	if chainId != 31336 {
		t.Fatalf("ChainId is not equal 31336, it is %i", chainId)
	}
}
