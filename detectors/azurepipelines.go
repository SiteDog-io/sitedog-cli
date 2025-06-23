package detectors

import (
	"os"
	"path/filepath"
)

// AzurePipelinesDetector detects Azure Pipelines configuration
type AzurePipelinesDetector struct{}

func (a *AzurePipelinesDetector) Name() string {
	return "azure-pipelines"
}

func (a *AzurePipelinesDetector) Description() string {
	return "Azure Pipelines CI/CD detector"
}

func (a *AzurePipelinesDetector) ShouldRun() bool {
	// Check for azure-pipelines.yml, azure-pipelines.yaml, or .azure-pipelines/azure-pipelines.yml
	configPaths := []string{
		"azure-pipelines.yml",
		"azure-pipelines.yaml",
		".azure-pipelines/azure-pipelines.yml",
		".azure-pipelines/azure-pipelines.yaml",
	}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

func (a *AzurePipelinesDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Try to find Azure Pipelines config file
	configPaths := []string{
		"azure-pipelines.yml",
		"azure-pipelines.yaml",
		".azure-pipelines/azure-pipelines.yml",
		".azure-pipelines/azure-pipelines.yaml",
	}
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
			Value:       "https://dev.azure.com/your-org/" + currentDir + "/_build",
			Description: "Azure Pipelines configuration detected",
			Confidence:  0.95,
		})
	}

	return results, nil
}