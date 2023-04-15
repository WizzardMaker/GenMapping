package commands

import (
	"AutoMapper/generator/mappings"
	"strings"
)

type Command interface {
	Read(text string)
}

type OverrideCommand interface {
	Command
	IsOverrideTarget(node *mappings.MappingNode, path string) bool
	OverrideSource() string
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

type Tag string

type CommandMap map[Tag]func(text string) []Command

var commandMap = CommandMap{
	MapperTag:      NewCommand[MapperCommand](MapperTag),
	MappingTag:     NewCommand[MappingCommand](MappingTag),
	TranslationTag: NewCommand[TranslationCommand](TranslationTag),
	ExpressionTag:  NewCommand[ExpressionCommand](ExpressionTag),
}
