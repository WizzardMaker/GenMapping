package commands

import "GenMapping/generator/mappings"

const MappingTag Tag = "@mapping"

type MappingCommand struct {
}

func (c *MappingCommand) Read(text string) string {
	_, readText := getCommandText(text, MappingTag)

	return readText
}

// TranslationTag
// from="X",to="Y"
const TranslationTag Tag = "@translate"

type TranslationCommand struct {
	From Attribute
	To   Attribute
}

func (c *TranslationCommand) Read(text string) string {
	commandText, readText := getCommandText(text, TranslationTag)
	c.From.FromCommandText("from", commandText)
	c.To.FromCommandText("to", commandText)
	return readText
}

func (c *TranslationCommand) IsOverrideTarget(node *mappings.MappingNode, path string) bool {
	return c.To.Value == path
}

func (c *TranslationCommand) OverrideSource(*mappings.MappingNode, string) string {
	return c.From.Value
}

// ExpressionTag
// target="X",expression="Y",isType*="false"
const ExpressionTag Tag = "@expression"

type ExpressionCommand struct {
	Target     Attribute
	Expression Attribute
	IsType     Attribute //Optional - default: false
}

func (c *ExpressionCommand) Read(text string) string {
	commandText, readText := getCommandText(text, ExpressionTag)
	c.Target.FromCommandText("target", commandText)
	c.Expression.FromCommandText("expression", commandText)
	c.IsType.Default = "false"
	c.IsType.FromCommandText("isType", commandText)

	return readText
}

func (c *ExpressionCommand) IsOverrideTarget(node *mappings.MappingNode, path string) bool {
	if c.IsType.Bool() {
		return c.Target.Value == node.TargetType.GetTypeName()
	}

	return c.Target.Value == path
}

func (c *ExpressionCommand) OverrideSource(*mappings.MappingNode, string) string {
	return c.Expression.Value
}

var PerMappingTags = []Tag{MappingTag, TranslationTag, ExpressionTag}
