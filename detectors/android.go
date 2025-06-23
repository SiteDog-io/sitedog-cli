package detectors

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// AndroidDetector detects Android-specific configurations and services
type AndroidDetector struct{}

func (a *AndroidDetector) Name() string {
	return "android"
}

func (a *AndroidDetector) Description() string {
	return "Android platform and services detector"
}

func (a *AndroidDetector) ShouldRun() bool {
	// Check for Android project structure
	androidPaths := []string{
		"android/app/build.gradle",                    // Flutter Android
		"android/app/src/main/AndroidManifest.xml",   // Android manifest
		"app/build.gradle",                            // Native Android
		"build.gradle",                                // Android root build
		"AndroidManifest.xml",                         // Native manifest
		"gradle.properties",                           // Gradle properties
		"settings.gradle",                             // Gradle settings
	}

	for _, path := range androidPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

func (a *AndroidDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult
	var allContent string

	// Read Android configuration files
	filesToCheck := []string{
		"android/app/build.gradle",
		"android/app/src/main/AndroidManifest.xml",
		"android/build.gradle",
		"android/gradle.properties",
		"app/build.gradle",
		"build.gradle",
		"AndroidManifest.xml",
		"gradle.properties",
		"settings.gradle",
	}

	for _, file := range filesToCheck {
		if data, err := ioutil.ReadFile(file); err == nil {
			allContent += string(data) + "\n"
		}
	}

	// Also check for Android-specific directories and files
	if _, err := os.Stat("android"); err == nil {
		// Look for additional Android files in android directory
		err := filepath.Walk("android", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if strings.HasSuffix(path, ".gradle") ||
			   strings.HasSuffix(path, ".xml") ||
			   strings.HasSuffix(path, ".properties") ||
			   strings.HasSuffix(path, ".json") {
				if data, readErr := ioutil.ReadFile(path); readErr == nil {
					allContent += string(data) + "\n"
				}
			}
			return nil
		})
		if err != nil {
			// Continue even if walk fails
		}
	}

	if allContent == "" {
		return results, nil
	}

	content := strings.ToLower(allContent)

	// Detect Android-specific services and configurations
	services := map[string]map[string]interface{}{
		// Google Play Store and Distribution (AUTO-ADD - high confidence)
		"play_store": {
			"patterns": []string{"google play", "play console", "play store", "application id", "versioncode", "versionname"},
			"name": "Google Play Console",
			"url": "https://play.google.com/console",
			"key": "android_distribution",
			"confidence": 0.95, // Auto-add for Android platform
		},
		"google_developer": {
			"patterns": []string{"applicationid", "android", "compilesdk"},
			"name": "Google Play Console",
			"url": "https://play.google.com/console",
			"key": "google_play_console",
			"confidence": 0.95, // Auto-add for Android platform
		},
		"play_app_signing": {
			"patterns": []string{"play app signing", "signing config", "keystore", "key alias"},
			"name": "Play App Signing",
			"url": "https://play.google.com/console",
			"key": "android_signing",
			"confidence": 0.90, // Conditional auto-add
		},

		// Firebase Android
		"firebase_android": {
			"patterns": []string{"firebase-bom", "firebase-analytics", "firebase-messaging", "google-services.json", "google-services"},
			"name": "Firebase Android",
			"url": "https://console.firebase.google.com",
			"key": "cloud",
		},

		// Push Notifications
		"fcm": {
			"patterns": []string{"firebase-messaging", "fcm", "firebase cloud messaging", "push notification"},
			"name": "Firebase Cloud Messaging",
			"url": "https://console.firebase.google.com",
			"key": "notifications",
		},
		"onesignal_android": {
			"patterns": []string{"onesignal-android", "onesignal", "com.onesignal"},
			"name": "OneSignal Android",
			"url": "https://onesignal.com",
			"key": "notifications",
		},

		// Analytics
		"firebase_analytics_android": {
			"patterns": []string{"firebase-analytics", "google-analytics", "firebase analytics"},
			"name": "Firebase Analytics Android",
			"url": "https://console.firebase.google.com",
			"key": "analytics",
		},
		"amplitude_android": {
			"patterns": []string{"amplitude-android", "com.amplitude"},
			"name": "Amplitude Android",
			"url": "https://amplitude.com",
			"key": "analytics",
		},
		"mixpanel_android": {
			"patterns": []string{"mixpanel-android", "com.mixpanel.android"},
			"name": "Mixpanel Android",
			"url": "https://mixpanel.com",
			"key": "analytics",
		},
		"google_analytics_android": {
			"patterns": []string{"google-analytics", "analytics-android", "com.google.android.gms:play-services-analytics"},
			"name": "Google Analytics Android",
			"url": "https://analytics.google.com",
			"key": "analytics",
		},

		// Error Tracking and Monitoring
		"crashlytics_android": {
			"patterns": []string{"firebase-crashlytics", "crashlytics", "fabric crashlytics"},
			"name": "Firebase Crashlytics Android",
			"url": "https://console.firebase.google.com",
			"key": "monitoring",
		},
		"sentry_android": {
			"patterns": []string{"sentry-android", "io.sentry:sentry-android"},
			"name": "Sentry Android",
			"url": "https://sentry.io",
			"key": "monitoring",
		},
		"bugsnag_android": {
			"patterns": []string{"bugsnag-android", "com.bugsnag:bugsnag-android"},
			"name": "Bugsnag Android",
			"url": "https://bugsnag.com",
			"key": "monitoring",
		},

		// Social and Authentication
		"google_signin_android": {
			"patterns": []string{"play-services-auth", "google sign in", "google-signin", "gms.auth"},
			"name": "Google Sign-In Android",
			"url": "https://console.cloud.google.com",
			"key": "auth",
		},
		"facebook_android": {
			"patterns": []string{"facebook-android-sdk", "facebook-login", "com.facebook.android"},
			"name": "Facebook Android SDK",
			"url": "https://developers.facebook.com",
			"key": "auth",
		},
		"twitter_android": {
			"patterns": []string{"twitter-android-sdk", "twitter-kit", "twitter4j"},
			"name": "Twitter Android SDK",
			"url": "https://developer.twitter.com",
			"key": "auth",
		},

		// Payments and In-App Billing
		"google_play_billing": {
			"patterns": []string{"play-services-billing", "billing", "in-app billing", "iab", "sku"},
			"name": "Google Play Billing",
			"url": "https://play.google.com/console",
			"key": "payments",
		},
		"stripe_android": {
			"patterns": []string{"stripe-android", "com.stripe:stripe-android"},
			"name": "Stripe Android",
			"url": "https://stripe.com",
			"key": "payments",
		},
		"paypal_android": {
			"patterns": []string{"paypal-android", "paypal-sdk", "com.paypal"},
			"name": "PayPal Android",
			"url": "https://paypal.com",
			"key": "payments",
		},

		// Maps and Location
		"google_maps_android": {
			"patterns": []string{"play-services-maps", "google-maps", "maps-android", "com.google.android.gms:play-services-maps"},
			"name": "Google Maps Android",
			"url": "https://console.cloud.google.com",
			"key": "maps",
		},
		"google_location_android": {
			"patterns": []string{"play-services-location", "location", "gps", "com.google.android.gms:play-services-location"},
			"name": "Google Location Services",
			"url": "https://console.cloud.google.com",
			"key": "location",
		},

		// Cloud Services
		"google_cloud_android": {
			"patterns": []string{"google-cloud-storage", "cloud-storage", "gcs"},
			"name": "Google Cloud Storage Android",
			"url": "https://console.cloud.google.com",
			"key": "cloud",
		},
		"aws_android": {
			"patterns": []string{"aws-android-sdk", "amazonaws", "aws-sdk-android"},
			"name": "AWS Android SDK",
			"url": "https://console.aws.amazon.com",
			"key": "cloud",
		},

		// Development and CI/CD
		"gradle_play_publisher": {
			"patterns": []string{"gradle-play-publisher", "play-publisher", "com.github.triplet.play"},
			"name": "Gradle Play Publisher",
			"url": "https://play.google.com/console",
			"key": "ci",
		},
	}

		// Check each service and collect all found results
	serviceOrder := []string{
		"play_store", "google_developer", "play_app_signing",
		"firebase_android", "fcm", "onesignal_android",
		"firebase_analytics_android", "amplitude_android", "mixpanel_android", "google_analytics_android",
		"crashlytics_android", "sentry_android", "bugsnag_android",
		"google_signin_android", "facebook_android", "twitter_android",
		"google_play_billing", "stripe_android", "paypal_android",
		"google_maps_android", "google_location_android",
		"google_cloud_android", "aws_android",
		"gradle_play_publisher",
	}

	for _, serviceKey := range serviceOrder {
		serviceInfo := services[serviceKey]
		patterns := serviceInfo["patterns"].([]string)

		found := false
		for _, pattern := range patterns {
			// Use case-insensitive regex matching
			escapedPattern := regexp.QuoteMeta(pattern)
			regex := regexp.MustCompile(`(?i)` + escapedPattern)

			if regex.MatchString(content) {
				// Use custom confidence if specified, otherwise default to 0.85
				confidence := 0.85
				if customConf, ok := serviceInfo["confidence"].(float64); ok {
					confidence = customConf
				}

				results = append(results, &DetectionResult{
					Key:         serviceInfo["key"].(string),
					Value:       serviceInfo["url"].(string),
					Description: fmt.Sprintf("%s detected in Android project", serviceInfo["name"].(string)),
					Confidence:  confidence,
				})
				found = true
				break
			}
		}

		if found {
			continue // Only add each service once
		}
	}

	return results, nil
}