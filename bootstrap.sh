#!/bin/bash

KEY="key"
CHAINID="hekas_9899-1"
MONIKER="first"
TOTAL="100000000000000000000000000"
STAKE="100000000000000000000"
# Remember to change to other types of keyring like 'file' in-case exposing to outside world,
# otherwise your balance will be wiped quickly
# The keyring test does not require private key to steal tokens from you
KEYRING="os"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
# Set dedicated home directory for the hekasd instance
HOMEDIR="$HOME/.hekas"
# to trace evm
#TRACE="--trace"
TRACE=""

# Path variables
CONFIG=$HOMEDIR/config/config.toml
APP_TOML=$HOMEDIR/config/app.toml
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json

# used to exit on first error (any non-zero exit code)
set -e

# User prompt if an existing local node configuration is found.
if [ -d "$HOMEDIR" ]; then
	printf "\nAn existing folder at '%s' was found. You can choose to delete this folder and start a new local node with new keys from genesis. When declined, the existing local node is started. \n" "$HOMEDIR"
	echo "Overwrite the existing configuration and start a new local node? [y/n]"
	read -r overwrite
else
	overwrite="Y"
fi


# Setup local node if overwrite is set to Yes, otherwise skip setup
if [[ $overwrite == "y" || $overwrite == "Y" ]]; then
	# Remove the previous folder
	rm -rf "$HOMEDIR"

	# Set client config
	hekasd config keyring-backend $KEYRING --home "$HOMEDIR"
	hekasd config chain-id $CHAINID --home "$HOMEDIR"

	# If keys exist they should be deleted
	hekasd keys add "$KEY" --keyring-backend $KEYRING --algo $KEYALGO --home "$HOMEDIR"


	# Set moniker and chain-id for Hekas (Moniker can be anything, chain-id must be an integer)
	hekasd init $MONIKER -o --chain-id $CHAINID --home "$HOMEDIR"

	# Change parameter token denominations to ahekas
	jq '.app_state["staking"]["params"]["bond_denom"]="ahekas"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq '.app_state["crisis"]["constant_fee"]["denom"]="ahekas"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]="ahekas"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

	# Set gas limit in genesis
	jq '.consensus_params["block"]["max_gas"]="10000000"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

	# Disable claims
	jq '.app_state["claims"]["params"]["enable_claims"]=false' >"$TMP_GENESIS" "$GENESIS" && mv "$TMP_GENESIS" "$GENESIS"


	# Allocate genesis accounts (cosmos formatted addresses)
	hekasd add-genesis-account "$KEY" "$TOTAL"ahekas --keyring-backend $KEYRING --home "$HOMEDIR"


	jq -r --arg total_supply "$TOTAL" '.app_state["bank"]["supply"][0]["amount"]=$total_supply' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

	# Sign genesis transaction
	hekasd gentx "$KEY" "$STAKE"ahekas --keyring-backend $KEYRING --chain-id $CHAINID --home "$HOMEDIR"

	# Collect genesis tx
	hekasd collect-gentxs --home "$HOMEDIR"

	# Run this to ensure everything worked and that the genesis file is setup correctly
	hekasd validate-genesis --home "$HOMEDIR"

fi

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
#hekasd start --metrics "$TRACE" --log_level $LOGLEVEL --minimum-gas-prices=0.0001ahekas --json-rpc.api eth,txpool,personal,net,debug,web3 --api.enable --home "$HOMEDIR"
