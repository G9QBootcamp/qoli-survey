package logging

func mapToZapParams(extraKeys map[ExtraKey]interface{}) []interface{} {
	params := []interface{}{}
	if extraKeys == nil {
		extraKeys = map[ExtraKey]interface{}{}
	}
	for k, v := range extraKeys {
		params = append(params, string(k), v)

	}
	return params
}
