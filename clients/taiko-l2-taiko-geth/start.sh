#!/bin/bash
set -e

source /env_vars.sh

geth=/usr/local/bin/geth
FLAGS="--pcscdpath=\"\""

if [ "$HIVE_LOGLEVEL" != "" ]; then
  FLAGS="$FLAGS --verbosity=$HIVE_LOGLEVEL"
fi

# It doesn't make sense to dial out, use only a pre-set bootnode.
#FLAGS="$FLAGS --bootnodes=$HIVE_BOOTNODE"

if [ "$HIVE_SKIP_POW" != "" ]; then
  FLAGS="$FLAGS --fakepow"
fi

# If a specific network ID is requested, use that
if [ "$HIVE_NETWORK_ID" != "" ]; then
  FLAGS="$FLAGS --networkid $HIVE_NETWORK_ID"
else
  # Unless otherwise specified by hive, we try to avoid mainnet networkid. If geth detects mainnet network id,
  # then it tries to bump memory quite a lot
  FLAGS="$FLAGS --networkid 1337"
fi

# If the client is to be run in testnet mode, flag it as such
if [ "$HIVE_TESTNET" == "1" ]; then
  FLAGS="$FLAGS --testnet"
fi

# Handle any client mode or operation requests
if [ "$HIVE_NODETYPE" == "archive" ]; then
  FLAGS="$FLAGS --syncmode full --gcmode archive"
fi
if [ "$HIVE_NODETYPE" == "full" ]; then
  FLAGS="$FLAGS --syncmode full"
fi
if [ "$HIVE_NODETYPE" == "light" ]; then
  FLAGS="$FLAGS --syncmode light"
fi
if [ "$HIVE_NODETYPE" == "snap" ]; then
  FLAGS="$FLAGS --syncmode snap"
fi
if [ -z "$HIVE_NODETYPE" ]; then
  FLAGS="$FLAGS --syncmode snap"
fi

# Import clique signing key.
#if [ -n "$HIVE_CLIQUE_PRIVATEKEY" ]; then
#     Create password file.
#    echo "Importing clique key..."
#    echo "secret" >/geth-password-file.txt
#    $geth --nousb account import --password /geth-password-file.txt <(echo "$HIVE_CLIQUE_PRIVATEKEY")
#
#     Ensure password file is used when running geth in mining mode.
#    if [ -n "$HIVE_MINER" ]; then
#        FLAGS="$FLAGS --password /geth-password-file.txt --unlock $HIVE_MINER --allow-insecure-unlock"
#    fi
#fi

# Configure any mining operation
#if [ -n "$HIVE_MINER" ] && [ "$HIVE_NODETYPE" != "light" ]; then
#FLAGS="$FLAGS --mine --miner.threads 1 --miner.etherbase $HIVE_MINER"
#FLAGS="$FLAGS --mine --miner.etherbase 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
#fi
if [ -n "$HIVE_MINER_EXTRA" ]; then
  FLAGS="$FLAGS --miner.extradata $HIVE_MINER_EXTRA"
fi
FLAGS="$FLAGS --miner.gasprice 16000000000"

# Configure LES.
if [ "$HIVE_LES_SERVER" == "1" ]; then
  FLAGS="$FLAGS --light.serve 50 --light.nosyncserve"
fi

# Configure RPC.
FLAGS="$FLAGS --http --http.addr=0.0.0.0 --http.port=8545 --http.vhosts=* --http.api=admin,debug,eth,miner,net,personal,txpool,web3,taiko"
FLAGS="$FLAGS --ws --ws.addr=0.0.0.0 --ws.origins=* --ws.api=admin,debug,eth,miner,net,personal,txpool,web3,taiko"

# if [ "$HIVE_TERMINAL_TOTAL_DIFFICULTY" != "" ]; then
#echo "$HIVE_TAIKO_JWT_SECRET" >/jwtsecret
#FLAGS="$FLAGS --authrpc.addr=0.0.0.0 --authrpc.port=8551 --authrpc.vhosts=* --authrpc.jwtsecret=/jwtsecret"
# fi

# Configure GraphQL.
if [ -n "$HIVE_GRAPHQL_ENABLED" ]; then
  FLAGS="$FLAGS --graphql"
fi
# used for the graphql to allow submission of unprotected tx
if [ -n "$HIVE_ALLOW_UNPROTECTED_TX" ]; then
  FLAGS="$FLAGS --rpc.allow-unprotected-txs"
fi

# Run the go-ethereum implementation with the requested flags.
FLAGS="$FLAGS --nat=none"

# taiko part start:
FLAGS="$FLAGS --taiko --allow-insecure-unlock"
# taiko part end

FLAGS="$FLAGS --datadir /taiko-l2-network/node"

echo "Running taiko-geth with flags $FLAGS"
$geth $FLAGS | tee /taiko-geth.log &

IP_ADDRESS=127.0.0.1
PORT=8545
echo "IP Address: $IP_ADDRESS"
echo "Port: $PORT"

while ! nc -z "$IP_ADDRESS" "$PORT"; do
  echo "Waiting for TCP connection on $IP_ADDRESS:$PORT..."
  sleep 1
done

echo "TCP connection on $IP_ADDRESS:$PORT is available."

TAIKO_L2_GENESIS_HASH=$(geth --exec 'eth.getBlock(0).hash' attach /taiko-l2-network/node/geth.ipc)
export TAIKO_L2_GENESIS_HASH
echo "TAIKO_L2_GENESIS_HASH=$TAIKO_L2_GENESIS_HASH"
env | grep TAIKO >/saved_env.txt

wait
