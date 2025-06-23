package main

import (
	"fmt"
	"os"
	"path/filepath"

	"sitedog-cli/detectors"
)

// Test configuration for each detector
type TestConfig struct {
	Name      string
	Detector  detectors.Detector
	TestDir   string
	TestFiles []string
}

func main() {
	// Define test configurations
	testConfigs := []TestConfig{
		{
			Name:      "Go Modules",
			Detector:  &detectors.GoModDetector{},
			TestDir:   "gomod",
			TestFiles: []string{"go.mod"},
		},
		{
			Name:      "Python Requirements",
			Detector:  &detectors.RequirementsDetector{},
			TestDir:   "requirements",
			TestFiles: []string{"requirements.txt"},
		},
		{
			Name:      "Node.js Package JSON",
			Detector:  &detectors.PackageJSONDetector{},
			TestDir:   "package_json",
			TestFiles: []string{"package.json"},
		},
		{
			Name:      "Ruby Gemfile",
			Detector:  &detectors.GemfileDetector{},
			TestDir:   "gemfile",
			TestFiles: []string{"Gemfile"},
		},
		{
			Name:      "PHP Composer",
			Detector:  &detectors.ComposerDetector{},
			TestDir:   "composer",
			TestFiles: []string{"composer.json"},
		},
		{
			Name:      "Rust Cargo",
			Detector:  &detectors.CargoDetector{},
			TestDir:   "cargo",
			TestFiles: []string{"Cargo.toml"},
		},
		{
			Name:      ".NET/C#",
			Detector:  &detectors.DotNetDetector{},
			TestDir:   "dotnet",
			TestFiles: []string{"MyApp.csproj", "packages.config", "project.json"},
		},
		{
			Name:      "Java",
			Detector:  &detectors.JavaDetector{},
			TestDir:   "java",
			TestFiles: []string{"pom.xml", "build.gradle", "build.gradle.kts"},
		},
		{
			Name:      "Flutter/Dart",
			Detector:  &detectors.DartDetector{},
			TestDir:   "dart",
			TestFiles: []string{"pubspec.yaml"},
		},
		{
			Name:      "iOS Platform",
			Detector:  &detectors.IOSDetector{},
			TestDir:   "ios",
			TestFiles: []string{"Runner/Info.plist", "Podfile"},
		},
		{
			Name:      "Android Platform",
			Detector:  &detectors.AndroidDetector{},
			TestDir:   "android",
			TestFiles: []string{"app/build.gradle", "app/src/main/AndroidManifest.xml"},
		},
		{
			Name:      "Chrome Extension",
			Detector:  &detectors.ChromeExtensionDetector{},
			TestDir:   "chrome_extension",
			TestFiles: []string{"manifest.json"},
		},
		{
			Name:      "VS Code Extension",
			Detector:  &detectors.VSCodeExtensionDetector{},
			TestDir:   "vscode_extension",
			TestFiles: []string{"package.json"},
		},
		{
			Name:      "CircleCI",
			Detector:  &detectors.CircleCIDetector{},
			TestDir:   "circleci",
			TestFiles: []string{".circleci/config.yml"},
		},
		{
			Name:      "Travis CI",
			Detector:  &detectors.TravisCIDetector{},
			TestDir:   "travis",
			TestFiles: []string{".travis.yml"},
		},
		{
			Name:      "Azure Pipelines",
			Detector:  &detectors.AzurePipelinesDetector{},
			TestDir:   "azurepipelines",
			TestFiles: []string{"azure-pipelines.yml"},
		},
		{
			Name:      "Jenkins",
			Detector:  &detectors.JenkinsDetector{},
			TestDir:   "jenkins",
			TestFiles: []string{"Jenkinsfile"},
		},
		{
			Name:      "Bitbucket Pipelines",
			Detector:  &detectors.BitbucketPipelinesDetector{},
			TestDir:   "bitbucket",
			TestFiles: []string{"bitbucket-pipelines.yml"},
		},
		{
			Name:      "Vercel",
			Detector:  &detectors.VercelDetector{},
			TestDir:   "hosting",
			TestFiles: []string{"vercel.json"},
		},
		{
			Name:      "Netlify",
			Detector:  &detectors.NetlifyDetector{},
			TestDir:   "hosting",
			TestFiles: []string{"netlify.toml"},
		},
		{
			Name:      "Heroku",
			Detector:  &detectors.HerokuDetector{},
			TestDir:   "hosting",
			TestFiles: []string{"Procfile"},
		},
		{
			Name:      "Firebase",
			Detector:  &detectors.FirebaseDetector{},
			TestDir:   "hosting",
			TestFiles: []string{"firebase.json"},
		},
		{
			Name:      "Railway",
			Detector:  &detectors.RailwayDetector{},
			TestDir:   "hosting",
			TestFiles: []string{"railway.json"},
		},
		{
			Name:      "Render",
			Detector:  &detectors.RenderDetector{},
			TestDir:   "hosting",
			TestFiles: []string{"render.yaml"},
		},
		{
			Name:      "Fly.io",
			Detector:  &detectors.FlyIODetector{},
			TestDir:   "hosting",
			TestFiles: []string{"fly.toml"},
		},
		{
			Name:      "WordPress",
			Detector:  &detectors.WordPressDetector{},
			TestDir:   "wordpress",
			TestFiles: []string{"wp-config.php", "style.css"},
		},
		{
			Name:      "Container Registry",
			Detector:  &detectors.ContainerRegistryDetector{},
			TestDir:   "container_registry",
			TestFiles: []string{"Dockerfile", ".gitlab-ci.yml", ".github/workflows/docker.yml", "Makefile", "package.json"},
		},
		{
			Name:      "I18n/Translation Services",
			Detector:  &detectors.I18nDetector{},
			TestDir:   "i18n",
			TestFiles: []string{"package.json", "crowdin.yml", ".lokaliserc", ".phraseapp.yml", ".tx/config", "tolgee.json", "locales/en/common.json"},
		},
		{
			Name:      "AI/ML Services",
			Detector:  &detectors.AIServicesDetector{},
			TestDir:   "ai_services",
			TestFiles: []string{"package.json", "requirements.txt", "env_example", "src/ai_client.py", "src/ai_client.js"},
		},
		{
			Name:      "Search Services",
			Detector:  &detectors.SearchServicesDetector{},
			TestDir:   "search_services",
			TestFiles: []string{"package.json", "requirements.txt", "env_example"},
		},
		{
			Name:      "Maps & Location Services",
			Detector:  &detectors.MapsServicesDetector{},
			TestDir:   "maps_services",
			TestFiles: []string{"package.json", "requirements.txt", "env_example"},
		},
		{
			Name:      "Push Notifications",
			Detector:  &detectors.PushNotificationsDetector{},
			TestDir:   "push_notifications",
			TestFiles: []string{"package.json", "requirements.txt", "env_example"},
		},
	}

	// Get current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	fmt.Printf("Running tests from: %s\n\n", currentDir)

	// Test each detector
	for _, config := range testConfigs {
		fmt.Printf("=== Testing %s Detector ===\n", config.Name)
		testDetector(config, currentDir)
		fmt.Println()
	}
}

func testDetector(config TestConfig, baseDir string) {
	// Change to test directory
	testPath := filepath.Join(baseDir, config.TestDir)
	err := os.Chdir(testPath)
	if err != nil {
		fmt.Printf("Error changing to test directory %s: %v\n", testPath, err)
		return
	}

	fmt.Printf("Test directory: %s\n", testPath)
	fmt.Printf("Detector: %s\n", config.Detector.Name())
	fmt.Printf("Description: %s\n", config.Detector.Description())

	// Check if detector should run
	shouldRun := config.Detector.ShouldRun()
	fmt.Printf("Should run: %v\n", shouldRun)

	if !shouldRun {
		fmt.Printf("Detector should not run in this directory\n")
		// Change back to base directory
		os.Chdir(baseDir)
		return
	}

	// Run detection
	results, err := config.Detector.Detect()
	if err != nil {
		fmt.Printf("Error during detection: %v\n", err)
		// Change back to base directory
		os.Chdir(baseDir)
		return
	}

	// Display results
	fmt.Printf("Found %d services:\n", len(results))
	for i, result := range results {
		fmt.Printf("  %d. %s: %s (confidence: %.2f)\n",
			i+1, result.Key, result.Value, result.Confidence)
		fmt.Printf("     %s\n", result.Description)
	}

	// Change back to base directory
	os.Chdir(baseDir)
}