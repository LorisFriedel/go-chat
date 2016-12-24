package console

import (
	"fmt"
	"regexp"
	"strings"
)

type IParser interface {
	// le parser, on lui donne l'input et il nous permettra de récup la valeur de l'input, cad si
	// cest un simple message le message, sinon le nom de la commande et les arguments, avec surement
	// un getInt, getString, etc..
	// implicitement, écrire un message c'est comme marquer !message "Le message blabla"

	Parse(input string) (IProvider, error)
}

type CmdParser struct {
	prefix string
	// TODO
}

func NewParser(prefix string) *CmdParser {
	return &CmdParser{prefix}
}

func (p *CmdParser) Parse(input string) (IProvider, error) {
	// Special case for message, user can omit the message command
	if !strings.HasPrefix(input, p.prefix) {
		return NewProvider("message", []string{input}), nil
	}

	reg, err := regexp.Compile(fmt.Sprintf("%s([a-zA-Z]+)(?: (.+))?", p.prefix))
	if err != nil {
		return nil, err
	}

	regResult := reg.FindStringSubmatch(input)
	if regResult == nil || len(regResult) < 2 {
		return nil, fmt.Errorf("CmdParser.Parse: missing command in: %s", input)
	}

	cmdName := regResult[1]
	argv := strings.Split(regResult[2], " ")

	return NewProvider(cmdName, argv), nil
}
