package detectors

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// AIServicesDetector detects AI/ML services used in projects
type AIServicesDetector struct{}

func (a *AIServicesDetector) Name() string {
	return "ai-services"
}

func (a *AIServicesDetector) Description() string {
	return "AI/ML services and APIs detector"
}

func (a *AIServicesDetector) ShouldRun() bool {
	// Check for common files that might contain AI integrations
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
			// Quick check for AI-related keywords
			aiKeywords := []string{
				"openai", "anthropic", "claude", "gpt", "huggingface", "replicate",
				"cohere", "stability", "midjourney", "dall-e", "whisper", "langchain",
			}
			for _, keyword := range aiKeywords {
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

func (a *AIServicesDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Files to check for AI service references - only configuration files
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

	// Deep search mode - also check source code directories for AI imports/usage
	if DeepSearchMode {
		srcDirs := []string{"src", "lib", "app", "components", "pages", "api", "services", "utils", "models"}
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

	// Define AI services with their patterns and dashboards
	services := map[string]map[string]interface{}{
		"openai": {
			"patterns": []string{
				"openai", "sk-", "openai_api_key", "openai_api_token",
				"gpt-3", "gpt-4", "gpt-3.5", "dall-e", "whisper", "davinci", "curie",
				"api.openai.com", "chat.openai.com", "platform.openai.com",
				"from openai import", "import openai", "OpenAI()",
			},
			"name": "OpenAI API",
			"url":  "https://platform.openai.com/account/billing",
			"key":  "ai_service",
		},

		"anthropic": {
			"patterns": []string{
				"anthropic", "claude", "anthropic_api_key", "anthropic_api_token",
				"claude-1", "claude-2", "claude-3", "claude-instant",
				"api.anthropic.com", "console.anthropic.com",
				"from anthropic import", "import anthropic", "@anthropic-ai/sdk",
			},
			"name": "Anthropic Claude API",
			"url":  "https://console.anthropic.com/account/billing",
			"key":  "ai_service",
		},

		"huggingface": {
			"patterns": []string{
				"huggingface", "hf_", "hugging face", "transformers",
				"huggingface_hub", "hf_api_token", "hf_token",
				"api-inference.huggingface.co", "huggingface.co/models",
				"from transformers import", "import transformers", "pipeline(",
				"AutoTokenizer", "AutoModel", "AutoModelForCausalLM",
			},
			"name": "Hugging Face",
			"url":  "https://huggingface.co/settings/billing",
			"key":  "ai_service",
		},

		"replicate": {
			"patterns": []string{
				"replicate", "replicate_api_token", "replicate.com",
				"r8.im/", "api.replicate.com",
				"from replicate import", "import replicate", "replicate.run",
			},
			"name": "Replicate",
			"url":  "https://replicate.com/account/billing",
			"key":  "ai_service",
		},

		"cohere": {
			"patterns": []string{
				"cohere", "cohere_api_key", "cohere_api_token", "co.api_key",
				"api.cohere.ai", "dashboard.cohere.ai",
				"from cohere import", "import cohere", "cohere.Client",
			},
			"name": "Cohere",
			"url":  "https://dashboard.cohere.com/billing",
			"key":  "ai_service",
		},

		"stability": {
			"patterns": []string{
				"stability", "stability_api_key", "stable diffusion", "stablediffusion",
				"api.stability.ai", "platform.stability.ai",
				"stability-sdk", "stabilityai",
			},
			"name": "Stability AI",
			"url":  "https://platform.stability.ai/account/billing",
			"key":  "ai_service",
		},

		"pinecone": {
			"patterns": []string{
				"pinecone", "pinecone_api_key", "pinecone_environment",
				"api.pinecone.io", "app.pinecone.io",
				"from pinecone import", "import pinecone", "pinecone.init",
			},
			"name": "Pinecone Vector Database",
			"url":  "https://app.pinecone.io/billing",
			"key":  "ai_service",
		},

		"elevenlabs": {
			"patterns": []string{
				"elevenlabs", "eleven labs", "elevenlabs_api_key",
				"api.elevenlabs.io", "elevenlabs.io",
				"elevenlabs-python", "elevenlabslib",
			},
			"name": "ElevenLabs",
			"url":  "https://elevenlabs.io/subscription",
			"key":  "ai_service",
		},

		"langchain": {
			"patterns": []string{
				"langchain", "langchain_api_key", "langsmith",
				"api.smith.langchain.com", "smith.langchain.com",
				"from langchain import", "import langchain", "@langchain/",
			},
			"name": "LangChain/LangSmith",
			"url":  "https://smith.langchain.com/settings",
			"key":  "ai_service",
		},

		"together": {
			"patterns": []string{
				"together_api_key", "together.ai", "together_ai",
				"api.together.xyz", "together-python", "together-ai",
				"from together import", "import together", "@together-ai/",
				"TOGETHER_API_KEY", "TOGETHER_AI_API_KEY",
			},
			"name": "Together AI",
			"url":  "https://api.together.xyz/settings/billing",
			"key":  "ai_service",
		},

		"perplexity": {
			"patterns": []string{
				"perplexity", "perplexity_api_key", "pplx-",
				"api.perplexity.ai", "perplexity.ai",
			},
			"name": "Perplexity AI",
			"url":  "https://www.perplexity.ai/settings/api",
			"key":  "ai_service",
		},

		"mistral": {
			"patterns": []string{
				"mistral", "mistral_api_key", "mistralai",
				"api.mistral.ai", "console.mistral.ai",
				"@mistralai/", "mistralai/client",
			},
			"name": "Mistral AI",
			"url":  "https://console.mistral.ai/billing",
			"key":  "ai_service",
		},

		"groq": {
			"patterns": []string{
				"groq", "groq_api_key", "groq-sdk",
				"api.groq.com", "console.groq.com",
				"from groq import", "import groq",
			},
			"name": "Groq",
			"url":  "https://console.groq.com/settings/billing",
			"key":  "ai_service",
		},
	}

	// Check for specific AI services in order of popularity
	serviceOrder := []string{
		"openai", "anthropic", "huggingface", "langchain", "pinecone",
		"replicate", "cohere", "stability", "elevenlabs", "together",
		"perplexity", "mistral", "groq",
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