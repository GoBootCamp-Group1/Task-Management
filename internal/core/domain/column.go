package domain

type Column struct {
	ID            uint
	BoardID       uint
	Name          string
	OrderPosition uint
	IsFinal       bool
	CreatedBy     uint
}
