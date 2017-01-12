#!/bin/sh

BIN_PATH=./bin

if [ ! -f $BIN_PATH/go-chat ]; then
    ./build.sh
fi

$BIN_PATH/chat "$@"
