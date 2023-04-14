package main

import (
	"AutoMapper/generator"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	fmt.Println(args)

	if len(args) != 2 {
		fmt.Println("Missing args")
		return
	}

	GenerateMapper(args[0], args[1])

	return
}

func GenerateMapper(mapperInterfaceName, mainPath string) string {
	project, err := generator.ParseProject(mainPath)
	if err != nil {
		return ""
	}

	fmt.Println(project)

	return ""
}
