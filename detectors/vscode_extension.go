package detectors

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// VSCodeExtensionDetector detects VS Code Extension projects
type VSCodeExtensionDetector struct{}

func (v *VSCodeExtensionDetector) Name() string {
	return "vscode-extension"
}

func (v *VSCodeExtensionDetector) Description() string {
	return "VS Code Extension platform detector"
}

func (v *VSCodeExtensionDetector) ShouldRun() bool {
	// Check for VS Code extension package.json
	if _, err := os.Stat("package.json"); err == nil {
		// Read and verify it's a VS Code extension
		if data, readErr := ioutil.ReadFile("package.json"); readErr == nil {
			var packageJson map[string]interface{}
			if json.Unmarshal(data, &packageJson) == nil {
				// Check for VS Code extension specific fields
				if engines, hasEngines := packageJson["engines"].(map[string]interface{}); hasEngines {
					if _, hasVSCode := engines["vscode"]; hasVSCode {
						return true
					}
				}
				if contributes, hasContributes := packageJson["contributes"]; hasContributes {
					_ = contributes // VS Code extensions have contributes field
					return true
				}
				if categories, hasCategories := packageJson["categories"].([]interface{}); hasCategories {
					for _, category := range categories {
						if catStr, ok := category.(string); ok {
							catLower := strings.ToLower(catStr)
							if strings.Contains(catLower, "extension") ||
							   catLower == "other" || catLower == "snippets" ||
							   catLower == "themes" || catLower == "language packs" {
								return true
							}
						}
					}
				}
			}
		}
	}

	// Also check for .vscodeignore file
	if _, err := os.Stat(".vscodeignore"); err == nil {
		return true
	}

	return false
}

func (v *VSCodeExtensionDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Read package.json
	data, err := ioutil.ReadFile("package.json")
	if err != nil {
		return results, nil
	}

	var packageJson map[string]interface{}
	if err := json.Unmarshal(data, &packageJson); err != nil {
		return results, nil
	}

	content := strings.ToLower(string(data))

		// Auto-add VS Code extension platform links (high confidence)
	results = append(results, &DetectionResult{
		Key:         "vscode_marketplace",
		Value:       "https://marketplace.visualstudio.com/manage",
		Description: "VS Code Marketplace Publisher Portal detected for VS Code Extension",
		Confidence:  0.98, // Auto-add
	})

	results = append(results, &DetectionResult{
		Key:         "vscode_developer",
		Value:       "https://code.visualstudio.com/api",
		Description: "VS Code Extension API Documentation detected",
		Confidence:  0.98, // Auto-add
	})

	// Check if using Azure DevOps for publishing (common pattern)
	if strings.Contains(content, "vsce") || strings.Contains(content, "azure devops") {
		results = append(results, &DetectionResult{
			Key:         "azure_devops",
			Value:       "https://dev.azure.com",
			Description: "Azure DevOps detected for VS Code Extension publishing",
			Confidence:  0.90, // Conditional auto-add
		})
	}

	// Detect specific VS Code Extension services and APIs
	services := map[string]map[string]interface{}{
		// Analytics for extensions
		"telemetry": {
			"patterns": []string{"telemetry", "vscode.env.telemetrylevel", "application insights"},
			"name": "VS Code Telemetry",
			"url": "https://code.visualstudio.com/api/extension-guides/telemetry",
			"key": "analytics",
		},

		// External analytics
		"google_analytics_vscode": {
			"patterns": []string{"google analytics", "gtag", "ga.js", "analytics.js"},
			"name": "Google Analytics for Extensions",
			"url": "https://analytics.google.com",
			"key": "analytics",
		},

		// VS Code APIs usage
		"vscode_workspace": {
			"patterns": []string{"vscode.workspace", "workspace api"},
			"name": "VS Code Workspace API",
			"url": "https://code.visualstudio.com/api/references/vscode-api#workspace",
			"key": "vscode_api_workspace",
		},

		"vscode_window": {
			"patterns": []string{"vscode.window", "window api"},
			"name": "VS Code Window API",
			"url": "https://code.visualstudio.com/api/references/vscode-api#window",
			"key": "vscode_api_window",
		},

		"vscode_commands": {
			"patterns": []string{"vscode.commands", "commands api", "registercommand"},
			"name": "VS Code Commands API",
			"url": "https://code.visualstudio.com/api/references/vscode-api#commands",
			"key": "vscode_api_commands",
		},

		// Language Server Protocol
		"language_server": {
			"patterns": []string{"language server", "lsp", "languageclient"},
			"name": "Language Server Protocol",
			"url": "https://microsoft.github.io/language-server-protocol",
			"key": "language_server",
		},

		// External services commonly used in extensions
		"github_api": {
			"patterns": []string{"github api", "octokit", "github.com/api"},
			"name": "GitHub API",
			"url": "https://docs.github.com/en/rest",
			"key": "github_integration",
		},

		"openai_api": {
			"patterns": []string{"openai", "gpt", "chatgpt", "api.openai.com"},
			"name": "OpenAI API",
			"url": "https://platform.openai.com",
			"key": "ai_integration",
		},

		"firebase_vscode": {
			"patterns": []string{"firebase", "firestore", "firebase.google.com"},
			"name": "Firebase for Extensions",
			"url": "https://console.firebase.google.com",
			"key": "cloud",
		},

		"sentry_vscode": {
			"patterns": []string{"sentry", "@sentry/node"},
			"name": "Sentry for Extensions",
			"url": "https://sentry.io",
			"key": "monitoring",
		},

		// CI/CD for extensions
		"github_actions_vscode": {
			"patterns": []string{"github actions", "actions/checkout", ".github/workflows"},
			"name": "GitHub Actions for VS Code Extensions",
			"url": "https://github.com/features/actions",
			"key": "ci",
		},
	}

	// Check for specific services
	serviceOrder := []string{
		"telemetry", "google_analytics_vscode",
		"vscode_workspace", "vscode_window", "vscode_commands",
		"language_server", "github_api", "openai_api",
		"firebase_vscode", "sentry_vscode", "github_actions_vscode",
	}

	for _, serviceKey := range serviceOrder {
		serviceInfo := services[serviceKey]
		patterns := serviceInfo["patterns"].([]string)

		for _, pattern := range patterns {
			if strings.Contains(content, pattern) {
				results = append(results, &DetectionResult{
					Key:         serviceInfo["key"].(string),
					Value:       serviceInfo["url"].(string),
					Description: fmt.Sprintf("%s detected in VS Code Extension", serviceInfo["name"].(string)),
					Confidence:  0.85,
				})
				break // Only add each service once
			}
		}
	}

	return results, nil
}