package generator

import (
	"AutoMapper/generator/commands"
	"AutoMapper/generator/mappings"
	"go/ast"
)

type Project struct {
	Packages         map[string]*ast.Package
	MapperInterfaces []Mapper
	Structs          []Structure
	Imports          []mappings.Import
	GlobalCommands   []commands.Command
}
