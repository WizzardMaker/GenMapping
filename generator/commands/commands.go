package commands

type Command interface {
	Read(text string)
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

const MapperTag string = "@mapper"

type MapperCommand struct {
	TargetFile Attribute
}

func (c *MapperCommand) Read(text string) {

}

const MappingTag string = "@mapping"

type MappingCommand struct {
}

func (c *MappingCommand) Read(text string) {

}

// TranslationTag
// from="X",to="Y"
const TranslationTag string = "@translate"

type TranslationCommand struct {
	From Attribute
	To   Attribute
}

func (c *TranslationCommand) Read(text string) {

}

// ExpressionTag
// target="X",expression="Y"
const ExpressionTag string = "@expression"

type ExpressionCommand struct {
	Target     Attribute
	Expression Attribute
}

func (c *ExpressionCommand) Read(text string) {

}
