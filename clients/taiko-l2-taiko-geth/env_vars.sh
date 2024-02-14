#!/bin/bash
config_output=""
if [ -e "/deployments/config_output" ]; then
  config_output="/deployments/config_output"
else
  if [ -e "./deployments/config_output" ]; then
    config_output="./deployments/config_output"
  fi
fi
bridge_impl=$(grep "BridgeImpl: '" "$config_output" | cut -d "'" -f 2)
erc20_vault_impl=$(grep "ERC20VaultImpl: '" "$config_output" | cut -d "'" -f 2)
erc721_vault_impl=$(grep "ERC721VaultImpl: '" "$config_output" | cut -d "'" -f 2)
erc1155_vault_impl=$(grep "ERC1155VaultImpl: '" "$config_output" | cut -d "'" -f 2)
signal_service_impl=$(grep "SignalServiceImpl: '" "$config_output" | cut -d "'" -f 2)
shared_address_manager_impl=$(grep "SharedAddressManagerImpl: '" "$config_output" | cut -d "'" -f 2)
bridged_erc20_impl=$(grep "BridgedERC20Impl: '" "$config_output" | cut -d "'" -f 2)
bridged_erc721_impl=$(grep "BridgedERC721Impl: '" "$config_output" | cut -d "'" -f 2)
bridged_erc1155_impl=$(grep "BridgedERC1155Impl: '" "$config_output" | cut -d "'" -f 2)
taiko_l2_impl=$(grep "TaikoL2Impl: '" "$config_output" | cut -d "'" -f 2)
rollup_address_manager_impl=$(grep "RollupAddressManagerImpl: '" "$config_output" | cut -d "'" -f 2)
bridge=$(grep "Bridge: '" "$config_output" | cut -d "'" -f 2)
erc20_vault=$(grep "ERC20Vault: '" "$config_output" | cut -d "'" -f 2)
erc721_vault=$(grep "ERC721Vault: '" "$config_output" | cut -d "'" -f 2)
erc1155_vault=$(grep "ERC1155Vault: '" "$config_output" | cut -d "'" -f 2)
signal_service=$(grep "SignalService: '" "$config_output" | cut -d "'" -f 2)
shared_address_manager=$(grep "SharedAddressManager: '" "$config_output" | cut -d "'" -f 2)
taiko_l2=$(grep "TaikoL2: '" "$config_output" | cut -d "'" -f 2)
rollup_address_manager=$(grep "RollupAddressManager: '" "$config_output" | cut -d "'" -f 2)

export TAIKO_L2_ADDRESS=$taiko_l2_impl
export TAIKO_L2_SIGNAL_SERVICE=$signal_service_impl
export TAIKO_TOKEN_NAME="Taiko Token Katla"
export TAIKO_TOKEN_SYMBOL=TTKOk
export TAIKO_SHARED_ADDRESS_MANAGER=$shared_address_manager
export TAIKO_L2_GENESIS_HASH=0xa378ed591ada87cce719d2b43ce0c7e632cc798d33451d69eaa384181b9281d9
export TAIKO_BRIDGE_IMPL="$bridge_impl"
export TAIKO_ERC20_VAULT_IMPL="$erc20_vault_impl"
export TAIKO_ERC721_VAULT_IMPL="$erc721_vault_impl"
export TAIKO_ERC1155_VAULT_IMPL="$erc1155_vault_impl"
export TAIKO_SIGNAL_SERVICE_IMPL="$signal_service_impl"
export TAIKO_SHARED_ADDRESS_MANAGER_IMPL="$shared_address_manager_impl"
export TAIKO_BRIDGED_ERC20_IMPL="$bridged_erc20_impl"
export TAIKO_BRIDGED_ERC721_IMPL="$bridged_erc721_impl"
export TAIKO_BRIDGED_ERC1155_IMPL="$bridged_erc1155_impl"
export TAIKO_L2_IMPL="$taiko_l2_impl"
export TAIKO_ROLLUP_ADDRESS_MANAGER_IMPL="$rollup_address_manager_impl"
export TAIKO_BRIDGE="$bridge"
export TAIKO_ERC20_VAULT="$erc20_vault"
export TAIKO_ERC721_VAULT="$erc721_vault"
export TAIKO_ERC1155_VAULT="$erc1155_vault"
export TAIKO_SIGNAL_SERVICE="$signal_service"
export TAIKO_SHARED_ADDRESS_MANAGER="$shared_address_manager"
export TAIKO_L2="$taiko_l2"
export TAIKO_ROLLUP_ADDRESS_MANAGER="$rollup_address_manager"

echo "export L2_GENESIS_HASH=$TAIKO_L2_GENESIS_HASH"
echo "export TAIKO_L2_ADDRESS=$TAIKO_L2_ADDRESS"
echo "export L2_SIGNAL_SERVICE=$TAIKO_L2_SIGNAL_SERVICE"
echo "export TAIKO_TOKEN_NAME=\"$TAIKO_TOKEN_NAME\""
echo "export TAIKO_TOKEN_SYMBOL=$TAIKO_TOKEN_SYMBOL"
echo "export BRIDGE_IMPL=$TAIKO_BRIDGE_IMPL"
echo "export ERC20_VAULT_IMPL=$TAIKO_ERC20_VAULT_IMPL"
echo "export ERC721_VAULT_IMPL=$TAIKO_ERC721_VAULT_IMPL"
echo "export ERC1155_VAULT_IMPL=$TAIKO_ERC1155_VAULT_IMPL"
echo "export SIGNAL_SERVICE_IMPL=$TAIKO_SIGNAL_SERVICE_IMPL"
echo "export SHARED_ADDRESS_MANAGER_IMPL=$TAIKO_SHARED_ADDRESS_MANAGER_IMPL"
echo "export BRIDGED_ERC20_IMPL=$TAIKO_BRIDGED_ERC20_IMPL"
echo "export BRIDGED_ERC721_IMPL=$TAIKO_BRIDGED_ERC721_IMPL"
echo "export BRIDGED_ERC1155_IMPL=$TAIKO_BRIDGED_ERC1155_IMPL"
echo "export TAIKO_L2_IMPL=$TAIKO_L2_IMPL"
echo "export ROLLUP_ADDRESS_MANAGER_IMPL=$TAIKO_ROLLUP_ADDRESS_MANAGER_IMPL"
echo "export BRIDGE=$TAIKO_BRIDGE"
echo "export ERC20_VAULT=$TAIKO_ERC20_VAULT"
echo "export ERC721_VAULT=$TAIKO_ERC721_VAULT"
echo "export ERC1155_VAULT=$TAIKO_ERC1155_VAULT"
echo "export SIGNAL_SERVICE=$TAIKO_SIGNAL_SERVICE"
#echo "export SHARED_ADDRESS_MANAGER=$TAIKO_SHARED_ADDRESS_MANAGER"
echo "export TAIKO_L2=$TAIKO_L2"
echo "export ROLLUP_ADDRESS_MANAGER=$TAIKO_ROLLUP_ADDRESS_MANAGER"

echo "export TAIKO_TOKEN_PREMINT_RECIPIENT=0xa0Ee7A142d267C1f36714E4a8F75612F20a79720"
echo "export PROTOCOL_LAYER=L1"
echo "export SECURITY_COUNCIL=0x60997970C51812dc3A010C7d01b50e0d17dc79C8"
echo "export PROPOSER_ONE=0x0000000000000000000000000000000000000000"
echo "export PROPOSER=0x0000000000000000000000000000000000000000"
echo "export PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
echo "export MIN_GUARDIANS=3"
echo "export LOG_LEVEL=debug"
echo "export GUARDIAN_PROVERS=0x1000777700000000000000000000000000000001,0x1000777700000000000000000000000000000002,0x1000777700000000000000000000000000000003,0x1000777700000000000000000000000000000004,0x1000777700000000000000000000000000000005"
echo "export TIMELOCK_CONTROLLER=0x0000000000000000000000000000000000000000"
echo "export SHARED_ADDRESS_MANAGER=0x0000000000000000000000000000000000000000"
