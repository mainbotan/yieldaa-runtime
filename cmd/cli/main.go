// main.go
package main

import (
	"fmt"
	"os"
	"yielda/runtime/internal/preset"
)

func main() {
	pkg, processed, fatalErrs := preset.LoadAndProcessPreset(
		"./../../sources/presets/package4/", 100)

	preset.PrintResults(pkg, processed, fatalErrs)

	if len(fatalErrs) > 0 || preset.HasValidationErrors(processed) {
		os.Exit(1)
	}

	if len(processed) > 0 {
		outputPath := "./output/entities.json"
		if err := preset.SaveEntitiesToJSON(processed, outputPath); err != nil {
			fmt.Printf("Failed to save JSON: %v\n", err)
		}
	}

	if len(fatalErrs) > 0 || preset.HasValidationErrors(processed) {
		os.Exit(1)
	}
}
