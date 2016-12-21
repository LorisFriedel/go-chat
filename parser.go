package main

type Parser interface {
	// le parser, on lui donne l'input et il nous permettra de récup la valeur de l'input, cad si
	// cest un simple message le message, sinon le nom de la commande et les arguments, avec surement
	// un getInt, getString, etc..
	// implicitement, écrire un message c'est comme marquer !message "Le message blabla"
}
