SOURCEDIR=.
SOURCES=$(shell find ${SOURCEDIR} -name '*.go')

BINARY=chat

.DEFAULT_GOAL: install

${BINARY}: ${SOURCES}
	go build -o ./bin/${BINARY} main.go

.PHONY: install
install:
	go install ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
