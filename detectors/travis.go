package detectors

import (
	"os"
	"path/filepath"
)

// TravisCIDetector detects Travis CI configuration
type TravisCIDetector struct{}

func (t *TravisCIDetector) Name() string {
	return "travis-ci"
}

func (t *TravisCIDetector) Description() string {
	return "Travis CI detector"
}

func (t *TravisCIDetector) ShouldRun() bool {
	// Check for .travis.yml or .travis.yaml
	configPaths := []string{".travis.yml", ".travis.yaml"}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

func (t *TravisCIDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Try to find Travis CI config file
	configPaths := []string{".travis.yml", ".travis.yaml"}
	var foundConfig string

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			foundConfig = path
			break
		}
	}

	if foundConfig != "" {
		// Get project name from current directory
		currentDir, err := os.Getwd()
		if err != nil {
			currentDir = "project"
		} else {
			currentDir = filepath.Base(currentDir)
		}

		results = append(results, &DetectionResult{
			Key:         "ci",
			Value:       "https://app.travis-ci.com/github/your-org/" + currentDir,
			Description: "Travis CI configuration detected",
			Confidence:  0.95,
		})
	}

	return results, nil
}