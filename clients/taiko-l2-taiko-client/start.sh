#!/bin/bash
source /saved_env.txt
set -e

FLAGS=""
FLAGS="--l1.ws $HIVE_TAIKO_L1_WS_ENDPOINT"
FLAGS="$FLAGS --taikoL1 $HIVE_TAIKO_L1_ROLLUP_ADDRESS  --taikoL2 $HIVE_TAIKO_L2_ROLLUP_ADDRESS"
FLAGS="$FLAGS --verbosity $HIVE_LOGLEVEL"

case $HIVE_TAIKO_ROLE in
"driver")
  echo "$HIVE_TAIKO_JWT_SECRET" >/jwtsecret
  FLAGS="$FLAGS --l2.ws $HIVE_TAIKO_L2_WS_ENDPOINT"
  FLAGS="$FLAGS --l2.auth $HIVE_TAIKO_L2_ENGINE_ENDPOINT"
  FLAGS="$FLAGS --l2.throwawayBlockBuilderPrivKey=$HIVE_TAIKO_THROWAWAY_BLOCK_BUILDER_PRIVATE_KEY"
  FLAGS="$FLAGS --jwtSecret /jwtsecret"
  if [ "$HIVE_TAIKO_ENABLE_L2_P2P" != "" ]; then
    FLAGS="$FLAGS --p2p.syncVerifiedBlocks"
  fi
  ;;
"prover")
  FLAGS="$FLAGS --l2.ws $HIVE_TAIKO_L2_WS_ENDPOINT"
  FLAGS="$FLAGS --l1.http $HIVE_TAIKO_L1_HTTP_ENDPOINT"
  FLAGS="$FLAGS --l2.http $HIVE_TAIKO_L2_HTTP_ENDPOINT"
  FLAGS="$FLAGS --zkevmRpcdEndpoint=ws://127.0.0.1:18545"
  FLAGS="$FLAGS --zkevmRpcdParamsPath=12345"
  FLAGS="$FLAGS --l1.proverPrivKey=$HIVE_TAIKO_PROVER_PRIVATE_KEY"
  FLAGS="$FLAGS --dummy"
  ;;
"proposer")
  FLAGS="$FLAGS --l2.http $HIVE_TAIKO_L2_HTTP_ENDPOINT"
  FLAGS="$FLAGS --l1.proposerPrivKey=$HIVE_TAIKO_PROPOSER_PRIVATE_KEY"
  FLAGS="$FLAGS --l2.suggestedFeeRecipient=$HIVE_TAIKO_SUGGESTED_FEE_RECIPIENT"
  FLAGS="$FLAGS --proposeInterval=$HIVE_TAIKO_PROPOSE_INTERVAL"
  if [ "$HIVE_TAIKO_PRODUCE_INVALID_BLOCKS_INTERVAL" != "" ]; then
    FLAGS="$FLAGS --produceInvalidBlocks"
    FLAGS="$FLAGS --produceInvalidBlocksInterval=$HIVE_TAIKO_PRODUCE_INVALID_BLOCKS_INTERVAL"
  fi
  ;;
esac

# Run the go-ethereum implementation with the requested flags.
echo "Running $HIVE_TAIKO_ROLE with flags $FLAGS"
taiko-client $HIVE_TAIKO_ROLE $FLAGS
