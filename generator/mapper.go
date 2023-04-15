package generator

import (
	"AutoMapper/generator/commands"
	"fmt"
	"go/ast"
	"go/types"
)

type Mapper struct {
	Interface ast.InterfaceType
	Name      string
	Methods   []Method
	Imports   []Import
	Commands  []commands.Command
}

func NewMethods(methodList *ast.FieldList, currentPackage string) (methods []Method) {
	for _, field := range methodList.List {
		errorHappened := false

		funcType, ok := field.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		var commandList []commands.Command

		errorHandling := false

		var target Type
		for i, t := range NewTypes(funcType.Results, currentPackage) {
			if t.Name == "error" {
				if i == 0 {
					fmt.Printf("Error handling is not allowed as first argument | Method:%s", field.Names[0].Name)
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
			Name:          field.Names[0].Name,
			Target:        target,
			Params:        NewTypes(funcType.Params, currentPackage),
			Commands:      commandList,
			ErrorHandling: errorHandling,
		})
	}

	return
}

func NewTypes(decl *ast.FieldList, currentPackage string) (types []Type) {
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

		types = append(types, Type{
			ArgumentName: argumentName,
			Name:         typeExpr.Sel.Name,
			Package:      packageName,
		})
	}

	return
}

func NewStructure(spec *ast.TypeSpec, structType *ast.StructType, info *types.Info, currentPackage string) Structure {
	//structInfo, _ :=.(*types.Struct)
	var structInfo *types.Struct
	for expr, value := range info.Defs {
		if expr.Obj != nil && expr.Obj.Name == spec.Name.Name {
			typ := value.Type().(*types.Named)
			structInfo = typ.Underlying().(*types.Struct)
			break
		}
	}

	if structInfo == nil {
		fmt.Printf("Failed to find %s in info map\n", spec.Name.Name)
		return Structure{}
	}

	fieldCount := structInfo.NumFields()

	//var fields []mappings.Field
	//for _, fieldItem := range structType.Fields.List {
	//	var field mappings.Field
	//	expr := fieldItem.Type
	//	switch expr.(type) {
	//	case *ast.ArrayType:
	//		field = mappings.NewArrayField(fieldItem, currentPackage)
	//		break
	//	case *ast.Ident: //Direct fieldItem without package declaration
	//		field = mappings.NewObjectField(fieldItem.Type, fieldItem.Names[0].Name, currentPackage)
	//		break
	//	case *ast.SelectorExpr: //fieldItem with package declaration
	//		field = mappings.NewObjectField(fieldItem.Type, fieldItem.Names[0].Name, currentPackage)
	//		break
	//	}
	//
	//	fields = append(fields, field)
	//}

	var fields []*types.Var
	for i := 0; i < fieldCount; i++ {
		field := structInfo.Field(i)
		fields = append(fields, field)
	}

	return Structure{
		Package: currentPackage,
		Name:    spec.Name.Name,
		Fields:  fields,
	}
}

type Structure struct {
	Package string
	Name    string
	Fields  []*types.Var
}
