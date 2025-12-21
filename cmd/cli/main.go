// main.go
package main

import (
	"os"
	"yielda/runtime/internal/preset"
)

func main() {
	pkg, processed, fatalErrs := preset.LoadAndProcessPreset(
		"./../../sources/presets/package3/", 5)

	preset.PrintResults(pkg, processed, fatalErrs)

	if len(fatalErrs) > 0 || preset.HasValidationErrors(processed) {
		os.Exit(1)
	}
}
