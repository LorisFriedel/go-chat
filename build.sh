#!/bin/sh

./clean.sh

start=`date +%s`

mkdir -p bin
make

end=`date +%s`

runtime=$((end-start))
echo "Build done. ("${runtime}" s)"
