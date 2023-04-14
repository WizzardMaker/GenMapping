package generator

import (
	"AutoMapper/generator/commands"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"os"
	"path/filepath"
	"strings"
)

func ParseProject(projectRoot string) (*Project, error) {
	goModPath := filepath.Join(projectRoot, "go.mod")
	_, err := os.Stat(goModPath)
	if err != nil {
		return nil, fmt.Errorf("go.mod file not found")
	}

	files := token.FileSet{}
	var parsedProject map[string]*ast.Package

	const loadMode = packages.NeedName |
		packages.NeedFiles |
		packages.NeedCompiledGoFiles |
		packages.NeedSyntax |
		packages.NeedImports |
		packages.NeedDeps |
		packages.NeedTypes |
		packages.NeedTypesInfo

	var mappers []Mapper
	var structs []Structure
	var imports []Import

	err = filepath.WalkDir(projectRoot, func(filePath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Check if the current entry is a directory
		if d.IsDir() {
			parsed, err := parser.ParseDir(&files, filePath, nil, parser.ParseComments)
			if err != nil {
				return err
			}
			parsedProject = merge(parsed, parsedProject)

		}
		return nil
	})

	//info := &types.Info{
	//	Defs:  make(map[*ast.Ident]types.Object),
	//	Types: make(map[ast.Expr]types.TypeAndValue),
	//}

	loadConfig := new(packages.Config)
	loadConfig.Dir = projectRoot
	loadConfig.Mode = loadMode
	loadConfig.Fset = &files
	packs, err := packages.Load(loadConfig, "./...")
	fmt.Println(packs, err)

	for _, a := range parsedProject {
		var info *types.Info
		for _, pack := range packs {
			if pack.Name == a.Name {
				info = pack.TypesInfo
				break
			}

		}
		
		mappersInPackage, structures, importsInPackage := FindMappersInPackage(a, info)
		mappers = append(mappers, mappersInPackage...)
		structs = append(structs, structures...)
		imports = append(imports, importsInPackage...)
	}

	if err != nil {
		return nil, err
	}

	proj := Project{
		Packages: parsedProject,
		//MapperInterfaces: mapper,
	}

	return &proj, nil
}

func merge[T any](ms ...map[string]T) map[string]T {
	res := map[string]T{}
	for _, m := range ms {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
}

func FindMappersInPackage(pack *ast.Package, info *types.Info) ([]Mapper, []Structure, []Import) {
	var mappers []Mapper
	var structs []Structure

	var packageImports []Import

	ast.Inspect(pack, func(node ast.Node) bool {
		inter, ok := node.(*ast.GenDecl)
		if ok {
			for _, spec := range inter.Specs {
				importSpec, ok := spec.(*ast.ImportSpec)
				if ok {
					imp := Import{
						Path: strings.Trim(importSpec.Path.Value, "\""),
					}

					if importSpec.Name != nil {
						imp.Name = importSpec.Name.Name
					} else {
						paths := strings.Split(imp.Path, "/")
						imp.Name = paths[len(paths)-1]
					}

					packageImports = append(packageImports, imp)
				}

				typeSpec, ok := spec.(*ast.TypeSpec)
				if ok {
					var doc string = ""
					if inter.Doc != nil {
						doc = inter.Doc.Text()
					}
					if typeSpec.Doc != nil {
						doc += typeSpec.Doc.Text()
					}

					interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
					if ok {
						//only interfaces should be mapper
						if !strings.Contains(doc, commands.MapperTag) {
							return true
						}

						methods := NewMethods(interfaceType.Methods, pack.Name)

						var neededImports []Import
						for _, method := range methods {
							types := append(method.Params, method.Return...)

							for _, t := range types {
								for _, packageImport := range packageImports {
									if packageImport.Name == t.Package {
										neededImports = append(neededImports, packageImport)
									}
								}
							}
						}

						mapper := Mapper{
							Interface: *interfaceType,
							Name:      typeSpec.Name.Name,
							Methods:   methods,
							Imports:   neededImports,
						}
						mappers = append(mappers, mapper)
					}

					structType, ok := typeSpec.Type.(*ast.StructType)
					if ok {
						structs = append(structs, NewStructure(typeSpec, structType, info, pack.Name))
					}
				}
			}

			return false
		}

		return true
	})

	return mappers, structs, packageImports
}
