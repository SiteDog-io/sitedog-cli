package detectors

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// HerokuDetector detects Heroku hosting configuration
type HerokuDetector struct{}

func (h *HerokuDetector) Name() string {
	return "heroku"
}

func (h *HerokuDetector) Description() string {
	return "Heroku hosting detector"
}

func (h *HerokuDetector) ShouldRun() bool {
	// Check for heroku git remote (most reliable)
	if hasHerokuRemote() {
		return true
	}

	// Check for app.json (Heroku app configuration)
	if _, err := os.Stat("app.json"); err == nil {
		return true
	}

	// Don't trigger on Procfile alone (too generic)
	return false
}

func (h *HerokuDetector) Detect() (*DetectionResult, error) {
	confidence := 0.0
	var appName string
	isHeroku := false

	// 1. Check for heroku git remote (most reliable indicator)
	if herokuAppName := getHerokuAppName(); herokuAppName != "" {
		appName = herokuAppName
		confidence = 1.0
		isHeroku = true
	}

	// 2. Check for app.json (quite specific to Heroku)
	if _, err := os.Stat("app.json"); err == nil {
		if data, err := ioutil.ReadFile("app.json"); err == nil {
			var config map[string]interface{}
			if err := json.Unmarshal(data, &config); err == nil {
				if name, ok := config["name"]; ok {
					if nameStr, ok := name.(string); ok {
						appName = nameStr
					}
				}
			}
		}
		isHeroku = true
		if confidence < 0.9 {
			confidence = 0.9
		}
	}

		// 3. Check for Procfile (only as confidence booster, not trigger)
	if _, err := os.Stat("Procfile"); err == nil && isHeroku {
		// Procfile + other indicators = slightly higher confidence
		if confidence < 0.95 {
			confidence = 0.95
		}
	}

	// Only proceed if we have reliable Heroku indicators
	if !isHeroku {
		return nil, nil
	}

	// 4. Fallback to git repo name for app name
	if appName == "" {
		if originURL, err := getGitOriginURL(); err == nil && originURL != "" {
			repoURL := convertToHTTPSURL(originURL)
			parts := strings.Split(repoURL, "/")
			if len(parts) >= 2 {
				appName = parts[len(parts)-1]
			}
		}
	}

	// Construct Heroku dashboard URL
	var hostingURL string
	if appName != "" {
		hostingURL = fmt.Sprintf("https://dashboard.heroku.com/apps/%s", appName)
	} else {
		hostingURL = "https://dashboard.heroku.com/apps"
	}

	return &DetectionResult{
		Key:         "hosting",
		Value:       hostingURL,
		Description: "Heroku hosting configuration detected",
		Confidence:  confidence,
	}, nil
}

// hasHerokuRemote checks if there's a heroku git remote
func hasHerokuRemote() bool {
	cmd := exec.Command("git", "remote", "-v")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), "heroku.com")
}

// getHerokuAppName extracts app name from heroku git remote
func getHerokuAppName() string {
	cmd := exec.Command("git", "remote", "-v")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "heroku.com") && strings.Contains(line, "(push)") {
			// Extract app name from URL like: heroku	https://git.heroku.com/myapp.git (push)
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				url := parts[1]
				// Extract app name from https://git.heroku.com/APPNAME.git
				if strings.Contains(url, "git.heroku.com/") {
					start := strings.Index(url, "git.heroku.com/") + len("git.heroku.com/")
					end := strings.Index(url[start:], ".git")
					if end > 0 {
						return url[start : start+end]
					}
				}
			}
		}
	}

	return ""
}