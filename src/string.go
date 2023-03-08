package main

import (
	"strconv"
	"strings"
)

var stringCompile = makeRegex("(( *)\"((\\\\([a-z\\\"'`]))|[^\\\"])*\"( *))|(( *)'((\\\\([a-z\\'\"`]))|[^\\'])*'( *))")

func isString(code UNPARSEcode) bool {
	return stringCompile.MatchString(code.code)
}

func unquoted(
	str string,
) (string, error) {
	str = strings.Trim(str, " ")
	if str[0] == '\'' {
		str = strings.Replace(str, "\\\"", "\"", -1)
		str = strings.Replace(str, "\"", "\\\"", -1)
	}
	str = str[1 : len(str)-1]
	str = strings.Replace(str, "\\'", "'", -1)
	str = "\"" + str + "\""
	return strconv.Unquote(str)
}

// returns translateString, success, error
func parseString(code UNPARSEcode) (string, bool, ArErr, int) {
	trim := strings.Trim(code.code, " ")

	unquoted, err := unquoted(trim)
	if err != nil {
		return "", false, ArErr{"Syntax Error", "invalid string", code.line, code.path, code.realcode, true}, 1
	}

	return unquoted, true, ArErr{}, 1
}