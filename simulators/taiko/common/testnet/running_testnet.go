package testnet

import (
	//"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/kr/pretty"
	"github.com/taikoxyz/hive-taiko-clients/clients/taiko/driver"
	"github.com/taikoxyz/hive-taiko-clients/clients/taiko/proposer"
	"github.com/taikoxyz/hive-taiko-clients/clients/taiko/prover"
	"math/big"
	"net"
	//"sync"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	//"github.com/pkg/errors"

	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon/altair"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	//"github.com/protolambda/zrnt/eth2/beacon/phase0"
	"github.com/protolambda/zrnt/eth2/util/math"
	"github.com/protolambda/ztyp/tree"

	"github.com/ethereum/hive/hivesim"
	execution_config "github.com/taikoxyz/hive/simulators/taiko/common/config/execution"

	//beacon_client "github.com/taikoxyz/hive-taiko-clients/clients/beacon"
	exec_client "github.com/taikoxyz/hive-taiko-clients/clients/execution"
	"github.com/taikoxyz/hive-taiko-clients/clients/taiko/node"
	protocol_deployer_client "github.com/taikoxyz/hive-taiko-clients/clients/taiko/protocol_deployer"
)

const (
	MAX_PARTICIPATION_SCORE = 7
)

var (
	EMPTY_EXEC_HASH = ethcommon.Hash{}
	EMPTY_TREE_ROOT = tree.Root{}
	JWT_SECRET, _   = hex.DecodeString(
		"7365637265747365637265747365637265747365637265747365637265747365",
	)
)

type Testnet struct {
	*hivesim.T
	node.TaikoNodes

	genesisTime           common.Timestamp
	genesisValidatorsRoot common.Root

	// Consensus chain configuration
	spec *common.Spec
	// Execution chain configuration and genesis info
	L1ExecutionClientGenesis *execution_config.ExecutionGenesis
	L2ExecutionClientGenesis *execution_config.ExecutionGenesis
	// Consensus genesis state
	//eth2GenesisState common.BeaconState

	// Test configuration
	maxConsecutiveErrorsOnWaits int
}

type ActiveSpec struct {
	*common.Spec
}

const slotsTolerance common.Slot = 2

func (spec *ActiveSpec) EpochTimeoutContext(
	parent context.Context,
	epochs common.Epoch,
) (context.Context, context.CancelFunc) {
	return context.WithTimeout(
		parent,
		time.Duration(
			uint64((spec.SLOTS_PER_EPOCH*common.Slot(epochs))+slotsTolerance)*
				uint64(spec.SECONDS_PER_SLOT),
		)*time.Second,
	)
}

func (spec *ActiveSpec) SlotTimeoutContext(
	parent context.Context,
	slots common.Slot,
) (context.Context, context.CancelFunc) {
	return context.WithTimeout(
		parent,
		time.Duration(
			uint64(slots+slotsTolerance)*
				uint64(spec.SECONDS_PER_SLOT))*time.Second,
	)
}

func (spec *ActiveSpec) EpochsTimeout(epochs common.Epoch) <-chan time.Time {
	return time.After(
		time.Duration(
			uint64(
				spec.SLOTS_PER_EPOCH*common.Slot(epochs),
			)*uint64(
				spec.SECONDS_PER_SLOT,
			),
		) * time.Second,
	)
}

func (spec *ActiveSpec) SlotsTimeout(slots common.Slot) <-chan time.Time {
	return time.After(
		time.Duration(
			uint64(slots)*uint64(spec.SECONDS_PER_SLOT),
		) * time.Second,
	)
}

func (t *Testnet) Spec() *ActiveSpec {
	return &ActiveSpec{
		Spec: t.spec,
	}
}

func (t *Testnet) GenesisTime() common.Timestamp {
	// return time.Unix(int64(t.genesisTime), 0)
	return t.genesisTime
}

func (t *Testnet) GenesisTimeUnix() time.Time {
	return time.Unix(int64(t.genesisTime), 0)
}

//func (t *Testnet) GenesisBeaconState() common.BeaconState {
//	return t.eth2GenesisState
//}

func (t *Testnet) GenesisValidatorsRoot() common.Root {
	return t.genesisValidatorsRoot
}

func (t *Testnet) ExecutionGenesis() *core.Genesis {
	return t.L1ExecutionClientGenesis.Genesis
}

func StartTestnet(
	parentCtx context.Context,
	t *hivesim.T,
	env *Environment,
	config *Config,
) *Testnet {
	var (
		prep        = prepareTestnet(t, env, config)
		testnet     = prep.createTestnet(t)
		genesisTime = testnet.GenesisTimeUnix()
	)
	t.Logf("Config: %+v", fmt.Sprintf("%# v", pretty.Formatter(config)))
	t.Logf("PreparedTestnet: %+v", fmt.Sprintf("%# v", pretty.Formatter(prep)))
	t.Logf("created Testnet: %+v", fmt.Sprintf("%# v", pretty.Formatter(testnet)))
	t.Logf(
		"Created new testnet, genesis at %s (%s from now)",
		genesisTime,
		time.Until(genesisTime),
	)

	var simulatorIP net.IP
	if simIPStr, err := t.Sim.ContainerNetworkIP(
		testnet.T.SuiteID,
		"bridge",
		"simulation",
	); err != nil {
		panic(err)
	} else {
		simulatorIP = net.ParseIP(simIPStr)
	}

	testnet.TaikoNodes = make(node.TaikoNodes, len(config.NodeDefinitions))

	// Init all client bundles
	for nodeIndex := range testnet.TaikoNodes {
		testnet.TaikoNodes[nodeIndex] = new(node.TaikoNode)
	}

	// For each key partition, we start a client bundle that consists of:
	// - 1 execution client
	// - 1 beacon client
	// - 1 validator client,
	for nodeIndex, node := range config.NodeDefinitions {
		// Prepare clients for this node
		var (
			nodeClient = testnet.TaikoNodes[nodeIndex]

			L1executionDef = env.Clients.ClientByNameAndRole(
				node.L1ExecutionClientName(),
				"l1_client",
			)
			L2executionDef = env.Clients.ClientByNameAndRole(
				node.L2ExecutionClientName(),
				"l2_client",
			)
			protocolDeployerDef = env.Clients.ClientByNameAndRole(
				node.L1L2ProtocolDeployerClientName(),
				"l1l2_protocol_deployer",
			)
			driverDef = env.Clients.ClientByNameAndRole(
				node.L2DriverClientName(),
				"l2_driver",
			)
			proposerDef = env.Clients.ClientByNameAndRole(
				node.L2ProposerClientName(),
				"l2_proposer",
			)
			proverDef = env.Clients.ClientByNameAndRole(
				node.L2ProverClientName(),
				"l2_prover",
			)
			//beaconDef = env.Clients.ClientByNameAndRole(
			//	node.ConsensusClientName(),
			//	"beacon",
			//)
			//validatorDef = env.Clients.ClientByNameAndRole(
			//	node.ValidatorClientName(),
			//	"validator",
			//)
			executionTTD = int64(0)
			//beaconTTD    = int64(0)
		)

		if L1executionDef == nil {
			//if L1executionDef == nil || beaconDef == nil || validatorDef == nil {
			t.Fatalf("FAIL: Unable to get L1 execution client")
		}
		if node.L1ExecutionClientTTD != nil {
			executionTTD = node.L1ExecutionClientTTD.Int64()
		} else if testnet.L1ExecutionClientGenesis.Genesis.Config.TerminalTotalDifficulty != nil {
			executionTTD = testnet.L1ExecutionClientGenesis.Genesis.Config.TerminalTotalDifficulty.Int64()
		}
		//if node.BeaconNodeTTD != nil {
		//	beaconTTD = node.BeaconNodeTTD.Int64()
		//} else if testnet.L1ExecutionClientGenesis.Genesis.Config.TerminalTotalDifficulty != nil {
		//	beaconTTD = testnet.L1ExecutionClientGenesis.Genesis.Config.TerminalTotalDifficulty.Int64()
		//}

		// Prepare the client objects with all the information necessary to
		// eventually start
		nodeClient.L1ExecutionClient = prep.prepareL1ExecutionNode(
			parentCtx,
			testnet,
			L1executionDef,
			config.L1ExecutionConsensus,
			node.L1Chain,
			exec_client.ExecutionClientConfig{
				ClientIndex:             nodeIndex,
				TerminalTotalDifficulty: executionTTD,
				Subnet:                  node.GetL1Subnet(),
				JWTSecret:               JWT_SECRET,
				ProxyConfig: &exec_client.ExecutionProxyConfig{
					Host:                   simulatorIP,
					Port:                   exec_client.PortEngineRPC + nodeIndex,
					TrackForkchoiceUpdated: true,
					LogEngineCalls:         env.LogEngineCalls,
				},
			},
		)

		if L2executionDef == nil {
			//if L2executionDef == nil || beaconDef == nil || validatorDef == nil {
			t.Fatalf("FAIL: Unable to get L2 execution client")
		}
		//if node.L2ExecutionClientTTD != nil {
		//	executionTTD = node.L2ExecutionClientTTD.Int64()
		//} else if testnet.L2ExecutionClientGenesis.Genesis.Config.TerminalTotalDifficulty != nil {
		//	executionTTD = testnet.L2ExecutionClientGenesis.Genesis.Config.TerminalTotalDifficulty.Int64()
		//}
		//if node.BeaconNodeTTD != nil {
		//	beaconTTD = node.BeaconNodeTTD.Int64()
		//} else if testnet.L2ExecutionClientGenesis.Genesis.Config.TerminalTotalDifficulty != nil {
		//	beaconTTD = testnet.L2ExecutionClientGenesis.Genesis.Config.TerminalTotalDifficulty.Int64()
		//}

		// Prepare the client objects with all the information necessary to
		// eventually start
		nodeClient.L2ExecutionClient = prep.prepareL2ExecutionNode(
			parentCtx,
			testnet,
			L2executionDef,
			config.L2ExecutionConsensus,
			node.L2Chain,
			exec_client.ExecutionClientConfig{
				ClientIndex:             nodeIndex,
				TerminalTotalDifficulty: executionTTD,
				Subnet:                  node.GetL2Subnet(),
				JWTSecret:               JWT_SECRET,
				ProxyConfig: &exec_client.ExecutionProxyConfig{
					Host:                   simulatorIP,
					Port:                   exec_client.PortEngineRPC + nodeIndex,
					TrackForkchoiceUpdated: true,
					LogEngineCalls:         env.LogEngineCalls,
				},
			},
			nodeClient.L1ExecutionClient,
		)

		if node.L1L2ProtocolDeployerClient != "" && protocolDeployerDef == nil {
			t.Fatalf("FAIL: Unable to get protocol deployer client")
		}

		if node.L1L2ProtocolDeployerClient != "" {
			nodeClient.L1L2ProtocolDeployerClient = prep.prepareL1L2ProtocolDeployerNode(
				parentCtx,
				testnet,
				protocolDeployerDef,
				protocol_deployer_client.ProtocolDeployerClientConfig{
					ClientIndex: nodeIndex,
					Subnet:      node.GetL1Subnet(),
				},
				nodeClient.L1ExecutionClient,
				nil,
			)
		}

		if node.L2DriverClient != "" && driverDef == nil {
			t.Fatalf("FAIL: Unable to get driver client")
		}

		if node.L2DriverClient != "" {
			nodeClient.L2DriverClient = prep.prepareL2DriverNode(
				parentCtx,
				testnet,
				driverDef,
				driver.DriverClientConfig{
					ClientIndex: nodeIndex,
					Subnet:      node.GetL2Subnet(),
				},
				nodeClient.L1ExecutionClient,
				nodeClient.L2ExecutionClient,
			)
		}

		if node.L2ProposerClient != "" && proposerDef == nil {
			t.Fatalf("FAIL: Unable to get proposer client")
		}

		if node.L2ProposerClient != "" {
			nodeClient.L2ProposerClient = prep.prepareL2ProposerNode(
				parentCtx,
				testnet,
				proposerDef,
				proposer.ProposerClientConfig{
					ClientIndex: nodeIndex,
					Subnet:      node.GetL2Subnet(),
				},
				nodeClient.L1ExecutionClient,
				nodeClient.L2ExecutionClient,
			)
		}

		if node.L2ProverClient != "" && proverDef == nil {
			t.Fatalf("FAIL: Unable to get prover client")
		}

		if node.L2ProverClient != "" {
			nodeClient.L2ProverClient = prep.prepareL2ProverNode(
				parentCtx,
				testnet,
				proverDef,
				prover.ProverClientConfig{
					ClientIndex: nodeIndex,
					Subnet:      node.GetL2Subnet(),
				},
				nodeClient.L1ExecutionClient,
				nodeClient.L2ExecutionClient,
			)
		}
		// Add rest of properties
		nodeClient.Logging = t
		nodeClient.Index = nodeIndex
		//nodeClient.Verification = node.TestVerificationNode
		// Start the node clients if specified so
		//if !node.DisableStartup {
		//jsonData, err := json.Marshal(nodeClient.L1ExecutionClient.)
		//if err != nil {
		//	t.Fatal("JSON marshal error", err)
		//}
		//t.Logf("NodeClient: %v", string(jsonData))
		if err := nodeClient.Start(); err != nil {
			t.Fatalf("FAIL: Unable to start node %d: %v", nodeIndex, err)
		}
		//}

		//t.nodeClient.L1L2ProtocolDeployerClient.(clients.HiveManagedClient).HiveClient.Container
	}

	return testnet
}

//TODO:
//func (t *Testnet) Stop() {
//	for _, p := range t.Proxies().Running() {
//		p.Cancel()
//	}
//	for _, b := range t.BeaconClients() {
//		if b.Builder != nil {
//			if builder, ok := b.Builder.(builder_types.Builder); ok {
//				builder.Cancel()
//			}
//		}
//	}
//}

//func (t *Testnet) ValidatorClientIndex(pk [48]byte) (int, error) {
//	for i, v := range t.ValidatorClients() {
//		if v.ContainsKey(pk) {
//			return i, nil
//		}
//	}
//	return 0, fmt.Errorf("key not found in any validator client")
//}

// Wait until the beacon chain genesis happens.
func (t *Testnet) WaitForGenesis(ctx context.Context) {
	genesis := t.GenesisTimeUnix()
	select {
	case <-ctx.Done():
	case <-time.After(time.Until(genesis)):
	}
}

// Wait a certain amount of slots while printing the current status.
//func (t *Testnet) WaitSlots(ctx context.Context, slots common.Slot) error {
//	for s := common.Slot(0); s < slots; s++ {
//		t.BeaconClients().Running().PrintStatus(ctx)
//		select {
//		case <-time.After(time.Duration(t.spec.SECONDS_PER_SLOT) * time.Second):
//		case <-ctx.Done():
//			return ctx.Err()
//		}
//	}
//	return nil
//}

// WaitForFork blocks until a beacon client reaches specified fork,
// or context finalizes, whichever happens first.
//func (t *Testnet) WaitForFork(ctx context.Context, fork string) error {
//	var (
//		genesis      = t.GenesisTimeUnix()
//		slotDuration = time.Duration(t.spec.SECONDS_PER_SLOT) * time.Second
//		timer        = time.NewTicker(slotDuration)
//		runningNodes = t.VerificationNodes().Running()
//		results      = makeResults(runningNodes, t.maxConsecutiveErrorsOnWaits)
//	)
//
//	for {
//		select {
//		case <-ctx.Done():
//			return ctx.Err()
//		case tim := <-timer.C:
//			// start polling after first slot of genesis
//			if tim.Before(genesis.Add(slotDuration)) {
//				t.Logf("Time till genesis: %s", genesis.Sub(tim))
//				continue
//			}
//
//			// new slot, log and check status of all beacon nodes
//			var (
//				wg        sync.WaitGroup
//				clockSlot = t.spec.TimeToSlot(
//					common.Timestamp(time.Now().Unix()),
//					t.GenesisTime(),
//				)
//			)
//			results.Clear()
//
//			for i, n := range runningNodes {
//				wg.Add(1)
//				go func(
//					ctx context.Context,
//					n *node.TaikoNode,
//					r *result,
//				) {
//					defer wg.Done()
//
//					b := n.BeaconClient
//
//					checkpoints, err := b.BlockFinalityCheckpoints(
//						ctx,
//						eth2api.BlockHead,
//					)
//					if err != nil {
//						r.err = errors.Wrap(
//							err,
//							"failed to poll finality checkpoint",
//						)
//						return
//					}
//
//					versionedBlock, err := b.BlockV2(
//						ctx,
//						eth2api.BlockHead,
//					)
//					if err != nil {
//						r.err = errors.Wrap(err, "failed to retrieve block")
//						return
//					}
//
//					execution := ethcommon.Hash{}
//					if executionPayload, err := versionedBlock.ExecutionPayload(); err == nil {
//						execution = executionPayload.BlockHash
//					}
//
//					slot := versionedBlock.Slot()
//					if clockSlot > slot &&
//						(clockSlot-slot) >= t.spec.SLOTS_PER_EPOCH {
//						r.fatal = fmt.Errorf(
//							"unable to sync for an entire epoch: clockSlot=%d, slot=%d",
//							clockSlot,
//							slot,
//						)
//						return
//					}
//
//					r.msg = fmt.Sprintf(
//						"fork=%s, clock_slot=%s, slot=%d, head=%s, exec_payload=%s, justified=%s, finalized=%s",
//						versionedBlock.Version,
//						clockSlot,
//						slot,
//						utils.Shorten(versionedBlock.Root().String()),
//						utils.Shorten(execution.String()),
//						utils.Shorten(checkpoints.CurrentJustified.String()),
//						utils.Shorten(checkpoints.Finalized.String()),
//					)
//
//					if versionedBlock.Version == fork {
//						r.done = true
//					}
//				}(ctx, n, results[i])
//			}
//			wg.Wait()
//
//			if err := results.CheckError(); err != nil {
//				return err
//			}
//			results.PrintMessages(t.Logf)
//			if results.AllDone() {
//				return nil
//			}
//		}
//	}
//}

// WaitForFinality blocks until a beacon client reaches finality,
// or timeoutSlots have passed, whichever happens first.
//func (t *Testnet) WaitForFinality(ctx context.Context) (
//	common.Checkpoint, error,
//) {
//	var (
//		genesis      = t.GenesisTimeUnix()
//		slotDuration = time.Duration(t.spec.SECONDS_PER_SLOT) * time.Second
//		timer        = time.NewTicker(slotDuration)
//		runningNodes = t.VerificationNodes().Running()
//		results      = makeResults(runningNodes, t.maxConsecutiveErrorsOnWaits)
//	)
//
//	for {
//		select {
//		case <-ctx.Done():
//			return common.Checkpoint{}, ctx.Err()
//		case tim := <-timer.C:
//			// start polling after first slot of genesis
//			if tim.Before(genesis.Add(slotDuration)) {
//				t.Logf("Time till genesis: %s", genesis.Sub(tim))
//				continue
//			}
//
//			// new slot, log and check status of all beacon nodes
//			var (
//				wg        sync.WaitGroup
//				clockSlot = t.spec.TimeToSlot(
//					common.Timestamp(time.Now().Unix()),
//					t.GenesisTime(),
//				)
//			)
//			results.Clear()
//
//			for i, n := range runningNodes {
//				wg.Add(1)
//				go func(ctx context.Context, n *node.TaikoNode, r *result) {
//					defer wg.Done()
//
//					b := n.BeaconClient
//
//					checkpoints, err := b.BlockFinalityCheckpoints(
//						ctx,
//						eth2api.BlockHead,
//					)
//					if err != nil {
//						r.err = errors.Wrap(
//							err,
//							"failed to poll finality checkpoint",
//						)
//						return
//					}
//
//					versionedBlock, err := b.BlockV2(
//						ctx,
//						eth2api.BlockHead,
//					)
//					if err != nil {
//						r.err = errors.Wrap(err, "failed to retrieve block")
//						return
//					}
//					execution := ethcommon.Hash{}
//					if executionPayload, err := versionedBlock.ExecutionPayload(); err == nil {
//						execution = executionPayload.BlockHash
//					}
//
//					slot := versionedBlock.Slot()
//					if clockSlot > slot &&
//						(clockSlot-slot) >= t.spec.SLOTS_PER_EPOCH {
//						r.fatal = fmt.Errorf(
//							"unable to sync for an entire epoch: clockSlot=%d, slot=%d",
//							clockSlot,
//							slot,
//						)
//						return
//					}
//
//					health, _ := GetHealth(ctx, b, t.spec, slot)
//
//					r.msg = fmt.Sprintf(
//						"fork=%s, clock_slot=%d, slot=%d, head=%s, "+
//							"health=%.2f, exec_payload=%s, justified=%s, "+
//							"finalized=%s",
//						versionedBlock.Version,
//						clockSlot,
//						slot,
//						utils.Shorten(versionedBlock.Root().String()),
//						health,
//						utils.Shorten(execution.String()),
//						utils.Shorten(checkpoints.CurrentJustified.String()),
//						utils.Shorten(checkpoints.Finalized.String()),
//					)
//
//					if (checkpoints.Finalized != common.Checkpoint{}) {
//						r.done = true
//						r.result = checkpoints.Finalized
//					}
//				}(ctx, n, results[i])
//			}
//			wg.Wait()
//
//			if err := results.CheckError(); err != nil {
//				return common.Checkpoint{}, err
//			}
//			results.PrintMessages(t.Logf)
//			if results.AllDone() {
//				if cp, ok := results[0].result.(common.Checkpoint); ok {
//					return cp, nil
//				}
//			}
//		}
//	}
//}

// WaitForExecutionFinality blocks until a beacon client reaches finality
// and the finality checkpoint contains an execution payload,
// or timeoutSlots have passed, whichever happens first.
//func (t *Testnet) WaitForExecutionFinality(
//	ctx context.Context,
//) (common.Checkpoint, error) {
//	var (
//		genesis      = t.GenesisTimeUnix()
//		slotDuration = time.Duration(t.spec.SECONDS_PER_SLOT) * time.Second
//		timer        = time.NewTicker(slotDuration)
//		runningNodes = t.VerificationNodes().Running()
//		results      = makeResults(runningNodes, t.maxConsecutiveErrorsOnWaits)
//	)
//
//	for {
//		select {
//		case <-ctx.Done():
//			return common.Checkpoint{}, ctx.Err()
//		case tim := <-timer.C:
//			// start polling after first slot of genesis
//			if tim.Before(genesis.Add(slotDuration)) {
//				t.Logf("Time till genesis: %s", genesis.Sub(tim))
//				continue
//			}
//
//			// new slot, log and check status of all beacon nodes
//			var (
//				wg        sync.WaitGroup
//				clockSlot = t.spec.TimeToSlot(
//					common.Timestamp(time.Now().Unix()),
//					t.GenesisTime(),
//				)
//			)
//			results.Clear()
//
//			for i, n := range runningNodes {
//				wg.Add(1)
//				go func(ctx context.Context, n *node.TaikoNode, r *result) {
//					defer wg.Done()
//					var (
//						b             = n.BeaconClient
//						finalizedFork string
//					)
//
//					headBlock, err := b.BlockV2(ctx, eth2api.BlockHead)
//					if err != nil {
//						r.err = errors.Wrap(err, "failed to poll head")
//						return
//					}
//					slot := headBlock.Slot()
//					if clockSlot > slot &&
//						(clockSlot-slot) >= t.spec.SLOTS_PER_EPOCH {
//						r.fatal = fmt.Errorf(
//							"unable to sync for an entire epoch: clockSlot=%d, slot=%d",
//							clockSlot,
//							slot,
//						)
//						return
//					}
//
//					checkpoints, err := b.BlockFinalityCheckpoints(
//						ctx,
//						eth2api.BlockHead,
//					)
//					if err != nil {
//						r.err = errors.Wrap(
//							err,
//							"failed to poll finality checkpoint",
//						)
//						return
//					}
//
//					execution := ethcommon.Hash{}
//					if exeuctionPayload, err := headBlock.ExecutionPayload(); err == nil {
//						execution = exeuctionPayload.BlockHash
//					}
//
//					finalizedExecution := ethcommon.Hash{}
//					if (checkpoints.Finalized != common.Checkpoint{}) {
//						if finalizedBlock, err := b.BlockV2(
//							ctx,
//							eth2api.BlockIdRoot(checkpoints.Finalized.Root),
//						); err != nil {
//							r.err = errors.Wrap(
//								err,
//								"failed to retrieve block",
//							)
//							return
//						} else {
//							finalizedFork = finalizedBlock.Version
//							if exeuctionPayload, err := finalizedBlock.ExecutionPayload(); err == nil {
//								finalizedExecution = exeuctionPayload.BlockHash
//							}
//						}
//					}
//
//					r.msg = fmt.Sprintf(
//						"fork=%s, finalized_fork=%s, clock_slot=%s, slot=%d, head=%s, "+
//							"exec_payload=%s, finalized_exec_payload=%s, justified=%s, finalized=%s",
//						headBlock.Version,
//						finalizedFork,
//						clockSlot,
//						slot,
//						utils.Shorten(headBlock.Root().String()),
//						utils.Shorten(execution.Hex()),
//						utils.Shorten(finalizedExecution.Hex()),
//						utils.Shorten(checkpoints.CurrentJustified.String()),
//						utils.Shorten(checkpoints.Finalized.String()),
//					)
//
//					if !bytes.Equal(
//						finalizedExecution[:],
//						EMPTY_EXEC_HASH[:],
//					) {
//						r.done = true
//						r.result = checkpoints.Finalized
//					}
//				}(
//					ctx,
//					n,
//					results[i],
//				)
//			}
//			wg.Wait()
//
//			if err := results.CheckError(); err != nil {
//				return common.Checkpoint{}, err
//			}
//			results.PrintMessages(t.Logf)
//			if results.AllDone() {
//				if cp, ok := results[0].result.(common.Checkpoint); ok {
//					return cp, nil
//				}
//			}
//		}
//	}
//}

// Waits for the current epoch to be finalized, or timeoutSlots have passed, whichever happens first.
//func (t *Testnet) WaitForCurrentEpochFinalization(
//	ctx context.Context,
//) (common.Checkpoint, error) {
//	var (
//		genesis      = t.GenesisTimeUnix()
//		slotDuration = time.Duration(
//			t.spec.SECONDS_PER_SLOT,
//		) * time.Second
//		timer        = time.NewTicker(slotDuration)
//		runningNodes = t.VerificationNodes().Running()
//		results      = makeResults(
//			runningNodes,
//			t.maxConsecutiveErrorsOnWaits,
//		)
//		epochToBeFinalized = t.spec.SlotToEpoch(t.spec.TimeToSlot(
//			common.Timestamp(time.Now().Unix()),
//			t.GenesisTime(),
//		))
//	)
//
//	for {
//		select {
//		case <-ctx.Done():
//			return common.Checkpoint{}, ctx.Err()
//		case tim := <-timer.C:
//			// start polling after first slot of genesis
//			if tim.Before(genesis.Add(slotDuration)) {
//				t.Logf("Time till genesis: %s", genesis.Sub(tim))
//				continue
//			}
//
//			// new slot, log and check status of all beacon nodes
//			var (
//				wg        sync.WaitGroup
//				clockSlot = t.spec.TimeToSlot(
//					common.Timestamp(time.Now().Unix()),
//					t.GenesisTime(),
//				)
//			)
//			results.Clear()
//
//			for i, n := range runningNodes {
//				i := i
//				wg.Add(1)
//				go func(ctx context.Context, n *node.TaikoNode, r *result) {
//					defer wg.Done()
//
//					b := n.BeaconClient
//
//					headInfo, err := b.BlockHeader(ctx, eth2api.BlockHead)
//					if err != nil {
//						r.err = errors.Wrap(err, "failed to poll head")
//						return
//					}
//
//					slot := headInfo.Header.Message.Slot
//					if clockSlot > slot &&
//						(clockSlot-slot) >= t.spec.SLOTS_PER_EPOCH {
//						r.fatal = fmt.Errorf(
//							"unable to sync for an entire epoch: clockSlot=%d, slot=%d",
//							clockSlot,
//							slot,
//						)
//						return
//					}
//
//					checkpoints, err := b.BlockFinalityCheckpoints(
//						ctx,
//						eth2api.BlockHead,
//					)
//					if err != nil {
//						r.err = errors.Wrap(
//							err,
//							"failed to poll finality checkpoint",
//						)
//						return
//					}
//
//					r.msg = fmt.Sprintf(
//						"clock_slot=%d, slot=%d, head=%s justified=%s, "+
//							"finalized=%s, epoch_to_finalize=%d",
//						clockSlot,
//						slot,
//						utils.Shorten(headInfo.Root.String()),
//						utils.Shorten(checkpoints.CurrentJustified.String()),
//						utils.Shorten(checkpoints.Finalized.String()),
//						epochToBeFinalized,
//					)
//
//					if checkpoints.Finalized != (common.Checkpoint{}) &&
//						checkpoints.Finalized.Epoch >= epochToBeFinalized {
//						r.done = true
//						r.result = checkpoints.Finalized
//					}
//				}(ctx, n, results[i])
//
//			}
//			wg.Wait()
//
//			if err := results.CheckError(); err != nil {
//				return common.Checkpoint{}, err
//			}
//			results.PrintMessages(t.Logf)
//			if results.AllDone() {
//				t.Logf("INFO: Epoch %d finalized", epochToBeFinalized)
//				if cp, ok := results[0].result.(common.Checkpoint); ok {
//					return cp, nil
//				}
//			}
//		}
//	}
//}

// Waits for any execution payload to be available included in a beacon block (merge),
// or timeoutSlots have passed, whichever happens first.
//func (t *Testnet) WaitForExecutionPayload(
//	ctx context.Context,
//) (ethcommon.Hash, error) {
//	var (
//		genesis      = t.GenesisTimeUnix()
//		slotDuration = time.Duration(t.spec.SECONDS_PER_SLOT) * time.Second
//		timer        = time.NewTicker(slotDuration)
//		runningNodes = t.VerificationNodes().Running()
//		results      = makeResults(
//			runningNodes,
//			t.maxConsecutiveErrorsOnWaits,
//		)
//		executionClient = runningNodes[0].ExecutionClient
//		ttdReached      = false
//	)
//
//	for {
//		select {
//		case <-ctx.Done():
//			return ethcommon.Hash{}, ctx.Err()
//		case tim := <-timer.C:
//			// start polling after first slot of genesis
//			if tim.Before(genesis.Add(slotDuration)) {
//				t.Logf("Time till genesis: %s", genesis.Sub(tim))
//				continue
//			}
//
//			if !ttdReached {
//				// Check if TTD has been reached
//				if td, err := executionClient.TotalDifficultyByNumber(ctx, nil); err == nil {
//					if td.Cmp(
//						t.L1ExecutionClientGenesis.Genesis.Config.TerminalTotalDifficulty,
//					) >= 0 {
//						ttdReached = true
//					} else {
//						continue
//					}
//				} else {
//					t.Logf("Error querying eth1 for TTD: %v", err)
//				}
//			}
//
//			// new slot, log and check status of all beacon nodes
//			var (
//				wg        sync.WaitGroup
//				clockSlot = t.spec.TimeToSlot(
//					common.Timestamp(time.Now().Unix()),
//					t.GenesisTime(),
//				)
//			)
//			results.Clear()
//
//			for i, n := range runningNodes {
//				wg.Add(1)
//				go func(ctx context.Context, n *node.TaikoNode, r *result) {
//					defer wg.Done()
//
//					b := n.BeaconClient
//
//					versionedBlock, err := b.BlockV2(
//						ctx,
//						eth2api.BlockHead,
//					)
//					if err != nil {
//						r.err = errors.Wrap(err, "failed to retrieve block")
//						return
//					}
//
//					slot := versionedBlock.Slot()
//					if clockSlot > slot &&
//						(clockSlot-slot) >= t.spec.SLOTS_PER_EPOCH {
//						r.fatal = fmt.Errorf(
//							"unable to sync for an entire epoch: clockSlot=%d, slot=%d",
//							clockSlot,
//							slot,
//						)
//						return
//					}
//
//					executionHash := ethcommon.Hash{}
//					if executionPayload, err := versionedBlock.ExecutionPayload(); err == nil {
//						executionHash = executionPayload.BlockHash
//					}
//
//					health, _ := GetHealth(ctx, b, t.spec, slot)
//
//					r.msg = fmt.Sprintf(
//						"fork=%s, clock_slot=%d, slot=%d, "+
//							"head=%s, health=%.2f, exec_payload=%s",
//						versionedBlock.Version,
//						clockSlot,
//						slot,
//						utils.Shorten(versionedBlock.Root().String()),
//						health,
//						utils.Shorten(executionHash.Hex()),
//					)
//
//					if !bytes.Equal(executionHash[:], EMPTY_EXEC_HASH[:]) {
//						r.done = true
//						r.result = executionHash
//					}
//				}(ctx, n, results[i])
//			}
//			wg.Wait()
//
//			if err := results.CheckError(); err != nil {
//				return ethcommon.Hash{}, err
//			}
//			results.PrintMessages(t.Logf)
//			if results.AllDone() {
//				if h, ok := results[0].result.(ethcommon.Hash); ok {
//					return h, nil
//				}
//			}
//
//		}
//	}
//}

//func GetHealth(
//	parentCtx context.Context,
//	bn *beacon_client.BeaconClient,
//	spec *common.Spec,
//	slot common.Slot,
//) (float64, error) {
//	var health float64
//	stateInfo, err := bn.BeaconStateV2(parentCtx, eth2api.StateIdSlot(slot))
//	if err != nil {
//		return 0, fmt.Errorf("failed to retrieve state: %v", err)
//	}
//	currentEpochParticipation := stateInfo.CurrentEpochParticipation()
//	if currentEpochParticipation != nil {
//		// Altair and after
//		health = calcHealth(currentEpochParticipation)
//	} else {
//		if stateInfo.Version != "phase0" {
//			return 0, fmt.Errorf("calculate participation")
//		}
//		state := stateInfo.Data.(*phase0.BeaconState)
//		epoch := spec.SlotToEpoch(slot)
//		validatorIds := make([]eth2api.ValidatorId, 0, len(state.Validators))
//		for id, validator := range state.Validators {
//			if epoch >= validator.ActivationEligibilityEpoch &&
//				epoch < validator.ExitEpoch &&
//				!validator.Slashed {
//				validatorIds = append(
//					validatorIds,
//					eth2api.ValidatorIdIndex(id),
//				)
//			}
//		}
//		var (
//			beforeEpoch = 0
//			afterEpoch  = spec.SlotToEpoch(slot)
//		)
//
//		// If it's genesis, keep before also set to 0.
//		if afterEpoch != 0 {
//			beforeEpoch = int(spec.SlotToEpoch(slot)) - 1
//		}
//		balancesBefore, err := bn.StateValidatorBalances(
//			parentCtx,
//			eth2api.StateIdSlot(beforeEpoch*int(spec.SLOTS_PER_EPOCH)),
//			validatorIds,
//		)
//		if err != nil {
//			return 0, fmt.Errorf(
//				"failed to retrieve validator balances: %v",
//				err,
//			)
//		}
//		balancesAfter, err := bn.StateValidatorBalances(
//			parentCtx,
//			eth2api.StateIdSlot(int(afterEpoch)*int(spec.SLOTS_PER_EPOCH)),
//			validatorIds,
//		)
//		if err != nil {
//			return 0, fmt.Errorf(
//				"failed to retrieve validator balances: %v",
//				err,
//			)
//		}
//		health = legacyCalcHealth(spec, balancesBefore, balancesAfter)
//	}
//	return health, nil
//}

func calcHealth(p altair.ParticipationRegistry) float64 {
	sum := 0
	for _, p := range p {
		sum += int(p)
	}
	avg := float64(sum) / float64(len(p))
	return avg / float64(MAX_PARTICIPATION_SCORE)
}

// legacyCalcHealth calculates the health of the network based on balances at
// the beginning of an epoch versus the balances at the end.
//
// NOTE: this isn't strictly the most correct way of doing things, but it is
// quite accurate and doesn't require implementing the attestation processing
// logic here.
func legacyCalcHealth(
	spec *common.Spec,
	before, after []eth2api.ValidatorBalanceResponse,
) float64 {
	sum_before := big.NewInt(0)
	sum_after := big.NewInt(0)
	for i := range before {
		sum_before.Add(sum_before, big.NewInt(int64(before[i].Balance)))
		sum_after.Add(sum_after, big.NewInt(int64(after[i].Balance)))
	}
	count := big.NewInt(int64(len(before)))
	avg_before := big.NewInt(0).Div(sum_before, count).Uint64()
	avg_after := sum_after.Div(sum_after, count).Uint64()
	reward := avg_before * uint64(
		spec.BASE_REWARD_FACTOR,
	) / math.IntegerSquareRootPrysm(
		sum_before.Uint64(),
	) / uint64(
		spec.HYSTERESIS_QUOTIENT,
	)
	return float64(
		avg_after-avg_before,
	) / float64(
		reward*common.BASE_REWARDS_PER_EPOCH,
	)
}
