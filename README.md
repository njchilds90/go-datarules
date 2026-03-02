# go-datarules

Declarative, deterministic data validation and transformation engine for Go.

## Features

- Zero dependencies
- Deterministic rule evaluation
- Structured machine-readable errors
- Field-level validation
- Default value support
- Transformation hooks
- Context support
- Pure functional behavior

## Installation

```bash
go get github.com/njchilds90/go-datarules

package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/njchilds90/go-datarules"
)

func main() {
	schema := datarules.New().
		Required("name").
		String("name").
		MinLength("name", 3).
		Integer("age").
		Default("active", true).
		Transform("name", func(v any) any {
			return strings.ToUpper(v.(string))
		})

	input := map[string]any{
		"name": "john",
		"age":  30,
	}

	result, err := schema.Validate(context.Background(), input)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}

{
  "errors": [
    {
      "field": "name",
      "code": "min_length",
      "message": "minimum length is 3"
    }
  ]
}
