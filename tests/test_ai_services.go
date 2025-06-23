package main

import (
	"fmt"
	"os"
	"path/filepath"

	"./detectors"
)

func testAIServices() {
	fmt.Println("Testing AI Services detector...")

	// Create test directory structure
	testDir := "tests/ai_services"
	originalDir, _ := os.Getwd()

	// Change to test directory
	if err := os.Chdir(testDir); err != nil {
		fmt.Printf("Error changing to test directory: %v\n", err)
		return
	}

	defer func() {
		// Change back to original directory
		os.Chdir(originalDir)
	}()

	// Initialize detector
	detector := &detectors.AIServicesDetector{}

	// Test ShouldRun
	fmt.Printf("ShouldRun: %v\n", detector.ShouldRun())

	// Test detection
	results, err := detector.Detect()
	if err != nil {
		fmt.Printf("Error during detection: %v\n", err)
		return
	}

	fmt.Printf("Found %d AI services:\n", len(results))
	for _, result := range results {
		fmt.Printf("- %s: %s (confidence: %.2f)\n",
			result.Description, result.Value, result.Confidence)
	}

	// Expected services based on test files:
	// OpenAI, Anthropic, LangChain, Hugging Face, Replicate, Cohere, Groq, Mistral, Pinecone
	expectedServices := []string{
		"OpenAI API",
		"Anthropic Claude API",
		"LangChain/LangSmith",
		"Hugging Face",
		"Replicate",
		"Cohere",
		"Groq",
		"Mistral AI",
		"Pinecone Vector Database",
	}

	fmt.Printf("\nExpected services: %v\n", expectedServices)

	// Verify we found the expected services
	foundServices := make(map[string]bool)
	for _, result := range results {
		for _, expected := range expectedServices {
			if result.Description == expected+" detected in project" {
				foundServices[expected] = true
			}
		}
	}

	fmt.Printf("Successfully detected %d out of %d expected services\n",
		len(foundServices), len(expectedServices))

	for service := range foundServices {
		fmt.Printf("✓ %s\n", service)
	}
}