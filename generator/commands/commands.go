package commands

type Command interface {
	Read(text string)
}

const MapperTag string = "@mapper"

type MapperCommand struct {
	TargetFile Attribute
}

const MappingTag string = "@mapping"

type MappingCommand struct {
}

// TranslationTag
// from="X",to="Y"
const TranslationTag string = "@translate"

type TranslationCommand struct {
	From Attribute
	To   Attribute
}

// ExpressionTag
// target="X",expression="Y"
const ExpressionTag string = "@expression"

type ExpressionCommand struct {
	Target     Attribute
	Expression Attribute
}
