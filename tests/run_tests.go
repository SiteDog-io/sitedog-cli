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