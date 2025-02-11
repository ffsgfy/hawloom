package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Value[T any] struct {
	V T
}

func (v *Value[T]) UnmarshalJSON(data []byte) error {
	// Check for environment variable ("${VAR}")
	if len(data) > 5 {
		if data[0] == '"' && data[1] == '$' && data[2] == '{' && data[len(data)-2] == '}' && data[len(data)-1] == '"' {
			varName := string(data[3 : len(data)-2])
			varValue := os.Getenv(varName)
			switch any(v.V).(type) {
			case string:
				var err error
				data, err = json.Marshal(varValue)
				if err != nil {
					return fmt.Errorf("failed to encode value of $%s: %w", varName, err)
				}
			default:
				data = []byte(varValue)
			}
		}
	}

	return json.Unmarshal(data, &v.V)
}
