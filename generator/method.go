package generator

import (
	"AutoMapper/generator/commands"
	mappings "AutoMapper/generator/mappings"
	"fmt"
	"go/types"
)

type Method struct {
	Name          string
	Target        mappings.Type
	ErrorHandling bool
	Params        []mappings.Type
	Commands      []commands.Command
}

func (m Method) HasErrorHandling() bool {
	return false //TODO: check if method has return
}

func (m Method) GenerateMapping(project Project) string {
	//config := commands.FilterCommand[*commands.MappingCommand](m.Commands)
	var output string

	var targetStruct Structure

	type InputStructure struct {
		Structure
		ArgumentName string
	}

	var sourceStructs []InputStructure

	var params []mappings.Type

	//Iterate over all known structs
	for _, structure := range project.Structs {
		sName := structure.Package + structure.Name
		if sName == m.Target.Package+m.Target.Name {
			targetStruct = structure
			continue
		}

		for _, param := range m.Params {
			if sName == param.Package+param.Name {
				sourceStructs = append(sourceStructs, InputStructure{
					Structure:    structure,
					ArgumentName: param.ArgumentName,
				})
				break
			}
		}
	}

	for _, param := range m.Params {
		isStruct := false
		for _, str := range sourceStructs {
			sName := str.Package + "." + str.Name
			if param.GetTypeName() == sName {
				isStruct = true
				break
			}
		}

		if !isStruct {
			params = append(params, param)
		}
	}

	overrides := commands.FilterCommands[commands.OverrideCommand](m.Commands)
	globalOverrides := commands.FilterCommands[commands.OverrideCommand](project.GlobalCommands)
	overrides = append(overrides, globalOverrides...)

	fields := targetStruct.Fields
	targetMappings := StructFieldsToMappings(fields)

	type SourceMapping struct {
		parent InputStructure
		nodes  []*mappings.MappingNode
	}

	var sourceMappings []SourceMapping
	for _, sourceStruct := range sourceStructs {
		sourceMappings = append(sourceMappings, SourceMapping{parent: sourceStruct, nodes: StructFieldsToMappings(sourceStruct.Fields)})
	}

	for _, override := range overrides {
		mapping, path := mappings.Find("", targetMappings, override.IsOverrideTarget)
		if mapping == nil {
			fmt.Println("Failed to find override target, ")
			continue
		}
		mapping.Source.Mapped = true
		mapping.Source.Source = override.OverrideSource(mapping, path)
	}

	for _, target := range targetMappings {
		if target.Source.Mapped {
			continue
		}
		target.Inspect("target.", func(targetFullPath string, targetNode *mappings.MappingNode) bool {
			for _, sourceMapping := range sourceMappings {
				for _, field := range sourceMapping.nodes {
					foundFinalMapping := false

					field.Inspect(sourceMapping.parent.ArgumentName+".", func(sourceFullPath string, sourceNode *mappings.MappingNode) bool {
						foundMapper, nodeMapper, nodeMapping := project.MapperInterfaces.GetFittingMapper(sourceNode.TargetType, targetNode.TargetType)
						if foundMapper {
							targetNode.Source.Mapped = true
							targetNode.Source.Source = fmt.Sprintf("%s.%s(%s)", "%*__*%"+nodeMapper.PackagePath(&project), nodeMapping.Name+"%*__*%", sourceFullPath)
							foundFinalMapping = true
							return false
						}

						if targetNode.TargetType.GetFullName() == sourceNode.TargetType.GetFullName() {
							targetNode.Source.Mapped = true
							targetNode.Source.Source = sourceFullPath
							foundFinalMapping = true
							return false
						}

						return true
					})

					if foundFinalMapping {
						return false
					}
				}
			}
			return true
		})
	}

	for _, target := range targetMappings {
		mappingCount := 0
		mappings.InspectStacked(target, "", nil, MappingCreateInspect(&output, &mappingCount))
	}

	return output
}

type StackContext struct {
	arrayDepth    int
	arrayFullPath string
}

func MappingCreateInspect(output *string, mappingCount *int) mappings.InspectionStackedFunc[StackContext] {
	return func(fullPath string, node *mappings.MappingNode, c *StackContext) bool {
		switch node.TargetType.Underlying {
		case mappings.PointerType:
			fallthrough
		case mappings.DefaultType:
			if !node.Source.Mapped {
				*output += fmt.Sprintf("\n\t//target.%s is not mapped", fullPath)
				return true
			}
			*mappingCount++
			*output += fmt.Sprintf("\n\ttarget.%s = %s", fullPath, node.Source.Source)
			break
		case mappings.ArrayType:
			arrayIndex := fmt.Sprintf("i%d", c.arrayDepth)
			c.arrayDepth++

			arrayMappingCount := 0
			arrayOutput := ""

			arrayOutput += fmt.Sprintf("\n\tfor %s := range target.%s {", arrayIndex, fullPath)

			if node.Source.Mapped {
				arrayMappingCount++
				arrayOutput += fmt.Sprintf("\n\ttarget.%s[%s] = %s", fullPath, arrayIndex, node.Source.Source)
			}

			for _, arrayChild := range node.Children {
				mappings.InspectStacked(arrayChild, fullPath+"["+arrayIndex+"].", c, MappingCreateInspect(&arrayOutput, &arrayMappingCount))
			}

			arrayOutput += fmt.Sprintf("\n\t}")

			if arrayMappingCount != 0 {
				*output += arrayOutput
				*mappingCount += arrayMappingCount
			} else {
				*output += fmt.Sprintf("\n\t//target.%s[%s] is not mapped", fullPath, arrayIndex)
			}

			return false
		}

		return true
	}
}

func StructFieldsToMappings(fields []*types.Var) []*mappings.MappingNode {
	var mappings []*mappings.MappingNode
	for _, field := range fields {
		if field.Exported() {
			mappings = append(mappings, NewMapping(field))
		}
	}

	return mappings
}

func NewMapping(field *types.Var) *mappings.MappingNode {
	t := field.Type()
	return MappingFromType(field, t)
}

func MappingFromType(field *types.Var, t types.Type) *mappings.MappingNode {
	var result mappings.MappingNode
	switch t.(type) {
	case *types.Pointer:
		pResult := MappingFromType(field, t.(*types.Pointer).Elem())
		pResult.TargetType.Underlying = mappings.PointerType
		return pResult
	case *types.Slice:
		pResult := MappingFromType(field, t.(*types.Slice).Elem())
		pResult.TargetType.Underlying = mappings.ArrayType
		return pResult
	case *types.Basic:
		basic := t.(*types.Basic)

		result.TargetType = mappings.Type{
			ArgumentName: field.Name(),
			Name:         basic.Name(),
			Package:      "--go--",
		}
		break
	case *types.Struct:
		struc := t.(*types.Struct)
		result.Children = StructFieldsToMappings(GetStructureFields(struc))
		break
	case *types.Named:
		named := t.(*types.Named)
		result = *MappingFromType(field, named.Underlying())

		object := named.Obj()
		result.TargetType = mappings.Type{
			ArgumentName: field.Name(),
			Name:         object.Name(),
			Package:      object.Pkg().Path(),
		}
		break
	default:
		fmt.Printf("Unknown type detected (%s)\n", t.String())
		break
	}

	return &result
}
