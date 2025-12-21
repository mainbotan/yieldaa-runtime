package preset

import (
	"os"
	"strings"
)

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
