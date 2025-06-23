package detectors

import (
	"os"
	"path/filepath"
)

// CircleCIDetector detects CircleCI configuration
type CircleCIDetector struct{}

func (c *CircleCIDetector) Name() string {
	return "circleci"
}

func (c *CircleCIDetector) Description() string {
	return "CircleCI CI/CD detector"
}

func (c *CircleCIDetector) ShouldRun() bool {
	// Check for .circleci/config.yml or .circleci/config.yaml
	configPaths := []string{".circleci/config.yml", ".circleci/config.yaml"}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

func (c *CircleCIDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Try to find CircleCI config file
	configPaths := []string{".circleci/config.yml", ".circleci/config.yaml"}
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
			Value:       "https://app.circleci.com/pipelines/github/your-org/" + currentDir,
			Description: "CircleCI configuration detected",
			Confidence:  0.95,
		})
	}

	return results, nil
}