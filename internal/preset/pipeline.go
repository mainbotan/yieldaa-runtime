package preset

// centralized entrypoint
func LoadAndProcessPreset(dir string, workers int) (*Package, []ProcessedEntity, []error) {
	pkg, err := LoadPreset(dir)
	if err != nil {
		return nil, nil, []error{err}
	}

	if pkg.EntitiesCount == 0 {
		return pkg, []ProcessedEntity{}, nil
	}

	processed, fatalErrors := ProcessEntities(pkg.EntitiesFiles, workers)
	return pkg, processed, fatalErrors
}
