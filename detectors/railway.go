package detectors

import (
	"os"
	"path/filepath"
)

// RailwayDetector detects Railway hosting configuration
type RailwayDetector struct{}

func (r *RailwayDetector) Name() string {
	return "railway"
}

func (r *RailwayDetector) Description() string {
	return "Railway hosting detector"
}

func (r *RailwayDetector) ShouldRun() bool {
	// Check for railway.json, railway.toml, or .railway directory
	configPaths := []string{"railway.json", "railway.toml", ".railway"}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

func (r *RailwayDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Try to find Railway config file or directory
	configPaths := []string{"railway.json", "railway.toml", ".railway"}
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
			Key:         "hosting",
			Value:       "https://railway.app/project/" + currentDir,
			Description: "Railway hosting configuration detected",
			Confidence:  0.95,
		})
	}

	return results, nil
}