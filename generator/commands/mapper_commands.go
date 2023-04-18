package commands

const MapperTag Tag = "@mapper"

type MapperCommand struct {
	TargetFile Attribute
}

func (c *MapperCommand) Read(text string) string {
	commandText, readText := getCommandText(text, MapperTag)

	c.TargetFile.FromCommandText("targetFile", commandText)

	return readText
}

var PerMapperTags = []Tag{MapperTag}
