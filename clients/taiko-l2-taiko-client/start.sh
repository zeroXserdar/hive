#!/bin/bash
check_tcp_connection() {
    local url="$1"

    # Extract host and port from the URL
    host=$(echo "$url" | cut -d':' -f2 | cut -d'/' -f3)
    port=$(echo "$url" | cut -d':' -f3 | cut -d'/' -f1)

    # Check if the connection is open
    if nc -z "$host" "$port"; then
        echo "Connection to $host on port $port is open"
    else
        echo "Connection to $host on port $port is closed"
    fi
}

source /saved_env.txt
#set -e
TAIKO_L1_HTTP_ENDPOINT="${HIVE_MAINNET_URL}"
TAIKO_L1_WS_ENDPOINT="${HIVE_MAINNET_URL%5}6"
TAIKO_L2_HTTP_ENDPOINT="${HIVE_TAIKO_URL}"
TAIKO_L2_WS_ENDPOINT="${HIVE_TAIKO_URL%5}6"

check_tcp_connection $TAIKO_L1_HTTP_ENDPOINT
check_tcp_connection $TAIKO_L1_WS_ENDPOINT
check_tcp_connection $TAIKO_L2_HTTP_ENDPOINT
check_tcp_connection $TAIKO_L2_WS_ENDPOINT

TAIKO_L1_WS_ENDPOINT="${TAIKO_L1_WS_ENDPOINT/http/ws}"
TAIKO_L2_WS_ENDPOINT="${TAIKO_L2_WS_ENDPOINT/http/ws}"

TAIKO_L1_ROLLUP_ADDRESS="${HIVE_TAIKO_L1_ADDRESS}"
TAIKO_L2_ROLLUP_ADDRESS="${TAIKO_L2_ADDRESS}"

TAIKO_L1_TOKEN_ADDRESS="${HIVE_TAIKO_L1_TOKEN_ADDRESS}"

echo "${HIVE_TAIKO_L2_GETH_JWT_SECRET}" > "/root/jwtsecret"
TAIKO_L2_GETH_JWT_SECRET="/root/jwtsecret"

//TODO: odpalowac
HIVE_TAIKO_PROVER_PRIVATE_KEY=5599a7be5589682da3e0094806840e8510dae6493665a701b06c59cbe9d97968
PROVER_CAPACITY=1

FLAGS=""
FLAGS="--l1.ws $TAIKO_L1_WS_ENDPOINT"
FLAGS="$FLAGS --taikoL1 $TAIKO_L1_ROLLUP_ADDRESS  --taikoL2 $TAIKO_L2_ROLLUP_ADDRESS"
#FLAGS="$FLAGS --verbosity $HIVE_LOGLEVEL"
FLAGS="$FLAGS --verbosity 5"

case $HIVE_TAIKO_ROLE in
"driver")
  echo "$HIVE_TAIKO_JWT_SECRET" >/jwtsecret
  FLAGS="$FLAGS --l2.ws $TAIKO_L2_WS_ENDPOINT"
  FLAGS="$FLAGS --l2.auth $TAIKO_L2_HTTP_ENDPOINT"
#  FLAGS="$FLAGS --l2.throwawayBlockBuilderPrivKey=$HIVE_TAIKO_THROWAWAY_BLOCK_BUILDER_PRIVATE_KEY"
  FLAGS="$FLAGS --jwtSecret $TAIKO_L2_GETH_JWT_SECRET"
  if [ "$HIVE_TAIKO_ENABLE_L2_P2P" != "" ]; then
    FLAGS="$FLAGS --p2p.syncVerifiedBlocks"
  fi
  ;;
"prover")
  FLAGS="$FLAGS --l2.ws $TAIKO_L2_WS_ENDPOINT"
  FLAGS="$FLAGS --l1.http $TAIKO_L1_HTTP_ENDPOINT"
  FLAGS="$FLAGS --l2.http $TAIKO_L2_HTTP_ENDPOINT"
  FLAGS="$FLAGS --zkevmRpcdEndpoint=ws://127.0.0.1:18545"
  FLAGS="$FLAGS --zkevmRpcdParamsPath=12345"
  FLAGS="$FLAGS --l1.proverPrivKey=$HIVE_TAIKO_PROVER_PRIVATE_KEY"
  FLAGS="$FLAGS --dummy"
  FLAGS="$FLAGS --taikoToken $TAIKO_L1_TOKEN_ADDRESS"
  FLAGS="$FLAGS --prover.capacity $PROVER_CAPACITY"
  ;;
"proposer")
  FLAGS="$FLAGS --l2.http $TAIKO_L2_HTTP_ENDPOINT"
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
taiko-client $HIVE_TAIKO_ROLE $FLAGS | tee /root/taiko-client.logs &

nc -l -p 8545 &

while true; do
  echo "Sleeping for 60 sec"
  sleep 60
done