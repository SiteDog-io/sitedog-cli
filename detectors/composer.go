package detectors

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// ComposerDetector detects external services in PHP composer.json
type ComposerDetector struct{}

func (c *ComposerDetector) Name() string {
	return "composer"
}

func (c *ComposerDetector) Description() string {
	return "PHP composer.json external services detector"
}

func (c *ComposerDetector) ShouldRun() bool {
	_, err := os.Stat("composer.json")
	return err == nil
}

func (c *ComposerDetector) Detect() ([]*DetectionResult, error) {
	data, err := ioutil.ReadFile("composer.json")
	if err != nil {
		return nil, err
	}

	var composerData map[string]interface{}
	if err := json.Unmarshal(data, &composerData); err != nil {
		return nil, err
	}

	// Get all dependencies (both regular and dev dependencies)
	allDeps := make(map[string]bool)

	if deps, ok := composerData["require"].(map[string]interface{}); ok {
		for pkg := range deps {
			allDeps[strings.ToLower(pkg)] = true
		}
	}

	if devDeps, ok := composerData["require-dev"].(map[string]interface{}); ok {
		for pkg := range devDeps {
			allDeps[strings.ToLower(pkg)] = true
		}
	}

	// Map of composer package names to service info
	services := map[string]map[string]interface{}{
		// Monitoring and Error Tracking
		"sentry": {
			"patterns": []string{"sentry/sentry", "sentry/sentry-laravel", "sentry/sentry-symfony"},
			"name": "Sentry",
			"url": "https://sentry.io",
			"key": "monitoring",
		},
		"bugsnag": {
			"patterns": []string{"bugsnag/bugsnag", "bugsnag/bugsnag-laravel", "bugsnag/bugsnag-symfony"},
			"name": "Bugsnag",
			"url": "https://bugsnag.com",
			"key": "monitoring",
		},
		"rollbar": {
			"patterns": []string{"rollbar/rollbar", "rollbar/rollbar-laravel"},
			"name": "Rollbar",
			"url": "https://rollbar.com",
			"key": "monitoring",
		},
		"airbrake": {
			"patterns": []string{"airbrake/phpbrake"},
			"name": "Airbrake",
			"url": "https://airbrake.io",
			"key": "monitoring",
		},
		"datadog": {
			"patterns": []string{"datadog/dd-trace", "datadog/php-datadogstatsd"},
			"name": "Datadog",
			"url": "https://datadog.com",
			"key": "monitoring",
		},
		"newrelic": {
			"patterns": []string{"newrelic/monolog-enricher"},
			"name": "New Relic",
			"url": "https://newrelic.com",
			"key": "monitoring",
		},

		// Analytics
		"mixpanel": {
			"patterns": []string{"mixpanel/mixpanel-php"},
			"name": "Mixpanel",
			"url": "https://mixpanel.com",
			"key": "analytics",
		},
		"amplitude": {
			"patterns": []string{"amplitude/analytics"},
			"name": "Amplitude",
			"url": "https://amplitude.com",
			"key": "analytics",
		},
		"segment": {
			"patterns": []string{"segmentio/analytics-php"},
			"name": "Segment",
			"url": "https://segment.com",
			"key": "analytics",
		},

		// Payments
		"stripe": {
			"patterns": []string{"stripe/stripe-php"},
			"name": "Stripe",
			"url": "https://stripe.com",
			"key": "payments",
		},
		"paypal": {
			"patterns": []string{"paypal/rest-api-sdk-php", "paypal/paypal-checkout-sdk"},
			"name": "PayPal",
			"url": "https://paypal.com",
			"key": "payments",
		},

		// Email Services
		"sendgrid": {
			"patterns": []string{"sendgrid/sendgrid", "sendgrid/php-http-client"},
			"name": "SendGrid",
			"url": "https://sendgrid.com",
			"key": "email_delivery",
		},
		"mailgun": {
			"patterns": []string{"mailgun/mailgun-php"},
			"name": "Mailgun",
			"url": "https://mailgun.com",
			"key": "email_delivery",
		},
		"postmark": {
			"patterns": []string{"wildbit/postmark-php"},
			"name": "Postmark",
			"url": "https://postmarkapp.com",
			"key": "email_delivery",
		},

		// Cloud Platforms
		"aws": {
			"patterns": []string{"aws/aws-sdk-php"},
			"name": "AWS",
			"url": "https://console.aws.amazon.com",
			"key": "cloud",
		},
		"gcp": {
			"patterns": []string{"google/cloud", "google/cloud-storage", "google/cloud-firestore"},
			"name": "Google Cloud",
			"url": "https://console.cloud.google.com",
			"key": "cloud",
		},
		"azure": {
			"patterns": []string{"microsoft/azure-storage", "microsoft/azure-storage-blob"},
			"name": "Azure",
			"url": "https://portal.azure.com",
			"key": "cloud",
		},
	}

	// Check each service and collect all found results
	var results []*DetectionResult
	serviceOrder := []string{
		"sentry", "bugsnag", "rollbar", "airbrake", "datadog", "newrelic",
		"mixpanel", "amplitude", "segment",
		"stripe", "paypal",
		"sendgrid", "mailgun", "postmark",
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
					Description: fmt.Sprintf("%s service detected in composer.json", serviceInfo["name"].(string)),
					Confidence:  0.9,
				})
				break // Only add each service once
			}
		}
	}

	return results, nil
}