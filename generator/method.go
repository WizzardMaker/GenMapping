package generator

import (
	"AutoMapper/generator/commands"
	mappings2 "AutoMapper/generator/mappings"
	"fmt"
	"go/types"
)

type Method struct {
	Name          string
	Target        mappings2.Type
	ErrorHandling bool
	Params        []mappings2.Type
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

	overrides := commands.FilterCommands[commands.OverrideCommand](m.Commands)

	fields := targetStruct.Fields
	targetMappings := StructFieldsToMappings(fields)

	for _, override := range overrides {
		mapping := mappings2.Find("target", targetMappings, override.IsOverrideTarget)
		if mapping == nil {
			fmt.Println("Failed to find override target, ")
			continue
		}
		mapping.Source.Mapped = true
		mapping.Source.Source = override.OverrideSource()
	}

	for _, target := range targetMappings {
		if target.Source.Mapped {
			continue
		}
		target.Inspect("target", func(fullPath string, node *mappings2.MappingNode) bool {

			return true
		})
	}

	for _, target := range targetMappings {
		target.Inspect("target", func(fullPath string, node *mappings2.MappingNode) bool {
			if !node.Source.Mapped {
				output += fmt.Sprintf("\n\t//target.%s is not mapped", fullPath)
				return true
			}

			output += fmt.Sprintf("\n\ttarget.%s = %s", fullPath, node.Source.Source)
			return true
		})
	}

	return output
}

func StructFieldsToMappings(fields []*types.Var) []*mappings2.MappingNode {
	var mappings []*mappings2.MappingNode
	for _, field := range fields {
		if field.Exported() {
			mappings = append(mappings, NewMapping(field))
		}
	}

	return mappings
}

func NewMapping(field *types.Var) *mappings2.MappingNode {
	t := field.Type()
	return MappingFromType(field, t)
}

func MappingFromType(field *types.Var, t types.Type) *mappings2.MappingNode {
	var result mappings2.MappingNode
	switch t.(type) {
	case *types.Pointer:
		pResult := MappingFromType(field, t.(*types.Pointer).Elem())
		pResult.TargetType.Name = "*" + pResult.TargetType.Name
		return pResult
	case *types.Slice:
		pResult := MappingFromType(field, t.(*types.Slice).Elem())
		pResult.TargetType.Name = "[]" + pResult.TargetType.Name
		return pResult
	case *types.Basic:
		basic := t.(*types.Basic)

		result.TargetType = mappings2.Type{
			ArgumentName: field.Name(),
			Name:         basic.Name(),
			Package:      "--go--",
		}
		break
	case *types.Named:
		named := t.(*types.Named)
		struc, ok := named.Underlying().(*types.Struct)
		if !ok {
			fmt.Println("Unknown type!")
		}

		object := named.Obj()

		result.TargetType = mappings2.Type{
			ArgumentName: field.Name(),
			Name:         object.Name(),
			Package:      object.Pkg().Path(),
		}

		result.Children = StructFieldsToMappings(GetStructureFields(struc))

		break
	default:
		fmt.Printf("Unknown type detected (%s)\n", t.String())
		break
	}

	return &result
}
