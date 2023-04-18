package generator

import (
	"AutoMapper/generator/commands"
	"AutoMapper/generator/mappings"
	"fmt"
	"go/ast"
	"go/types"
	"path/filepath"
	"strings"
)

type Mapper struct {
	Interface  ast.InterfaceType
	Name       string
	outputPath string
	Methods    []Method
	Imports    []mappings.Import
	Commands   []commands.Command
}

func (m Mapper) OutputPath(project *Project) string {
	outputDirName := "mapper"
	//TODO: Add global commands
	//if project != nil {
	//	projectConfig := commands.FilterCommand[*commands.MapperCommand](project.GlobalCommands)
	//	if projectConfig != nil && projectConfig.
	//}
	config := commands.FilterCommand[*commands.MapperCommand](m.Commands)

	if config != nil {
		if config.TargetFile.Value != "" {
			return config.TargetFile.Value
		}
	}

	return m.outputPath + "/" + outputDirName + "/" + m.Name + "/mapper" + ".go"
}

func (m Mapper) PackagePath(project *Project) string {
	//C:/src/module/X/mapper
	path := m.OutputPath(project)

	//X/mapper
	path = strings.TrimPrefix(path, project.BasePath)
	path = filepath.Dir(path)

	//ModuleName/X/mapper
	return strings.Replace(project.ModuleName+path, "\\", "/", -1)
}

func NewMethods(methodList *ast.FieldList, currentPackage string, info *types.Info) (methods []Method) {
	for _, method := range methodList.List {
		errorHappened := false

		funcType, ok := method.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		test := info.TypeOf(funcType)
		fmt.Println(test)

		commandList := commands.FromText(method.Doc.Text(), commands.PerMappingTags...)

		errorHandling := false

		var target mappings.Type
		for i, t := range NewTypes(funcType.Results, currentPackage) {
			if t.Name == "error" {
				if i == 0 {
					fmt.Printf("Error handling is not allowed as first argument | Method:%s", method.Names[0].Name)
					errorHappened = true
					break
				}

				errorHandling = true
				continue
			}

			if t.ArgumentName == "" {
				target = t
				target.ArgumentName = "target"
			}
		}
		if errorHappened {
			continue
		}

		methods = append(methods, Method{
			Name:          method.Names[0].Name,
			Target:        target,
			Params:        NewTypes(funcType.Params, currentPackage),
			Commands:      commandList,
			ErrorHandling: errorHandling,
		})
	}

	return
}

type Mappers []Mapper

func (m Mappers) GetFittingMapper(from, to mappings.Type) (bool, Mapper, Method) {
	for _, mapper := range m {
		for _, method := range mapper.Methods {
			if len(method.Params) == 1 {
				if method.Params[0].GetTypeName() == from.GetTypeName() && method.Target.GetTypeName() == to.GetTypeName() {
					return true, mapper, method
				}
			}
		}
	}

	return false, Mapper{}, Method{}
}

func NewTypes(decl *ast.FieldList, currentPackage string) (types []mappings.Type) {
	for _, field := range decl.List {
		argumentName := ""
		if len(field.Names) != 0 {
			argumentName = field.Names[0].Name
		}
		packageName := currentPackage

		typeExpr, ok := field.Type.(*ast.SelectorExpr)
		if !ok {
			continue
		}

		packageType, ok := typeExpr.X.(*ast.Ident)
		if ok {
			packageName = packageType.Name
		}

		types = append(types, mappings.Type{
			ArgumentName: argumentName,
			Name:         typeExpr.Sel.Name,
			Package:      packageName,
		})
	}

	return
}

func NewStructure(spec *ast.TypeSpec, info *types.Info, currentPackage string) Structure {
	//structInfo, _ :=.(*types.Struct)
	var structInfo *types.Struct
	for expr, value := range info.Defs {
		if expr.Obj != nil && expr.Obj.Name == spec.Name.Name {
			typ, ok := value.Type().(*types.Named)
			if !ok {
				continue
			}
			structInfo = typ.Underlying().(*types.Struct)
			currentPackage = typ.Obj().Pkg().Path()
			break
		}
	}

	if structInfo == nil {
		fmt.Printf("Failed to find %s in info map\n", spec.Name.Name)
		return Structure{}
	}

	fields := GetStructureFields(structInfo)

	return Structure{
		Package: currentPackage,
		Name:    spec.Name.Name,
		Fields:  fields,
	}
}

func GetStructureFields(structInfo *types.Struct) []*types.Var {
	fieldCount := structInfo.NumFields()
	var fields []*types.Var
	for i := 0; i < fieldCount; i++ {
		field := structInfo.Field(i)
		fields = append(fields, field)
	}
	return fields
}

type Structure struct {
	Package string
	Name    string
	Fields  []*types.Var
}
