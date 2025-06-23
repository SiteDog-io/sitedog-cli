package detectors

import (
	"os"
	"path/filepath"
)

// BitbucketPipelinesDetector detects Bitbucket Pipelines configuration
type BitbucketPipelinesDetector struct{}

func (b *BitbucketPipelinesDetector) Name() string {
	return "bitbucket-pipelines"
}

func (b *BitbucketPipelinesDetector) Description() string {
	return "Bitbucket Pipelines CI/CD detector"
}

func (b *BitbucketPipelinesDetector) ShouldRun() bool {
	// Check for bitbucket-pipelines.yml or bitbucket-pipelines.yaml
	configPaths := []string{"bitbucket-pipelines.yml", "bitbucket-pipelines.yaml"}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

func (b *BitbucketPipelinesDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Try to find Bitbucket Pipelines config file
	configPaths := []string{"bitbucket-pipelines.yml", "bitbucket-pipelines.yaml"}
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
			Value:       "https://bitbucket.org/your-workspace/" + currentDir + "/addon/pipelines/home",
			Description: "Bitbucket Pipelines configuration detected",
			Confidence:  0.95,
		})
	}

	return results, nil
}