package detectors

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// NetlifyDetector detects Netlify hosting configuration
type NetlifyDetector struct{}

func (n *NetlifyDetector) Name() string {
	return "netlify"
}

func (n *NetlifyDetector) Description() string {
	return "Netlify hosting detector"
}

func (n *NetlifyDetector) ShouldRun() bool {
	// Check for netlify.toml config file
	if _, err := os.Stat("netlify.toml"); err == nil {
		return true
	}

	// Check for _redirects file (Netlify specific)
	if _, err := os.Stat("_redirects"); err == nil {
		return true
	}

	// Check for .netlify directory
	if _, err := os.Stat(".netlify"); err == nil {
		return true
	}

	// Check for public/_redirects (common location)
	if _, err := os.Stat("public/_redirects"); err == nil {
		return true
	}

	return false
}

func (n *NetlifyDetector) Detect() (*DetectionResult, error) {
	// Try to get site name from various sources
	var siteName string

	// Try to read .netlify/state.json for site info
	if _, err := os.Stat(".netlify/state.json"); err == nil {
		// This would contain site ID and other info, but it's complex to parse
		// For now, we'll use git repo name as fallback
	}

	// Try to get git repo name as site name
	if originURL, err := getGitOriginURL(); err == nil && originURL != "" {
		repoURL := convertToHTTPSURL(originURL)
		// Extract repo name from URL like https://github.com/user/repo
		parts := strings.Split(repoURL, "/")
		if len(parts) >= 2 {
			siteName = parts[len(parts)-1]
		}
	}

	// Check if netlify.toml contains build settings to increase confidence
	confidence := 0.8
	if _, err := os.Stat("netlify.toml"); err == nil {
		data, err := ioutil.ReadFile("netlify.toml")
		if err == nil {
			content := string(data)
			// If netlify.toml contains build configuration, higher confidence
			if strings.Contains(content, "[build]") || strings.Contains(content, "command") {
				confidence = 0.95
			}
		}
	}

	// Construct Netlify site URL
	var hostingURL string
	if siteName != "" {
		// Try to construct site URL - Netlify uses various patterns
		// Most common is sitename.netlify.app or custom domain
		hostingURL = fmt.Sprintf("https://app.netlify.com/sites/%s", siteName)
	} else {
		hostingURL = "https://app.netlify.com/sites"
	}

	return &DetectionResult{
		Key:         "hosting",
		Value:       hostingURL,
		Description: "Netlify hosting configuration detected",
		Confidence:  confidence,
	}, nil
}