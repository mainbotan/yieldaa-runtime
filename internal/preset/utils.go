package preset

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/cespare/xxhash/v2"
)

func calculateContentHash(content []byte) string {
	hash := xxhash.Sum64(content)
	return fmt.Sprintf("%016x", hash)
}

func calculateStructureHash(files []EntityFile) uint32 {
	hash := xxhash.New()
	for _, f := range files {
		relPath := strings.TrimPrefix(f.Path, filepath.Dir(f.Path)+string(os.PathSeparator))
		hash.Write([]byte(relPath))
		binary.Write(hash, binary.LittleEndian, f.Size)
		binary.Write(hash, binary.LittleEndian, f.ModTime.Unix())
	}

	result := binary.LittleEndian.Uint32(hash.Sum(nil)[:4])
	return result
}

func GetField(data map[string]any, key string) any {
	if val, ok := data[key]; ok {
		return val
	}
	return nil
}

func GetFieldString(data map[string]any, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func EntityKey(data map[string]any) string {
	return strings.Join([]string{
		GetFieldString(data, "module"),
		GetFieldString(data, "object"),
		GetFieldString(data, "property"),
		GetFieldString(data, "code"),
	}, ".")
}

func ShortPath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}
	// Показываем последнюю часть пути
	parts := strings.Split(path, string(os.PathSeparator))
	result := ""
	for i := len(parts) - 1; i >= 0; i-- {
		if len(result)+len(parts[i])+1 > maxLen {
			break
		}
		if result == "" {
			result = parts[i]
		} else {
			result = parts[i] + "/" + result
		}
	}
	if result != path && len(result) < maxLen-3 {
		result = "…/" + result
	}
	return result
}

func normalizePattern(pattern string) string {
	if pattern == "YYYY-MM-DD" {
		return `^\d{4}-\d{2}-\d{2}$`
	}

	re := regexp.MustCompile(`\\\\([0-9])`)
	return re.ReplaceAllString(pattern, `\$1`)
}

func getFloat(field map[string]any, key string) (float64, bool) {
	if val, ok := field[key]; ok {
		switch v := val.(type) {
		case float64:
			return v, true
		case int:
			return float64(v), true
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f, true
			}
		}
	}
	return 0, false
}

func getNumberValue(field map[string]any, key string) *float64 {
	if val, ok := field[key]; ok {
		switch v := val.(type) {
		case float64:
			return &v
		case int:
			f := float64(v)
			return &f
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return &f
			}
		}
	}
	return nil
}

func validatePattern(pattern string) (bool, string) {
	if pattern == "" {
		return true, ""
	}

	if pattern == "YYYY-MM-DD" {
		return true, ""
	}

	if _, err := regexp.Compile(pattern); err != nil {
		return false, fmt.Sprintf("invalid regex pattern: %v", err)
	}

	return true, ""
}

func validateMinMax(min, max *float64, fieldType string) []string {
	var errors []string

	if min != nil && max != nil && *min > *max {
		errors = append(errors, "min cannot be greater than max")
	}

	if fieldType == "string" && min != nil && *min < 0 {
		errors = append(errors, "min cannot be negative for string")
	}

	return errors
}

func normalizePatternForSchema(pattern string) string {
	if pattern == "YYYY-MM-DD" {
		return "^\\d{4}-\\d{2}-\\d{2}$"
	}
	return pattern
}

func isValidType(t string) bool {
	switch t {
	case "string", "number", "integer", "boolean", "enum":
		return true
	default:
		return false
	}
}

func getFieldCodeAndType(field map[string]any) (code, fieldType string) {
	code, _ = field["code"].(string)
	fieldType, _ = field["type"].(string)
	return
}

func validateEnumValues(values []any) (bool, string) {
	if values == nil || len(values) == 0 {
		return false, "enum requires values array"
	}

	for i, val := range values {
		if _, ok := val.(string); !ok {
			return false, fmt.Sprintf("enum value at index %d is not a string", i)
		}
	}

	return true, ""
}
