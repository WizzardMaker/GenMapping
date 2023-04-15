package generator

import (
	"AutoMapper/generator/commands"
	"fmt"
	"go/types"
)

type Method struct {
	Name          string
	Target        Type
	ErrorHandling bool
	Params        []Type
	Commands      []commands.Command
}

func (m Method) HasErrorHandling() bool {
	return false //TODO: check if method has return
}

func (m Method) GenerateMapping(project Project) string {
	//config := commands.FilterCommand[*commands.MappingCommand](m.Commands)
	var output string

	var targetStruct Structure
	var sourceStructs []Structure

	for _, structure := range project.Structs {
		sName := structure.Package + structure.Name
		if sName == m.Target.Package+m.Target.Name {
			targetStruct = structure
			continue
		}

		for _, param := range m.Params {
			if sName == param.Package+param.Name {
				sourceStructs = append(sourceStructs, structure)
				break
			}
		}
	}

	translations := commands.FilterCommands[*commands.TranslationCommand](m.Commands)

	for _, field := range targetStruct.Fields {
		t := field.Type()
		switch t.(type) {
		case *types.Basic:
			//basic := t.(*types.Basic)

			trans := findTranslation(translations, field.Name())
			if len(trans) != 0 {
				output += fmt.Sprintf("\n\t%s = %s", trans[0].From.Value, m.Target.ArgumentName+"."+field.Name())
			} else {
				source := findSource(m.Params, field.Name())
				output += fmt.Sprintf("\n\t%s = %s", source.Name, m.Target.ArgumentName+"."+field.Name())
			}

			break
		}
	}

	return output
}

func findTranslation(translations []*commands.TranslationCommand, to string) []*commands.TranslationCommand {
	var output []*commands.TranslationCommand
	for _, translation := range translations {
		if translation.To.Value == to {
			output = append(output, translation)
		}
	}

	return output
}
func findSource(types []Type, to string) Type {
	for _, sourceType := range types {
		if sourceType.Name == to {
			return sourceType
		}
	}

	return Type{}
}
