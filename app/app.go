package app

type Application struct {
	Models []*Model
}

type Model struct {
	Name       string
	Fields     []*Field
	Pagination *Pagination
}

type Field struct {
	Name string
	Type FieldType
}

type FieldType int

const (
	STRING FieldType = iota
	INT
	DOUBLE
	DATE
	DATETIME
)

type Pagination struct {
	PerPage    int
	MaxPerPage int
}
