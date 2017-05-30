#!/usr/bin/env bash

cd $(dirname $0)

./clean.sh

start=`date +%s`

go get github.com/LorisFriedel/go-chat
godep save
sudo docker run --rm -it -v "$GOPATH":/gopath -v "$(pwd)":/app -e "GOPATH=/gopath" -w /app golang:1.8.1 sh -c 'CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags="-s" -o bin/go-chat'

end=`date +%s`

runtime=$((end-start))
echo "Build done. ("${runtime}" s)"
