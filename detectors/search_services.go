package detectors

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// SearchServicesDetector detects search and indexing services used in projects
type SearchServicesDetector struct{}

func (s *SearchServicesDetector) Name() string {
	return "search-services"
}

func (s *SearchServicesDetector) Description() string {
	return "Search and indexing services detector"
}

func (s *SearchServicesDetector) ShouldRun() bool {
	// Check for common files that might contain search integrations
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
			// Quick check for search-related keywords
			searchKeywords := []string{
				"algolia", "elasticsearch", "solr", "opensearch", "meilisearch",
				"typesense", "swiftype", "searchkit", "instantsearch",
			}
			for _, keyword := range searchKeywords {
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

func (s *SearchServicesDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Files to check for search service references - only configuration files
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
	}

	// Store file contents with metadata
	type FileContent struct {
		Path    string
		Content string
		Lines   []string
	}
	var fileContents []FileContent

	// Read individual files
	for _, file := range files {
		if data, err := ioutil.ReadFile(file); err == nil {
			content := string(data)
			fileContents = append(fileContents, FileContent{
				Path:    file,
				Content: strings.ToLower(content),
				Lines:   strings.Split(content, "\n"),
			})
		}
	}

	// Deep search mode - also check source code directories for search imports/usage
	if DeepSearchMode {
		srcDirs := []string{"src", "lib", "app", "components", "pages", "api", "services", "utils", "search"}
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
							if data, readErr := ioutil.ReadFile(path); readErr == nil {
								content := string(data)
								fileContents = append(fileContents, FileContent{
									Path:    path,
									Content: strings.ToLower(content),
									Lines:   strings.Split(content, "\n"),
								})
							}
						}
					}
					return nil
				})
			}
		}
	}

	// Define search services with their patterns and dashboards
	services := map[string]map[string]interface{}{
		"algolia": {
			"patterns": []string{
				"algolia", "algolia_api_key", "algolia_app_id", "algolia_admin_key",
				"algoliasearch", "algolia-search", "@algolia/client-search",
				"api.algolia.com", "dashboard.algolia.com",
				"from algoliasearch import", "import algoliasearch", "algoliasearch.Client",
				"InstantSearch", "react-instantsearch", "vue-instantsearch",
				"algolia/algoliasearch-client-go", "algolia/algoliasearch-client-php",
			},
			"name": "Algolia Search",
			"url":  "https://dashboard.algolia.com/billing",
			"key":  "search_service",
		},

		"elasticsearch": {
			"patterns": []string{
				"elasticsearch", "elasticsearch_url", "elastic_url",
				"@elastic/elasticsearch", "elasticsearch-py", "elasticsearch-dsl",
				"elastic.co", "cloud.elastic.co", "elasticsearch.org",
				"from elasticsearch import", "import elasticsearch", "Elasticsearch(",
				"elastic/elasticsearch", "olivere/elastic", "ELASTICSEARCH_URL",
				"ELASTIC_SEARCH_URL", "ES_HOST", "ES_URL",
			},
			"name": "Elasticsearch",
			"url":  "https://cloud.elastic.co/deployments",
			"key":  "search_service",
		},

		"opensearch": {
			"patterns": []string{
				"opensearch", "opensearch_url", "opensearch-py", "opensearch-js",
				"opensearch.org", "aws.amazon.com/opensearch",
				"from opensearchpy import", "import opensearchpy", "OpenSearch(",
				"opensearch-project/opensearch-go",
			},
			"name": "OpenSearch",
			"url":  "https://console.aws.amazon.com/opensearch",
			"key":  "search_service",
		},

		"meilisearch": {
			"patterns": []string{
				"meilisearch", "meilisearch_url", "meilisearch_api_key",
				"meilisearch-js", "meilisearch-python", "meilisearch-go",
				"meilisearch.com", "cloud.meilisearch.com",
				"from meilisearch import", "import meilisearch", "meilisearch.Client",
				"react-instantsearch-hooks-web", "instantsearch.js",
			},
			"name": "MeiliSearch",
			"url":  "https://cloud.meilisearch.com/projects",
			"key":  "search_service",
		},

		"typesense": {
			"patterns": []string{
				"typesense", "typesense_url", "typesense_api_key",
				"typesense-js", "typesense-python", "typesense-go",
				"typesense.org", "cloud.typesense.org",
				"from typesense import", "import typesense", "typesense.Client",
				"typesense/typesense-js", "typesense/typesense-go",
			},
			"name": "Typesense",
			"url":  "https://cloud.typesense.org/clusters",
			"key":  "search_service",
		},

		"swiftype": {
			"patterns": []string{
				"swiftype", "swiftype_api_key", "swiftype_engine_key",
				"swiftype-search-jquery", "swiftype-autocomplete-jquery",
				"api.swiftype.com", "swiftype.com",
				"elastic.co/products/site-search",
			},
			"name": "Swiftype Site Search",
			"url":  "https://app.swiftype.com/engines",
			"key":  "search_service",
		},

		"solr": {
			"patterns": []string{
				"solr", "apache solr", "solr_url", "solr_core",
				"pysolr", "solrj", "solr-node-client",
				"solr.apache.org", "lucene.apache.org/solr",
				"from pysolr import", "import pysolr", "solr.SolrClient",
				"org.apache.solr", "solrj-lib",
			},
			"name": "Apache Solr",
			"url":  "https://solr.apache.org/guide/solr/latest/deployment-guide/solr-admin-ui.html",
			"key":  "search_service",
		},
	}

	// Check for specific search services in order of popularity
	serviceOrder := []string{
		"algolia", "elasticsearch", "opensearch", "meilisearch",
		"typesense", "swiftype", "solr",
	}

	for _, serviceKey := range serviceOrder {
		serviceInfo := services[serviceKey]
		patterns := serviceInfo["patterns"].([]string)

		// Check each pattern in each file
		for _, pattern := range patterns {
			for _, fileContent := range fileContents {
				if strings.Contains(fileContent.Content, pattern) {
					// Find the line number where the pattern was found
					lineNum := 0
					sourceLine := ""
					for i, line := range fileContent.Lines {
						if strings.Contains(strings.ToLower(line), pattern) {
							lineNum = i + 1
							sourceLine = strings.TrimSpace(line)
							break
						}
					}

					results = append(results, &DetectionResult{
						Key:         serviceInfo["key"].(string),
						Value:       serviceInfo["url"].(string),
						Description: serviceInfo["name"].(string) + " detected in project",
						Confidence:  0.90,
						DebugInfo:   "Found pattern '" + pattern + "' in " + fileContent.Path,
						SourceFile:  fileContent.Path,
						SourceLine:  lineNum,
						SourceText:  maskSecrets(sourceLine),
					})
					goto nextService // Only add each service once
				}
			}
		}
		nextService:
	}

	return results, nil
}