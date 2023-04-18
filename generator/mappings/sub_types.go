package mappings

type UnderlyingType int

const (
	DefaultType UnderlyingType = iota
	ArrayType
	PointerType
)

type Type struct {
	ArgumentName string
	Name         string
	Package      string
	Underlying   UnderlyingType
}

func (t Type) GetFullName() string {
	return t.ArgumentName + " " + t.Package + "." + t.Name
}

func (t Type) GetTypeName() string {
	return t.Package + "." + t.Name
}

type Import struct {
	Name string
	Path string
}
