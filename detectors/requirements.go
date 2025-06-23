package detectors

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// RequirementsDetector detects external services in Python requirements.txt
type RequirementsDetector struct{}

func (r *RequirementsDetector) Name() string {
	return "requirements"
}

func (r *RequirementsDetector) Description() string {
	return "Python requirements.txt external services detector"
}

func (r *RequirementsDetector) ShouldRun() bool {
	_, err := os.Stat("requirements.txt")
	return err == nil
}

func (r *RequirementsDetector) Detect() ([]*DetectionResult, error) {
	data, err := ioutil.ReadFile("requirements.txt")
	if err != nil {
		return nil, err
	}

	content := strings.ToLower(string(data))

	// Map of package names to service info
	services := map[string]map[string]interface{}{
		// Monitoring and Error Tracking
		"sentry": {
			"patterns": []string{"sentry-sdk"},
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
		"airbrake": {
			"patterns": []string{"airbrake-pybrake"},
			"name": "Airbrake",
			"url": "https://airbrake.io",
			"key": "monitoring",
		},
		"datadog": {
			"patterns": []string{"datadog"},
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
			"patterns": []string{"amplitude-analytics"},
			"name": "Amplitude",
			"url": "https://amplitude.com",
			"key": "analytics",
		},
		"segment": {
			"patterns": []string{"segment-analytics-python"},
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
			"patterns": []string{"paypal-sdk", "paypalrestsdk"},
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
			"patterns": []string{"boto3", "botocore"},
			"name": "AWS",
			"url": "https://console.aws.amazon.com",
			"key": "cloud",
		},
		"gcp": {
			"patterns": []string{"google-cloud-", "google-api-python-client"},
			"name": "Google Cloud",
			"url": "https://console.cloud.google.com",
			"key": "cloud",
		},
		"azure": {
			"patterns": []string{"azure-", "azure-mgmt-"},
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
			if strings.Contains(content, pattern) {
				results = append(results, &DetectionResult{
					Key:         serviceInfo["key"].(string),
					Value:       serviceInfo["url"].(string),
					Description: fmt.Sprintf("%s service detected in requirements.txt", serviceInfo["name"].(string)),
					Confidence:  0.9,
				})
				break // Only add each service once
			}
		}
	}

	return results, nil
}