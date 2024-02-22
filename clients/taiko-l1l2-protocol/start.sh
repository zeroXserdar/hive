#!/bin/bash
source /saved_env.txt
touch /root/sleep
TAIKO_URL=$HIVE_TAIKO_URL
URL="${TAIKO_URL/http:\/\//}"

#HIVE_MAINNET_URL="http://172.17.0.2:8545"
MAINNET_URL=$HIVE_MAINNET_URL
URL="${MAINNET_URL/http:\/\//}"

# Extract the IP address (assuming it comes before the first ":")
IP_ADDRESS="${URL%%:*}"

# Extract the port (assuming it comes after the last ":")
PORT="${URL##*:}"

echo "IP Address: $IP_ADDRESS"
echo "Port: $PORT"

while ! nc -z "$IP_ADDRESS" "$PORT"; do
  echo "Waiting for TCP connection on $IP_ADDRESS:$PORT..."
  sleep 1
done

echo "TCP connection on $IP_ADDRESS:$PORT is available."

#env
export FORGE_VERBOSE=true
#PRIVATE_KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
#source /shit.txt
#env
/root/.foundry/bin/forge script script/DeployOnL1.s.sol:DeployOnL1 \
  --fork-url $MAINNET_URL \
  --broadcast \
  --ffi \
  -vvvv \
  --private-key $PRIVATE_KEY \
  --block-gas-limit 100000000 |
  tee /root/deploy.log

nc -l -p 8545 &

while true; do
  echo "Sleeping for 60 sec"
  sleep 60
done