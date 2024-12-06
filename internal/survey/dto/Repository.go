package dto

type RepositoryFilter struct {
	Field    string
	Operator string
	Value    string
}
type RepositorySort struct {
	Field    string
	SortType string
}

type RepositoryRequest struct {
	Filters []*RepositoryFilter
	Sorts   []*RepositorySort
	Limit   uint
	Offset  uint
	With    string
}
