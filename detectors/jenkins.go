package detectors

import (
	"os"
	"path/filepath"
)

// JenkinsDetector detects Jenkins configuration
type JenkinsDetector struct{}

func (j *JenkinsDetector) Name() string {
	return "jenkins"
}

func (j *JenkinsDetector) Description() string {
	return "Jenkins CI/CD detector"
}

func (j *JenkinsDetector) ShouldRun() bool {
	// Check for Jenkinsfile in various locations
	configPaths := []string{
		"Jenkinsfile",
		"jenkins/Jenkinsfile",
		".jenkins/Jenkinsfile",
		"ci/Jenkinsfile",
	}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

func (j *JenkinsDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Try to find Jenkins config file
	configPaths := []string{
		"Jenkinsfile",
		"jenkins/Jenkinsfile",
		".jenkins/Jenkinsfile",
		"ci/Jenkinsfile",
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
			Value:       "https://jenkins.your-domain.com/job/" + currentDir,
			Description: "Jenkins configuration detected",
			Confidence:  0.95,
		})
	}

	return results, nil
}