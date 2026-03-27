package main

type token struct {
	kind  string
	value string
}

func tokenize(input string) []token {

	input += "\n"

	now := 0

	tokens := []token{}
	for now < len([]rune(input)) {
		char := string(input[now])

		if char == "(" { // Checks for left parenthesis
			tokens = append(tokens, token{kind: "lparen", value: char})

			now++
			continue
		}

		if char == ")" { // Checks for right parenthesis
			tokens = append(tokens, token{kind: "rparen", value: char})

			now++
			continue
		}

		if char == "\n" || char == " " || char == "\t" { // Checks for whitespace, then ignore
			now++
			continue
		}

		if isNumber(char) { // Checks for numbers, then adds them to the token list
			num := ""

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
