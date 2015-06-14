package app

type Application struct {
	Models []*Model
}

type Model struct {
	Name       string
	Fields     []*Field
	Pagination *Pagination
	Actions    *RestfulActions
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

type RestfulActions struct {
	Index   bool
	Create  bool
	Show    bool
	Update  bool
	Destroy bool
}

func DefaultRestfulActions() *RestfulActions {
	return &RestfulActions{
		Index:   true,
		Create:  true,
		Show:    true,
		Update:  true,
		Destroy: true,
	}
}
