package commands

type Attribute struct {
	Name  string
	Value string
}

func (a *Attribute) FromCommandText(text string) {
	a.Value = text //TODO: Name="Value"
}
