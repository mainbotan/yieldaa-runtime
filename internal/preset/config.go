package preset

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"unicode/utf8"

	"github.com/ghodss/yaml"
)

// loadConfig - try to load package.yml
func loadConfig(dir string) (*Package, error) {
	configPath := filepath.Join(dir, "package.yml")

	config, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read package.yml: %w", err)
	}

	var pkg Package
	if err := yaml.Unmarshal(config, &pkg); err != nil {
		return nil, fmt.Errorf("invalid YAML in package.yml: %w", err)
	}

	if err := validateConfig(pkg); err != nil {
		return nil, fmt.Errorf("error validating package.yml: %w", err)
	}

	return &pkg, nil
}

// validateConfig - validate config fields
func validateConfig(pkg Package) error {
	// name
	if pkg.Name == "" {
		return fmt.Errorf("'name' is required")
	}

	nameLen := utf8.RuneCountInString(pkg.Name)
	if nameLen < 4 || nameLen > 32 {
		return fmt.Errorf(
			"'name' must be 4-32 characters, got %d (%s)",
			nameLen, pkg.Name)
	}

	// version
	if pkg.Version == "" {
		return fmt.Errorf("'version' is required")
	}

	versionRegex := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	if !versionRegex.MatchString(pkg.Version) {
		return fmt.Errorf(
			"'version' must be X.Y.Z format (e.g. 0.0.1), got %s",
			pkg.Version)
	}

	// region
	if pkg.Region != "" {
		regionLen := utf8.RuneCountInString(pkg.Region)
		if regionLen < 2 || regionLen > 3 {
			return fmt.Errorf(
				"'region' must be 2-3 characters (e.g. 'ru'), got %s",
				pkg.Region)
		}
	}

	return nil
}
