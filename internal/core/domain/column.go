package domain

type Column struct {
	ID            uint
	BoardID       uint
	Name          string
	OrderPosition int
	IsFinal       bool
	CreatedBy     uint
}

type ColumnUpdate struct {
	ID   uint
	Name string
}

type ColumnMove struct {
	ID            uint
	OrderPosition int
}
