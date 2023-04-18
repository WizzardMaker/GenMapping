package commands

import (
	"fmt"
	"strconv"
)

type Attribute struct {
	Name    string
	Value   string
	Default string
}

func (a *Attribute) FromCommandText(name, text string) {
	a.Name = name

	pattern := fmt.Sprintf("%s=\"([\\w\\W]*?)\"", a.Name)
	a.Value, _ = getMultilineRegexResult(text, pattern)

	if a.Value == "" {
		a.Value = a.Default
	}
}

func (a *Attribute) Bool() bool {
	b, ok := strconv.ParseBool(a.Value)

	return ok == nil && b
}
