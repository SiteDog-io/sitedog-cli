package detectors

import (
	"os"
	"path/filepath"
)

// RenderDetector detects Render hosting configuration
type RenderDetector struct{}

func (r *RenderDetector) Name() string {
	return "render"
}

func (r *RenderDetector) Description() string {
	return "Render hosting detector"
}

func (r *RenderDetector) ShouldRun() bool {
	// Check for render.yaml or render.yml
	configPaths := []string{"render.yaml", "render.yml"}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

func (r *RenderDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Try to find Render config file
	configPaths := []string{"render.yaml", "render.yml"}
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
			Value:       "https://dashboard.render.com/web/srv-" + currentDir,
			Description: "Render hosting configuration detected",
			Confidence:  0.95,
		})
	}

	return results, nil
}