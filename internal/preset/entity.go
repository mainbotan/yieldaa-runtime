package preset

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"os"

	"github.com/ghodss/yaml"
)

func ProcessEntity(file EntityFile) ProcessedEntity {
	result := ProcessedEntity{File: file}

	data, err := os.ReadFile(file.Path)
	if err != nil {
		result.FatalError = fmt.Errorf("read: %w", err)
		return result
	}

	result.ContentHash = crc32.ChecksumIEEE(data)

	// .yml -> json without inter go struct
	jsonData, err := yaml.YAMLToJSON(data)
	if err != nil {
		result.FatalError = fmt.Errorf("YAML→JSON: %w", err)
		return result
	}
	result.JSONData = jsonData

	// validation round 1
	var parsed map[string]any
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		result.FatalError = fmt.Errorf("invalid JSON: %w", err)
		return result
	}
	result.ParsedData = parsed

	// validation round 2
	if errs := validateStructure(parsed); len(errs) > 0 {
		result.Errors = append(result.Errors, errs...)
	}

	// validation round 3
	if errs := validateFieldsDirectly(parsed); len(errs) > 0 {
		result.Errors = append(result.Errors, errs...)
	}

	// generate schema
	if result.FatalError == nil && len(result.Errors) == 0 {
		schema, err := GenerateJSONSchema(result.ParsedData)
		if err != nil {
			result.Errors = append(result.Errors,
				fmt.Sprintf("JSON Schema generation failed: %v", err))
		} else {
			result.Schema = schema
		}
	}

	return result
}

// structure validation
func validateStructure(data map[string]any) []string {
	var errors []string

	required := []string{"module", "object", "property", "code", "name", "fields"}
	for _, field := range required {
		if _, exists := data[field]; !exists {
			errors = append(errors, fmt.Sprintf("missing: %s", field))
		}
	}

	if fields, ok := data["fields"].([]any); !ok {
		errors = append(errors, "fields must be array")
	} else if len(fields) == 0 {
		errors = append(errors, "fields array empty")
	}

	return errors
}

// fields validation
func validateFieldsDirectly(data map[string]any) []string {
	var errors []string

	fields, ok := data["fields"].([]any)
	if !ok {
		return []string{"fields is not array"}
	}

	seenCodes := make(map[string]bool)

	for i, fieldAny := range fields {
		field, ok := fieldAny.(map[string]any)
		if !ok {
			errors = append(errors, fmt.Sprintf("field[%d]: not object", i))
			continue
		}

		code, typeStr := getFieldCodeAndType(field)
		if code == "" {
			errors = append(errors, fmt.Sprintf("field[%d]: missing code", i))
			continue
		}

		if seenCodes[code] {
			errors = append(errors, fmt.Sprintf("duplicate field code: %s", code))
		}
		seenCodes[code] = true

		if !isValidType(typeStr) {
			errors = append(errors, fmt.Sprintf("field %s: invalid type '%s'", code, typeStr))
			continue
		}

		if pattern, ok := field["pattern"].(string); ok && pattern != "" {
			if valid, errMsg := validatePattern(pattern); !valid {
				errors = append(errors, fmt.Sprintf("field %s: %s", code, errMsg))
			}
		}

		if min := getNumberValue(field, "min"); min != nil {
			max := getNumberValue(field, "max")
			if minMaxErrs := validateMinMax(min, max, typeStr); len(minMaxErrs) > 0 {
				for _, err := range minMaxErrs {
					errors = append(errors, fmt.Sprintf("field %s: %s", code, err))
				}
			}
		}

		if typeStr == "enum" {
			if values, ok := field["values"].([]any); ok {
				if valid, errMsg := validateEnumValues(values); !valid {
					errors = append(errors, fmt.Sprintf("field %s: %s", code, errMsg))
				}
			} else {
				errors = append(errors, fmt.Sprintf("field %s: enum requires values array", code))
			}
		}
	}

	return errors
}
