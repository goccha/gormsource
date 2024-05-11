package foundations

import (
	"context"
)

type UpdateField func(ctx context.Context, m map[string]interface{})

func UpdateValue(name string, value interface{}) UpdateField {
	return func(ctx context.Context, m map[string]interface{}) {
		m[name] = value
	}
}

func Updates(ctx context.Context, fields ...UpdateField) map[string]interface{} {
	m := make(map[string]interface{})
	for _, field := range fields {
		field(ctx, m)
	}
	return m
}
