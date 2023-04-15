package generator

import (
	"AutoMapper/generator/commands"
	"fmt"
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
	//config := commands.HasCommand[*commands.MappingCommand](m.Commands)
	var output string

	var targetStruct Structure

	for _, structure := range project.Structs {
		if structure.Package+structure.Name == m.Target.Package+m.Target.Name {
			targetStruct = structure
		}
	}

	for _, field := range targetStruct.Fields {
		fmt.Print(field)
	}

	return output
}
