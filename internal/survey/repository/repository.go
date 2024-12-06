package repository

import (
	"errors"

	"github.com/G9QBootcamp/qoli-survey/internal/survey/dto"
	"gorm.io/gorm"
)

func GetRecords[T any](db *gorm.DB, request *dto.RepositoryRequest) ([]T, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}

	query := db.Model(new(T))

	if request.With != "" {
		query = query.Preload(request.With)
	}

	// Apply filters
	for _, filter := range request.Filters {

		switch filter.Operator {
		case "=":
			query = query.Where(filter.Field+" = ?", filter.Value)
		case "!=":
			query = query.Where(filter.Field+" != ?", filter.Value)
		case ">":
			query = query.Where(filter.Field+" > ?", filter.Value)
		case "<":
			query = query.Where(filter.Field+" < ?", filter.Value)
		case ">=":
			query = query.Where(filter.Field+" >= ?", filter.Value)
		case "<=":
			query = query.Where(filter.Field+" <= ?", filter.Value)
		case "LIKE":
			query = query.Where(filter.Field+" LIKE ?", "%"+filter.Value+"%")
		default:
			return nil, errors.New("unsupported filter operator: " + filter.Operator)
		}
	}

	// Apply sorts
	for _, sort := range request.Sorts {
		if sort.SortType != "asc" && sort.SortType != "desc" {
			return nil, errors.New("invalid sort type: " + sort.SortType)
		}
		query = query.Order(sort.Field + " " + sort.SortType)
	}

	// Apply limit and offset
	if request.Limit > 0 {
		query = query.Limit(int(request.Limit))
	}
	query = query.Offset(int(request.Offset))

	// Execute the query
	var records []T
	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}
