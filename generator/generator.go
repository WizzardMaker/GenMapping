package generator

import (
	_ "embed"
	"html/template"
	"strings"
)

//go:embed templates/mapper.tmpl
var MapperTemplate string

//go:embed templates/mapping.tmpl
var MappingTemplate string

//go:embed templates/mapper.tmpl
var DirectMapping string

func GenerateMapper(mapper Mapper, project *Project) string {
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

	return output.String()
}
