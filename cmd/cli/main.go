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

	fmt.Printf("%s v%s | entities: %d | hash: %08x | size: %d bytes\n",
		preset.Name, preset.Version,
		preset.EntitiesCount,
		preset.EntitiesStructureHash,
		preset.EntitiesTotalSize)
}
