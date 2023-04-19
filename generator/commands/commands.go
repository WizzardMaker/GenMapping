package commands

import (
	"GenMapping/generator/mappings"
	"strings"
)

type Command interface {
	Read(text string) string
}

type OverrideCommand interface {
	Command
	IsOverrideTarget(node *mappings.MappingNode, path string) bool
	OverrideSource(node *mappings.MappingNode, nodePath string) string
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

type CommandParser func(text string) []Command

func NewCommand[T any, PT interface {
	Read(text string) string
	*T
}](tag Tag) CommandParser {
	return func(text string) []Command {
		commands := strings.Count(text, string(tag))

		var result []Command
		for i := 0; i < commands; i++ {
			var command T
			var commandInterface PT
			commandInterface = &command
			readText := commandInterface.Read(text)
			text = strings.Replace(text, readText, "", 1)
			result = append(result, commandInterface)
		}

		return result
	}
}

type Tag string

type CommandMap map[Tag]CommandParser

var commandMap = CommandMap{
	MapperTag:      NewCommand[MapperCommand](MapperTag),
	MappingTag:     NewCommand[MappingCommand](MappingTag),
	TranslationTag: NewCommand[TranslationCommand](TranslationTag),
	ExpressionTag:  NewCommand[ExpressionCommand](ExpressionTag),
}
