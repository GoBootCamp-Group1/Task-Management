package response

type PaginateResponseFromService[T any] struct {
	Data     T
	Page     int
	PageSize int
	Total    int64
}
