#!/bin/bash
set -eou pipefail

DIR=$(
  cd $(dirname ${BASH_SOURCE[0]})
  pwd
)
echo $DIR

rm -rf $DIR/deployments || true

# Get the Node.js version
node_version=$(node --version)

# Desired Node.js version
desired_version="v20.10.0"

# Check if the Node.js version matches the desired version
if [ "$node_version" != "$desired_version" ]; then
  echo "Error: Node.js version is not $desired_version"
  exit 1
fi

# If the Node.js version matches, continue with the rest of the script
echo "Node.js version is $desired_version. Continuing with the script..."

# Get the npm version
npm_version=$(npm --version)

# Desired npm version
desired_version="10.2.3"

# Check if the npm version matches the desired version
if [ "$npm_version" != "$desired_version" ]; then
  echo "Error: npm version is not $desired_version"
  exit 1
fi

# If the npm version matches, continue with the rest of the script
echo "npm version is $desired_version. Continuing with the script..."

mkdir deployments || true
rm -rf genesis_tmp || true
mkdir genesis_tmp
trap "rm -rf $DIR/genesis_tmp" EXIT INT KILL ERR

cd genesis_tmp

git clone git@github.com:taikoxyz/k8s-configs.git
cd k8s-configs
#git checkout update-a6-internal
#git checkout update-l2-genesis-2
#cp internal-devnet/docker/blockscout/internal-l2a-genesis.json ../../deployments/
#jq --sort-keys "." ../../deployments/internal-l2a-genesis.json >../../deployments/sorted_internal-l2a-genesis.json
cd ..
git clone git@github.com:taikoxyz/taiko-mono.git
cd taiko-mono/packages/protocol
git checkout alpha-6
PROJ_DIR=$(pwd)
#cp ../../../../deployments/genesis.json deployments/
# Input JSON file
#input_file="./deployments/genesis.json"

# Extract contents of 'alloc' field and write to alloc_data.json
#jq '.alloc' "$input_file" > "./deployments/genesis_alloc.json"

# Remove 'alloc' field and write the rest of the data to rest_of_data.json
#jq 'del(.alloc)' "$input_file" > "./deployments/genesis_config.json"

npm install
npm run compile
npm run generate:genesis ../../../k8s-configs/internal-devnet/config/genesis/internal_l2a.js | tee ./deployments/config_output

if ! command -v docker &>/dev/null 2>&1; then
  echo "ERROR: $(docker) command not found"
  exit 1
fi

if ! docker info >/dev/null 2>&1; then
  echo "ERROR: docker daemon isn't running"
  exit 1
fi

GENESIS_JSON=./deployments/genesis.json
TESTNET_CONFIG=$PROJ_DIR/genesis/testnet/docker-compose.yml

touch $GENESIS_JSON

# generate complete genesis json
echo -e '{\n "config": ' >$GENESIS_JSON
cat ./deployments/genesis_config.json | jq >>$GENESIS_JSON
#
echo -e ',\n "alloc": \n' >>$GENESIS_JSON
cat ./deployments/genesis_alloc.json >>$GENESIS_JSON
#
echo '}' >>$GENESIS_JSON

additional_config='{
  "byzantiumBlock": 0,
  "berlinBlock": 0,
  "constantinopleBlock": 0,
  "eip150Block": 0,
  "eip155Block": 0,
  "eip158Block": 0,
  "homesteadBlock": 0,
  "istanbulBlock": 0,
  "londonBlock": 0,
  "petersburgBlock": 0,
  "shanghaiTime": 0,
  "taiko": true,
  "terminalTotalDifficulty": 0,
  "terminalTotalDifficultyPassed": true
}'
jq --argjson additional_config "$additional_config" '.config += $additional_config' $GENESIS_JSON >temp.json && mv temp.json $GENESIS_JSON

additional_data='{
  "difficulty": "0x0",
  "blobGasUsed": null,
  "excessBlobGas": null,
  "extraData": "0x",
  "gasLimit": "0x5b8d80",
  "gasUsed": "0x0",
  "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "nonce": "0x0",
  "number": "0x0",
  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "timestamp": "0x0"
}'
jq --argjson additional_data "$additional_data" '. + $additional_data' $GENESIS_JSON >temp.json && mv temp.json $GENESIS_JSON

jq --sort-keys "." ./deployments/genesis.json >./deployments/sorted_genesis.json

cp -r ./deployments ../../../../

echo "Starting generate_genesis tests..."

# start a geth instance and init with the output genesis json
echo ""
echo "Start docker compose network..."

docker compose -f $TESTNET_CONFIG down -v --remove-orphans &>/dev/null
docker compose -f $TESTNET_CONFIG up -d

trap "docker compose -f $TESTNET_CONFIG down -v" EXIT INT KILL ERR

echo ""
echo "Start testing..."

function waitTestNode {
  echo "Waiting for test node: $1"
  # Wait till the test node fully started
  RETRIES=120
  i=0
  until curl \
    --silent \
    --fail \
    --noproxy localhost \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0","id":0,"method":"eth_chainId","params":[]}' \
    $1; do
    sleep 1
    if [ $i -eq $RETRIES ]; then
      echo 'Timed out waiting for test node'
      exit 1
    fi
    ((i = i + 1))
  done
}

waitTestNode http://localhost:18545

pwd

sed -i "s/test = 'test'/test = 'genesis'/" foundry.toml

forge test \
  -vvv \
  --gas-report \
  --fork-url http://localhost:18545 \
  --fork-retry-backoff 120 \
  --no-storage-caching \
  --evm-version cancun \
  --match-path ./genesis/GenerateGenesis.g.sol \
  --block-gas-limit 1000000000 ||
  true

docker compose -f $TESTNET_CONFIG down -v --remove-orphans

cd ../../../../
./env_vars.sh > ../taiko-l1l2-protocol/saved_env.txt
./env_vars.sh > ../taiko-l2-taiko-client/saved_env.txt


rm -rf $DIR/genesis_tmp
