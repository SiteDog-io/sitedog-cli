package detectors

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// IOSDetector detects iOS-specific configurations and services
type IOSDetector struct{}

func (i *IOSDetector) Name() string {
	return "ios"
}

func (i *IOSDetector) Description() string {
	return "iOS platform and services detector"
}

func (i *IOSDetector) ShouldRun() bool {
	// Check for iOS project structure
	iosPaths := []string{
		"ios/Runner/Info.plist",     // Flutter iOS
		"ios/Podfile",               // CocoaPods
		"Info.plist",                // Native iOS
		"Podfile",                   // CocoaPods in root
		"ios/Runner.xcodeproj",      // Flutter Xcode project
	}

	// Check direct paths first
	for _, path := range iosPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	// Check for Xcode projects with glob patterns
	xcodePatterns := []string{"*.xcodeproj", "*.xcworkspace", "ios/*.xcworkspace"}
	for _, pattern := range xcodePatterns {
		matches, _ := filepath.Glob(pattern)
		if len(matches) > 0 {
			return true
		}
	}

	return false
}

func (i *IOSDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult
	var allContent string

	// Read iOS configuration files
	filesToCheck := []string{
		"ios/Runner/Info.plist",
		"ios/Podfile",
		"Info.plist",
		"Podfile",
	}

	for _, file := range filesToCheck {
		if data, err := ioutil.ReadFile(file); err == nil {
			allContent += string(data) + "\n"
		}
	}

	// Also check for iOS-specific directories and files
	if _, err := os.Stat("ios"); err == nil {
		// Look for additional iOS files in ios directory
		err := filepath.Walk("ios", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if strings.HasSuffix(path, ".plist") ||
			   strings.HasSuffix(path, ".xcconfig") ||
			   strings.HasSuffix(path, "Podfile") {
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

	// Detect iOS-specific services and configurations
	services := map[string]map[string]interface{}{
		// App Store and Distribution (AUTO-ADD - high confidence)
		"app_store": {
			"patterns": []string{"app store", "appstoreconnect", "itunes connect", "bundle identifier", "cfbundleidentifier"},
			"name": "App Store Connect",
			"url": "https://appstoreconnect.apple.com",
			"key": "ios_distribution",
			"confidence": 0.95, // Auto-add for iOS platform
		},
		"apple_developer": {
			"patterns": []string{"cfbundleidentifier", "ios", "platform :ios"},
			"name": "Apple Developer",
			"url": "https://developer.apple.com",
			"key": "apple_developer",
			"confidence": 0.95, // Auto-add for iOS platform
		},
		"testflight": {
			"patterns": []string{"testflight", "beta testing", "external testing"},
			"name": "TestFlight",
			"url": "https://appstoreconnect.apple.com",
			"key": "ios_testing",
			"confidence": 0.90, // Conditional auto-add
		},

		// Push Notifications
		"apns": {
			"patterns": []string{"push notification", "apns", "apple push notification", "usernotifications"},
			"name": "Apple Push Notifications",
			"url": "https://developer.apple.com/notifications",
			"key": "notifications",
		},

		// Firebase iOS
		"firebase_ios": {
			"patterns": []string{"firebase/core", "firebase_core_ios", "googleservice-info.plist", "firebase_messaging"},
			"name": "Firebase iOS",
			"url": "https://console.firebase.google.com",
			"key": "cloud",
		},

		// Analytics
		"firebase_analytics_ios": {
			"patterns": []string{"firebase/analytics", "firebase_analytics_ios", "google analytics"},
			"name": "Firebase Analytics iOS",
			"url": "https://console.firebase.google.com",
			"key": "analytics",
		},
		"amplitude_ios": {
			"patterns": []string{"amplitude-ios", "amplitude_ios", "pod 'amplitude'"},
			"name": "Amplitude iOS",
			"url": "https://amplitude.com",
			"key": "analytics",
		},
		"mixpanel_ios": {
			"patterns": []string{"mixpanel-ios", "mixpanel_ios", "pod 'mixpanel'"},
			"name": "Mixpanel iOS",
			"url": "https://mixpanel.com",
			"key": "analytics",
		},

		// Error Tracking
		"sentry_ios": {
			"patterns": []string{"sentry-cocoa", "sentry_ios", "pod 'sentry'"},
			"name": "Sentry iOS",
			"url": "https://sentry.io",
			"key": "monitoring",
		},
		"crashlytics_ios": {
			"patterns": []string{"firebase/crashlytics", "crashlytics_ios", "fabric crashlytics"},
			"name": "Firebase Crashlytics iOS",
			"url": "https://console.firebase.google.com",
			"key": "monitoring",
		},

		// Social and Authentication
		"facebook_ios": {
			"patterns": []string{"fbsdkcorekit", "facebook-ios-sdk", "facebook login"},
			"name": "Facebook iOS SDK",
			"url": "https://developers.facebook.com",
			"key": "auth",
		},
		"google_signin_ios": {
			"patterns": []string{"googlesignin", "google-signin-ios", "google sign in"},
			"name": "Google Sign-In iOS",
			"url": "https://console.cloud.google.com",
			"key": "auth",
		},
		"twitter_ios": {
			"patterns": []string{"twitterkit", "twitter-ios-sdk"},
			"name": "Twitter iOS SDK",
			"url": "https://developer.twitter.com",
			"key": "auth",
		},

		// Payments
		"stripe_ios": {
			"patterns": []string{"stripe-ios", "stripe_ios", "pod 'stripe'"},
			"name": "Stripe iOS",
			"url": "https://stripe.com",
			"key": "payments",
		},
		"in_app_purchase": {
			"patterns": []string{"storekit", "in-app purchase", "iap", "subscription"},
			"name": "App Store In-App Purchases",
			"url": "https://appstoreconnect.apple.com",
			"key": "payments",
		},

		// Maps and Location
		"mapkit": {
			"patterns": []string{"mapkit", "core location", "corelocation"},
			"name": "Apple MapKit",
			"url": "https://developer.apple.com/maps",
			"key": "maps",
		},
		"google_maps_ios": {
			"patterns": []string{"googlemaps", "google-maps-ios", "google maps sdk"},
			"name": "Google Maps iOS",
			"url": "https://console.cloud.google.com",
			"key": "maps",
		},

		// Development Tools
		"xcode_cloud": {
			"patterns": []string{"xcode cloud", "ci_post_clone", "ci_pre_xcodebuild"},
			"name": "Xcode Cloud",
			"url": "https://developer.apple.com/xcode-cloud",
			"key": "ci",
		},
	}

		// Check each service and collect all found results
	serviceOrder := []string{
		"app_store", "apple_developer", "testflight",
		"apns",
		"firebase_ios", "firebase_analytics_ios", "amplitude_ios", "mixpanel_ios",
		"sentry_ios", "crashlytics_ios",
		"facebook_ios", "google_signin_ios", "twitter_ios",
		"stripe_ios", "in_app_purchase",
		"mapkit", "google_maps_ios",
		"xcode_cloud",
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
					Description: fmt.Sprintf("%s detected in iOS project", serviceInfo["name"].(string)),
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