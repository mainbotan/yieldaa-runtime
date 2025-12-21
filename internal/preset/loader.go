package preset

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
	"strings"
)

// preset loader
func LoadPreset(dir string) (*Package, error) {
	// scan package.yml
	packageData, err := loadConfig(dir)
	if err != nil {
		return nil, fmt.Errorf("package load failed: %w", err)
	}

	// scan entities
	entitiesDir := filepath.Join(dir, "entities")
	if _, err := os.Stat(entitiesDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("package does not have entities in the /entities directory")
	}
	entityFiles, err := ScanEntities(entitiesDir)
	if err != nil {
		return nil, fmt.Errorf("entities scan failed: %w", err)
	}

	// entities meta to pkg struct
	packageData.EntitiesFiles = entityFiles
	packageData.EntitiesCount = len(entityFiles)

	// control sum
	totalSize := int64(0)
	for _, f := range entityFiles {
		totalSize += f.Size
	}
	packageData.EntitiesTotalSize = totalSize
	packageData.EntitiesStructureHash = calculateStructureHash(entityFiles)

	return packageData, nil
}

// structure hash
func calculateStructureHash(files []EntityFile) uint32 {
	hash := crc32.NewIEEE()
	for _, f := range files {
		relPath := strings.TrimPrefix(f.Path, filepath.Dir(f.Path)+string(os.PathSeparator))
		hash.Write([]byte(relPath))
		binary.Write(hash, binary.LittleEndian, f.Size)
		binary.Write(hash, binary.LittleEndian, f.ModTime.Unix())
	}
	return hash.Sum32()
}
