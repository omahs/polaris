#!/bin/bash

geth=~/go/bin/polard
if [ "$HIVE_LOGLEVEL" != "" ]; then
    FLAGS="$FLAGS --verbosity=$HIVE_LOGLEVEL"
fi

# Run the go-ethereum implementation with the requested flags.
FLAGS="$FLAGS --nat=none"
echo "Running go-ethereum with flags $FLAGS"
$geth $FLAGS