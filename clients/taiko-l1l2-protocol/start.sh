#!/bin/bash
config_output=/config.output
bridge_impl=$(grep -oP "BridgeImpl: '(.*)'" "$config_output" | awk -F"'" '{print $2}')
erc20_vault_impl=$(grep -oP "ERC20VaultImpl: '(.*)'" "$config_output" | awk -F"'" '{print $2}')
erc721_vault_impl=$(grep -oP "ERC721VaultImpl: '(.*)'" "$config_output" | awk -F"'" '{print $2}')
erc1155_vault_impl=$(grep -oP "ERC1155VaultImpl: '(.*)'" "$config_output" | awk -F"'" '{print $2}')
signal_service_impl=$(grep -oP "SignalServiceImpl: '(.*)'" "$config_output" | awk -F"'" '{print $2}')
shared_address_manager_impl=$(grep -oP "SharedAddressManagerImpl: '(.*)'" "$config_output" | awk -F"'" '{print $2}')
bridged_erc20_impl=$(grep -oP "BridgedERC20Impl: '(.*)'" "$config_output" | awk -F"'" '{print $2}')
bridged_erc721_impl=$(grep -oP "BridgedERC721Impl: '(.*)'" "$config_output" | awk -F"'" '{print $2}')
bridged_erc1155_impl=$(grep -oP "BridgedERC1155Impl: '(.*)'" "$config_output" | awk -F"'" '{print $2}')
taiko_l2_impl=$(grep -oP "TaikoL2Impl: '(.*)'" "$config_output" | awk -F"'" '{print $2}')
rollup_address_manager_impl=$(grep -oP "RollupAddressManagerImpl: '(.*)'" "$config_output" | awk -F"'" '{print $2}')
bridge=$(grep -oP "Bridge: '(.*)'" "$config_output" | awk -F"'" '{print $2}')
erc20_vault=$(grep -oP "ERC20Vault: '(.*)'" "$config_output" | awk -F"'" '{print $2}')
erc721_vault=$(grep -oP "ERC721Vault: '(.*)'" "$config_output" | awk -F"'" '{print $2}')
erc1155_vault=$(grep -oP "ERC1155Vault: '(.*)'" "$config_output" | awk -F"'" '{print $2}')
signal_service=$(grep -oP "SignalService: '(.*)'" "$config_output" | awk -F"'" '{print $2}')
shared_address_manager=$(grep -oP "SharedAddressManager: '(.*)'" "$config_output" | awk -F"'" '{print $2}')
taiko_l2=$(grep -oP "TaikoL2: '(.*)'" "$config_output" | awk -F"'" '{print $2}')
rollup_address_manager=$(grep -oP "RollupAddressManager: '(.*)'" "$config_output" | awk -F"'" '{print $2}')

#export TAIKO_L2_ADDRESS=$taiko_l2_impl
#export L2_SIGNAL_SERVICE=$signal_service_impl
#ENV TAIKO_TOKEN_NAME="Taiko Token Katla"
#ENV TAIKO_TOKEN_SYMBOL=TTKOk
#export SHARED_ADDRESS_MANAGER=$shared_address_manager
#ENV L2_GENESIS_HASH=0xee1950562d42f0da28bd4550d88886bc90894c77c9c9eaefef775d4c8223f259
#export BRIDGE_IMPL="$bridge_impl"
#export ERC20_VAULT_IMPL="$erc20_vault_impl"
#export ERC721_VAULT_IMPL="$erc721_vault_impl"
#export ERC1155_VAULT_IMPL="$erc1155_vault_impl"
#export SIGNAL_SERVICE_IMPL="$signal_service_impl"
#export SHARED_ADDRESS_MANAGER_IMPL="$shared_address_manager_impl"
#export BRIDGED_ERC20_IMPL="$bridged_erc20_impl"
#export BRIDGED_ERC721_IMPL="$bridged_erc721_impl"
#export BRIDGED_ERC1155_IMPL="$bridged_erc1155_impl"
#export TAIKO_L2_IMPL="$taiko_l2_impl"
#export ROLLUP_ADDRESS_MANAGER_IMPL="$rollup_address_manager_impl"
#export BRIDGE="$bridge"
#export ERC20_VAULT="$erc20_vault"
#export ERC721_VAULT="$erc721_vault"
#export ERC1155_VAULT="$erc1155_vault"
#export SIGNAL_SERVICE="$signal_service"
#export SHARED_ADDRESS_MANAGER="$shared_address_manager"
#export TAIKO_L2="$taiko_l2"
#export ROLLUP_ADDRESS_MANAGER="$rollup_address_manager"

echo "BRIDGE_IMPL=$BRIDGE_IMPL"
echo "ERC20_VAULT_IMPL=$ERC20_VAULT_IMPL"
echo "ERC721_VAULT_IMPL=$ERC721_VAULT_IMPL"
echo "ERC1155_VAULT_IMPL=$ERC1155_VAULT_IMPL"
echo "SIGNAL_SERVICE_IMPL=$SIGNAL_SERVICE_IMPL"
echo "SHARED_ADDRESS_MANAGER_IMPL=$SHARED_ADDRESS_MANAGER_IMPL"
echo "BRIDGED_ERC20_IMPL=$BRIDGED_ERC20_IMPL"
echo "BRIDGED_ERC721_IMPL=$BRIDGED_ERC721_IMPL"
echo "BRIDGED_ERC1155_IMPL=$BRIDGED_ERC1155_IMPL"
echo "TAIKO_L2_IMPL=$TAIKO_L2_IMPL"
echo "ROLLUP_ADDRESS_MANAGER_IMPL=$ROLLUP_ADDRESS_MANAGER_IMPL"
echo "BRIDGE=$BRIDGE"
echo "ERC20_VAULT=$ERC20_VAULT"
echo "ERC721_VAULT=$ERC721_VAULT"
echo "ERC1155_VAULT=$ERC1155_VAULT"
echo "SIGNAL_SERVICE=$SIGNAL_SERVICE"
echo "SHARED_ADDRESS_MANAGER=$SHARED_ADDRESS_MANAGER"
echo "TAIKO_L2=$TAIKO_L2"
echo "ROLLUP_ADDRESS_MANAGER=$ROLLUP_ADDRESS_MANAGER"
# Inherited $PRIVATE_KEY
# Inherited $MAINNET_URL
#    SHARED_ADDRESS_MANAGER=$SHARED_ADDRESS_MANAGER \
#    GUARDIAN_PROVERS=$GUARDIAN_PROVERS \
#    TAIKO_L2_ADDRESS=$TAIKO_L2_ADDRESS \
#    L2_SIGNAL_SERVICE=$L2_SIGNAL_SERVICE \
#    TIMELOCK_CONTROLLER=$TIMELOCK_CONTROLLER \
#    PROPOSER=0x0000000000000000000000000000000000000000 \
#    PROPOSER_ONE=0x0000000000000000000000000000000000000000 \
#    TAIKO_TOKEN_PREMINT_RECIPIENT=$TAIKO_TOKEN_PREMINT_RECIPIENT \
#    TAIKO_TOKEN_NAME="Taiko Token Test" \
#    TAIKO_TOKEN_SYMBOL=TTKOt \
#    L2_GENESIS_HASH=$L2_GENESIS_HASH \
#    SECURITY_COUNCIL=$CONTRACTS_OWNER \
#    MIN_GUARDIANS=4 \
    export FORGE_VERBOSE=true
    /root/.foundry/bin/forge script script/DeployOnL1.s.sol:DeployOnL1 \
    --fork-url $MAINNET_URL \
    --broadcast \
    --ffi \
    -vvvv \
    --private-key $PRIVATE_KEY \
    --block-gas-limit 100000000 \
    | tee /root/deploy.log
