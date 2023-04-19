package main

import (
	"GenMapping/generator"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
)

func main() {
	args := os.Args[1:]

	fmt.Println(args)

	if len(args) != 1 {
		fmt.Println("Missing args")
		return
	}

	GenerateMappers(args[0])

	return
}

func GenerateMappers(mainPath string) {
	project, err := generator.ParseProject(mainPath)
	if err != nil {
		return
	}

	type Output struct {
		name string
		text string
	}

	result := make(map[string]*Output)

	for _, mapper := range project.MapperInterfaces {
		fmt.Printf("Mapper: %s\n", mapper.Name)
		for _, method := range mapper.Methods {
			fmt.Printf("- Mapping Method: %s\n", method.Name)

			fmt.Printf("-- From: %v\n", method.Params)
			fmt.Printf("-- To: %v\n", method.Target)
		}

		output := result[mapper.OutputPath(project)]
		if output == nil {
			output = new(Output)
			result[mapper.OutputPath(project)] = output
		}

		output.text += generator.GenerateMapper(mapper, project).Methods
		if output.name != "" {
			output.name += "_"
		}
		output.name += mapper.Name
	}

	for mapperPath, mapperOutput := range result {
		finalOutput := generator.ProcessImports(mapperOutput.text)
		finalOutput = fmt.Sprintf("package %s\n", mapperOutput.name) + finalOutput

		err := os.MkdirAll(filepath.Dir(mapperPath), 0)
		if err != nil {
			panic(err)
		}

		f, err := os.Create(mapperPath)
		if err != nil {
			panic(err)
		}

		formattedOutput, err := format.Source([]byte(finalOutput))
		if err != nil {
			panic(err)
		}

		_, err = f.Write(formattedOutput)
		if err != nil {
			panic(err)
		}
	}

	return
}
