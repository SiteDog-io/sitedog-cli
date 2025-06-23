package detectors

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// ChromeExtensionDetector detects Chrome Extension projects
type ChromeExtensionDetector struct{}

func (c *ChromeExtensionDetector) Name() string {
	return "chrome-extension"
}

func (c *ChromeExtensionDetector) Description() string {
	return "Chrome Extension platform detector"
}

func (c *ChromeExtensionDetector) ShouldRun() bool {
	// Check for Chrome Extension manifest
	if _, err := os.Stat("manifest.json"); err == nil {
		// Read and verify it's a Chrome extension manifest
		if data, readErr := ioutil.ReadFile("manifest.json"); readErr == nil {
			var manifest map[string]interface{}
			if json.Unmarshal(data, &manifest) == nil {
				// Check for Chrome extension specific fields
				if _, hasManifestVersion := manifest["manifest_version"]; hasManifestVersion {
					return true
				}
				if _, hasPermissions := manifest["permissions"]; hasPermissions {
					return true
				}
				if _, hasContentScripts := manifest["content_scripts"]; hasContentScripts {
					return true
				}
			}
		}
	}
	return false
}

func (c *ChromeExtensionDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Read manifest.json
	data, err := ioutil.ReadFile("manifest.json")
	if err != nil {
		return results, nil
	}

	var manifest map[string]interface{}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return results, nil
	}

	content := strings.ToLower(string(data))

		// Auto-add Chrome extension platform links (high confidence)
	results = append(results, &DetectionResult{
		Key:         "chrome_webstore",
		Value:       "https://chrome.google.com/webstore/devconsole",
		Description: "Chrome Web Store Developer Console detected for Chrome Extension",
		Confidence:  0.98, // Auto-add
	})

	results = append(results, &DetectionResult{
		Key:         "chrome_developer",
		Value:       "https://developer.chrome.com/docs/extensions",
		Description: "Chrome Extensions Developer Documentation detected",
		Confidence:  0.98, // Auto-add
	})

	// Detect specific Chrome Extension services and APIs
	services := map[string]map[string]interface{}{
		// Analytics for extensions
		"google_analytics_extension": {
			"patterns": []string{"google analytics", "gtag", "ga.js", "analytics.js"},
			"name": "Google Analytics for Extensions",
			"url": "https://analytics.google.com",
			"key": "analytics",
		},

		// Chrome APIs usage
		"chrome_storage": {
			"patterns": []string{"chrome.storage", "storage api"},
			"name": "Chrome Storage API",
			"url": "https://developer.chrome.com/docs/extensions/reference/storage",
			"key": "chrome_api_storage",
		},

		"chrome_tabs": {
			"patterns": []string{"chrome.tabs", "tabs api"},
			"name": "Chrome Tabs API",
			"url": "https://developer.chrome.com/docs/extensions/reference/tabs",
			"key": "chrome_api_tabs",
		},

		"chrome_notifications": {
			"patterns": []string{"chrome.notifications", "notifications api"},
			"name": "Chrome Notifications API",
			"url": "https://developer.chrome.com/docs/extensions/reference/notifications",
			"key": "chrome_api_notifications",
		},

		// External services commonly used in extensions
		"firebase_extension": {
			"patterns": []string{"firebase", "firestore", "firebase.google.com"},
			"name": "Firebase for Extensions",
			"url": "https://console.firebase.google.com",
			"key": "cloud",
		},

		"sentry_extension": {
			"patterns": []string{"sentry", "@sentry/browser"},
			"name": "Sentry for Extensions",
			"url": "https://sentry.io",
			"key": "monitoring",
		},

		// Payment services for extensions
		"stripe_extension": {
			"patterns": []string{"stripe", "stripe.com"},
			"name": "Stripe for Extensions",
			"url": "https://stripe.com",
			"key": "payments",
		},

		// Extension stores and distribution
		"edge_addons": {
			"patterns": []string{"edge", "microsoft edge", "edge addons"},
			"name": "Microsoft Edge Addons",
			"url": "https://partner.microsoft.com/en-us/dashboard/microsoftedge",
			"key": "edge_distribution",
		},

		"firefox_addons": {
			"patterns": []string{"firefox", "mozilla", "firefox addon"},
			"name": "Firefox Add-ons",
			"url": "https://addons.mozilla.org/developers",
			"key": "firefox_distribution",
		},
	}

	// Check for specific services
	serviceOrder := []string{
		"google_analytics_extension", "chrome_storage", "chrome_tabs", "chrome_notifications",
		"firebase_extension", "sentry_extension", "stripe_extension",
		"edge_addons", "firefox_addons",
	}

	for _, serviceKey := range serviceOrder {
		serviceInfo := services[serviceKey]
		patterns := serviceInfo["patterns"].([]string)

		for _, pattern := range patterns {
			if strings.Contains(content, pattern) {
				results = append(results, &DetectionResult{
					Key:         serviceInfo["key"].(string),
					Value:       serviceInfo["url"].(string),
					Description: fmt.Sprintf("%s detected in Chrome Extension", serviceInfo["name"].(string)),
					Confidence:  0.85,
				})
				break // Only add each service once
			}
		}
	}

	return results, nil
}