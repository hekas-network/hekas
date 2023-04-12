#!/bin/bash

KEY="dev0"
CHAINID="hekas_9898-1"
MONIKER="mymoniker"
DATA_DIR=$(mktemp -d -t hekas-datadir.XXXXX)

echo "create and add new keys"
./hekasd keys add $KEY --home $DATA_DIR --no-backup --chain-id $CHAINID --algo "eth_secp256k1" --keyring-backend test
echo "init Hekas with moniker=$MONIKER and chain-id=$CHAINID"
./hekasd init $MONIKER --chain-id $CHAINID --home $DATA_DIR
echo "prepare genesis: Allocate genesis accounts"
./hekasd add-genesis-account \
"$(./hekasd keys show $KEY -a --home $DATA_DIR --keyring-backend test)" 1000000000000000000ahekas,1000000000000000000stake \
--home $DATA_DIR --keyring-backend test
echo "prepare genesis: Sign genesis transaction"
./hekasd gentx $KEY 1000000000000000000stake --keyring-backend test --home $DATA_DIR --keyring-backend test --chain-id $CHAINID
echo "prepare genesis: Collect genesis tx"
./hekasd collect-gentxs --home $DATA_DIR
echo "prepare genesis: Run validate-genesis to ensure everything worked and that the genesis file is setup correctly"
./hekasd validate-genesis --home $DATA_DIR

echo "starting hekas node $i in background ..."
./hekasd start --pruning=nothing --rpc.unsafe \
--keyring-backend test --home $DATA_DIR \
>$DATA_DIR/node.log 2>&1 & disown

echo "started hekas node"
tail -f /dev/null