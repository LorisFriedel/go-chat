#!/usr/bin/env bash

cd $(dirname $0)
BIN_PATH=./bin

./build.sh
docker build -t lorisfriedel/go-chat .
./clean.sh
