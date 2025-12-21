package preset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func SaveEntitiesToJSON(processed []ProcessedEntity, outputPath string) error {
	output := make([]EntityOutput, 0, len(processed))

	for _, entity := range processed {
		entityOutput := convertToEntityOutput(entity)
		output = append(output, entityOutput)
	}

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(output); err != nil {
		return fmt.Errorf("write JSON: %w", err)
	}

	fmt.Printf("Saved %d entities to %s\n", len(output), outputPath)
	return nil
}

// ProcessedEntity -> EntityOutput
func convertToEntityOutput(pe ProcessedEntity) EntityOutput {
	metadata := EntityMetadata{
		SourceFile:  pe.File.Path,
		FileSize:    pe.File.Size,
		ModTime:     pe.File.ModTime,
		ContentHash: fmt.Sprintf("%08x", pe.ContentHash),
		ProcessedAt: time.Now(),
	}

	if pe.ParsedData != nil {
		if module, ok := pe.ParsedData["module"].(string); ok {
			metadata.Module = module
		}
		if object, ok := pe.ParsedData["object"].(string); ok {
			metadata.Object = object
		}
		if property, ok := pe.ParsedData["property"].(string); ok {
			metadata.Property = property
		}
		if code, ok := pe.ParsedData["code"].(string); ok {
			metadata.Code = code
		}
		if name, ok := pe.ParsedData["name"].(string); ok {
			metadata.Name = name
		}
	}

	var jsonDataStr string
	if pe.JSONData != nil {
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, pe.JSONData, "", "  "); err == nil {
			jsonDataStr = prettyJSON.String()
		} else {
			jsonDataStr = string(pe.JSONData)
		}
	}

	// validation result
	validation := ValidationResult{
		IsValid:    pe.FatalError == nil && len(pe.Errors) == 0,
		HasFatal:   pe.FatalError != nil,
		ErrorCount: len(pe.Errors),
		Errors:     pe.Errors,
	}

	if pe.FatalError != nil {
		validation.FatalError = pe.FatalError.Error()
	}

	return EntityOutput{
		Metadata:   metadata,
		ParsedData: pe.ParsedData,
		JSONData:   jsonDataStr,
		Schema:     pe.Schema,
		Validation: validation,
	}
}
