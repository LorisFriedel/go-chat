#!/bin/sh

BIN_PATH=./bin

if [ ! -f $BIN_PATH/chat ]; then
    ./build.sh
fi

$BIN_PATH/chat "$@"
