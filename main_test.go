package main

import (
	"os"
	"testing"
)

func TestGenerateMapper(t *testing.T) {
	projectPath, _ := os.Getwd()

	//testInterface := "SensorMapper"
	//testInterfaceFile := "/test/testInterface.go"

	GenerateMappers(projectPath)
}
