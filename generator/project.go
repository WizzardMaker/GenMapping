package generator

import (
	"go/ast"
)

type Project struct {
	Packages         map[string]*ast.Package
	MapperInterfaces []Mapper
	Structs          []Structure
	Imports          []Import
}
