package generator

type Type struct {
	ArgumentName string
	Name         string
	Package      string
}

type Field struct {
	Name string
}

type Import struct {
	Name string
	Path string
}
