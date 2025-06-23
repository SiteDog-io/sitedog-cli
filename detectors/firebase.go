package detectors

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// FirebaseDetector detects Firebase hosting configuration
type FirebaseDetector struct{}

func (f *FirebaseDetector) Name() string {
	return "firebase"
}

func (f *FirebaseDetector) Description() string {
	return "Firebase hosting detector"
}

func (f *FirebaseDetector) ShouldRun() bool {
	// Check for firebase.json config file
	if _, err := os.Stat("firebase.json"); err == nil {
		return true
	}

	// Check for .firebaserc project file
	if _, err := os.Stat(".firebaserc"); err == nil {
		return true
	}

	// Check for firebase directory (functions)
	if _, err := os.Stat("firebase"); err == nil {
		return true
	}

	// Check for functions directory
	if _, err := os.Stat("functions"); err == nil {
		// Check if it's actually Firebase functions
		if _, err := os.Stat("functions/package.json"); err == nil {
			return true
		}
	}

	return false
}

func (f *FirebaseDetector) Detect() ([]*DetectionResult, error) {
	confidence := 0.0
	var projectID string

	// 1. Check for .firebaserc (contains project ID)
	if _, err := os.Stat(".firebaserc"); err == nil {
		if data, err := ioutil.ReadFile(".firebaserc"); err == nil {
			var config map[string]interface{}
			if err := json.Unmarshal(data, &config); err == nil {
				if projects, ok := config["projects"]; ok {
					if projectsMap, ok := projects.(map[string]interface{}); ok {
						// Try to get default project
						if defaultProject, ok := projectsMap["default"]; ok {
							if projectStr, ok := defaultProject.(string); ok {
								projectID = projectStr
								confidence = 1.0 // .firebaserc is very reliable
							}
						}
					}
				}
			}
		}
	}

	// 2. Check for firebase.json (hosting configuration)
	if _, err := os.Stat("firebase.json"); err == nil {
		if data, err := ioutil.ReadFile("firebase.json"); err == nil {
			var config map[string]interface{}
			if err := json.Unmarshal(data, &config); err == nil {
				// Check if hosting is configured
				if _, hasHosting := config["hosting"]; hasHosting {
					if confidence < 0.95 {
						confidence = 0.95
					}
				}
				// Check for other Firebase services
				if _, hasFunctions := config["functions"]; hasFunctions {
					if confidence < 0.9 {
						confidence = 0.9
					}
				}
				if _, hasFirestore := config["firestore"]; hasFirestore {
					if confidence < 0.9 {
						confidence = 0.9
					}
				}
			}
		}
		if confidence == 0.0 {
			confidence = 0.8 // firebase.json exists but no specific config found
		}
	}

	// 3. Check for Firebase functions
	if _, err := os.Stat("functions/package.json"); err == nil {
		if data, err := ioutil.ReadFile("functions/package.json"); err == nil {
			content := strings.ToLower(string(data))
			if strings.Contains(content, "firebase-functions") {
				if confidence < 0.85 {
					confidence = 0.85
				}
			}
		}
	}

	// Construct Firebase console URL
	var hostingURL string
	if projectID != "" {
		hostingURL = fmt.Sprintf("https://console.firebase.google.com/project/%s", projectID)
	} else {
		hostingURL = "https://console.firebase.google.com"
	}

	return []*DetectionResult{{
		Key:         "hosting",
		Value:       hostingURL,
		Description: "Firebase hosting configuration detected",
		Confidence:  confidence,
	}}, nil
}