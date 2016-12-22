package provider

import (
	"errors"
	"fmt"
	"strconv"
)

// Arguments provider type

type IProvider interface {
	HasMore() bool
	CommandName() string
	GetString() (string, error)
	GetInt() (int, error)
}

type ArgProvider struct {
	cmdName string
	args    []string
	index   int
	argSize int
}

// TODO reduce slice each NextXX ?

func New(cmdName string, args []string) *ArgProvider {
	return &ArgProvider{cmdName, args, -1, len(args)}
}

func (p *ArgProvider) CommandName() string {
	return p.cmdName
}

func (p *ArgProvider) GetString() (string, error) {
	p.index++

	if p.HasMore() {
		return p.args[p.index], nil
	}
	return "", errors.New("ArgProvider.GetString: no more value")
}

func (p *ArgProvider) GetInt() (int, error) {
	p.index++

	if p.HasMore() {
		if result, err := strconv.Atoi(p.args[p.index]); err == nil {
			return result, nil
		}
		return 0, fmt.Errorf("ArgProvider.GetInt: %v cannot be converted to int", p.args[p.index])
	}
	return 0, errors.New("ArgProvider.GetInt: no more value") // TODO
}

func (p *ArgProvider) HasMore() bool {
	return p.index < p.argSize
}
