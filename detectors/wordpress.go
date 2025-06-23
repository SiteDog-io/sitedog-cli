package detectors

import (
	"io/ioutil"
	"os"
	"strings"
)

// WordPressDetector detects WordPress projects
type WordPressDetector struct{}

func (w *WordPressDetector) Name() string {
	return "wordpress"
}

func (w *WordPressDetector) Description() string {
	return "WordPress platform detector"
}

func (w *WordPressDetector) ShouldRun() bool {
	// Check for WordPress core configuration file - 100% confidence
	if _, err := os.Stat("wp-config.php"); err == nil {
		return true
	}

	// Check for WordPress theme structure (style.css with Theme Name header)
	if _, err := os.Stat("style.css"); err == nil {
		if data, readErr := ioutil.ReadFile("style.css"); readErr == nil {
			content := strings.ToLower(string(data))
			if strings.Contains(content, "theme name:") {
				return true
			}
		}
	}

	return false
}

func (w *WordPressDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Auto-add WordPress platform links (high confidence)
	results = append(results, &DetectionResult{
		Key:         "wordpress_admin",
		Value:       "https://wordpress.com/wp-admin",
		Description: "WordPress Admin Dashboard detected for site management",
		Confidence:  0.98, // Auto-add
	})

	// Read all relevant files to detect services
	var projectContent strings.Builder

	// Common WordPress files to check
	files := []string{"wp-config.php", "style.css", "functions.php", "index.php", "composer.json", "package.json", "readme.txt"}

	for _, file := range files {
		if data, err := ioutil.ReadFile(file); err == nil {
			projectContent.WriteString(strings.ToLower(string(data)))
		}
	}

	content := projectContent.String()

	// Detect WordPress hosting and services - only real services with dashboards
	services := map[string]map[string]interface{}{
		"wp_engine": {
			"patterns": []string{"wpengine", "wp engine", "wpengine.com"},
			"name": "WP Engine Hosting",
			"url": "https://my.wpengine.com",
			"key": "hosting",
		},

		"kinsta": {
			"patterns": []string{"kinsta", "kinsta.com"},
			"name": "Kinsta Hosting",
			"url": "https://my.kinsta.com",
			"key": "hosting",
		},

		"siteground": {
			"patterns": []string{"siteground", "siteground.com"},
			"name": "SiteGround Hosting",
			"url": "https://tools.siteground.com",
			"key": "hosting",
		},

		"cloudflare": {
			"patterns": []string{"cloudflare", "cf-ray", "cloudflare.com"},
			"name": "Cloudflare CDN",
			"url": "https://dash.cloudflare.com",
			"key": "cdn",
		},

		"jetpack": {
			"patterns": []string{"jetpack", "jetpack.com"},
			"name": "Jetpack by WordPress.com",
			"url": "https://wordpress.com/jetpack",
			"key": "jetpack",
		},

		"woocommerce": {
			"patterns": []string{"woocommerce", "wc-", "woo commerce"},
			"name": "WooCommerce",
			"url": "https://woocommerce.com/my-account",
			"key": "ecommerce",
		},

		"yoast_seo": {
			"patterns": []string{"yoast", "yoast seo", "wordpress seo"},
			"name": "Yoast SEO",
			"url": "https://my.yoast.com",
			"key": "seo",
		},

		"elementor": {
			"patterns": []string{"elementor", "elementor.com"},
			"name": "Elementor Page Builder",
			"url": "https://my.elementor.com",
			"key": "page_builder",
		},

		"wpml": {
			"patterns": []string{"wpml", "wpml.org"},
			"name": "WPML Translation",
			"url": "https://wpml.org/account",
			"key": "translation",
		},

		"mailchimp": {
			"patterns": []string{"mailchimp", "mc4wp", "mailchimp.com"},
			"name": "Mailchimp",
			"url": "https://mailchimp.com",
			"key": "email_marketing",
		},

		"google_analytics_wp": {
			"patterns": []string{"google analytics", "gtag", "ga.js", "analytics.js"},
			"name": "Google Analytics",
			"url": "https://analytics.google.com",
			"key": "analytics",
		},

		"google_search_console": {
			"patterns": []string{"google search console", "search console", "webmaster"},
			"name": "Google Search Console",
			"url": "https://search.google.com/search-console",
			"key": "seo_tools",
		},

		"stripe_wp": {
			"patterns": []string{"stripe", "stripe.com", "pk_test", "sk_test"},
			"name": "Stripe Payments",
			"url": "https://dashboard.stripe.com",
			"key": "payments",
		},

		"paypal_wp": {
			"patterns": []string{"paypal", "paypal.com"},
			"name": "PayPal",
			"url": "https://paypal.com",
			"key": "payments",
		},

		"akismet": {
			"patterns": []string{"akismet", "akismet.com"},
			"name": "Akismet Anti-Spam",
			"url": "https://akismet.com/account",
			"key": "security",
		},

		"wordfence": {
			"patterns": []string{"wordfence", "wordfence.com"},
			"name": "Wordfence Security",
			"url": "https://www.wordfence.com/central",
			"key": "security",
		},

		"updraftplus": {
			"patterns": []string{"updraftplus", "updraft plus"},
			"name": "UpdraftPlus Backup",
			"url": "https://updraftplus.com/my-account",
			"key": "backup",
		},

		"wp_rocket": {
			"patterns": []string{"wp rocket", "wp-rocket", "wprocket"},
			"name": "WP Rocket Caching",
			"url": "https://wp-rocket.me/account",
			"key": "performance",
		},
	}

	// Check for specific services
	serviceOrder := []string{
		"wp_engine", "kinsta", "siteground", "cloudflare",
		"jetpack", "woocommerce", "yoast_seo", "elementor", "wpml",
		"mailchimp", "google_analytics_wp", "google_search_console",
		"stripe_wp", "paypal_wp", "akismet", "wordfence", "updraftplus", "wp_rocket",
	}

	for _, serviceKey := range serviceOrder {
		serviceInfo := services[serviceKey]
		patterns := serviceInfo["patterns"].([]string)

		for _, pattern := range patterns {
			if strings.Contains(content, pattern) {
				results = append(results, &DetectionResult{
					Key:         serviceInfo["key"].(string),
					Value:       serviceInfo["url"].(string),
					Description: serviceInfo["name"].(string) + " detected in WordPress project",
					Confidence:  0.85,
				})
				break // Only add each service once
			}
		}
	}

	return results, nil
}