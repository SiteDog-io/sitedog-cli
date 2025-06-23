package detectors

import (
	"io/ioutil"
	"os"
	"strings"
)

// ShopifyDetector detects Shopify projects (apps and themes)
type ShopifyDetector struct{}

func (s *ShopifyDetector) Name() string {
	return "shopify"
}

func (s *ShopifyDetector) Description() string {
	return "Shopify app and theme detector"
}

func (s *ShopifyDetector) ShouldRun() bool {
	// Check for Shopify app configuration
	if _, err := os.Stat("shopify.app.toml"); err == nil {
		return true
	}

	// Check for Shopify CLI config
	if _, err := os.Stat(".shopify-cli.yml"); err == nil {
		return true
	}

	// Check for Shopify theme structure
	if _, err := os.Stat("config.yml"); err == nil {
		// Check if it's in a typical Shopify theme structure
		if _, err := os.Stat("templates"); err == nil {
			return true
		}
		if _, err := os.Stat("sections"); err == nil {
			return true
		}
		if _, err := os.Stat("snippets"); err == nil {
			return true
		}
	}

	// Check for package.json with Shopify dependencies
	if _, err := os.Stat("package.json"); err == nil {
		if data, readErr := ioutil.ReadFile("package.json"); readErr == nil {
			content := strings.ToLower(string(data))
			if strings.Contains(content, "@shopify/") ||
			   strings.Contains(content, "shopify-cli") ||
			   strings.Contains(content, "theme-kit") {
				return true
			}
		}
	}

	// Check for Gemfile with Shopify gems
	if _, err := os.Stat("Gemfile"); err == nil {
		if data, readErr := ioutil.ReadFile("Gemfile"); readErr == nil {
			content := strings.ToLower(string(data))
			if strings.Contains(content, "shopify_api") ||
			   strings.Contains(content, "shopify_app") ||
			   strings.Contains(content, "theme_kit") {
				return true
			}
		}
	}

	return false
}

func (s *ShopifyDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Auto-add Shopify platform links (high confidence)
	results = append(results, &DetectionResult{
		Key:         "shopify_partners",
		Value:       "https://partners.shopify.com",
		Description: "Shopify Partners Dashboard detected for Shopify development",
		Confidence:  0.98, // Auto-add
	})

	results = append(results, &DetectionResult{
		Key:         "shopify_admin",
		Value:       "https://admin.shopify.com",
		Description: "Shopify Admin detected for store management",
		Confidence:  0.98, // Auto-add
	})

	// Check what type of Shopify project this is
	var projectContent strings.Builder

	// Read shopify.app.toml if exists
	if data, err := ioutil.ReadFile("shopify.app.toml"); err == nil {
		projectContent.WriteString(strings.ToLower(string(data)))
	}

	// Read .shopify-cli.yml if exists
	if data, err := ioutil.ReadFile(".shopify-cli.yml"); err == nil {
		projectContent.WriteString(strings.ToLower(string(data)))
	}

	// Read package.json if exists
	if data, err := ioutil.ReadFile("package.json"); err == nil {
		projectContent.WriteString(strings.ToLower(string(data)))
	}

	// Read Gemfile if exists
	if data, err := ioutil.ReadFile("Gemfile"); err == nil {
		projectContent.WriteString(strings.ToLower(string(data)))
	}

	content := projectContent.String()

	// Detect specific Shopify services and tools - only real services with dashboards
	services := map[string]map[string]interface{}{
		"polaris": {
			"patterns": []string{"@shopify/polaris", "polaris-react", "polaris"},
			"name": "Shopify Polaris Design System",
			"url": "https://polaris.shopify.com",
			"key": "polaris",
		},

		"hydrogen": {
			"patterns": []string{"@shopify/hydrogen", "hydrogen"},
			"name": "Shopify Hydrogen",
			"url": "https://hydrogen.shopify.dev",
			"key": "hydrogen",
		},

		"shopify_payments": {
			"patterns": []string{"shopify payments", "payments api"},
			"name": "Shopify Payments",
			"url": "https://www.shopify.com/payments",
			"key": "payments",
		},

		"shopify_plus": {
			"patterns": []string{"shopify plus", "plus api"},
			"name": "Shopify Plus",
			"url": "https://www.shopify.com/plus",
			"key": "shopify_plus",
		},

		"ngrok": {
			"patterns": []string{"ngrok", "tunnel"},
			"name": "ngrok for Shopify development",
			"url": "https://ngrok.com",
			"key": "development_tunnel",
		},
	}

	// Check for specific services
	serviceOrder := []string{
		"polaris", "hydrogen", "shopify_payments", "shopify_plus", "ngrok",
	}

	for _, serviceKey := range serviceOrder {
		serviceInfo := services[serviceKey]
		patterns := serviceInfo["patterns"].([]string)

		for _, pattern := range patterns {
			if strings.Contains(content, pattern) {
				results = append(results, &DetectionResult{
					Key:         serviceInfo["key"].(string),
					Value:       serviceInfo["url"].(string),
					Description: serviceInfo["name"].(string) + " detected in Shopify project",
					Confidence:  0.85,
				})
				break // Only add each service once
			}
		}
	}

	// Check for specific file patterns to determine project type
	if _, err := os.Stat("shopify.app.toml"); err == nil {
		results = append(results, &DetectionResult{
			Key:         "shopify_app_store",
			Value:       "https://apps.shopify.com/partners",
			Description: "Shopify App Store detected for app publishing",
			Confidence:  0.90, // Conditional auto-add
		})
	}

	// Check for theme structure
	if _, err := os.Stat("templates"); err == nil {
		if _, err := os.Stat("sections"); err == nil {
			results = append(results, &DetectionResult{
				Key:         "shopify_themes",
				Value:       "https://themes.shopify.com/services/themes",
				Description: "Shopify Theme Store detected for theme development",
				Confidence:  0.90, // Conditional auto-add
			})
		}
	}

	return results, nil
}