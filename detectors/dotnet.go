package detectors

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// DotNetDetector detects external services in .NET projects
type DotNetDetector struct{}

func (d *DotNetDetector) Name() string {
	return "dotnet"
}

func (d *DotNetDetector) Description() string {
	return ".NET/C# external services detector"
}

func (d *DotNetDetector) ShouldRun() bool {
	// Check for .csproj files, packages.config, or project.json
	if d.findCsprojFiles() {
		return true
	}

	configPaths := []string{"packages.config", "project.json"}
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

func (d *DotNetDetector) findCsprojFiles() bool {
	// Look for *.csproj files in current directory
	files, err := filepath.Glob("*.csproj")
	if err != nil {
		return false
	}
	return len(files) > 0
}

func (d *DotNetDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult
	var allContent string

	// Read .csproj files
	csprojFiles, err := filepath.Glob("*.csproj")
	if err == nil {
		for _, file := range csprojFiles {
			if data, err := ioutil.ReadFile(file); err == nil {
				allContent += string(data) + "\n"
			}
		}
	}

	// Read packages.config if exists
	if data, err := ioutil.ReadFile("packages.config"); err == nil {
		allContent += string(data) + "\n"
	}

	// Read project.json if exists
	if data, err := ioutil.ReadFile("project.json"); err == nil {
		allContent += string(data) + "\n"
	}

	if allContent == "" {
		return results, nil
	}

	content := strings.ToLower(allContent)

	// Map of .NET package names to service info
	services := map[string]map[string]interface{}{
		// Monitoring and Error Tracking
		"sentry": {
			"patterns": []string{"sentry", "sentry.aspnetcore", "sentry.extensions.logging"},
			"name": "Sentry",
			"url": "https://sentry.io",
			"key": "monitoring",
		},
		"bugsnag": {
			"patterns": []string{"bugsnag", "bugsnag.aspnet", "bugsnag.aspnetcore"},
			"name": "Bugsnag",
			"url": "https://bugsnag.com",
			"key": "monitoring",
		},
		"rollbar": {
			"patterns": []string{"rollbar", "rollbar.netcore.aspnet"},
			"name": "Rollbar",
			"url": "https://rollbar.com",
			"key": "monitoring",
		},
		"airbrake": {
			"patterns": []string{"sharpbrake"},
			"name": "Airbrake",
			"url": "https://airbrake.io",
			"key": "monitoring",
		},
		"datadog": {
			"patterns": []string{"datadog.trace", "dogstatsd-csharp-client"},
			"name": "Datadog",
			"url": "https://datadog.com",
			"key": "monitoring",
		},
		"newrelic": {
			"patterns": []string{"newrelic.agent", "newrelic.agent.api"},
			"name": "New Relic",
			"url": "https://newrelic.com",
			"key": "monitoring",
		},

		// Analytics
		"mixpanel": {
			"patterns": []string{"mixpanel"},
			"name": "Mixpanel",
			"url": "https://mixpanel.com",
			"key": "analytics",
		},
		"amplitude": {
			"patterns": []string{"amplitude"},
			"name": "Amplitude",
			"url": "https://amplitude.com",
			"key": "analytics",
		},
		"segment": {
			"patterns": []string{"analytics.net", "segment.analytics"},
			"name": "Segment",
			"url": "https://segment.com",
			"key": "analytics",
		},

		// Payments
		"stripe": {
			"patterns": []string{"stripe.net"},
			"name": "Stripe",
			"url": "https://stripe.com",
			"key": "payments",
		},
		"paypal": {
			"patterns": []string{"paypal"},
			"name": "PayPal",
			"url": "https://paypal.com",
			"key": "payments",
		},

		// Email Services
		"sendgrid": {
			"patterns": []string{"sendgrid"},
			"name": "SendGrid",
			"url": "https://sendgrid.com",
			"key": "email_delivery",
		},
		"mailgun": {
			"patterns": []string{"mailgun", "restsharp.portable.mailgun"},
			"name": "Mailgun",
			"url": "https://mailgun.com",
			"key": "email_delivery",
		},
		"postmark": {
			"patterns": []string{"postmark", "postmark.net"},
			"name": "Postmark",
			"url": "https://postmarkapp.com",
			"key": "email_delivery",
		},

		// Cloud Platforms
		"aws": {
			"patterns": []string{"awssdk", "amazon.", "aws."},
			"name": "AWS",
			"url": "https://console.aws.amazon.com",
			"key": "cloud",
		},
		"azure": {
			"patterns": []string{"azure.", "microsoft.azure", "windowsazure.storage"},
			"name": "Azure",
			"url": "https://portal.azure.com",
			"key": "cloud",
		},
		"gcp": {
			"patterns": []string{"google.cloud", "google.apis"},
			"name": "Google Cloud",
			"url": "https://console.cloud.google.com",
			"key": "cloud",
		},
	}

	// Check each service and collect all found results
	serviceOrder := []string{
		"sentry", "bugsnag", "rollbar", "airbrake", "datadog", "newrelic",
		"mixpanel", "amplitude", "segment",
		"stripe", "paypal",
		"sendgrid", "mailgun", "postmark",
		"aws", "azure", "gcp",
	}

	for _, serviceKey := range serviceOrder {
		serviceInfo := services[serviceKey]
		patterns := serviceInfo["patterns"].([]string)

		for _, pattern := range patterns {
			// Use regex to match package references in various .NET formats
			// PackageReference, package id, dependencies, etc.
			escapedPattern := regexp.QuoteMeta(pattern)
			regex := regexp.MustCompile(`(?i)(packagereference|package\s+id|"` + escapedPattern + `"|>` + escapedPattern + `<)`)

			if regex.MatchString(content) {
				results = append(results, &DetectionResult{
					Key:         serviceInfo["key"].(string),
					Value:       serviceInfo["url"].(string),
					Description: fmt.Sprintf("%s service detected in .NET project", serviceInfo["name"].(string)),
					Confidence:  0.9,
				})
				break // Only add each service once
			}
		}
	}

	return results, nil
}