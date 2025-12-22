package preset

import (
	"bytes"         // ← ДОБАВИТЬ
	"encoding/json" // ← ДОБАВИТЬ
	"fmt"
	"strconv"
)

func GenerateJSONSchema(parsed map[string]any) (map[string]any, error) {
	// fields from parsed
	fieldsAny, ok := parsed["fields"].([]any)
	if !ok {
		return nil, fmt.Errorf("fields not found or not an array")
	}

	// convert into struct
	fields := make([]map[string]any, 0, len(fieldsAny))
	for _, f := range fieldsAny {
		if field, ok := f.(map[string]any); ok {
			fields = append(fields, field)
		}
	}

	// generate json_schema
	schema := map[string]any{
		"$schema": JSONSchemaDraft,
		"$id": fmt.Sprintf("/%s/%s/%s/schema.json",
			parsed["module"], parsed["object"], parsed["code"]),
		"type":                 "object",
		"title":                parsed["name"],
		"additionalProperties": false,
		"properties":           make(map[string]any),
		"required":             []string{},
	}

	properties := schema["properties"].(map[string]any)
	required := make([]string, 0)

	for _, field := range fields {
		fieldCode, _ := field["code"].(string)
		fieldName, _ := field["name"].(string)
		isRequired, _ := field["required"].(bool)

		// schema for field
		fieldSchema := generateFieldJSONSchema(field)

		// + title
		fieldSchema["title"] = fieldName

		if description, ok := field["description"].(string); ok && description != "" {
			fieldSchema["description"] = description
		}

		if examples, ok := field["examples"].([]any); ok && len(examples) > 0 {
			fieldSchema["examples"] = examples
		}

		properties[fieldCode] = fieldSchema

		if isRequired {
			required = append(required, fieldCode)
		}
	}

	schema["required"] = required

	if rootExamples, ok := parsed["examples"].([]any); ok && len(rootExamples) > 0 {
		schema["examples"] = rootExamples
	}

	return schema, nil
}

func generateFieldJSONSchema(field map[string]any) map[string]any {
	fieldType, _ := field["type"].(string)
	schema := make(map[string]any)

	switch fieldType {
	case "string":
		schema["type"] = "string"

		if pattern, ok := field["pattern"].(string); ok && pattern != "" {
			normalized := normalizePatternForSchema(pattern)
			if pattern == "YYYY-MM-DD" {
				schema["format"] = "date"
			}
			schema["pattern"] = normalized
		}

		if min := getNumberValue(field, "min"); min != nil {
			schema["minLength"] = int(*min)
		}
		if max := getNumberValue(field, "max"); max != nil {
			schema["maxLength"] = int(*max)
		}

		// default value
		if def, ok := field["default"].(string); ok && def != "" {
			schema["default"] = def
		}

	case "number":
		schema["type"] = "number"

		// range
		if min := getNumberValue(field, "min"); min != nil {
			schema["minimum"] = *min
		}
		if max := getNumberValue(field, "max"); max != nil {
			schema["maximum"] = *max
		}

		// multiplicity
		if multiple := getNumberValue(field, "multiple_of"); multiple != nil {
			schema["multipleOf"] = *multiple
		}
		if multiple := getNumberValue(field, "multipleOf"); multiple != nil {
			schema["multipleOf"] = *multiple
		}

		// default
		if def, ok := field["default"].(float64); ok {
			schema["default"] = def
		}
		if def, ok := field["default"].(int); ok {
			schema["default"] = float64(def)
		}

	case "integer":
		schema["type"] = "integer"

		// range
		if min := getNumberValue(field, "min"); min != nil {
			schema["minimum"] = int(*min)
		}
		if max := getNumberValue(field, "max"); max != nil {
			schema["maximum"] = int(*max)
		}

		// multiplicity
		if multiple := getNumberValue(field, "multiple_of"); multiple != nil {
			schema["multipleOf"] = *multiple
		}
		if multiple := getNumberValue(field, "multipleOf"); multiple != nil {
			schema["multipleOf"] = *multiple
		}

		// default
		if def, ok := field["default"].(float64); ok {
			schema["default"] = int(def)
		}
		if def, ok := field["default"].(int); ok {
			schema["default"] = def
		}
		if def, ok := field["default"].(string); ok && def != "" {
			if intVal, err := strconv.Atoi(def); err == nil {
				schema["default"] = intVal
			}
		}

	case "boolean":
		schema["type"] = "boolean"

		// default
		if def, ok := field["default"].(bool); ok {
			schema["default"] = def
		}
		if def, ok := field["default"].(string); ok && def != "" {
			if def == "true" {
				schema["default"] = true
			} else if def == "false" {
				schema["default"] = false
			}
		}

	case "enum":
		schema["type"] = "string"

		// allowable
		if values, ok := field["values"].([]any); ok && len(values) > 0 {
			schema["enum"] = values
		}

		// default
		if def, ok := field["default"].(string); ok && def != "" {
			schema["default"] = def
		}

	default:
		// undefined -> string
		schema["type"] = "string"
	}

	return schema
}

func formatJSONData(jsonData []byte) string {
	if jsonData == nil || len(jsonData) == 0 {
		return ""
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, jsonData, "", "  "); err != nil {
		return string(jsonData)
	}
	return prettyJSON.String()
}
