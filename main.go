package main

import (
	"fmt"
	"log"
	"strings"
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

		for token.kind != "paren" && token.value == ")" {
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

/* Transformer */

func transform(a ast) ast {
	newAst := ast{
		body: []node{},
		kind: "program",
	}

	a.context = &newAst.body

	traverse(a, map[string]func(node *node, parent node){
		"number": func(n *node, parent node) {
			*parent.context = append(*parent.context, node{
				kind:  "number",
				value: n.value,
			})
		},

		"call": func(n *node, parent node) {
			expression := node{
				kind:      "call",
				callee:    &node{name: n.name, kind: "identifier"},
				arguments: new([]node),
			}

			n.context = expression.arguments

			if parent.kind != "call" {

				es := node{
					kind:       "expression",
					expression: &expression,
				}

				*parent.context = append(*parent.context, es)
			} else {
				*parent.context = append(*parent.context, expression)
			}

		},
	})

	return newAst
}

/* Code Generator */

func codeGenerator(node node) string {
	switch node.kind {
	case "program":
		var r []string
		for _, n := range node.body {
			r = append(r, codeGenerator(n))
		}
		return strings.Join(r, "\n")
	case "expression":
		return codeGenerator(*node.expression) + ";"
	case "call":
		var ra []string
		c := codeGenerator(*node.callee)
		for _, a := range *node.arguments {
			ra = append(ra, codeGenerator(a))
		}

		r := strings.Join(ra, ", ")
		return c + "(" + r + ")"

	case "identifier":
		return node.name
	case "number":
		return node.value
	default:
		log.Fatal(node.kind + " is not a valid node")
		return ""
	}
}

/* Compiler */

func compile(input string) string {
	tokens := tokenize(input)
	ast := parse(tokens)
	newAst := transform(ast)
	output := codeGenerator(node(newAst))

	return output
}

func main() {
	program := "(add 10 (subtract 4 2))"
	output := compile(program)
	fmt.Println(output)
}
