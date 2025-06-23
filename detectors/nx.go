package detectors

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

// NxDetector detects Nx monorepo projects
type NxDetector struct{}

func (n *NxDetector) Name() string {
	return "nx"
}

func (n *NxDetector) Description() string {
	return "Nx monorepo detector"
}

func (n *NxDetector) ShouldRun() bool {
	// Check for Nx configuration files
	nxFiles := []string{"nx.json", "workspace.json", "project.json"}
	for _, file := range nxFiles {
		if _, err := os.Stat(file); err == nil {
			return true
		}
	}

	// Check for Nx in package.json
	if _, err := os.Stat("package.json"); err == nil {
		if data, readErr := ioutil.ReadFile("package.json"); readErr == nil {
			content := strings.ToLower(string(data))
			if strings.Contains(content, "@nrwl/") || strings.Contains(content, "@nx/") {
				return true
			}
		}
	}

	return false
}

func (n *NxDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Auto-add Nx platform links (high confidence)
	results = append(results, &DetectionResult{
		Key:         "nx_console",
		Value:       "https://nx.app",
		Description: "Nx Console detected for Nx monorepo management",
		Confidence:  0.98, // Auto-add
	})

	// Check for Nx Cloud access token
	if data, err := ioutil.ReadFile("nx.json"); err == nil {
		var nxConfig map[string]interface{}
		if json.Unmarshal(data, &nxConfig) == nil {
			if _, hasToken := nxConfig["nxCloudAccessToken"]; hasToken {
				results = append(results, &DetectionResult{
					Key:         "nx_cloud",
					Value:       "https://cloud.nx.app",
					Description: "Nx Cloud detected for distributed caching and CI",
					Confidence:  0.90, // Conditional auto-add
				})
			}
		}
	}

	// For Nx, we focus only on real services with dashboards, not documentation links
	// The main value is Nx Console and Nx Cloud for actual project management

	return results, nil
}