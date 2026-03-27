package main

import (
	"log"
)

type token struct {
	kind  string
	value string
}

/*
 Lexer
*/

func tokenize(input string) []token {

	input += "\n"

	now := 0

	tokens := []token{}
	for now < len([]rune(input)) {
		char := string(input[now])

		if char == "(" { // Checks for left parenthesis
			tokens = append(tokens, token{kind: "paren", value: char})

			now++
			continue
		}

		if char == ")" { // Checks for right parenthesis
			tokens = append(tokens, token{kind: "paren", value: char})

			now++
			continue
		}

		if char == "\n" || char == " " || char == "\t" { // Checks for whitespace, then ignore
			now++
			continue
		}

		if isNumber(char) { // Checks for numbers, then adds them to the token list
			value := ""

			for isNumber(char) {
				value += char
				now++
				char = string([]rune(input)[now])
			}

			tokens = append(tokens, token{
				kind:  "number",
				value: value,
			})

			continue
		}

		if isLetter(char) { // Checks for symbols, then adds them to the token list
			value := ""

			for isLetter(char) {
				value += char
				now++
				char = string([]rune(input)[now])
			}

			tokens = append(tokens, token{
				kind:  "symbol",
				value: value,
			})
			continue
		}
		break
	}

	return tokens // Return the list of tokens
}

func isNumber(char string) bool {
	if char == "" {
		return false
	}
	n := []rune(char)[0]
	if n >= '0' && n <= '9' {
		return true
	}
	return false
}

func isLetter(char string) bool {
	if char == "" {
		return false
	}
	n := []rune(char)[0]
	if n >= 'a' && n <= 'z' || n >= 'A' && n <= 'Z' {
		return true
	}
	return false
}

/*
 Parser
*/

type node struct {
	kind       string
	value      string
	name       string
	callee     *node
	expression *node
	body       []node
	parameters []node
	arguments  *[]node
	context    *[]node
}

type ast node

var pc int
var pt []token

func parse(tokens []token) ast {
	pc = 0
	pt = tokens

	ast := ast{
		body: []node{},
		kind: "program",
	}

	for pc < len(pt) {
		ast.body = append(ast.body, parseExpression())
	}

	return ast
}

func parseExpression() node {
	token := pt[pc]

	if token.kind == "number" {
		pc++
		return node{
			kind:  "number",
			value: token.value,
		}
	}

	if token.kind == "paren" && token.value == "(" {

		pc++
		token = pt[pc]

		n := node{
			kind:       "call",
			name:       token.value,
			parameters: []node{},
		}
		pc++
		token = pt[pc]

		for token.kind != "paren" || (token.kind == "paren" && token.value == ")") {
			{
				n.parameters = append(n.parameters, parseExpression())
				token = pt[pc]
			}

			pc++
		}

		return n
	}

	log.Fatal(token.kind + " is not a valid token")
	return node{}
}

/* Traverser */

type visitor map[string]func(node *node, parent node)

func traverse(ast ast, visitor visitor) {

	traverseNode(node(ast), node{}, visitor)
}

func traverseArray(a []node, parent node, visitor visitor) {

	for _, child := range a {
		traverseNode(child, parent, visitor)
	}
}

func traverseNode(node node, parent node, visitor visitor) {

	for k, va := range visitor {
		if node.kind == k {
			va(&node, parent)
		}
	}

	switch node.kind {
	case "program":
		traverseArray(node.body, node, visitor)
		break

	case "call":
		traverseArray(node.parameters, node, visitor)
		break

	case "number":
		break

	default:
		log.Fatal(node.kind + " is not a valid node")
	}
}
