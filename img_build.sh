#!/usr/bin/env bash

cd $(dirname $0)
BIN_PATH=./bin

if [ ! -f $BIN_PATH/go-chat ]; then
    ./build.sh
fi

docker build -t lorisfriedel/go-chat .
# ./clean.sh
