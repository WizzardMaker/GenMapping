package commands

import (
	"AutoMapper/generator/mappings"
	"strings"
)

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

func FromText(text string, allowedTags ...Tag) []Command {
	var commands []Command

	if len(allowedTags) == 0 {
		for _, ctor := range commandMap {
			commands = append(commands, ctor(text)...)
		}
	} else {
		for _, tag := range allowedTags {
			commands = append(commands, commandMap[tag](text)...)
		}
	}
	return commands
}

type Tag string

const MapperTag Tag = "@mapper"

type MapperCommand struct {
	TargetFile Attribute
}

func (c *MapperCommand) Read(text string) {
	_ = getCommandText(text, MapperTag)
}

const MappingTag Tag = "@mapping"

type MappingCommand struct {
}

func (c *MappingCommand) Read(text string) {
	_ = getCommandText(text, MappingTag)
}

type OverrideCommand interface {
	Command
	IsOverrideTarget(node *mappings.MappingNode, path string) bool
	OverrideSource() string
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

func NewCommand[T any, PT interface {
	Read(text string)
	*T
}](tag Tag) func(text string) []Command {
	return func(text string) []Command {
		commands := strings.Count(text, string(tag))

		var result []Command
		for i := 0; i < commands; i++ {
			var n T
			var pn PT
			pn = &n
			pn.Read(text)
			result = append(result, pn)
		}

		return result
	}
}

type CommandMap map[Tag]func(text string) []Command

var commandMap CommandMap = CommandMap{
	MapperTag:      NewCommand[MapperCommand](MapperTag),
	MappingTag:     NewCommand[MappingCommand](MappingTag),
	TranslationTag: NewCommand[TranslationCommand](TranslationTag),
	ExpressionTag:  NewCommand[ExpressionCommand](ExpressionTag),
}
