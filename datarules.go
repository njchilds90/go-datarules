package datarules

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ValidationError struct {
	Errors []FieldError `json:"errors"`
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

func (e *ValidationError) Add(field, code, message string) {
	e.Errors = append(e.Errors, FieldError{
		Field:   field,
		Code:    code,
		Message: message,
	})
}

func (e *ValidationError) HasErrors() bool {
	return len(e.Errors) > 0
}

type rule func(ctx context.Context, data map[string]any, verr *ValidationError)

type Schema struct {
	rules []rule
}

func New() *Schema {
	return &Schema{}
}

func (s *Schema) addRule(r rule) *Schema {
	s.rules = append(s.rules, r)
	return s
}

func (s *Schema) Required(field string) *Schema {
	return s.addRule(func(ctx context.Context, data map[string]any, verr *ValidationError) {
		if _, ok := data[field]; !ok {
			verr.Add(field, "required", "field is required")
		}
	})
}

func (s *Schema) String(field string) *Schema {
	return s.addRule(func(ctx context.Context, data map[string]any, verr *ValidationError) {
		if v, ok := data[field]; ok {
			if _, ok := v.(string); !ok {
				verr.Add(field, "type_string", "must be a string")
			}
		}
	})
}

func (s *Schema) Integer(field string) *Schema {
	return s.addRule(func(ctx context.Context, data map[string]any, verr *ValidationError) {
		if v, ok := data[field]; ok {
			kind := reflect.TypeOf(v).Kind()
			if kind != reflect.Int && kind != reflect.Int64 && kind != reflect.Int32 {
				verr.Add(field, "type_integer", "must be an integer")
			}
		}
	})
}

func (s *Schema) MinLength(field string, min int) *Schema {
	return s.addRule(func(ctx context.Context, data map[string]any, verr *ValidationError) {
		if v, ok := data[field]; ok {
			str, ok := v.(string)
			if !ok {
				return
			}
			if len(str) < min {
				verr.Add(field, "min_length", fmt.Sprintf("minimum length is %d", min))
			}
		}
	})
}

func (s *Schema) Default(field string, value any) *Schema {
	return s.addRule(func(ctx context.Context, data map[string]any, verr *ValidationError) {
		if _, ok := data[field]; !ok {
			data[field] = value
		}
	})
}

func (s *Schema) Transform(field string, fn func(any) any) *Schema {
	return s.addRule(func(ctx context.Context, data map[string]any, verr *ValidationError) {
		if v, ok := data[field]; ok {
			data[field] = fn(v)
		}
	})
}

func (s *Schema) Validate(ctx context.Context, input map[string]any) (map[string]any, error) {
	if ctx == nil {
		return nil, errors.New("context cannot be nil")
	}

	output := make(map[string]any)
	for k, v := range input {
		output[k] = v
	}

	verr := &ValidationError{}

	for _, r := range s.rules {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			r(ctx, output, verr)
		}
	}

	if verr.HasErrors() {
		return nil, verr
	}

	return output, nil
}
