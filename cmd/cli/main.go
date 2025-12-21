package main

import (
	"fmt"
	"log"
	"yielda/runtime/internal/preset"
)

func main() {
	preset, err := preset.LoadPreset("./../../sources/presets/package/")

	if err != nil {
		log.Fatal("Failed to load preset: ", err)
	}

	fmt.Printf("Loaded preset: %s v%s\n", preset.Name, preset.Version)
}
