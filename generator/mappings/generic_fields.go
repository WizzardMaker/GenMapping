package mappings

import (
	"AutoMapper/helper"
	"go/ast"
)

func NewArrayField(arrayType *ast.Field, currentPackage string) Field {
	typ := arrayType.Type.(*ast.ArrayType)

	return ArrayField{
		Name:         arrayType.Names[0].Name,
		MappingField: NewObjectField(typ.Elt, arrayType.Names[0].Name+"[%d]", currentPackage),
	}
}

type ArrayField struct {
	Name         string
	MappingField Field
}

func (f ArrayField) GetName() string {
	return f.Name
}

func (f ArrayField) GetType() string {
	return "[]" + f.MappingField.GetType()
}

func NewObjectField(typeExpr ast.Expr, name string, currentPackage string) Field {
	ident, ok := typeExpr.(*ast.Ident)

	if ok {
		field := DirectField{
			Package: currentPackage,
		}

		if ident.Obj != nil {
			field.Name = name
			field.Type = ident.Obj.Name
		} else {
			field.Type = ident.Name
			field.Name = name
		}

		if helper.IsBasicDataType(field.Type) {
			field.Package = "--go--"
		}

		return field
	}

	selector, ok := typeExpr.(*ast.SelectorExpr)
	if ok {
		return DirectField{
			Type:    selector.Sel.Name,
			Name:    name,
			Package: selector.X.(*ast.Ident).Name,
		}
	}

	return nil
}

type DirectField struct {
	Name    string
	Type    string
	Package string
}

func (f DirectField) GetName() string {
	return f.Name
}

func (f DirectField) GetType() string {
	return f.Type
}
