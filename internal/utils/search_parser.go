package utils

import (
	"fmt"
	"strings"
)

type SearchCondition struct {
	Field    string
	Value    string
	Operator string // Defaults to "contains" or you can support more
}

// ParseSearchQuery parses search parameter like "name:finan" or "name:finan,code:HR"
func ParseSearchQuery(searchParam string) ([]SearchCondition, error) {
	if searchParam == "" {
		return nil, nil
	}

	var conditions []SearchCondition

	// Split by comma for multiple conditions (if no searchJoin param)
	parts := strings.Split(searchParam, ";")

	for _, part := range parts {
		// Split by colon to separate field and value
		colonIndex := strings.Index(part, ":")
		if colonIndex == -1 {
			return nil, fmt.Errorf("invalid search format: %s", part)
		}

		field := part[:colonIndex]
		value := part[colonIndex+1:]

		if field == "" || value == "" {
			return nil, fmt.Errorf("empty field or value in search: %s", part)
		}

		conditions = append(conditions, SearchCondition{
			Field:    field,
			Value:    value,
			Operator: "LIKE", // Default operator
		})
	}

	return conditions, nil
}
