package generator

import (
	"GenMapping/generator/commands"
	"GenMapping/generator/mappings"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ParseProject(projectRoot string) (*Project, error) {
	goModPath := filepath.Join(projectRoot, "go.mod")
	_, err := os.Stat(goModPath)
	if err != nil {
		return nil, fmt.Errorf("go.mod file not found")
	}

	modContent, err := os.ReadFile(goModPath)
	if err != nil {
		return nil, err
	}

	r, err := regexp.Compile("module (.*)")
	if err != nil {
		return nil, err
	}
	moduleName := r.FindStringSubmatch(string(modContent))[1]

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
	var imports []mappings.Import

	err = filepath.WalkDir(projectRoot, func(filePath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Check if the current entry is a directory
		if d.IsDir() {
			parsed, err := parser.ParseDir(&files, filePath, nil, parser.ParseComments)
			if err != nil {
				fmt.Println("Error during parsing of package \""+d.Name()+"\"\n\t- Error:\n\t", err)
				return nil
			}
			parsedProject = merge(parsed, parsedProject)

		}
		return nil
	})

	loadConfig := new(packages.Config)
	loadConfig.Dir = projectRoot
	loadConfig.Mode = loadMode
	loadConfig.Fset = &files
	packs, err := packages.Load(loadConfig, "./...")

	var globalInfos []*types.Info
	var globalCommands []commands.Command

	for _, a := range parsedProject {
		var info *types.Info
		for _, pack := range packs {
			if len(pack.GoFiles) == 0 {
				continue
			}

			for fPath := range a.Files {
				if pack.GoFiles[0] != fPath {
					continue
				}

				info = pack.TypesInfo
				globalInfos = append(globalInfos, info)
				break
			}
		}

		if info == nil {
			fmt.Println("Couldn't find info for package: ", a.Name)
		}

		mappersInPackage, structures, importsInPackage := FindMappersInPackage(a, info)

		//for _, mapper := range mappersInPackage {
		//	for _, mappingMethod := range mapper.Methods {
		//		if len(mappingMethod.Params) == 1 {
		//			//TODO: check for error handling!
		//			globalCommands = append(globalCommands, &commands.ExpressionCommand{
		//				Target:     commands.Attribute{Value: mappingMethod.Params[0].GetTypeName()},
		//				Expression: commands.Attribute{Value: fmt.Sprintf("%s.%s(%s)", mapper.Name, mappingMethod.Name, "%s")},
		//				IsType:     commands.Attribute{Value: "true"},
		//			})
		//		}
		//	}
		//}

		mappers = append(mappers, mappersInPackage...)
		structs = append(structs, structures...)
		imports = append(imports, importsInPackage...)
	}

	if err != nil {
		return nil, err
	}

	proj := Project{
		Packages:         parsedProject,
		MapperInterfaces: mappers,
		Imports:          imports,
		Structs:          structs,
		GlobalTypes:      globalInfos,
		GlobalCommands:   globalCommands,
		ModuleName:       moduleName,
		BasePath:         projectRoot,
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

func FindMappersInPackage(pack *ast.Package, info *types.Info) ([]Mapper, []Structure, []mappings.Import) {
	var mappers []Mapper
	var structs []Structure

	var packageImports []mappings.Import

	ast.Inspect(pack, func(node ast.Node) bool {
		inter, ok := node.(*ast.GenDecl)
		if ok {
			for _, spec := range inter.Specs {
				importSpec, ok := spec.(*ast.ImportSpec)
				if ok {
					imp := mappings.Import{
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
						if !strings.Contains(doc, string(commands.MapperTag)) {
							return true
						}

						methods := NewMethods(interfaceType.Methods, pack.Name, info)

						var neededImports []mappings.Import
						for i, method := range methods {
							var methodTypes []*mappings.Type
							methodTypes = append(methodTypes, &method.Target)
							for i := range method.Params {
								methodTypes = append(methodTypes, &method.Params[i])
							}

							for _, t := range methodTypes {
								for _, packageImport := range packageImports {
									if packageImport.Name == t.Package {
										t.Package = packageImport.Path
										neededImports = append(neededImports, packageImport)
									}
								}
							}

							methods[i] = method
						}

						mapper := Mapper{
							Interface:  *interfaceType,
							Name:       typeSpec.Name.Name,
							Methods:    methods,
							Imports:    neededImports,
							Commands:   commands.FromText(doc, commands.PerMapperTags...),
							outputPath: GetBasePathOfPkg(pack),
						}
						mappers = append(mappers, mapper)
					}

					_, ok = typeSpec.Type.(*ast.StructType)
					if ok {
						structure := NewStructure(typeSpec, info, pack.Name)
						structs = append(structs, structure)
					}
				}
			}

			return false
		}

		return true
	})

	return mappers, structs, packageImports
}

func GetBasePathOfPkg(pack *ast.Package) string {
	for path := range pack.Files {
		return filepath.Dir(path)
	}

	return ""
}
