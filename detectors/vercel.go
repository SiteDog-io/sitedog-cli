package detectors

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// VercelDetector detects Vercel hosting configuration
type VercelDetector struct{}

func (v *VercelDetector) Name() string {
	return "vercel"
}

func (v *VercelDetector) Description() string {
	return "Vercel hosting detector"
}

func (v *VercelDetector) ShouldRun() bool {
	// Check for vercel.json config file
	if _, err := os.Stat("vercel.json"); err == nil {
		return true
	}

	// Check for .vercel directory
	if _, err := os.Stat(".vercel"); err == nil {
		return true
	}

	return false
}

func (v *VercelDetector) Detect() ([]*DetectionResult, error) {
	// Try to get project name and construct hosting URL
	var projectName string

	// First try to read vercel.json for project name
	if _, err := os.Stat("vercel.json"); err == nil {
		data, err := ioutil.ReadFile("vercel.json")
		if err == nil {
			var config map[string]interface{}
			if err := json.Unmarshal(data, &config); err == nil {
				if name, ok := config["name"]; ok {
					if nameStr, ok := name.(string); ok {
						projectName = nameStr
					}
				}
			}
		}
	}

	// Try to get git repo name if no project name found
	if projectName == "" {
		if originURL, err := getGitOriginURL(); err == nil && originURL != "" {
			repoURL := convertToHTTPSURL(originURL)
			// Extract repo name from URL like https://github.com/user/repo
			parts := strings.Split(repoURL, "/")
			if len(parts) >= 2 {
				projectName = parts[len(parts)-1]
			}
		}
	}

	// Construct Vercel dashboard URL
	var hostingURL string
	if projectName != "" {
		// Vercel dashboard URL format
		hostingURL = fmt.Sprintf("https://vercel.com/dashboard/%s", projectName)
	} else {
		hostingURL = "https://vercel.com/dashboard"
	}

	return []*DetectionResult{{
		Key:         "hosting",
		Value:       hostingURL,
		Description: "Vercel hosting configuration detected",
		Confidence:  0.9,
	}}, nil
}