package console

import "strings"

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

func NewCmdParser(prefix string) *CmdParser {
	return &CmdParser{prefix}
}

func (p *CmdParser) Parse(input string) (IProvider, error) {
	// Special case for message, user can omit the message command
	if !strings.HasPrefix(p.prefix, input) {
		input = p.prefix + "message " + input
	}

	return NewArgProvider("message", []interface{}{"Ceci est un message test"}), nil

	// TODO isolate command name

	// TODO put args in a slice of interface{}, when get -> assertion on type + iterator incrementation
	// TODO
	// TODO
	// TODO prefix + if not prefix put !message
	//return nil, nil // TODO
}