package detectors

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// GemfileDetector detects monitoring services in Ruby Gemfile
type GemfileDetector struct{}

func (g *GemfileDetector) Name() string {
	return "gemfile"
}

func (g *GemfileDetector) Description() string {
	return "Ruby Gemfile external services detector"
}

func (g *GemfileDetector) ShouldRun() bool {
	_, err := os.Stat("Gemfile")
	return err == nil
}

func (g *GemfileDetector) Detect() ([]*DetectionResult, error) {
	data, err := ioutil.ReadFile("Gemfile")
	if err != nil {
		return nil, err
	}

	content := strings.ToLower(string(data))

	// Map of gem names to service info
	services := map[string]map[string]interface{}{
		// Monitoring and Error Tracking
		"sentry": {
			"patterns": []string{"sentry-ruby", "sentry-raven", "\"sentry\""},
			"name": "Sentry",
			"url": "https://sentry.io",
			"key": "monitoring",
		},
		"appsignal": {
			"patterns": []string{"appsignal"},
			"name": "AppSignal",
			"url": "https://appsignal.com",
			"key": "monitoring",
		},
		"airbrake": {
			"patterns": []string{"airbrake"},
			"name": "Airbrake",
			"url": "https://airbrake.io",
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
		"honeybadger": {
			"patterns": []string{"honeybadger"},
			"name": "Honeybadger",
			"url": "https://honeybadger.io",
			"key": "monitoring",
		},
		"datadog": {
			"patterns": []string{"ddtrace", "dogstatsd-ruby"},
			"name": "Datadog",
			"url": "https://datadog.com",
			"key": "monitoring",
		},
		"newrelic": {
			"patterns": []string{"newrelic_rpm"},
			"name": "New Relic",
			"url": "https://newrelic.com",
			"key": "monitoring",
		},

		// Analytics
		"mixpanel": {
			"patterns": []string{"mixpanel-ruby"},
			"name": "Mixpanel",
			"url": "https://mixpanel.com",
			"key": "analytics",
		},
		"amplitude": {
			"patterns": []string{"amplitude-api"},
			"name": "Amplitude",
			"url": "https://amplitude.com",
			"key": "analytics",
		},
		"segment": {
			"patterns": []string{"analytics-ruby"},
			"name": "Segment",
			"url": "https://segment.com",
			"key": "analytics",
		},

		// Payments
		"stripe": {
			"patterns": []string{"stripe"},
			"name": "Stripe",
			"url": "https://stripe.com",
			"key": "payments",
		},
		"paypal": {
			"patterns": []string{"paypal-sdk-rest", "paypal-checkout-sdk"},
			"name": "PayPal",
			"url": "https://paypal.com",
			"key": "payments",
		},

		// Email Services
		"sendgrid": {
			"patterns": []string{"sendgrid-ruby"},
			"name": "SendGrid",
			"url": "https://sendgrid.com",
			"key": "email_delivery",
		},
		"mailgun": {
			"patterns": []string{"mailgun-ruby"},
			"name": "Mailgun",
			"url": "https://mailgun.com",
			"key": "email_delivery",
		},
		"postmark": {
			"patterns": []string{"postmark-rails", "mail-postmark"},
			"name": "Postmark",
			"url": "https://postmarkapp.com",
			"key": "email_delivery",
		},

		// Cloud Platforms
		"aws": {
			"patterns": []string{"aws-sdk", "aws-sdk-ruby", "aws-sdk-core"},
			"name": "AWS",
			"url": "https://console.aws.amazon.com",
			"key": "cloud",
		},
		"gcp": {
			"patterns": []string{"google-cloud", "google-api-client"},
			"name": "Google Cloud",
			"url": "https://console.cloud.google.com",
			"key": "cloud",
		},
		"azure": {
			"patterns": []string{"azure", "azure-storage"},
			"name": "Azure",
			"url": "https://portal.azure.com",
			"key": "cloud",
		},
	}

	// Check each service and collect all found results
	var results []*DetectionResult
	serviceOrder := []string{
		"sentry", "appsignal", "airbrake", "bugsnag", "rollbar", "honeybadger", "datadog", "newrelic",
		"mixpanel", "amplitude", "segment",
		"stripe", "paypal",
		"sendgrid", "mailgun", "postmark",
		"aws", "gcp", "azure",
	}

	for _, serviceKey := range serviceOrder {
		serviceInfo := services[serviceKey]
		patterns := serviceInfo["patterns"].([]string)

		for _, pattern := range patterns {
			if strings.Contains(content, pattern) {
				results = append(results, &DetectionResult{
					Key:         serviceInfo["key"].(string),
					Value:       serviceInfo["url"].(string),
					Description: fmt.Sprintf("%s service detected in Gemfile", serviceInfo["name"].(string)),
					Confidence:  0.9,
				})
				break // Only add each service once
			}
		}
	}

	return results, nil
}