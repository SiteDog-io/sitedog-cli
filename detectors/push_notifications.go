package detectors

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// PushNotificationsDetector detects push notification services used in projects
type PushNotificationsDetector struct{}

func (p *PushNotificationsDetector) Name() string {
	return "push-notifications"
}

func (p *PushNotificationsDetector) Description() string {
	return "Push notifications and messaging services detector"
}

func (p *PushNotificationsDetector) ShouldRun() bool {
	// Check for common files that might contain push notification integrations
	files := []string{
		"package.json",
		"requirements.txt",
		"go.mod",
		"Cargo.toml",
		"composer.json",
		"Gemfile",
		".env",
		".env.example",
		".env.local",
	}

	for _, file := range files {
		if data, err := ioutil.ReadFile(file); err == nil {
			content := strings.ToLower(string(data))
			// Quick check for push notification-related keywords
			pushKeywords := []string{
				"onesignal", "pusher", "firebase", "fcm", "apns", "push notification",
				"expo push", "sendbird", "pushy", "airship", "pushwoosh",
			}
			for _, keyword := range pushKeywords {
				if strings.Contains(content, keyword) {
					return true
				}
			}
		}
	}

	// Check for source code directories
	srcDirs := []string{"src", "lib", "app", "components", "pages", "api", "services"}
	for _, dir := range srcDirs {
		if _, err := os.Stat(dir); err == nil {
			return true
		}
	}

	return false
}

func (p *PushNotificationsDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Files to check for push notification service references
	files := []string{
		"package.json",
		"requirements.txt",
		"go.mod",
		"Cargo.toml",
		"composer.json",
		"Gemfile",
		"pom.xml",
		"build.gradle",
		".env",
		".env.example",
		".env.local",
		".env.production",
		"config.js",
		"config.json",
		"config.yaml",
		"README.md",
		"google-services.json",
		"GoogleService-Info.plist",
	}

	// Exclude lock files that often have false positives
	excludeFiles := []string{
		"package-lock.json",
		"yarn.lock",
		"poetry.lock",
		"Pipfile.lock",
		"Cargo.lock",
		"composer.lock",
	}

	// Store file contents with metadata
	type FileContent struct {
		Path    string
		Content string
		Lines   []string
	}

	var fileContents []FileContent

	// Helper function to check if file should be excluded
	isExcluded := func(filename string) bool {
		for _, excluded := range excludeFiles {
			if filename == excluded {
				return true
			}
		}
		return false
	}

	// Read individual files
	for _, file := range files {
		if !isExcluded(file) {
			if data, err := ioutil.ReadFile(file); err == nil {
				content := string(data)
				lines := strings.Split(content, "\n")
				fileContents = append(fileContents, FileContent{
					Path:    file,
					Content: strings.ToLower(content),
					Lines:   lines,
				})
			}
		}
	}

	// Also check source code directories for push notification imports/usage
	srcDirs := []string{"src", "lib", "app", "components", "pages", "api", "services", "utils", "notifications", "push"}
	for _, dir := range srcDirs {
		if _, err := os.Stat(dir); err == nil {
			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err == nil && !info.IsDir() {
					// Check common source file extensions
					ext := strings.ToLower(filepath.Ext(info.Name()))
					if ext == ".js" || ext == ".ts" || ext == ".jsx" || ext == ".tsx" ||
					   ext == ".py" || ext == ".go" || ext == ".php" || ext == ".rb" ||
					   ext == ".java" || ext == ".cs" || ext == ".rs" || ext == ".swift" ||
					   ext == ".kt" || ext == ".dart" {
						// Skip lock files even in subdirectories
						if !isExcluded(info.Name()) {
							if data, readErr := ioutil.ReadFile(path); readErr == nil {
								content := string(data)
								lines := strings.Split(content, "\n")
								fileContents = append(fileContents, FileContent{
									Path:    path,
									Content: strings.ToLower(content),
									Lines:   lines,
								})
							}
						}
					}
				}
				return nil
			})
		}
	}

	// Define push notification services with their patterns and dashboards
	services := map[string]map[string]interface{}{
		"onesignal": {
			"patterns": []string{
				"onesignal", "onesignal_app_id", "onesignal_api_key", "one signal",
				"react-onesignal", "onesignal-node", "onesignal-python-api",
				"api.onesignal.com", "onesignal.com",
				"OneSignal.init", "OneSignal.sendNotification", "from onesignal_sdk import",
				"onesignal/onesignal-go-api", "OneSignalClient",
			},
			"name": "OneSignal",
			"url":  "https://onesignal.com/billing",
			"key":  "push_service",
		},

		"firebase_messaging": {
			"patterns": []string{
				"firebase messaging", "firebase-messaging", "firebase/messaging",
				"@firebase/messaging", "firebase-admin", "google-services.json",
				"firebase.messaging", "getMessaging(", "onMessage(", "getToken(",
				"from firebase_admin import messaging", "firebase_admin.messaging",
				"FirebaseMessaging", "FCMNotification", "FirebaseCloudMessaging",
				"FIREBASE_SERVER_KEY", "FCM_SERVER_KEY", "firebase_server_key",
			},
			"name": "Firebase Cloud Messaging",
			"url":  "https://console.firebase.google.com/project/_/settings/cloudmessaging",
			"key":  "push_service",
		},

		"pusher": {
			"patterns": []string{
				"pusher_app_id", "pusher_key", "pusher_secret", "pusher_cluster",
				"pusher-js", "pusher-http-python", "pusher-http-node",
				"api.pusherapp.com", "pusher.com",
				"new Pusher(", "Pusher.trigger", "from pusher import Pusher",
				"pusher/pusher-http-go", "github.com/pusher/pusher-http-go",
				"PUSHER_APP_ID", "PUSHER_KEY", "PUSHER_SECRET", "PUSHER_CLUSTER",
			},
			"name": "Pusher Channels",
			"url":  "https://dashboard.pusher.com/accounts/sign_in",
			"key":  "push_service",
		},

		"expo_push": {
			"patterns": []string{
				"expo push", "expo-notifications", "expo-server-sdk", "expo push token",
				"@expo/server-sdk", "expo-server-sdk-python", "expo-server-sdk-node",
				"exp.host/--/api/v2/push", "expo.io",
				"Expo.sendPushNotificationsAsync", "expo.sendPushNotificationsAsync",
				"from exponent_server_sdk import", "ExponentPushClient",
			},
			"name": "Expo Push Notifications",
			"url":  "https://expo.dev/accounts/[account]/settings/api-tokens",
			"key":  "push_service",
		},

		"sendbird": {
			"patterns": []string{
				"sendbird", "sendbird_app_id", "sendbird_api_token", "sendbird_bot_token",
				"sendbird-uikit", "sendbird-chat", "@sendbird/chat",
				"api.sendbird.com", "sendbird.com",
				"SendBird.init", "new SendBird", "SendbirdChat.init",
				"sendbird-platform-sdk-python", "sendbird-platform-sdk-javascript",
			},
			"name": "SendBird",
			"url":  "https://dashboard.sendbird.com/settings/application",
			"key":  "push_service",
		},

		"pushy": {
			"patterns": []string{
				"pushy", "pushy_api_key", "pushy_secret_api_key",
				"pushy-react-native", "pushy-cordova", "pushy-flutter",
				"api.pushy.me", "pushy.me",
				"Pushy.listen", "Pushy.register", "import me.pushy.sdk",
			},
			"name": "Pushy",
			"url":  "https://dashboard.pushy.me/apps",
			"key":  "push_service",
		},

		"airship": {
			"patterns": []string{
				"airship", "urban airship", "airship_app_key", "airship_master_secret",
				"urbanairship", "@urbanairship/node-library", "urbanairship-python",
				"api.airship.com", "go.airship.com",
				"urbanairship.push", "airship.push", "import urbanairship",
				"UAirship", "Airship.takeOff",
			},
			"name": "Airship (Urban Airship)",
			"url":  "https://go.airship.com/accounts/login",
			"key":  "push_service",
		},

		"pushwoosh": {
			"patterns": []string{
				"pushwoosh", "pushwoosh_app_id", "pushwoosh_api_token",
				"pushwoosh-react-native-plugin", "pushwoosh-cordova-plugin",
				"api.pushwoosh.com", "pushwoosh.com",
				"Pushwoosh.init", "pushwoosh.init", "import pushwoosh",
				"PWMessaging", "Pushwoosh.getInstance",
			},
			"name": "Pushwoosh",
			"url":  "https://cp.pushwoosh.com/applications",
			"key":  "push_service",
		},

		"twilio_notify": {
			"patterns": []string{
				"twilio notify", "twilio.notify", "twilio_notify_service_sid",
				"twilio-node", "twilio-python", "@twilio-labs/serverless-api",
				"api.twilio.com", "console.twilio.com",
				"client.notify", "twilio.rest.notify", "from twilio.rest import Client",
				"TwilioNotifyService", "twilio/twilio-go",
			},
			"name": "Twilio Notify",
			"url":  "https://console.twilio.com/us1/develop/notify/services",
			"key":  "push_service",
		},

		"amazon_sns": {
			"patterns": []string{
				"amazon sns", "aws sns", "sns_topic_arn", "aws_sns_topic_arn",
				"@aws-sdk/client-sns", "sns.amazonaws.com", "console.aws.amazon.com/sns",
				"sns.publish(", "sns_client.publish", "SNSClient", "PublishCommand",
				"github.com/aws/aws-sdk-go/service/sns", "AWS_SNS_TOPIC_ARN",
				"sns_publish_message", "amazon-sns", "aws-sns",
			},
			"name": "Amazon SNS",
			"url":  "https://console.aws.amazon.com/sns/v3/home",
			"key":  "push_service",
		},

		"apple_push": {
			"patterns": []string{
				"apple push notification", "apns_key_id", "apns_team_id",
				"node-apn", "apns2", "PyAPNs2", "apns-http2",
				"api.push.apple.com", "developer.apple.com",
				"apn.send", "apns.send_notification", "from apns2.client import",
				"APNSClient", "github.com/sideshow/apns2", "APNS_KEY_ID", "APNS_TEAM_ID",
				"apns_private_key", "apns_certificate", "apple_push_notification",
			},
			"name": "Apple Push Notification Service",
			"url":  "https://developer.apple.com/account/resources/authkeys/list",
			"key":  "push_service",
		},

		"webpush": {
			"patterns": []string{
				"web-push", "webpush", "vapid_public_key", "vapid_private_key",
				"web-push-libs", "pywebpush", "web-push-go",
				"push.mozilla.org", "fcm.googleapis.com/fcm/send",
				"webpush.sendNotification", "webpush.send_notification", "import webpush",
				"WebPushClient", "ServiceWorkerRegistration.pushManager",
			},
			"name": "Web Push Protocol",
			"url":  "https://web.dev/push-notifications/",
			"key":  "push_service",
		},
	}

	// Check for specific push notification services in order of popularity
	serviceOrder := []string{
		"firebase_messaging", "onesignal", "pusher", "expo_push", "apple_push",
		"amazon_sns", "twilio_notify", "sendbird", "pushy", "airship",
		"pushwoosh", "webpush",
	}

	for _, serviceKey := range serviceOrder {
		serviceInfo := services[serviceKey]
		patterns := serviceInfo["patterns"].([]string)

		// Check each pattern against all file contents
		for _, pattern := range patterns {
			found := false
			for _, fileContent := range fileContents {
				if strings.Contains(fileContent.Content, pattern) {
					// Find the exact line where the pattern was found
					lineNum := 0
					sourceText := ""
					for i, line := range fileContent.Lines {
						if strings.Contains(strings.ToLower(line), pattern) {
							lineNum = i + 1 // 1-indexed
							sourceText = strings.TrimSpace(line)
							break
						}
					}

					results = append(results, &DetectionResult{
						Key:         serviceInfo["key"].(string),
						Value:       serviceInfo["url"].(string),
						Description: serviceInfo["name"].(string) + " detected in project",
						Confidence:  0.90,
						DebugInfo:   fmt.Sprintf("Found pattern '%s' in %s", pattern, fileContent.Path),
						SourceFile:  fileContent.Path,
						SourceLine:  lineNum,
						SourceText:  sourceText,
					})
					found = true
					break // Only add each service once
				}
			}
			if found {
				break // Only add each service once
			}
		}
	}

	return results, nil
}