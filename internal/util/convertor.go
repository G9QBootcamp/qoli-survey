package util

import (
	"encoding/json"
	"errors"

	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
)

func ConvertTypes[T any](logger logging.Logger, source interface{}, dest *T) error {
	// Marshal the source object to JSON
	b, err := json.Marshal(source)
	if err != nil {
		logger.Error(logging.Internal, logging.FailedConvertDto, "error in marshaling", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return errors.New("failed to marshal source object")
	}

	// Unmarshal the JSON into the destination object
	err = json.Unmarshal(b, dest)
	if err != nil {
		logger.Error(logging.Internal, logging.FailedConvertDto, "error in unmarshaling", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return errors.New("failed to unmarshal into destination object")
	}

	return nil
}
