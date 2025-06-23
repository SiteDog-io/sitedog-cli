package detectors

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// CargoDetector detects external services in Rust Cargo.toml
type CargoDetector struct{}

func (c *CargoDetector) Name() string {
	return "cargo"
}

func (c *CargoDetector) Description() string {
	return "Rust Cargo.toml external services detector"
}

func (c *CargoDetector) ShouldRun() bool {
	_, err := os.Stat("Cargo.toml")
	return err == nil
}

func (c *CargoDetector) Detect() ([]*DetectionResult, error) {
	data, err := ioutil.ReadFile("Cargo.toml")
	if err != nil {
		return nil, err
	}

	content := strings.ToLower(string(data))

	// Map of Rust crate names to service info
	services := map[string]map[string]interface{}{
		// Monitoring and Error Tracking
		"sentry": {
			"patterns": []string{"sentry", "sentry-core", "sentry-actix", "sentry-tower"},
			"name": "Sentry",
			"url": "https://sentry.io",
			"key": "monitoring",
		},
		"bugsnag": {
			"patterns": []string{"bugsnag"},
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
		"datadog": {
			"patterns": []string{"ddtrace", "datadog", "dogstatsd"},
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
			"patterns": []string{"segment", "analytics-rust"},
			"name": "Segment",
			"url": "https://segment.com",
			"key": "analytics",
		},
		"posthog": {
			"patterns": []string{"posthog"},
			"name": "PostHog",
			"url": "https://posthog.com",
			"key": "analytics",
		},

		// Payments
		"stripe": {
			"patterns": []string{"stripe-rust", "async-stripe"},
			"name": "Stripe",
			"url": "https://stripe.com",
			"key": "payments",
		},
		"paypal": {
			"patterns": []string{"paypal-rs"},
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
			"patterns": []string{"mailgun"},
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

		// Cloud Platforms
		"aws": {
			"patterns": []string{"aws-sdk-rust", "rusoto", "aws-config", "aws-types"},
			"name": "AWS",
			"url": "https://console.aws.amazon.com",
			"key": "cloud",
		},
		"gcp": {
			"patterns": []string{"google-cloud", "gcp", "tonic-gcp"},
			"name": "Google Cloud",
			"url": "https://console.cloud.google.com",
			"key": "cloud",
		},
		"azure": {
			"patterns": []string{"azure_core", "azure_storage", "azure-sdk"},
			"name": "Azure",
			"url": "https://portal.azure.com",
			"key": "cloud",
		},
	}

	// Check each service and collect all found results
	var results []*DetectionResult
	serviceOrder := []string{
		"sentry", "bugsnag", "rollbar", "datadog", "newrelic",
		"mixpanel", "amplitude", "segment", "posthog",
		"stripe", "paypal",
		"sendgrid", "mailgun", "postmark",
		"aws", "gcp", "azure",
	}

	for _, serviceKey := range serviceOrder {
		serviceInfo := services[serviceKey]
		patterns := serviceInfo["patterns"].([]string)

		for _, pattern := range patterns {
			// Use regex to match crate names in dependencies sections
			// Matches patterns like: crate_name = "version" or crate_name = { version = "..." }
			escapedPattern := regexp.QuoteMeta(pattern)
			regex := regexp.MustCompile(`(?m)^\s*` + escapedPattern + `\s*=`)

			if regex.MatchString(content) {
				results = append(results, &DetectionResult{
					Key:         serviceInfo["key"].(string),
					Value:       serviceInfo["url"].(string),
					Description: fmt.Sprintf("%s service detected in Cargo.toml", serviceInfo["name"].(string)),
					Confidence:  0.9,
				})
				break // Only add each service once
			}
		}
	}

	return results, nil
}