package detectors

import (
	"os"
	"path/filepath"
)

// FlyIODetector detects Fly.io hosting configuration
type FlyIODetector struct{}

func (f *FlyIODetector) Name() string {
	return "flyio"
}

func (f *FlyIODetector) Description() string {
	return "Fly.io hosting detector"
}

func (f *FlyIODetector) ShouldRun() bool {
	// Check for fly.toml
	_, err := os.Stat("fly.toml")
	return err == nil
}

func (f *FlyIODetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Check if fly.toml exists
	if _, err := os.Stat("fly.toml"); err == nil {
		// Get project name from current directory
		currentDir, err := os.Getwd()
		if err != nil {
			currentDir = "project"
		} else {
			currentDir = filepath.Base(currentDir)
		}

		results = append(results, &DetectionResult{
			Key:         "hosting",
			Value:       "https://fly.io/apps/" + currentDir,
			Description: "Fly.io hosting configuration detected",
			Confidence:  0.95,
		})
	}

	return results, nil
}