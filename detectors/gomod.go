package detectors

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// GoModDetector detects external services in Go go.mod
type GoModDetector struct{}

func (g *GoModDetector) Name() string {
	return "go-mod"
}

func (g *GoModDetector) Description() string {
	return "Go go.mod external services detector"
}

func (g *GoModDetector) ShouldRun() bool {
	_, err := os.Stat("go.mod")
	return err == nil
}

func (g *GoModDetector) Detect() ([]*DetectionResult, error) {
	data, err := ioutil.ReadFile("go.mod")
	if err != nil {
		return nil, err
	}

	content := string(data)

	// Extract all dependencies from go.mod using regex
	// Match lines like: github.com/some/package v1.2.3
	depRegex := regexp.MustCompile(`^\s*([^\s]+)\s+v[^\s]+`)
	lines := strings.Split(content, "\n")

	allDeps := make(map[string]bool)
	inRequireBlock := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check for require block start/end
		if strings.HasPrefix(line, "require (") {
			inRequireBlock = true
			continue
		}
		if inRequireBlock && line == ")" {
			inRequireBlock = false
			continue
		}

		// Parse dependency line
		if inRequireBlock || strings.HasPrefix(line, "require ") {
			// Remove "require " prefix if it exists
			depLine := strings.TrimPrefix(line, "require ")

			matches := depRegex.FindStringSubmatch(depLine)
			if len(matches) >= 2 {
				dep := strings.ToLower(matches[1])
				allDeps[dep] = true
			}
		}
	}

	// Map of Go module names to service info
	services := map[string]map[string]interface{}{
		// Monitoring and Error Tracking
		"sentry": {
			"patterns": []string{"github.com/getsentry/sentry-go"},
			"name": "Sentry",
			"url": "https://sentry.io",
			"key": "monitoring",
		},
		"bugsnag": {
			"patterns": []string{"github.com/bugsnag/bugsnag-go"},
			"name": "Bugsnag",
			"url": "https://bugsnag.com",
			"key": "monitoring",
		},
		"rollbar": {
			"patterns": []string{"github.com/rollbar/rollbar-go"},
			"name": "Rollbar",
			"url": "https://rollbar.com",
			"key": "monitoring",
		},
		"airbrake": {
			"patterns": []string{"github.com/airbrake/gobrake"},
			"name": "Airbrake",
			"url": "https://airbrake.io",
			"key": "monitoring",
		},
		"datadog": {
			"patterns": []string{"gopkg.in/datadog/dd-trace-go.v1", "github.com/datadog/datadog-go"},
			"name": "Datadog",
			"url": "https://datadog.com",
			"key": "monitoring",
		},
		"newrelic": {
			"patterns": []string{"github.com/newrelic/go-agent"},
			"name": "New Relic",
			"url": "https://newrelic.com",
			"key": "monitoring",
		},

		// Analytics
		"mixpanel": {
			"patterns": []string{"github.com/mixpanel/mixpanel-go"},
			"name": "Mixpanel",
			"url": "https://mixpanel.com",
			"key": "analytics",
		},
		"amplitude": {
			"patterns": []string{"github.com/amplitude/analytics-go"},
			"name": "Amplitude",
			"url": "https://amplitude.com",
			"key": "analytics",
		},
		"segment": {
			"patterns": []string{"github.com/segmentio/analytics-go"},
			"name": "Segment",
			"url": "https://segment.com",
			"key": "analytics",
		},
		"posthog": {
			"patterns": []string{"github.com/posthog/posthog-go"},
			"name": "PostHog",
			"url": "https://posthog.com",
			"key": "analytics",
		},

		// Payments
		"stripe": {
			"patterns": []string{"github.com/stripe/stripe-go"},
			"name": "Stripe",
			"url": "https://stripe.com",
			"key": "payments",
		},
		"paypal": {
			"patterns": []string{"github.com/plutov/paypal"},
			"name": "PayPal",
			"url": "https://paypal.com",
			"key": "payments",
		},

		// Email Services
		"sendgrid": {
			"patterns": []string{"github.com/sendgrid/sendgrid-go"},
			"name": "SendGrid",
			"url": "https://sendgrid.com",
			"key": "email_delivery",
		},
		"mailgun": {
			"patterns": []string{"github.com/mailgun/mailgun-go"},
			"name": "Mailgun",
			"url": "https://mailgun.com",
			"key": "email_delivery",
		},
		"postmark": {
			"patterns": []string{"github.com/mattbaird/gochimp", "github.com/keighl/postmark"},
			"name": "Postmark",
			"url": "https://postmarkapp.com",
			"key": "email_delivery",
		},

		// Cloud Platforms
		"aws": {
			"patterns": []string{"github.com/aws/aws-sdk-go", "github.com/aws/aws-sdk-go-v2"},
			"name": "AWS",
			"url": "https://console.aws.amazon.com",
			"key": "cloud",
		},
		"gcp": {
			"patterns": []string{"cloud.google.com/go", "google.golang.org/api"},
			"name": "Google Cloud",
			"url": "https://console.cloud.google.com",
			"key": "cloud",
		},
		"azure": {
			"patterns": []string{"github.com/azure/azure-sdk-for-go"},
			"name": "Azure",
			"url": "https://portal.azure.com",
			"key": "cloud",
		},
	}

	// Check each service and collect all found results
	var results []*DetectionResult
	serviceOrder := []string{
		"sentry", "bugsnag", "rollbar", "airbrake", "datadog", "newrelic",
		"mixpanel", "amplitude", "segment", "posthog",
		"stripe", "paypal",
		"sendgrid", "mailgun", "postmark",
		"aws", "gcp", "azure",
	}

	for _, serviceKey := range serviceOrder {
		serviceInfo := services[serviceKey]
		patterns := serviceInfo["patterns"].([]string)

		for _, pattern := range patterns {
			if allDeps[strings.ToLower(pattern)] {
				// Find the line in go.mod where this dependency appears
				lineNum := 0
				sourceText := ""
				for i, line := range lines {
					if strings.Contains(strings.ToLower(line), strings.ToLower(pattern)) {
						lineNum = i + 1
						sourceText = strings.TrimSpace(line)
						break
					}
				}

				results = append(results, &DetectionResult{
					Key:         serviceInfo["key"].(string),
					Value:       serviceInfo["url"].(string),
					Description: fmt.Sprintf("%s service detected in go.mod", serviceInfo["name"].(string)),
					Confidence:  0.9,
					DebugInfo:   fmt.Sprintf("Found Go module '%s' in go.mod", pattern),
					SourceFile:  "go.mod",
					SourceLine:  lineNum,
					SourceText:  sourceText,
				})
				break // Only add each service once
			}
		}
	}

	return results, nil
}