package detectors

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// ImageServicesDetector detects image processing and optimization services used in projects
type ImageServicesDetector struct{}

func (i *ImageServicesDetector) Name() string {
	return "image-services"
}

func (i *ImageServicesDetector) Description() string {
	return "Image processing, optimization and management services detector"
}

func (i *ImageServicesDetector) ShouldRun() bool {
	// Check for common files that might contain image service integrations
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
			// Quick check for image service-related keywords
			imageKeywords := []string{
				"cloudinary", "uploadcare", "imagekit", "tinypng", "kraken",
				"optimole", "shortpixel", "imageoptim", "compressor", "imagify",
				"unsplash", "pexels", "pixabay", "shutterstock",
			}
			for _, keyword := range imageKeywords {
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

func (i *ImageServicesDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Read all relevant files to detect image services
	var projectContent strings.Builder

	// Files to check for image service references
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
		"next.config.js",
		"nuxt.config.js",
		"gatsby-config.js",
		"webpack.config.js",
	}

	// Also check source code directories for image service imports/usage
	srcDirs := []string{"src", "lib", "app", "components", "pages", "api", "services", "utils", "images", "assets"}
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

	// Define image services with their patterns and dashboards
	services := map[string]map[string]interface{}{
		"cloudinary": {
			"patterns": []string{
				"cloudinary", "cloudinary_url", "cloudinary_cloud_name", "cloudinary_api_key",
				"cloudinary-core", "cloudinary-react", "cloudinary-vue", "@cloudinary/url-gen",
				"api.cloudinary.com", "cloudinary.com",
				"cloudinary.image", "cloudinary.uploader", "from cloudinary import",
				"CloudinaryImage", "github.com/cloudinary/cloudinary-go",
			},
			"name": "Cloudinary",
			"url":  "https://cloudinary.com/console/settings/account",
			"key":  "image_service",
		},

		"uploadcare": {
			"patterns": []string{
				"uploadcare", "uploadcare_public_key", "uploadcare_secret_key",
				"@uploadcare/upload-client", "@uploadcare/react-widget", "uploadcare-widget",
				"api.uploadcare.com", "uploadcare.com",
				"uploadcare.file", "uploadcare.upload", "from pyuploadcare import",
				"UploadcareClient", "github.com/uploadcare/uploadcare-go",
			},
			"name": "Uploadcare",
			"url":  "https://uploadcare.com/dashboard/",
			"key":  "image_service",
		},

		"imagekit": {
			"patterns": []string{
				"imagekit", "imagekit_public_key", "imagekit_private_key", "imagekit_url_endpoint",
				"imagekit-javascript", "imagekit-react", "imagekitio-vue",
				"api.imagekit.io", "imagekit.io",
				"imagekit.upload", "imagekit.url", "from imagekitio import",
				"ImageKitClient", "github.com/imagekit-developer/imagekit-go",
			},
			"name": "ImageKit",
			"url":  "https://imagekit.io/dashboard/developer/api-keys",
			"key":  "image_service",
		},

		"tinypng": {
			"patterns": []string{
				"tinypng", "tinify", "tinypng_api_key", "tinify_key",
				"tinify-node", "tinify-python", "tinypng-compress",
				"api.tinify.com", "tinypng.com",
				"tinify.fromdatabase", "tinify.fromfile", "import tinify",
				"TinifyClient", "github.com/tinify/tinify-go",
			},
			"name": "TinyPNG",
			"url":  "https://tinypng.com/developers",
			"key":  "image_service",
		},

		"kraken": {
			"patterns": []string{
				"kraken", "krakenio", "kraken_api_key", "kraken_api_secret",
				"kraken-io", "krakenio-node", "kraken-io-python",
				"api.kraken.io", "kraken.io",
				"kraken.upload", "kraken.url", "from krakenio import",
				"KrakenClient", "github.com/kraken-io/kraken-go",
			},
			"name": "Kraken.io",
			"url":  "https://kraken.io/account/api-credentials",
			"key":  "image_service",
		},

		"optimole": {
			"patterns": []string{
				"optimole", "optimole_api_key", "optimole_service_url",
				"optimole-sdk", "optimole-php", "optimole-js",
				"dashboard.optimole.com", "optimole.com",
				"optimole.generateurl", "optimole.upload", "from optimole import",
				"OptimoleClient", "Optimole\\OptiMole",
			},
			"name": "Optimole",
			"url":  "https://dashboard.optimole.com/",
			"key":  "image_service",
		},

		"shortpixel": {
			"patterns": []string{
				"shortpixel", "shortpixel_api_key", "short_pixel_api_key",
				"shortpixel-image-optimiser", "shortpixel-php", "shortpixel-node",
				"api.shortpixel.com", "shortpixel.com",
				"shortpixel.fromfile", "shortpixel.fromurl", "from shortpixel import",
				"ShortPixelClient", "ShortPixel\\ShortPixel",
			},
			"name": "ShortPixel",
			"url":  "https://shortpixel.com/login",
			"key":  "image_service",
		},

		"imageoptim": {
			"patterns": []string{
				"imageoptim", "imageoptim_api_username", "imageoptim_api_key",
				"imageoptim-cli", "imageoptim-node", "imageoptim-python",
				"im2.io", "imageoptim.com",
				"imageoptim.optimize", "imageoptim.upload", "from imageoptim import",
				"ImageOptimClient", "github.com/imageoptim/imageoptim-cli",
			},
			"name": "ImageOptim API",
			"url":  "https://imageoptim.com/api/register",
			"key":  "image_service",
		},

		"imagify": {
			"patterns": []string{
				"imagify", "imagify_api_key", "imagify_secret_key",
				"imagify-node", "imagify-php", "imagify-python",
				"api.imagify.io", "imagify.io",
				"imagify.optimize", "imagify.upload", "from imagify import",
				"ImagifyClient", "Imagify\\API\\Imagify",
			},
			"name": "Imagify",
			"url":  "https://app.imagify.io/api/",
			"key":  "image_service",
		},

		"unsplash": {
			"patterns": []string{
				"unsplash", "unsplash_access_key", "unsplash_secret_key", "unsplash_application_id",
				"unsplash-js", "python-unsplash", "unsplash-php",
				"api.unsplash.com", "unsplash.com",
				"unsplash.photos", "unsplash.search", "from unsplash.api import",
				"UnsplashApi", "github.com/unsplash/unsplash-go",
			},
			"name": "Unsplash API",
			"url":  "https://unsplash.com/developers",
			"key":  "image_service",
		},

		"pexels": {
			"patterns": []string{
				"pexels", "pexels_api_key", "pexels_authorization",
				"pexels-api", "pexels-js", "pexels-python",
				"api.pexels.com", "pexels.com",
				"pexels.search", "pexels.photos", "from pexels_api import",
				"PexelsApi", "github.com/pexels/pexels-go",
			},
			"name": "Pexels API",
			"url":  "https://www.pexels.com/api/",
			"key":  "image_service",
		},

		"pixabay": {
			"patterns": []string{
				"pixabay", "pixabay_api_key", "pixabay_key",
				"pixabay-api", "pixabay-js", "pixabay-python",
				"pixabay.com/api", "pixabay.com",
				"pixabay.search", "pixabay.images", "from pixabay import",
				"PixabayApi", "github.com/pixabay/pixabay-go",
			},
			"name": "Pixabay API",
			"url":  "https://pixabay.com/api/docs/",
			"key":  "image_service",
		},

		"shutterstock": {
			"patterns": []string{
				"shutterstock", "shutterstock_client_id", "shutterstock_client_secret",
				"shutterstock-api", "shutterstock-node", "shutterstock-python",
				"api.shutterstock.com", "shutterstock.com",
				"shutterstock.search", "shutterstock.images", "from shutterstock_api import",
				"ShutterstockApi", "github.com/shutterstock/shutterstock-go",
			},
			"name": "Shutterstock API",
			"url":  "https://www.shutterstock.com/account/developers/apps",
			"key":  "image_service",
		},

		"imgproxy": {
			"patterns": []string{
				"imgproxy", "imgproxy_key", "imgproxy_salt", "imgproxy_url",
				"imgproxy-node", "imgproxy-php", "imgproxy-python",
				"imgproxy.net", "imgproxy:latest", "darthsim/imgproxy",
				"imgproxy_base_url", "imgproxy_signature", "imgproxy_endpoint",
				"IMGPROXY_KEY", "IMGPROXY_SALT", "IMGPROXY_BASE_URL",
				"imgproxy.generate_url", "imgproxy_url_for", "from imgproxy import",
			},
			"name": "imgproxy",
			"url":  "https://imgproxy.net/",
			"key":  "image_service",
		},
	}

	// Check for specific image services in order of popularity
	serviceOrder := []string{
		"cloudinary", "uploadcare", "imagekit", "imgproxy", "tinypng", "kraken",
		"optimole", "shortpixel", "imageoptim", "imagify", "unsplash",
		"pexels", "pixabay", "shutterstock",
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