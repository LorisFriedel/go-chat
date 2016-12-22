package console

import (
	"errors"
	"fmt"
)

// Arguments provider type

type IProvider interface {
	HasMore() bool
	CommandName() string
	NextString() (string, error)
	NextInt() (int, error)
}

type ArgProvider struct {
	cmdName string
	args    []interface{}
	index   int
	argSize int
}

// TODO reduce slice each NextXX ?

func NewArgProvider(cmdName string, args []interface{}) *ArgProvider {
	return &ArgProvider{cmdName, args, -1, len(args)}
}

func (p *ArgProvider) CommandName() string { // no error here, command name is already known
	// If pas de point d'exclamation, alors c'est un message
	return p.cmdName // TODO
}

func (p *ArgProvider) NextString() (string, error) {
	p.index++

	if p.HasMore() {
		result, ok := p.args[p.index].(string)
		if ok {
			return result, nil
		}
		return "", errors.New(fmt.Sprintf("Non matching type: %T (expecting string)", result))
	}
	return "", errors.New("No more value")
}

func (p *ArgProvider) NextInt() (int, error) {
	p.index++

	if p.HasMore() {
		result, ok := p.args[p.index].(int)
		if ok {
			return result, nil
		}
		return 0, errors.New(fmt.Sprintf("Non matching type: %T (expecting int)", result))
	}
	return 0, errors.New("No more value") // TODO
}

func (p *ArgProvider) HasMore() bool {
	return p.index < p.argSize
}
