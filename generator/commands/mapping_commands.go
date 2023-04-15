package commands

import "AutoMapper/generator/mappings"

const MappingTag Tag = "@mapping"

type MappingCommand struct {
}

func (c *MappingCommand) Read(text string) {
	_ = getCommandText(text, MappingTag)
}

// TranslationTag
// from="X",to="Y"
const TranslationTag Tag = "@translate"

type TranslationCommand struct {
	From Attribute
	To   Attribute
}

func (c *TranslationCommand) Read(text string) {
	commandText := getCommandText(text, TranslationTag)
	c.From.FromCommandText("from", commandText)
	c.To.FromCommandText("to", commandText)
}

func (c *TranslationCommand) IsOverrideTarget(node *mappings.MappingNode, path string) bool {
	return c.To.Value == path
}

func (c *TranslationCommand) OverrideSource() string {
	return c.From.Value
}

// ExpressionTag
// target="X",expression="Y"
const ExpressionTag Tag = "@expression"

type ExpressionCommand struct {
	Target     Attribute
	Expression Attribute
}

func (c *ExpressionCommand) Read(text string) {
	commandText := getCommandText(text, ExpressionTag)
	c.Target.FromCommandText("target", commandText)
	c.Expression.FromCommandText("expression", commandText)
}

func (c *ExpressionCommand) IsOverrideTarget(node *mappings.MappingNode, path string) bool {
	return c.Target.Value == path
}

func (c *ExpressionCommand) OverrideSource() string {
	return c.Expression.Value
}

var PerMappingTags = []Tag{MappingTag, TranslationTag, ExpressionTag}
