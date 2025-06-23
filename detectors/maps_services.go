package detectors

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// MapsServicesDetector detects maps, geocoding, and location services used in projects
type MapsServicesDetector struct{}

func (m *MapsServicesDetector) Name() string {
	return "maps-services"
}

func (m *MapsServicesDetector) Description() string {
	return "Maps, geocoding, and location services detector"
}

func (m *MapsServicesDetector) ShouldRun() bool {
	// Check for common files that might contain maps/location integrations
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
			// Quick check for maps/location-related keywords
			mapsKeywords := []string{
				"google maps", "mapbox", "openstreetmap", "geocoding", "geolocation",
				"mapquest", "here maps", "tomtom", "bing maps", "leaflet",
			}
			for _, keyword := range mapsKeywords {
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

func (m *MapsServicesDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Read all relevant files to detect maps services
	var projectContent strings.Builder

	// Files to check for maps service references
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
		"package-lock.json",
		"yarn.lock",
		"poetry.lock",
		"Pipfile",
		"Pipfile.lock",
	}

	// Also check source code directories for maps imports/usage
	srcDirs := []string{"src", "lib", "app", "components", "pages", "api", "services", "utils", "maps", "location"}
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
							projectContent.WriteString(strings.ToLower(string(data)))
						}
					}
				}
				return nil
			})
		}
	}

	// Read individual files
	for _, file := range files {
		if data, err := ioutil.ReadFile(file); err == nil {
			projectContent.WriteString(strings.ToLower(string(data)))
		}
	}

	content := projectContent.String()

	// Define maps services with their patterns and dashboards
	services := map[string]map[string]interface{}{
		"google_maps": {
			"patterns": []string{
				"google maps", "googlemaps", "google_maps_api_key", "maps_api_key",
				"@googlemaps/js-api-loader", "@google/maps", "googlemaps-js-api",
				"maps.googleapis.com", "console.cloud.google.com/google/maps-apis",
				"google.maps", "new google.maps", "GoogleMap", "react-google-maps",
				"vue-google-maps", "google-maps-react", "googlemaps/google-maps-services-go",
			},
			"name": "Google Maps Platform",
			"url":  "https://console.cloud.google.com/google/maps-apis/credentials",
			"key":  "maps_service",
		},

		"mapbox": {
			"patterns": []string{
				"mapbox", "mapbox_api_key", "mapbox_access_token", "mapbox_token",
				"mapbox-gl", "mapbox-gl-js", "react-map-gl", "@mapbox/mapbox-gl-geocoder",
				"api.mapbox.com", "account.mapbox.com",
				"mapboxgl", "mapbox.map", "new mapboxgl.Map", "MapboxGL",
				"mapbox/mapbox-sdk-py", "mapbox/mapbox-java",
			},
			"name": "Mapbox",
			"url":  "https://account.mapbox.com/billing",
			"key":  "maps_service",
		},

		"here_maps": {
			"patterns": []string{
				"here maps", "here_api_key", "here_app_id", "here_app_code",
				"@here/maps-api-for-javascript", "heremaps",
				"developer.here.com", "here.com",
				"here.map", "H.Map", "HERE.Map", "here-js-api",
				"heremaps/here-location-services-python",
			},
			"name": "HERE Maps",
			"url":  "https://developer.here.com/projects",
			"key":  "maps_service",
		},

		"opencage": {
			"patterns": []string{
				"opencage", "opencage_api_key", "opencagedata",
				"opencage-geocoder", "python-opencage-geocoder",
				"api.opencagedata.com", "opencagedata.com",
				"opencage.geocode", "OpenCageGeocode",
			},
			"name": "OpenCage Geocoding",
			"url":  "https://opencagedata.com/dashboard",
			"key":  "maps_service",
		},

		"mapquest": {
			"patterns": []string{
				"mapquest", "mapquest_api_key", "mapquest_key",
				"mapquest-js", "mapquest-sdk",
				"developer.mapquest.com", "mapquestapi.com",
				"mapquest.map", "MapQuest", "L.mapquest",
			},
			"name": "MapQuest",
			"url":  "https://developer.mapquest.com/user/me/apps",
			"key":  "maps_service",
		},

		"tomtom": {
			"patterns": []string{
				"tomtom", "tomtom_api_key", "tomtom_key",
				"@tomtom-international/web-sdk-maps", "@tomtom-international/web-sdk-services",
				"developer.tomtom.com", "api.tomtom.com",
				"tomtom.map", "TomTomMap", "tt.map",
			},
			"name": "TomTom Maps",
			"url":  "https://developer.tomtom.com/user/me/apps",
			"key":  "maps_service",
		},

		"bing_maps": {
			"patterns": []string{
				"bing maps", "bingmaps", "bing_maps_api_key", "bing_maps_key",
				"bingmaps-js", "microsoft.maps",
				"dev.virtualearth.net", "bing.com/maps",
				"Microsoft.Maps", "new Microsoft.Maps.Map", "BingMapsAPI",
			},
			"name": "Bing Maps",
			"url":  "https://www.bingmapsportal.com/Application",
			"key":  "maps_service",
		},

		"locationiq": {
			"patterns": []string{
				"locationiq", "locationiq_api_key", "locationiq_token",
				"locationiq-js-client", "python-locationiq",
				"locationiq.com", "eu1.locationiq.com",
				"locationiq.geocode", "LocationIQ",
			},
			"name": "LocationIQ",
			"url":  "https://my.locationiq.com/dashboard",
			"key":  "maps_service",
		},

		"positionstack": {
			"patterns": []string{
				"positionstack", "positionstack_api_key", "positionstack_access_key",
				"api.positionstack.com", "positionstack.com",
				"positionstack.geocode", "PositionStack",
			},
			"name": "PositionStack",
			"url":  "https://positionstack.com/dashboard",
			"key":  "maps_service",
		},

		"geocodio": {
			"patterns": []string{
				"geocodio", "geocodio_api_key", "geocod.io",
				"pygeocodio", "geocodio-js", "geocodio-php",
				"api.geocod.io", "dash.geocod.io",
				"geocodio.geocode", "Geocodio",
			},
			"name": "Geocodio",
			"url":  "https://dash.geocod.io/usage",
			"key":  "maps_service",
		},

		"ipgeolocation": {
			"patterns": []string{
				"ipgeolocation", "ipgeolocation_api_key", "ipgeolocationapi",
				"ipgeolocation-js", "python-ipgeolocation",
				"api.ipgeolocation.io", "ipgeolocation.io",
				"ipgeolocation.locate", "IPGeolocation",
			},
			"name": "IPGeolocation",
			"url":  "https://ipgeolocation.io/dashboard",
			"key":  "maps_service",
		},

		"what3words": {
			"patterns": []string{
				"what3words", "what3words_api_key", "w3w_api_key",
				"@what3words/api", "what3words-python", "what3words-java",
				"api.what3words.com", "what3words.com",
				"what3words.convert", "What3Words", "w3w.convert",
			},
			"name": "What3Words",
			"url":  "https://accounts.what3words.com/billing",
			"key":  "maps_service",
		},
	}

	// Check for specific maps services in order of popularity
	serviceOrder := []string{
		"google_maps", "mapbox", "here_maps", "opencage", "mapquest",
		"tomtom", "bing_maps", "locationiq", "positionstack",
		"geocodio", "ipgeolocation", "what3words",
	}

	for _, serviceKey := range serviceOrder {
		serviceInfo := services[serviceKey]
		patterns := serviceInfo["patterns"].([]string)

		for _, pattern := range patterns {
			if strings.Contains(content, pattern) {
				results = append(results, &DetectionResult{
					Key:         serviceInfo["key"].(string),
					Value:       serviceInfo["url"].(string),
					Description: serviceInfo["name"].(string) + " detected in project",
					Confidence:  0.90,
				})
				break // Only add each service once
			}
		}
	}

	return results, nil
}