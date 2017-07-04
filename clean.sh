#!/usr/bin/env bash

cd $(dirname $0)
BIN_PATH=./bin
SUB_SRC_PATH=./src

sudo rm -Rf $BIN_PATH $SUB_SRC_PATH
echo "Clean done."
