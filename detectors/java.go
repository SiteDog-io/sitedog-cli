package detectors

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// JavaDetector detects external services in Java projects (Maven & Gradle)
type JavaDetector struct{}

func (j *JavaDetector) Name() string {
	return "java"
}

func (j *JavaDetector) Description() string {
	return "Java Maven/Gradle external services detector"
}

func (j *JavaDetector) ShouldRun() bool {
	// Check for pom.xml (Maven) or build.gradle* (Gradle) files
	configPaths := []string{"pom.xml", "build.gradle", "build.gradle.kts"}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

func (j *JavaDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult
	var allContent string

	// Read Maven pom.xml if exists
	if data, err := ioutil.ReadFile("pom.xml"); err == nil {
		allContent += string(data) + "\n"
	}

	// Read Gradle build files if exist
	gradleFiles := []string{"build.gradle", "build.gradle.kts"}
	for _, file := range gradleFiles {
		if data, err := ioutil.ReadFile(file); err == nil {
			allContent += string(data) + "\n"
		}
	}

	// Also check for multi-module projects
	if files, err := filepath.Glob("*/pom.xml"); err == nil {
		for _, file := range files {
			if data, err := ioutil.ReadFile(file); err == nil {
				allContent += string(data) + "\n"
			}
		}
	}

	if allContent == "" {
		return results, nil
	}

	content := strings.ToLower(allContent)

	// Map of Java library names to service info
	services := map[string]map[string]interface{}{
		// Monitoring and Error Tracking
		"sentry": {
			"patterns": []string{"io.sentry", "sentry-java", "sentry-spring", "sentry-spring-boot-starter"},
			"name": "Sentry",
			"url": "https://sentry.io",
			"key": "monitoring",
		},
		"bugsnag": {
			"patterns": []string{"com.bugsnag", "bugsnag-java", "bugsnag-spring"},
			"name": "Bugsnag",
			"url": "https://bugsnag.com",
			"key": "monitoring",
		},
		"rollbar": {
			"patterns": []string{"com.rollbar", "rollbar-java", "rollbar-spring-boot-webmvc"},
			"name": "Rollbar",
			"url": "https://rollbar.com",
			"key": "monitoring",
		},
		"datadog": {
			"patterns": []string{"com.datadoghq", "dd-trace-java", "datadog-slf4j"},
			"name": "Datadog",
			"url": "https://datadog.com",
			"key": "monitoring",
		},
		"newrelic": {
			"patterns": []string{"com.newrelic", "newrelic-agent", "newrelic-api"},
			"name": "New Relic",
			"url": "https://newrelic.com",
			"key": "monitoring",
		},
		"micrometer": {
			"patterns": []string{"io.micrometer", "micrometer-core", "micrometer-registry"},
			"name": "Micrometer",
			"url": "https://micrometer.io",
			"key": "monitoring",
		},

		// Analytics
		"mixpanel": {
			"patterns": []string{"com.mixpanel", "mixpanel-java"},
			"name": "Mixpanel",
			"url": "https://mixpanel.com",
			"key": "analytics",
		},
		"amplitude": {
			"patterns": []string{"com.amplitude", "amplitude-java"},
			"name": "Amplitude",
			"url": "https://amplitude.com",
			"key": "analytics",
		},
		"segment": {
			"patterns": []string{"com.segment.analytics.java", "analytics-java"},
			"name": "Segment",
			"url": "https://segment.com",
			"key": "analytics",
		},

		// Payments
		"stripe": {
			"patterns": []string{"com.stripe", "stripe-java"},
			"name": "Stripe",
			"url": "https://stripe.com",
			"key": "payments",
		},
		"paypal": {
			"patterns": []string{"com.paypal", "paypal-core", "checkout-sdk"},
			"name": "PayPal",
			"url": "https://paypal.com",
			"key": "payments",
		},

		// Email Services
		"sendgrid": {
			"patterns": []string{"com.sendgrid", "sendgrid-java"},
			"name": "SendGrid",
			"url": "https://sendgrid.com",
			"key": "email_delivery",
		},
		"mailgun": {
			"patterns": []string{"net.sargue.mailgun", "mailgun-java"},
			"name": "Mailgun",
			"url": "https://mailgun.com",
			"key": "email_delivery",
		},
		"postmark": {
			"patterns": []string{"com.wildbit.java", "postmark-java"},
			"name": "Postmark",
			"url": "https://postmarkapp.com",
			"key": "email_delivery",
		},

		// Cloud Platforms
		"aws": {
			"patterns": []string{"software.amazon.awssdk", "com.amazonaws", "aws-java-sdk"},
			"name": "AWS",
			"url": "https://console.aws.amazon.com",
			"key": "cloud",
		},
		"gcp": {
			"patterns": []string{"com.google.cloud", "google-cloud-java", "google-api-services"},
			"name": "Google Cloud",
			"url": "https://console.cloud.google.com",
			"key": "cloud",
		},
		"azure": {
			"patterns": []string{"com.azure", "com.microsoft.azure", "azure-sdk-bom"},
			"name": "Azure",
			"url": "https://portal.azure.com",
			"key": "cloud",
		},
	}

	// Check each service and collect all found results
	serviceOrder := []string{
		"sentry", "bugsnag", "rollbar", "datadog", "newrelic", "micrometer",
		"mixpanel", "amplitude", "segment",
		"stripe", "paypal",
		"sendgrid", "mailgun", "postmark",
		"aws", "gcp", "azure",
	}

	for _, serviceKey := range serviceOrder {
		serviceInfo := services[serviceKey]
		patterns := serviceInfo["patterns"].([]string)

		for _, pattern := range patterns {
			// Use regex to match dependencies in both Maven and Gradle formats
			// Maven: <groupId>pattern</groupId> or <artifactId>pattern</artifactId>
			// Gradle: implementation 'pattern:artifact:version' or compile 'pattern'
			escapedPattern := regexp.QuoteMeta(pattern)
			regex := regexp.MustCompile(`(?i)(groupid>` + escapedPattern + `<|artifactid>` + escapedPattern + `|implementation.*['"]` + escapedPattern + `|compile.*['"]` + escapedPattern + `)`)

			if regex.MatchString(content) {
				results = append(results, &DetectionResult{
					Key:         serviceInfo["key"].(string),
					Value:       serviceInfo["url"].(string),
					Description: fmt.Sprintf("%s service detected in Java project", serviceInfo["name"].(string)),
					Confidence:  0.9,
				})
				break // Only add each service once
			}
		}
	}

	return results, nil
}