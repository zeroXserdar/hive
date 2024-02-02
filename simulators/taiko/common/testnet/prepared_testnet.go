package testnet

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/holiman/uint256"
	"github.com/protolambda/ztyp/tree"
	"github.com/protolambda/ztyp/view"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/configs"

	"github.com/ethereum/hive/hivesim"
	execution_client "github.com/taikoxyz/hive-taiko-clients/clients/execution"
	protocol_deployer_client "github.com/taikoxyz/hive-taiko-clients/clients/taiko/protocol_deployer"
	"github.com/taikoxyz/hive/simulators/taiko/common/clients"
	execution "github.com/taikoxyz/hive/simulators/taiko/common/config/execution"
)

var (
	depositAddress                              common.Eth1Address
	DEFAULT_SAFE_SLOTS_TO_IMPORT_OPTIMISTICALLY = big.NewInt(128)
	DEFAULT_MAX_CONSECUTIVE_ERRORS_ON_WAITS     = 3
)

func init() {
	_ = depositAddress.UnmarshalText(
		[]byte("0x4242424242424242424242424242424242424242"),
	)
}

// PreparedTestnet has all the options for starting nodes, ready to build the network.
type PreparedTestnet struct {
	// Consensus chain configuration
	//spec *common.Spec

	// Execution chain configuration and genesis info
	L1ExecutionClientGenesis *execution.ExecutionGenesis
	L2ExecutionClientGenesis *execution.ExecutionGenesis

	// Consensus genesis state
	//eth2Genesis common.BeaconState
	// Secret keys of validators, to fabricate extra signed test messages with during testnet/
	// E.g. to test a slashable offence that would not otherwise happen.
	//keys *[]blsu.SecretKey

	// Configuration to apply to every node of the given type
	L1ExecutionOpts          hivesim.StartOption
	L2ExecutionOpts          hivesim.StartOption
	L1L2ProtocolDeployerOpts hivesim.StartOption
	L2DriverOpts             hivesim.StartOption
	L2ProposerOpts           hivesim.StartOption
	L2ProverOpts             hivesim.StartOption

	// A tranche is a group of validator keys to run on 1 node
	//keyTranches []cl.ValidatorDetailsMap
}

// Prepares the fork timestamps of post-merge forks based on the
// consensus layer genesis time and the fork epochs
func prepareExecutionForkConfig(
	eth2GenesisTime common.Timestamp,
	config *Config,
) *params.ChainConfig {
	chainConfig := params.ChainConfig{}
	//if config.CapellaForkEpoch != nil {
	shanghai := uint64(eth2GenesisTime)
	//if config.CapellaForkEpoch.Uint64() != 0 {
	//	shanghai += uint64(
	//		config.CapellaForkEpoch.Int64() * config.SlotTime.Int64() * 32,
	//	)
	//}
	chainConfig.ShanghaiTime = &shanghai
	//}
	return &chainConfig
}

// Build all artifacts require to start a testnet.
func prepareTestnet(
	t *hivesim.T,
	env *Environment,
	config *Config,
) *PreparedTestnet {
	l1ExecutionClientGenesisTime := common.Timestamp(time.Now().Unix())
	l2ExecutionClientGenesisTime := l1ExecutionClientGenesisTime + 30

	// Sanitize configuration according to the clients used
	if err := config.fillDefaults(); err != nil {
		t.Fatal(fmt.Errorf("FAIL: error filling defaults: %v", err))
	}

	if configJson, err := json.MarshalIndent(config, "", "  "); err != nil {
		panic(err)
	} else {
		t.Logf("Testnet config: %s", configJson)
	}

	// Generate genesis for execution clients
	eth1Genesis := execution.BuildExecutionGenesis(
		config.TerminalTotalDifficulty,
		uint64(l1ExecutionClientGenesisTime),
		config.L1ExecutionConsensus,
		prepareExecutionForkConfig(l2ExecutionClientGenesisTime, config),
		config.GenesisExecutionAccounts,
	)
	if config.InitialBaseFeePerGas != nil {
		eth1Genesis.Genesis.BaseFee = config.InitialBaseFeePerGas
	}
	eth1ConfigOpt := eth1Genesis.ToParams(depositAddress)
	eth1Bundle, err := execution.ExecutionBundle(eth1Genesis.Genesis)
	if err != nil {
		t.Fatal(err)
	}
	execNodeOpts := hivesim.Params{
		"HIVE_LOGLEVEL": os.Getenv("HIVE_LOGLEVEL"),
		"HIVE_NODETYPE": "full",
	}
	jwtSecret := hivesim.Params{"HIVE_JWTSECRET": "true"}
	executionOpts := hivesim.Bundle(
		eth1ConfigOpt,
		eth1Bundle,
		execNodeOpts,
		jwtSecret,
	)

	// Pre-generate PoW chains for L1 clients that require it
	for i := 0; i < len(config.NodeDefinitions); i++ {
		if config.NodeDefinitions[i].L1ChainGenerator != nil {
			config.NodeDefinitions[i].L1Chain, err = config.NodeDefinitions[i].L1ChainGenerator.Generate(
				eth1Genesis,
			)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("Generated chain for node %d:\n", i+1)
			for j, b := range config.NodeDefinitions[i].L1Chain {
				js, _ := json.MarshalIndent(b.Header(), "", "  ")
				fmt.Printf("Block %d: %s\n", j, js)
			}
		}
	}

	// Generate beacon spec
	//
	// TODO: specify build-target based on preset, to run clients in mainnet or minimal mode.
	// copy the default mainnet config, and make some minimal modifications for testnet usage
	specCpy := *configs.Mainnet
	spec := &specCpy
	spec.Config.DEPOSIT_CONTRACT_ADDRESS = depositAddress
	spec.Config.DEPOSIT_CHAIN_ID = view.Uint64View(
		eth1Genesis.Genesis.Config.ChainID.Uint64(),
	)
	spec.Config.DEPOSIT_NETWORK_ID = view.Uint64View(eth1Genesis.NetworkID)
	spec.Config.ETH1_FOLLOW_DISTANCE = 1

	// Alter versions to avoid conflicts with mainnet values
	//spec.Config.GENESIS_FORK_VERSION = common.Version{0x00, 0x00, 0x00, 0x0a}
	//if config.AltairForkEpoch != nil {
	//	spec.Config.ALTAIR_FORK_EPOCH = common.Epoch(
	//		config.AltairForkEpoch.Uint64(),
	//	)
	//}
	//spec.Config.ALTAIR_FORK_VERSION = common.Version{0x01, 0x00, 0x00, 0x0a}
	//if config.BellatrixForkEpoch != nil {
	//	spec.Config.BELLATRIX_FORK_EPOCH = common.Epoch(
	//		config.BellatrixForkEpoch.Uint64(),
	//	)
	//}
	//spec.Config.BELLATRIX_FORK_VERSION = common.Version{0x02, 0x00, 0x00, 0x0a}
	//if config.CapellaForkEpoch != nil {
	//	spec.Config.CAPELLA_FORK_EPOCH = common.Epoch(
	//		config.CapellaForkEpoch.Uint64(),
	//	)
	//}
	//spec.Config.CAPELLA_FORK_VERSION = common.Version{0x03, 0x00, 0x00, 0x0a}
	//spec.Config.DENEB_FORK_VERSION = common.Version{0x04, 0x00, 0x00, 0x0a}
	//if config.ValidatorCount == nil {
	//	t.Fatal(fmt.Errorf("ValidatorCount was not configured"))
	//}
	//spec.Config.MIN_GENESIS_ACTIVE_VALIDATOR_COUNT = view.Uint64View(
	//	config.ValidatorCount.Uint64(),
	//)
	//if config.SlotTime != nil {
	//	spec.Config.SECONDS_PER_SLOT = common.Timestamp(
	//		config.SlotTime.Uint64(),
	//	)
	//}
	tdd, _ := uint256.FromBig(config.TerminalTotalDifficulty)
	spec.Config.TERMINAL_TOTAL_DIFFICULTY = view.Uint256View(*tdd)
	if execution.IsEth1GenesisPostMerge(eth1Genesis.Genesis) {
		genesisBlock := eth1Genesis.Genesis.ToBlock()
		spec.Config.TERMINAL_BLOCK_HASH = tree.Root(
			genesisBlock.Hash(),
		)
		spec.Config.TERMINAL_BLOCK_HASH_ACTIVATION_EPOCH = common.Timestamp(0)
	}

	// Validators can exit immediately
	spec.Config.SHARD_COMMITTEE_PERIOD = 0
	spec.Config.CHURN_LIMIT_QUOTIENT = 2

	// Validators can withdraw immediately
	spec.Config.MIN_VALIDATOR_WITHDRAWABILITY_DELAY = 0

	spec.Config.PROPOSER_SCORE_BOOST = 40

	// Generate keys opts for validators
	//shares := config.NodeDefinitions.Shares()
	// ExtraShares defines an extra set of keys that none of the nodes will have.
	// E.g. to produce an environment where none of the nodes has 50%+ of the keys.
	//if config.ExtraShares != nil {
	//	shares = append(shares, config.ExtraShares.Uint64())
	//}
	//keyTranches := cl.KeyTranches(env.Keys, shares)

	//consensusConfigOpts, err := cl.ConsensusConfigsBundle(
	//	spec,
	//	L1ExecutionClientGenesis.Genesis,
	//config.ValidatorCount.Uint64(),
	//)
	//if err != nil {
	//	t.Fatal(err)
	//}

	//var optimisticSync hivesim.Params
	//if config.SafeSlotsToImportOptimistically == nil {
	//	config.SafeSlotsToImportOptimistically = DEFAULT_SAFE_SLOTS_TO_IMPORT_OPTIMISTICALLY
	//}
	//optimisticSync = optimisticSync.Set(
	//	"HIVE_ETH2_SAFE_SLOTS_TO_IMPORT_OPTIMISTICALLY",
	//	fmt.Sprintf("%d", config.SafeSlotsToImportOptimistically),
	//)

	// prepare genesis beacon state, with all the validators in it.
	//state, err := cl.BuildBeaconState(
	//	spec,
	//	L1ExecutionClientGenesis.Genesis,
	//	eth2GenesisTime,
	//	env.Keys,
	//)
	//if err != nil {
	//	t.Fatal(err)
	//}

	// Write info so that the genesis state can be generated by the client
	//stateOpt, err := cl.StateBundle(state)
	//if err != nil {
	//	t.Fatal(err)
	//}

	// Define additional start options for beacon chain
	//commonOpts := hivesim.Params{
	//	"HIVE_ETH2_BN_API_PORT": fmt.Sprintf(
	//		"%d",
	//		beacon_client.PortBeaconAPI,
	//	),
	//	"HIVE_ETH2_BN_GRPC_PORT": fmt.Sprintf(
	//		"%d",
	//		beacon_client.PortBeaconGRPC,
	//	),
	//	"HIVE_ETH2_METRICS_PORT": fmt.Sprintf(
	//		"%d",
	//		beacon_client.PortMetrics,
	//	),
	//	"HIVE_ETH2_CONFIG_DEPOSIT_CONTRACT_ADDRESS": depositAddress.String(),
	//	"HIVE_ETH2_DEPOSIT_DEPLOY_BLOCK_HASH": fmt.Sprintf(
	//		"%s",
	//		L1ExecutionClientGenesis.Genesis.ToBlock().Hash(),
	//	),
	//}
	//beaconOpts := hivesim.Bundle(
	//	commonOpts,
	//	hivesim.Params{
	//		"HIVE_CHECK_LIVE_PORT": fmt.Sprintf(
	//			"%d",
	//			beacon_client.PortBeaconAPI,
	//		),
	//		"HIVE_ETH2_MERGE_ENABLED": "1",
	//		"HIVE_ETH2_ETH1_GENESIS_TIME": fmt.Sprintf(
	//			"%d",
	//			L1ExecutionClientGenesis.Genesis.Timestamp,
	//		),
	//		"HIVE_ETH2_GENESIS_FORK": config.activeFork(),
	//	},
	//	stateOpt,
	//	consensusConfigOpts,
	//	optimisticSync,
	//)
	//
	//validatorOpts := hivesim.Bundle(
	//	commonOpts,
	//	hivesim.Params{
	//		"HIVE_CHECK_LIVE_PORT": "0",
	//	},
	//	consensusConfigOpts,
	//)

	return &PreparedTestnet{
		//spec:          spec,
		L1ExecutionClientGenesis: eth1Genesis,
		//eth2Genesis:   state,
		//keys:          env.Secrets,
		L1ExecutionOpts: executionOpts,
		//beaconOpts:    beaconOpts,
		//validatorOpts: validatorOpts,
		//keyTranches:   keyTranches,
	}
}

func (p *PreparedTestnet) createTestnet(t *hivesim.T) *Testnet {
	genesisTime := common.Timestamp(p.L1ExecutionClientGenesis.Genesis.Timestamp)
	//genesisValidatorsRoot, _ := p.eth2Genesis.GenesisValidatorsRoot()
	return &Testnet{
		T:           t,
		genesisTime: genesisTime,
		//genesisValidatorsRoot: genesisValidatorsRoot,
		//spec:                  p.spec,
		L1ExecutionClientGenesis: p.L1ExecutionClientGenesis,
		//eth2GenesisState:         p.eth2Genesis,

		// Testing
		maxConsecutiveErrorsOnWaits: DEFAULT_MAX_CONSECUTIVE_ERRORS_ON_WAITS,
	}
}

// Prepares an L1 execution client object with all the necessary information
// to start
func (p *PreparedTestnet) prepareL1ExecutionNode(
	parentCtx context.Context,
	testnet *Testnet,
	eth1Def *hivesim.ClientDefinition,
	consensus execution.ExecutionConsensus,
	chain []*types.Block,
	config execution_client.ExecutionClientConfig,
) *execution_client.ExecutionClient {
	testnet.Logf(
		"Preparing execution node: %s (%s)",
		eth1Def.Name,
		eth1Def.Version,
	)

	cm := &clients.HiveManagedClient{
		T:                    testnet.T,
		HiveClientDefinition: eth1Def,
	}

	// This method will return the options used to run the client.
	// Requires a method that returns the rest of the currently running
	// execution clients on the network at startup.
	cm.OptionsGenerator = func() ([]hivesim.StartOption, error) {
		opts := []hivesim.StartOption{p.L1ExecutionOpts}
		opts = append(opts, consensus.HiveParams(config.ClientIndex))

		currentlyRunningEcs := testnet.L1ExecutionClients().
			Running().
			Subnet(config.Subnet)
		if len(currentlyRunningEcs) > 0 {
			bootnode, err := currentlyRunningEcs.Enodes()
			if err != nil {
				return nil, err
			}

			// Make the client connect to the first eth1 node, as a bootnode for the eth1 net
			opts = append(opts, hivesim.Params{"HIVE_BOOTNODE": bootnode})
		}
		opts = append(
			opts,
			hivesim.Params{
				"HIVE_TERMINAL_TOTAL_DIFFICULTY": fmt.Sprintf(
					"%d",
					config.TerminalTotalDifficulty,
				),
			},
		)
		genesis := testnet.ExecutionGenesis().ToBlock()
		if config.TerminalTotalDifficulty <= genesis.Difficulty().Int64() {
			opts = append(
				opts,
				hivesim.Params{
					"HIVE_TERMINAL_BLOCK_HASH": fmt.Sprintf(
						"%s",
						genesis.Hash(),
					),
				},
			)
			opts = append(
				opts,
				hivesim.Params{
					"HIVE_TERMINAL_BLOCK_NUMBER": fmt.Sprintf(
						"%d",
						genesis.NumberU64(),
					),
				},
			)
		}

		if len(chain) > 0 {
			// Bundle the chain into the container
			chainParam, err := execution.ChainBundle(chain)
			if err != nil {
				return nil, err
			}
			opts = append(opts, chainParam)
		}
		return opts, nil
	}

	testnet.Logf(
		"Finished preparing execution node: %s (%s)",
		eth1Def.Name,
		eth1Def.Version,
	)

	return &execution_client.ExecutionClient{
		Client: cm,
		Logger: testnet.T,
		Config: config,
	}
}

// Prepares an L2 execution client object with all the necessary information
// to start
func (p *PreparedTestnet) prepareL2ExecutionNode(
	parentCtx context.Context,
	testnet *Testnet,
	eth1Def *hivesim.ClientDefinition,
	consensus execution.ExecutionConsensus,
	chain []*types.Block,
	config execution_client.ExecutionClientConfig,
) *execution_client.ExecutionClient {
	testnet.Logf(
		"Preparing execution node: %s (%s)",
		eth1Def.Name,
		eth1Def.Version,
	)

	cm := &clients.HiveManagedClient{
		T:                    testnet.T,
		HiveClientDefinition: eth1Def,
	}

	// This method will return the options used to run the client.
	// Requires a method that returns the rest of the currently running
	// execution clients on the network at startup.
	cm.OptionsGenerator = func() ([]hivesim.StartOption, error) {
		opts := []hivesim.StartOption{p.L1ExecutionOpts}
		opts = append(opts, consensus.HiveParams(config.ClientIndex))

		currentlyRunningEcs := testnet.L2ExecutionClients().
			Running().
			Subnet(config.Subnet)
		if len(currentlyRunningEcs) > 0 {
			bootnode, err := currentlyRunningEcs.Enodes()
			if err != nil {
				return nil, err
			}

			// Make the client connect to the first eth1 node, as a bootnode for the eth1 net
			opts = append(opts, hivesim.Params{"HIVE_BOOTNODE": bootnode})
		}
		opts = append(
			opts,
			hivesim.Params{
				"HIVE_TERMINAL_TOTAL_DIFFICULTY": fmt.Sprintf(
					"%d",
					config.TerminalTotalDifficulty,
				),
			},
		)
		genesis := testnet.ExecutionGenesis().ToBlock()
		if config.TerminalTotalDifficulty <= genesis.Difficulty().Int64() {
			opts = append(
				opts,
				hivesim.Params{
					"HIVE_TERMINAL_BLOCK_HASH": fmt.Sprintf(
						"%s",
						genesis.Hash(),
					),
				},
			)
			opts = append(
				opts,
				hivesim.Params{
					"HIVE_TERMINAL_BLOCK_NUMBER": fmt.Sprintf(
						"%d",
						genesis.NumberU64(),
					),
				},
			)
		}

		if len(chain) > 0 {
			// Bundle the chain into the container
			chainParam, err := execution.ChainBundle(chain)
			if err != nil {
				return nil, err
			}
			opts = append(opts, chainParam)
		}
		return opts, nil
	}

	testnet.Logf(
		"Finished preparing execution node: %s (%s)",
		eth1Def.Name,
		eth1Def.Version,
	)

	return &execution_client.ExecutionClient{
		Client: cm,
		Logger: testnet.T,
		Config: config,
	}
}

func (p *PreparedTestnet) prepareL1L2ProtocolDeployerNode(
	parentCtx context.Context,
	testnet *Testnet,
	protocolDeployerDef *hivesim.ClientDefinition,
	config protocol_deployer_client.ProtocolDeployerClientConfig,
	l1ExecutionClientEndpoint *execution_client.ExecutionClient,
	l2ExecutionClientEndpoint *execution_client.ExecutionClient,
) *protocol_deployer_client.ProtocolDeployerClient {
	testnet.Logf(
		"Preparing protocol deployer node: %s (%s)",
		protocolDeployerDef.Name,
		protocolDeployerDef.Version,
	)

	if l1ExecutionClientEndpoint == nil && l2ExecutionClientEndpoint == nil {
		panic(fmt.Errorf("at least 1 (l1/l2) execution client endpoint is required"))
	}

	cm := &clients.HiveManagedClient{
		T:                    testnet.T,
		HiveClientDefinition: protocolDeployerDef,
	}

	cl := &protocol_deployer_client.ProtocolDeployerClient{
		Client: cm,
		Logger: testnet.T,
		Config: config,
	}

	//if enableBuilders {
	//	simIP, err := testnet.T.Sim.ContainerNetworkIP(
	//		testnet.T.SuiteID,
	//		"bridge",
	//		"simulation",
	//	)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	options := []mock_builder.Option{
	//		mock_builder.WithExternalIP(net.ParseIP(simIP)),
	//		mock_builder.WithPort(
	//			mock_builder.DEFAULT_BUILDER_PORT + config.ClientIndex,
	//		),
	//		mock_builder.WithID(config.ClientIndex),
	//		mock_builder.WithBeaconGenesisTime(testnet.genesisTime),
	//		mock_builder.WithSpec(p.spec),
	//	}
	//
	//	if builderOptions != nil {
	//		options = append(options, builderOptions...)
	//	}
	//
	//	cl.Builder, err = mock_builder.NewMockBuilder(
	//		context.Background(),
	//		eth1Endpoints[0],
	//		cl,
	//		options...,
	//	)
	//	if err != nil {
	//		panic(err)
	//	}
	//}

	// This method will return the options used to run the client.
	// Requires a method that returns the rest of the currently running
	// beacon clients on the network at startup.
	cm.OptionsGenerator = func() ([]hivesim.StartOption, error) {
		opts := []hivesim.StartOption{p.L1L2ProtocolDeployerOpts}

		var deploymentTargetExecutionClient execution_client.ExecutionClient
		if l1ExecutionClientEndpoint != nil {
			deploymentTargetExecutionClient = *l1ExecutionClientEndpoint
		} else {
			deploymentTargetExecutionClient = *l2ExecutionClientEndpoint
		}

		if !deploymentTargetExecutionClient.IsRunning() || deploymentTargetExecutionClient.Proxy() == nil {
			return nil, fmt.Errorf(
				"attempted to start protocol deployment node when the execution client is not yet running",
			)
		}
		execNode := deploymentTargetExecutionClient.Proxy()
		userRPC, err := execNode.UserRPCAddress()
		if err != nil {
			return nil, fmt.Errorf(
				"execution client node used for protocol deployment without available RPC: %v",
				err,
			)
		}

		opts = append(opts, hivesim.Params{
			"MAINNET_URL": userRPC,
		})

		//opts = append(
		//	opts,
		//	hivesim.Params{
		//		"HIVE_TERMINAL_TOTAL_DIFFICULTY": fmt.Sprintf(
		//			"%d",
		//			config.TerminalTotalDifficulty,
		//		),
		//	},
		//)

		return opts, nil
	}

	return cl
}

// Prepares a beacon client object with all the necessary information
// to start
//func (p *PreparedTestnet) prepareBeaconNode(
//	parentCtx context.Context,
//	testnet *Testnet,
//	beaconDef *hivesim.ClientDefinition,
//	enableBuilders bool,
//	builderOptions []mock_builder.Option,
//	config beacon_client.BeaconClientConfig,
//	eth1Endpoints ...*execution_client.ExecutionClient,
//) *beacon_client.BeaconClient {
//	testnet.Logf(
//		"Preparing beacon node: %s (%s)",
//		beaconDef.Name,
//		beaconDef.Version,
//	)
//
//	if len(eth1Endpoints) == 0 {
//		panic(fmt.Errorf("at least 1 execution endpoint is required"))
//	}
//
//	cm := &prot_depl_clients.HiveManagedClient{
//		T:                    testnet.T,
//		HiveClientDefinition: beaconDef,
//	}
//
//	cl := &beacon_client.BeaconClient{
//		Client: cm,
//		Logger: testnet.T,
//		Config: config,
//	}
//
//	if enableBuilders {
//		simIP, err := testnet.T.Sim.ContainerNetworkIP(
//			testnet.T.SuiteID,
//			"bridge",
//			"simulation",
//		)
//		if err != nil {
//			panic(err)
//		}
//
//		options := []mock_builder.Option{
//			mock_builder.WithExternalIP(net.ParseIP(simIP)),
//			mock_builder.WithPort(
//				mock_builder.DEFAULT_BUILDER_PORT + config.ClientIndex,
//			),
//			mock_builder.WithID(config.ClientIndex),
//			mock_builder.WithBeaconGenesisTime(testnet.genesisTime),
//			mock_builder.WithSpec(p.spec),
//		}
//
//		if builderOptions != nil {
//			options = append(options, builderOptions...)
//		}
//
//		cl.Builder, err = mock_builder.NewMockBuilder(
//			context.Background(),
//			eth1Endpoints[0],
//			cl,
//			options...,
//		)
//		if err != nil {
//			panic(err)
//		}
//	}
//
//	// This method will return the options used to run the client.
//	// Requires a method that returns the rest of the currently running
//	// beacon clients on the network at startup.
//	cm.OptionsGenerator = func() ([]hivesim.StartOption, error) {
//		opts := []hivesim.StartOption{p.beaconOpts}
//
//		// Hook up beacon node to (maybe multiple) eth1 nodes
//		var addrs []string
//		var engineAddrs []string
//		for _, en := range eth1Endpoints {
//			if !en.IsRunning() || en.Proxy() == nil {
//				return nil, fmt.Errorf(
//					"attempted to start beacon node when the execution client is not yet running",
//				)
//			}
//			execNode := en.Proxy()
//			userRPC, err := execNode.UserRPCAddress()
//			if err != nil {
//				return nil, fmt.Errorf(
//					"eth1 node used for beacon without available RPC: %v",
//					err,
//				)
//			}
//			addrs = append(addrs, userRPC)
//			engineRPC, err := execNode.EngineRPCAddress()
//			if err != nil {
//				return nil, fmt.Errorf(
//					"eth1 node used for beacon without available RPC: %v",
//					err,
//				)
//			}
//			engineAddrs = append(engineAddrs, engineRPC)
//		}
//		opts = append(opts, hivesim.Params{
//			"HIVE_ETH2_ETH1_RPC_ADDRS":        strings.Join(addrs, ","),
//			"HIVE_ETH2_ETH1_ENGINE_RPC_ADDRS": strings.Join(engineAddrs, ","),
//			"HIVE_ETH2_BEACON_NODE_INDEX": fmt.Sprintf(
//				"%d",
//				config.ClientIndex,
//			),
//		})
//
//		currentlyRunningBcs := testnet.BeaconClients().
//			Running().
//			Subnet(config.Subnet)
//		if len(currentlyRunningBcs) > 0 {
//			if bootnodeENRs, err := currentlyRunningBcs.ENRs(parentCtx); err != nil {
//				return nil, fmt.Errorf(
//					"failed to get ENR as bootnode for every beacon node: %v",
//					err,
//				)
//			} else if bootnodeENRs != "" {
//				opts = append(opts, hivesim.Params{"HIVE_ETH2_BOOTNODE_ENRS": bootnodeENRs})
//			}
//
//			if staticPeers, err := currentlyRunningBcs.P2PAddrs(parentCtx); err != nil {
//				return nil, fmt.Errorf(
//					"failed to get p2paddr for every beacon node: %v",
//					err,
//				)
//			} else if staticPeers != "" {
//				opts = append(opts, hivesim.Params{"HIVE_ETH2_STATIC_PEERS": staticPeers})
//			}
//		}
//
//		opts = append(
//			opts,
//			hivesim.Params{
//				"HIVE_TERMINAL_TOTAL_DIFFICULTY": fmt.Sprintf(
//					"%d",
//					config.TerminalTotalDifficulty,
//				),
//			},
//		)
//
//		if cl.Builder != nil {
//			if builder, ok := cl.Builder.(builder_types.Builder); ok {
//				opts = append(opts, hivesim.Params{
//					"HIVE_ETH2_BUILDER_ENDPOINT": builder.Address(),
//				})
//			}
//		}
//
//		// TODO
//		//if p.configName != "mainnet" && hasBuildTarget(beaconDef, p.configName) {
//		//	opts = append(opts, hivesim.WithBuildTarget(p.configName))
//		//}
//
//		return opts, nil
//	}
//
//	return cl
//}

// Prepares a validator client object with all the necessary information
// to eventually start the client.
//func (p *PreparedTestnet) prepareValidatorClient(
//	parentCtx context.Context,
//	testnet *Testnet,
//	validatorDef *hivesim.ClientDefinition,
//	bn *beacon_client.BeaconClient,
//	keyIndex int,
//) *validator_client.ValidatorClient {
//	testnet.Logf(
//		"Preparing validator client: %s (%s)",
//		validatorDef.Name,
//		validatorDef.Version,
//	)
//	if keyIndex >= len(p.keyTranches) {
//		testnet.Fatalf(
//			"only have %d key tranches, cannot find index %d for VC",
//			len(p.keyTranches),
//			keyIndex,
//		)
//	}
//	keys := p.keyTranches[keyIndex]
//
//	cm := &prot_depl_clients.HiveManagedClient{
//		T:                    testnet.T,
//		HiveClientDefinition: validatorDef,
//	}
//
//	// This method will return the options used to run the client.
//	// Requires the beacon client object to which to connect.
//	cm.OptionsGenerator = func() ([]hivesim.StartOption, error) {
//		if !bn.IsRunning() {
//			return nil, fmt.Errorf(
//				"attempted to start a validator when the beacon node is not running",
//			)
//		}
//		// Hook up validator to beacon node
//		bnAPIOpt := hivesim.Params{
//			"HIVE_ETH2_BN_API_IP": bn.GetIP().String(),
//		}
//		keysOpt := cl.KeysBundle(keys)
//		opts := []hivesim.StartOption{p.validatorOpts, keysOpt, bnAPIOpt}
//
//		if bn.Builder != nil {
//			if builder, ok := bn.Builder.(builder_types.Builder); ok {
//				opts = append(opts, hivesim.Params{
//					"HIVE_ETH2_BUILDER_ENDPOINT": builder.Address(),
//				})
//			}
//		}
//
//		// TODO
//		//if p.configName != "mainnet" && hasBuildTarget(validatorDef, p.configName) {
//		//	opts = append(opts, hivesim.WithBuildTarget(p.configName))
//		//}
//		return opts, nil
//	}
//
//	return &validator_client.ValidatorClient{
//		Client:       cm,
//		Logger:       testnet.T,
//		ClientIndex:  keyIndex,
//		Keys:         keys.Keys(),
//		BeaconClient: bn,
//	}
//}
