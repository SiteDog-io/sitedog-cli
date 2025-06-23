package detectors

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// I18nDetector detects internationalization and translation services
type I18nDetector struct{}

func (i *I18nDetector) Name() string {
	return "i18n"
}

func (i *I18nDetector) Description() string {
	return "Internationalization and translation services detector"
}

func (i *I18nDetector) ShouldRun() bool {
	// Check for common i18n files and directories
	i18nPaths := []string{
		"locales",
		"translations",
		"lang",
		"i18n",
		"locale",
		"intl",
		"messages",
	}

	for _, path := range i18nPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	// Check for i18n configuration files
	i18nFiles := []string{
		"i18next.config.js",
		"next-i18next.config.js",
		"vue-i18n.config.js",
		"crowdin.yml",
		"crowdin.yaml",
		".crowdin.yml",
		".crowdin.yaml",
		"lokalise.cfg",
		"lokalise.yml",
		".lokaliserc",
		"phrase.yml",
		".phraseapp.yml",
		"transifex.yml",
		".tx/config",
		".poeditorrc",
		"poeditor.yml",
		"weblate.yml",
		"tolgee.json",
		"lingohub.yml",
	}

	for _, file := range i18nFiles {
		if _, err := os.Stat(file); err == nil {
			return true
		}
	}

	// Check for package.json with i18n dependencies
	if data, err := ioutil.ReadFile("package.json"); err == nil {
		content := strings.ToLower(string(data))
		i18nPackages := []string{
			"react-i18next", "i18next", "vue-i18n", "angular-i18n",
			"next-i18next", "gatsby-plugin-i18n", "svelte-i18n",
			"crowdin", "lokalise", "phrase", "transifex",
		}

		for _, pkg := range i18nPackages {
			if strings.Contains(content, pkg) {
				return true
			}
		}
	}

	return false
}

func (i *I18nDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Read all relevant files to detect translation services
	var projectContent strings.Builder

	// Files to check for translation service references
	files := []string{
		"package.json",
		"composer.json",
		"Gemfile",
		"requirements.txt",
		"go.mod",
		"Cargo.toml",
		"crowdin.yml",
		"crowdin.yaml",
		".crowdin.yml",
		".crowdin.yaml",
		"lokalise.cfg",
		"lokalise.yml",
		".lokaliserc",
		"phrase.yml",
		".phraseapp.yml",
		"transifex.yml",
		".tx/config",
		".poeditorrc",
		"poeditor.yml",
		"weblate.yml",
		"tolgee.json",
		"lingohub.yml",
		"i18next.config.js",
		"next-i18next.config.js",
		"vue-i18n.config.js",
		".github/workflows/i18n.yml",
		".github/workflows/translations.yml",
		".gitlab-ci.yml",
		"Makefile",
		"README.md",
		"package-lock.json",
		"yarn.lock",
	}

	// Also check CI/CD workflow files for translation automation
	if _, err := os.Stat(".github/workflows"); err == nil {
		filepath.Walk(".github/workflows", func(path string, info os.FileInfo, err error) error {
			if err == nil && (strings.HasSuffix(info.Name(), ".yml") || strings.HasSuffix(info.Name(), ".yaml")) {
				if data, readErr := ioutil.ReadFile(path); readErr == nil {
					projectContent.WriteString(strings.ToLower(string(data)))
				}
			}
			return nil
		})
	}

	// Read individual files
	for _, file := range files {
		if data, err := ioutil.ReadFile(file); err == nil {
			projectContent.WriteString(strings.ToLower(string(data)))
		}
	}

	content := projectContent.String()

	// Define translation services with their patterns and dashboards
	services := map[string]map[string]interface{}{
		"crowdin": {
			"patterns": []string{
				"crowdin.yml", "crowdin.yaml", ".crowdin.yml", ".crowdin.yaml",
				"crowdin.com", "crowdin-cli", "crowdin/cli", "crowdin-action",
				"CROWDIN_", "crowdin_project_id", "crowdin_api_token",
			},
			"name": "Crowdin",
			"url":  "https://crowdin.com/project",
			"key":  "translation_service",
		},

		"lokalise": {
			"patterns": []string{
				"lokalise.cfg", "lokalise.yml", ".lokaliserc", "lokalise-cli", "lokalise2",
				"lokalise.com", "lokalise-action", "@lokalise/", "LOKALISE_",
				"lokalise_project_id", "lokalise_api_token", "lokalise\":",
			},
			"name": "Lokalise",
			"url":  "https://app.lokalise.com",
			"key":  "translation_service",
		},

		"phrase": {
			"patterns": []string{
				"phrase.yml", ".phraseapp.yml", "phrase-cli", "phraseapp",
				"phrase.com", "PHRASE_", "phrase_project_id", "phrase_access_token",
			},
			"name": "Phrase",
			"url":  "https://app.phrase.com",
			"key":  "translation_service",
		},

		"transifex": {
			"patterns": []string{
				"transifex.yml", ".tx/config", "transifex-client", "tx pull", "tx push",
				"transifex.com", "TRANSIFEX_", "transifex_api_token",
			},
			"name": "Transifex",
			"url":  "https://www.transifex.com/dashboard",
			"key":  "translation_service",
		},

		"weblate": {
			"patterns": []string{
				"weblate.yml", "weblate.com", "weblate-cli", "WEBLATE_",
				"weblate_api_key", "weblate_project",
			},
			"name": "Weblate",
			"url":  "https://hosted.weblate.org",
			"key":  "translation_service",
		},

		"tolgee": {
			"patterns": []string{
				"tolgee.json", "tolgee.com", "@tolgee/", "tolgee-cli",
				"TOLGEE_", "tolgee_api_key", "tolgee_project_id",
			},
			"name": "Tolgee",
			"url":  "https://app.tolgee.io",
			"key":  "translation_service",
		},

		"localize": {
			"patterns": []string{
				"localize.com", "localize-cli", "LOCALIZE_",
				"localize_api_key", "localize_project_id",
			},
			"name": "Localize",
			"url":  "https://app.localize.com",
			"key":  "translation_service",
		},

		"onesky": {
			"patterns": []string{
				"onesky", "oneskyapp.com", "ONESKY_",
				"onesky_api_key", "onesky_project_id",
			},
			"name": "OneSky",
			"url":  "https://www.oneskyapp.com",
			"key":  "translation_service",
		},

		"lingohub": {
			"patterns": []string{
				"lingohub.yml", "lingohub.com", "lingohub-cli", "LINGOHUB_",
				"lingohub_api_key", "lingohub_project",
			},
			"name": "LingoHub",
			"url":  "https://app.lingohub.com",
			"key":  "translation_service",
		},

		"poeditor": {
			"patterns": []string{
				"poeditor", "poeditor.com", "POEDITOR_",
				"poeditor_api_token", "poeditor_project_id", "poeditor-cli",
			},
			"name": "POEditor",
			"url":  "https://poeditor.com/projects",
			"key":  "translation_service",
		},

		"smartling": {
			"patterns": []string{
				"smartling", "smartling.com", "SMARTLING_",
				"smartling_api_key", "smartling_project_id",
			},
			"name": "Smartling",
			"url":  "https://dashboard.smartling.com",
			"key":  "translation_service",
		},
	}

	// Check for specific translation services in order of popularity
	serviceOrder := []string{
		"crowdin", "lokalise", "phrase", "transifex", "poeditor",
		"weblate", "tolgee", "localize", "onesky", "lingohub", "smartling",
	}

	for _, serviceKey := range serviceOrder {
		serviceInfo := services[serviceKey]
		patterns := serviceInfo["patterns"].([]string)

		for _, pattern := range patterns {
			if strings.Contains(content, pattern) {
				results = append(results, &DetectionResult{
					Key:         serviceInfo["key"].(string),
					Value:       serviceInfo["url"].(string),
					Description: serviceInfo["name"].(string) + " translation service detected in project",
					Confidence:  0.90,
				})
				break // Only add each service once
			}
		}
	}

	return results, nil
}