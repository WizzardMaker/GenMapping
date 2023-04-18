package generator

import (
	"AutoMapper/generator/commands"
	"AutoMapper/generator/mappings"
	"go/ast"
	"go/types"
)

type Project struct {
	Packages         map[string]*ast.Package
	MapperInterfaces Mappers
	Structs          []Structure
	Imports          []mappings.Import
	GlobalCommands   []commands.Command
	GlobalTypes      []*types.Info
}
