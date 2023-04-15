package commands

const MapperTag Tag = "@mapper"

type MapperCommand struct {
	TargetFile Attribute
}

func (c *MapperCommand) Read(text string) {
	_ = getCommandText(text, MapperTag)
}

var PerMapperTags = []Tag{MapperTag}
