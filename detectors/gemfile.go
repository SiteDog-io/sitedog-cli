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
	return "Ruby Gemfile monitoring services detector"
}

func (g *GemfileDetector) ShouldRun() bool {
	_, err := os.Stat("Gemfile")
	return err == nil
}

func (g *GemfileDetector) Detect() (*DetectionResult, error) {
	data, err := ioutil.ReadFile("Gemfile")
	if err != nil {
		return nil, err
	}

	content := strings.ToLower(string(data))

	// Map of gem names to service info
	services := map[string]map[string]interface{}{
		"sentry": {
			"patterns": []string{"sentry-ruby", "sentry-raven", "\"sentry\""},
			"name": "Sentry",
			"url": "https://sentry.io",
		},
		"appsignal": {
			"patterns": []string{"appsignal"},
			"name": "AppSignal",
			"url": "https://appsignal.com",
		},
		"airbrake": {
			"patterns": []string{"airbrake"},
			"name": "Airbrake",
			"url": "https://airbrake.io",
		},
		"bugsnag": {
			"patterns": []string{"bugsnag"},
			"name": "Bugsnag",
			"url": "https://bugsnag.com",
		},
		"rollbar": {
			"patterns": []string{"rollbar"},
			"name": "Rollbar",
			"url": "https://rollbar.com",
		},
		"honeybadger": {
			"patterns": []string{"honeybadger"},
			"name": "Honeybadger",
			"url": "https://honeybadger.io",
		},
	}

	// Check each service
	for _, serviceInfo := range services {
		patterns := serviceInfo["patterns"].([]string)
		for _, pattern := range patterns {
			if strings.Contains(content, pattern) {
				return &DetectionResult{
					Key:         "monitoring",
					Value:       serviceInfo["url"].(string),
					Description: fmt.Sprintf("%s monitoring service detected", serviceInfo["name"].(string)),
					Confidence:  0.9,
				}, nil
			}
		}
	}

	return nil, nil
}