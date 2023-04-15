package commands

import "fmt"

type Attribute struct {
	Name  string
	Value string
}

func (a *Attribute) FromCommandText(name, text string) {
	a.Name = name

	pattern := fmt.Sprintf("%s=\"([\\w\\W]*?)\"", a.Name)
	a.Value = getMultilineRegexResult(text, pattern)
}
