package datarules

import (
	"context"
	"strings"
	"testing"
)

func TestSchemaValidation(t *testing.T) {
	ctx := context.Background()

	schema := New().
		Required("name").
		String("name").
		MinLength("name", 3).
		Integer("age").
		Default("active", true).
		Transform("name", func(v any) any {
			return strings.ToUpper(v.(string))
		})

	tests := []struct {
		name      string
		input     map[string]any
		expectErr bool
	}{
		{
			name: "valid input",
			input: map[string]any{
				"name": "john",
				"age":  30,
			},
			expectErr: false,
		},
		{
			name: "missing name",
			input: map[string]any{
				"age": 30,
			},
			expectErr: true,
		},
		{
			name: "short name",
			input: map[string]any{
				"name": "ab",
				"age":  30,
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		_, err := schema.Validate(ctx, tt.input)
		if tt.expectErr && err == nil {
			t.Fatalf("expected error for %s", tt.name)
		}
		if !tt.expectErr && err != nil {
			t.Fatalf("unexpected error for %s: %v", tt.name, err)
		}
	}
}
