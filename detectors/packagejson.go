package detectors

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// PackageJSONDetector detects external services in Node.js package.json
type PackageJSONDetector struct{}

func (p *PackageJSONDetector) Name() string {
	return "package-json"
}

func (p *PackageJSONDetector) Description() string {
	return "Node.js package.json external services detector"
}

func (p *PackageJSONDetector) ShouldRun() bool {
	_, err := os.Stat("package.json")
	return err == nil
}

func (p *PackageJSONDetector) Detect() ([]*DetectionResult, error) {
	data, err := ioutil.ReadFile("package.json")
	if err != nil {
		return nil, err
	}

	var packageData map[string]interface{}
	if err := json.Unmarshal(data, &packageData); err != nil {
		return nil, err
	}

	// Get all dependencies (both regular and dev dependencies)
	allDeps := make(map[string]bool)

	if deps, ok := packageData["dependencies"].(map[string]interface{}); ok {
		for pkg := range deps {
			allDeps[strings.ToLower(pkg)] = true
		}
	}

	if devDeps, ok := packageData["devDependencies"].(map[string]interface{}); ok {
		for pkg := range devDeps {
			allDeps[strings.ToLower(pkg)] = true
		}
	}

	// Map of npm package names to service info
	services := map[string]map[string]interface{}{
		// Monitoring and Error Tracking
		"sentry": {
			"patterns": []string{"@sentry/node", "@sentry/browser", "@sentry/react", "@sentry/nextjs"},
			"name": "Sentry",
			"url": "https://sentry.io",
			"key": "monitoring",
		},
		"bugsnag": {
			"patterns": []string{"@bugsnag/js", "@bugsnag/node", "bugsnag-js"},
			"name": "Bugsnag",
			"url": "https://bugsnag.com",
			"key": "monitoring",
		},
		"rollbar": {
			"patterns": []string{"rollbar"},
			"name": "Rollbar",
			"url": "https://rollbar.com",
			"key": "monitoring",
		},
		"airbrake": {
			"patterns": []string{"@airbrake/browser", "@airbrake/node"},
			"name": "Airbrake",
			"url": "https://airbrake.io",
			"key": "monitoring",
		},
		"datadog": {
			"patterns": []string{"@datadog/dd-trace", "dd-trace"},
			"name": "Datadog",
			"url": "https://datadog.com",
			"key": "monitoring",
		},
		"newrelic": {
			"patterns": []string{"newrelic"},
			"name": "New Relic",
			"url": "https://newrelic.com",
			"key": "monitoring",
		},

		// Analytics
		"mixpanel": {
			"patterns": []string{"mixpanel", "mixpanel-browser"},
			"name": "Mixpanel",
			"url": "https://mixpanel.com",
			"key": "analytics",
		},
		"amplitude": {
			"patterns": []string{"amplitude-js", "@amplitude/analytics-browser", "@amplitude/analytics-node"},
			"name": "Amplitude",
			"url": "https://amplitude.com",
			"key": "analytics",
		},
		"segment": {
			"patterns": []string{"analytics-node", "@segment/analytics-node"},
			"name": "Segment",
			"url": "https://segment.com",
			"key": "analytics",
		},
		"posthog": {
			"patterns": []string{"posthog-js", "posthog-node"},
			"name": "PostHog",
			"url": "https://posthog.com",
			"key": "analytics",
		},

		// Payments
		"stripe": {
			"patterns": []string{"stripe", "@stripe/stripe-js"},
			"name": "Stripe",
			"url": "https://stripe.com",
			"key": "payments",
		},
		"paypal": {
			"patterns": []string{"paypal-rest-sdk", "@paypal/checkout-server-sdk"},
			"name": "PayPal",
			"url": "https://paypal.com",
			"key": "payments",
		},

		// Email Services
		"sendgrid": {
			"patterns": []string{"@sendgrid/mail", "@sendgrid/client"},
			"name": "SendGrid",
			"url": "https://sendgrid.com",
			"key": "email_delivery",
		},
		"mailgun": {
			"patterns": []string{"mailgun-js", "mailgun.js"},
			"name": "Mailgun",
			"url": "https://mailgun.com",
			"key": "email_delivery",
		},
		"postmark": {
			"patterns": []string{"postmark"},
			"name": "Postmark",
			"url": "https://postmarkapp.com",
			"key": "email_delivery",
		},
		"resend": {
			"patterns": []string{"resend"},
			"name": "Resend",
			"url": "https://resend.com",
			"key": "email_delivery",
		},

		// Cloud Platforms
		"aws": {
			"patterns": []string{"aws-sdk", "@aws-sdk/client-s3", "@aws-sdk/lib-dynamodb"},
			"name": "AWS",
			"url": "https://console.aws.amazon.com",
			"key": "cloud",
		},
		"gcp": {
			"patterns": []string{"@google-cloud/storage", "@google-cloud/firestore", "googleapis"},
			"name": "Google Cloud",
			"url": "https://console.cloud.google.com",
			"key": "cloud",
		},
		"azure": {
			"patterns": []string{"@azure/storage-blob", "@azure/cosmos", "azure-storage"},
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
		"sendgrid", "mailgun", "postmark", "resend",
		"aws", "gcp", "azure",
	}

	for _, serviceKey := range serviceOrder {
		serviceInfo := services[serviceKey]
		patterns := serviceInfo["patterns"].([]string)

		for _, pattern := range patterns {
			if allDeps[strings.ToLower(pattern)] {
				results = append(results, &DetectionResult{
					Key:         serviceInfo["key"].(string),
					Value:       serviceInfo["url"].(string),
					Description: fmt.Sprintf("%s service detected in package.json", serviceInfo["name"].(string)),
					Confidence:  0.9,
				})
				break // Only add each service once
			}
		}
	}

	return results, nil
}