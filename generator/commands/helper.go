package commands

import (
	"fmt"
	"regexp"
	"strings"
)

func getCommandText(text string, command Tag) (string, string) {
	pattern := fmt.Sprintf("%s\\(([^@]*)\\)", command)
	return getMultilineRegexResult(text, pattern)
}

func getMultilineRegexResult(text string, pattern string) (string, string) {
	pattern = "(?m)" + pattern
	r, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}
	results := r.FindStringSubmatch(text)

	if len(results) < 2 {
		return "", ""
	}

	result := strings.Replace(strings.TrimSpace(results[1]), "//", "", -1)
	return result, results[0]
}

func FilterCommand[T Command](commands []Command) T {
	for _, command := range commands {
		if res, ok := command.(T); ok {
			return res
		}
	}

	var nothing T
	return nothing
}

func FilterCommands[T Command](commands []Command) []T {
	var result []T
	for _, command := range commands {
		if res, ok := command.(T); ok {
			result = append(result, res)
		}
	}

	return result
}
