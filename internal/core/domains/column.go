package domains

type Column struct {
	ID            uint   `json:"id,omitempty"`
	BoardID       uint   `json:"board_id,omitempty"`
	Name          string `json:"name,omitempty"`
	OrderPosition int    `json:"order_position,omitempty"`
	IsFinal       bool   `json:"is_final,omitempty"`
	CreatedBy     uint   `json:"created_by,omitempty"`
}

type ColumnUpdate struct {
	ID   uint   `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ColumnMove struct {
	ID            uint `json:"id,omitempty"`
	OrderPosition int  `json:"order_position,omitempty"`
}
