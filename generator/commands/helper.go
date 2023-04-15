package commands

import (
	"fmt"
	"regexp"
	"strings"
)

func getCommandText(text string, command Tag) string {
	pattern := fmt.Sprintf("%s\\(([\\w\\W]*?)\\)", command)
	return getMultilineRegexResult(text, pattern)
}

func getMultilineRegexResult(text string, pattern string) string {
	pattern = "(?m)" + pattern
	r, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}
	results := r.FindStringSubmatch(text)

	if len(results) < 2 {
		return ""
	}

	result := strings.Replace(strings.TrimSpace(results[1]), "//", "", -1)
	fmt.Println("Found command: ", result)
	return result
}
