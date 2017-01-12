SOURCEDIR=.
SOURCES=$(shell find ${SOURCEDIR} -name '*.go')

BINARY=go-chat

.DEFAULT_GOAL: install

${BINARY}: ${SOURCES}
	go get ${SOURCEDIR} && go build -o ./bin/${BINARY} main.go

.PHONY: install
install:
	go install ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
