package dto

type SortType string

const DescSort SortType = "desc"
const AscSort SortType = "asc"

type Sort struct {
	Type  SortType
	Field string
}

type GetNotifications struct {
	UserId uint
	Seen   *bool
	Limit  int
	Offset int
	Sort   *Sort
}
