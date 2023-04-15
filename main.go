package main

import (
	"AutoMapper/generator"
	"fmt"
	"os"
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

	for _, mapper := range project.MapperInterfaces {
		fmt.Printf("Mapper: %s\n", mapper.Name)
		for _, method := range mapper.Methods {
			fmt.Printf("- Mapping Method: %s\n", method.Name)

			fmt.Printf("-- From: %v\n", method.Params)
			fmt.Printf("-- To: %v\n", method.Target)
		}

		fmt.Println(generator.GenerateMapper(mapper, project))
	}

	return
}
