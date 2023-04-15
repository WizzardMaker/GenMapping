package generator

type Type struct {
	ArgumentName string
	Name         string
	Package      string
}

func (t Type) GetName() string {
	return t.ArgumentName + " " + t.Package + "." + t.Name
}

type Field struct {
	Name string
}

type Import struct {
	Name string
	Path string
}
