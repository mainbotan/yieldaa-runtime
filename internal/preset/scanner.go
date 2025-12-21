package preset

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ScanEntities(dir string) ([]EntityFile, error) {
	var files []EntityFile

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(strings.ToLower(info.Name()), ".yml") &&
			!strings.HasSuffix(strings.ToLower(info.Name()), ".yaml") {
			return nil
		}

		files = append(files, EntityFile{
			Path:    path,
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	return files, nil
}
