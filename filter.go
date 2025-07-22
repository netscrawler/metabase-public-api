package metabase

import (
	"errors"
	"reflect"
	"strings"
)

// Filter represents a query parameter for Metabase public API
type Filter struct {
	Type   string `json:"type"`
	Value  any    `json:"value"`
	Target any    `json:"target"`
}

// NewCategoryFilter builds category-type filter with template-tag target
func NewCategoryFilter(tagName string, value any) Filter {
	return Filter{
		Type:  "category",
		Value: value,
		Target: []any{
			"variable",
			[]any{"template-tag", tagName},
		},
	}
}

// GenerateFiltersFromStruct creates filters from struct fields using struct tags or field names
func GenerateFiltersFromStruct(input any) ([]Filter, error) {
	v := reflect.ValueOf(input)
	if v.Kind() != reflect.Struct {
		return nil, errors.New("input must be struct")
	}
	t := v.Type()
	var filters []Filter
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		name := field.Tag.Get("metabase")
		if name == "" {
			name = strings.ToLower(field.Name)
		}
		filters = append(filters, NewCategoryFilter(name, value))
	}
	return filters, nil
}
