package commands

const MapperTag Tag = "@mapper"

type MapperCommand struct {
	TargetFile Attribute
}

func (c *MapperCommand) Read(text string) string {
	_, readText := getCommandText(text, MapperTag)

	return readText
}

var PerMapperTags = []Tag{MapperTag}
