package generator

import (
	_ "embed"
	"fmt"
	"html/template"
	"regexp"
	"strings"
)

//go:embed templates/mapper.tmpl
var MapperTemplate string

//go:embed templates/mapping.tmpl
var MappingTemplate string

//go:embed templates/mapper.tmpl
var DirectMapping string

type MapperOutput struct {
	Methods string
}

func GenerateMapper(mapper Mapper, project *Project) (result MapperOutput) {
	templates := template.New("mapper")

	templates.Funcs(template.FuncMap{
		"getProject": func() Project { return *project },
	})

	mapperTmpl, err := templates.Parse(MapperTemplate)
	if err != nil {
		panic(err)
	}

	_, err = templates.Parse(MappingTemplate)
	if err != nil {
		panic(err)
	}

	output := new(strings.Builder)
	err = mapperTmpl.Execute(output, mapper)
	if err != nil {
		panic(err)
	}

	result.Methods = output.String()
	return
}

func ProcessImports(mappingFunctions string) string {
	output := mappingFunctions

	const IMPORT_PATTERN = "%\\*__\\*%"

	// html/template.X = 0->html/template 1->template 2->X
	r, err := regexp.Compile(IMPORT_PATTERN + "((?:[\\w\\d_/.]*?\\/)?([\\w\\d_-]*?))\\.(.*?)" + IMPORT_PATTERN)
	if err != nil {
		panic(err)
	}
	foundImports := r.FindAllStringSubmatch(mappingFunctions, -1)
	alreadyImports := make(map[string][][]int)

	//var imports []string

	// html/template.Alpha, Custom/template.Alpha
	// ->
	// html/template.Alpha, Custom/template2.Alpha

	for i, importItem := range foundImports {
		packetPath := importItem[1]
		packet := importItem[2]
		obj := importItem[3]

		fmt.Println(packetPath, " ", packet, " ", obj)

		found := false
		at := 0
		for i, imp := range alreadyImports[packet] {
			if foundImports[imp[0]][1] == packetPath {
				found = true
				at = i
			}
		}
		if !found {
			alreadyImports[packet] = append(alreadyImports[packet], []int{i})
		} else {
			alreadyImports[packet][at] = append(alreadyImports[packet][at], i)
		}
	}

	importText := ""

	for _, importIndices := range alreadyImports {
		for count, index := range importIndices {
			importItem := foundImports[index[0]]
			if strings.Contains(importItem[0], "--go--") {
				output = strings.Replace(output, importItem[0], importItem[3], -1)
				continue
			}

			importText += "\n\t"

			var packageOutput string
			if count != 0 {
				packageOutput = fmt.Sprintf("%s%d", importItem[2], count+1)
				importText += fmt.Sprintf("%s%d", importItem[2], count+1)
			} else {
				packageOutput = fmt.Sprintf("%s", importItem[2])
			}

			importText += fmt.Sprintf("\"%s\"", importItem[1])
			for removeIndex := range index {
				removeItem := foundImports[index[removeIndex]]
				output = strings.Replace(output, removeItem[0], packageOutput+"."+removeItem[3], -1)
			}
		}
	}

	return fmt.Sprintf("import (%s\n)\n", importText) + output
}
