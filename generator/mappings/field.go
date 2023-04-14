package mappings

type Field interface {
	GetName() string
	GetType() string
}
