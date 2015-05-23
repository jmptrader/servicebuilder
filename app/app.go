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
	Type int
}

type FieldType int

const (
	String FieldType = iota
	Text
	Integer
	Float
	Date
	DateTime
)

type Pagination struct {
	PerPage    int
	MaxPerPage int
}
