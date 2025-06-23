package detectors

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// DartDetector detects external services in Flutter/Dart projects
type DartDetector struct{}

func (d *DartDetector) Name() string {
	return "dart"
}

func (d *DartDetector) Description() string {
	return "Flutter/Dart pubspec.yaml external services detector"
}

func (d *DartDetector) ShouldRun() bool {
	// Check for pubspec.yaml file (standard Flutter/Dart configuration file)
	if _, err := os.Stat("pubspec.yaml"); err == nil {
		return true
	}
	return false
}

func (d *DartDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Read pubspec.yaml file
	data, err := ioutil.ReadFile("pubspec.yaml")
	if err != nil {
		return results, nil
	}

	content := strings.ToLower(string(data))

	// Map of Dart/Flutter package names to service info
	services := map[string]map[string]interface{}{
		// Monitoring and Error Tracking
		"sentry": {
			"patterns": []string{"sentry_flutter", "sentry_dart", "sentry"},
			"name": "Sentry",
			"url": "https://sentry.io",
			"key": "monitoring",
		},
		"bugsnag": {
			"patterns": []string{"bugsnag_flutter", "bugsnag_dart"},
			"name": "Bugsnag",
			"url": "https://bugsnag.com",
			"key": "monitoring",
		},
		"rollbar": {
			"patterns": []string{"rollbar_dart", "rollbar_flutter"},
			"name": "Rollbar",
			"url": "https://rollbar.com",
			"key": "monitoring",
		},
		"datadog": {
			"patterns": []string{"datadog_flutter_plugin", "datadog_sdk_flutter"},
			"name": "Datadog",
			"url": "https://datadog.com",
			"key": "monitoring",
		},
		"firebase_crashlytics": {
			"patterns": []string{"firebase_crashlytics", "crashlytics"},
			"name": "Firebase Crashlytics",
			"url": "https://console.firebase.google.com",
			"key": "monitoring",
		},

		// Analytics
		"firebase_analytics": {
			"patterns": []string{"firebase_analytics"},
			"name": "Firebase Analytics",
			"url": "https://console.firebase.google.com",
			"key": "analytics",
		},
		"mixpanel": {
			"patterns": []string{"mixpanel_flutter", "mixpanel_dart"},
			"name": "Mixpanel",
			"url": "https://mixpanel.com",
			"key": "analytics",
		},
		"amplitude": {
			"patterns": []string{"amplitude_flutter", "amplitude_dart"},
			"name": "Amplitude",
			"url": "https://amplitude.com",
			"key": "analytics",
		},
		"posthog": {
			"patterns": []string{"posthog_flutter", "posthog_dart"},
			"name": "PostHog",
			"url": "https://posthog.com",
			"key": "analytics",
		},
		"segment": {
			"patterns": []string{"segment_analytics", "analytics_flutter"},
			"name": "Segment",
			"url": "https://segment.com",
			"key": "analytics",
		},

		// Payments and In-App Purchases
		"stripe": {
			"patterns": []string{"stripe_payment", "flutter_stripe", "stripe_sdk"},
			"name": "Stripe",
			"url": "https://stripe.com",
			"key": "payments",
		},
		"paypal": {
			"patterns": []string{"paypal_flutter", "flutter_paypal", "paypal_sdk"},
			"name": "PayPal",
			"url": "https://paypal.com",
			"key": "payments",
		},
		"razorpay": {
			"patterns": []string{"razorpay_flutter", "flutter_razorpay"},
			"name": "Razorpay",
			"url": "https://razorpay.com",
			"key": "payments",
		},
		"square": {
			"patterns": []string{"square_in_app_payments", "flutter_square"},
			"name": "Square",
			"url": "https://squareup.com",
			"key": "payments",
		},

		// Cloud Services and Backend
		"firebase": {
			"patterns": []string{"firebase_core", "firebase_auth", "cloud_firestore", "firebase_storage"},
			"name": "Firebase",
			"url": "https://console.firebase.google.com",
			"key": "cloud",
		},
		"supabase": {
			"patterns": []string{"supabase_flutter", "supabase"},
			"name": "Supabase",
			"url": "https://supabase.com",
			"key": "cloud",
		},
		"aws": {
			"patterns": []string{"amplify_flutter", "aws_s3_api", "aws_lambda_dart_runtime"},
			"name": "AWS",
			"url": "https://console.aws.amazon.com",
			"key": "cloud",
		},
		"gcp": {
			"patterns": []string{"googleapis", "gcloud", "google_cloud"},
			"name": "Google Cloud",
			"url": "https://console.cloud.google.com",
			"key": "cloud",
		},

		// Push Notifications
		"onesignal": {
			"patterns": []string{"onesignal_flutter"},
			"name": "OneSignal",
			"url": "https://onesignal.com",
			"key": "notifications",
		},
		"pusher": {
			"patterns": []string{"pusher_channels_flutter", "pusher_flutter"},
			"name": "Pusher",
			"url": "https://pusher.com",
			"key": "notifications",
		},

		// Authentication
		"auth0": {
			"patterns": []string{"auth0_flutter", "flutter_auth0"},
			"name": "Auth0",
			"url": "https://auth0.com",
			"key": "auth",
		},

		// Storage and CDN
		"cloudinary": {
			"patterns": []string{"cloudinary_flutter", "flutter_cloudinary"},
			"name": "Cloudinary",
			"url": "https://cloudinary.com",
			"key": "storage",
		},
	}

	// Check each service and collect all found results
	serviceOrder := []string{
		"sentry", "bugsnag", "rollbar", "datadog", "firebase_crashlytics",
		"firebase_analytics", "mixpanel", "amplitude", "posthog", "segment",
		"stripe", "paypal", "razorpay", "square",
		"firebase", "supabase", "aws", "gcp",
		"onesignal", "pusher",
		"auth0",
		"cloudinary",
	}

	for _, serviceKey := range serviceOrder {
		serviceInfo := services[serviceKey]
		patterns := serviceInfo["patterns"].([]string)

		for _, pattern := range patterns {
			// Use regex to match dependencies in pubspec.yaml format
			// Looking for package names under dependencies: or dev_dependencies:
			// Format: "  package_name: ^version" or "  package_name:"
			escapedPattern := regexp.QuoteMeta(pattern)
			regex := regexp.MustCompile(`(?i)^\s*` + escapedPattern + `\s*:`)

			// Split content into lines and check each line
			lines := strings.Split(content, "\n")
			inDependenciesSection := false

			for _, line := range lines {
				trimmedLine := strings.TrimSpace(line)

				// Check if we're entering dependencies section
				if trimmedLine == "dependencies:" || trimmedLine == "dev_dependencies:" {
					inDependenciesSection = true
					continue
				}

				// Check if we're leaving dependencies section (new top-level key)
				if inDependenciesSection && len(trimmedLine) > 0 && !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") && strings.Contains(trimmedLine, ":") {
					inDependenciesSection = false
				}

				// If we're in dependencies section and find a match
				if inDependenciesSection && regex.MatchString(line) {
					results = append(results, &DetectionResult{
						Key:         serviceInfo["key"].(string),
						Value:       serviceInfo["url"].(string),
						Description: fmt.Sprintf("%s service detected in Flutter/Dart project", serviceInfo["name"].(string)),
						Confidence:  0.9,
					})
					break // Only add each service once
				}
			}
		}
	}

	return results, nil
}