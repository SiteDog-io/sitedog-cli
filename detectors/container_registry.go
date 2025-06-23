package detectors

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// ContainerRegistryDetector detects container registries used in projects
type ContainerRegistryDetector struct{}

func (c *ContainerRegistryDetector) Name() string {
	return "container-registry"
}

func (c *ContainerRegistryDetector) Description() string {
	return "Container Registry services detector"
}

func (c *ContainerRegistryDetector) ShouldRun() bool {
	// Check for common container-related files
	containerFiles := []string{
		"Dockerfile",
		"docker-compose.yml",
		"docker-compose.yaml",
		".dockerignore",
		"skaffold.yaml",
		"k8s.yaml",
		"deployment.yaml",
		"deployment.yml",
	}

	for _, file := range containerFiles {
		if _, err := os.Stat(file); err == nil {
			return true
		}
	}

	// Check for Kubernetes directory
	if _, err := os.Stat("k8s"); err == nil {
		return true
	}

	// Check for .kube directory
	if _, err := os.Stat(".kube"); err == nil {
		return true
	}

	// Check for helm charts
	if _, err := os.Stat("Chart.yaml"); err == nil {
		return true
	}

	// Check for CI/CD files that might contain registry info
	ciFiles := []string{
		".github/workflows",
		".gitlab-ci.yml",
		".circleci/config.yml",
		"azure-pipelines.yml",
		"Jenkinsfile",
		"bitbucket-pipelines.yml",
	}

	for _, file := range ciFiles {
		if _, err := os.Stat(file); err == nil {
			return true
		}
	}

	return false
}

func (c *ContainerRegistryDetector) Detect() ([]*DetectionResult, error) {
	var results []*DetectionResult

	// Read all relevant files to detect container registries where we PUSH images
	var projectContent strings.Builder

	// Files to check for registry references - focus on CI/CD and deployment configs
	files := []string{
		".gitlab-ci.yml",
		"azure-pipelines.yml",
		"Jenkinsfile",
		"bitbucket-pipelines.yml",
		"skaffold.yaml",
		"Chart.yaml",
		"values.yaml",
		"docker-compose.yml",
		"docker-compose.yaml",
		"package.json", // for scripts that might push images
		"Makefile",     // for build scripts
		".env",         // for registry credentials
		".env.example",
		"deploy.sh",
		"build.sh",
	}

	// Also check CI/CD workflow files
	if _, err := os.Stat(".github/workflows"); err == nil {
		filepath.Walk(".github/workflows", func(path string, info os.FileInfo, err error) error {
			if err == nil && strings.HasSuffix(info.Name(), ".yml") || strings.HasSuffix(info.Name(), ".yaml") {
				if data, readErr := ioutil.ReadFile(path); readErr == nil {
					projectContent.WriteString(strings.ToLower(string(data)))
				}
			}
			return nil
		})
	}

	// Check Kubernetes manifests
	k8sDirs := []string{"k8s", "kubernetes", "manifests", ".kube"}
	for _, dir := range k8sDirs {
		if _, err := os.Stat(dir); err == nil {
			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err == nil && (strings.HasSuffix(info.Name(), ".yml") || strings.HasSuffix(info.Name(), ".yaml")) {
					if data, readErr := ioutil.ReadFile(path); readErr == nil {
						projectContent.WriteString(strings.ToLower(string(data)))
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

	// Define container registry services with patterns focused on PUSH/DEPLOY operations
	registries := map[string]map[string]interface{}{
		"docker_hub": {
			"patterns": []string{
				"docker push", "docker.io/", "hub.docker.com/",
				"DOCKER_USERNAME", "DOCKER_PASSWORD", "dockerhub",
				"registry-1.docker.io", "index.docker.io",
			},
			"name": "Docker Hub",
			"url":  "https://hub.docker.com",
			"key":  "container_registry",
		},

		"ghcr": {
			"patterns": []string{
				"ghcr.io/", "docker.pkg.github.com/",
				"github.actor", "github.token", "CR_PAT",
				"push ghcr", "login ghcr.io",
			},
			"name": "GitHub Container Registry",
			"url":  "https://github.com/settings/packages",
			"key":  "container_registry",
		},

				"gitlab_registry": {
			"patterns": []string{
				"registry.gitlab.com/", "ci_registry", "ci_registry_image",
				"gitlab-ci-token", "ci_registry_password", "ci_registry_user",
				"docker login.*ci_registry", "push.*ci_registry",
			},
			"name": "GitLab Container Registry",
			"url":  "https://gitlab.com/-/packages/container_registry",
			"key":  "container_registry",
		},

		"gcr": {
			"patterns": []string{
				"gcr.io/", "eu.gcr.io/", "us.gcr.io/", "asia.gcr.io/",
				"gcloud docker", "docker push gcr.io",
				"GCR_", "GOOGLE_APPLICATION_CREDENTIALS",
			},
			"name": "Google Container Registry",
			"url":  "https://console.cloud.google.com/gcr",
			"key":  "container_registry",
		},

		"artifact_registry": {
			"patterns": []string{
				"pkg.dev/", "docker push.*pkg.dev",
				"gcloud auth configure-docker.*pkg.dev",
				"artifact-registry", "artifactregistry.googleapis.com",
			},
			"name": "Google Artifact Registry",
			"url":  "https://console.cloud.google.com/artifacts",
			"key":  "container_registry",
		},

		"ecr": {
			"patterns": []string{
				"amazonaws.com/", ".dkr.ecr.", "ecr get-login",
				"aws ecr", "ECR_REGISTRY", "AWS_ACCOUNT_ID",
				"docker push.*ecr", "ecr:GetAuthorizationToken",
			},
			"name": "Amazon Elastic Container Registry",
			"url":  "https://console.aws.amazon.com/ecr",
			"key":  "container_registry",
		},

		"acr": {
			"patterns": []string{
				"azurecr.io/", "az acr", "ACR_",
				"docker push.*azurecr.io", "az acr login",
				"containerregistry.azure.com", "ACR_REGISTRY",
			},
			"name": "Azure Container Registry",
			"url":  "https://portal.azure.com/#view/HubsExtension/BrowseResource/resourceType/Microsoft.ContainerRegistry%2Fregistries",
			"key":  "container_registry",
		},

		"quay": {
			"patterns": []string{
				"quay.io/", "quay.redhat.com/",
				"docker push quay.io", "QUAY_USERNAME", "QUAY_PASSWORD",
				"quay login", "push.*quay.io",
			},
			"name": "Red Hat Quay",
			"url":  "https://quay.io/repository",
			"key":  "container_registry",
		},

		"digitalocean_registry": {
			"patterns": []string{
				"registry.digitalocean.com/", "docr.io/",
				"doctl registry", "DIGITALOCEAN_ACCESS_TOKEN",
				"docker push.*registry.digitalocean.com",
			},
			"name": "DigitalOcean Container Registry",
			"url":  "https://cloud.digitalocean.com/registry",
			"key":  "container_registry",
		},

		"heroku_registry": {
			"patterns": []string{
				"registry.heroku.com/", "heroku container:push",
				"heroku container:release", "HEROKU_API_KEY",
				"docker push registry.heroku.com",
			},
			"name": "Heroku Container Registry",
			"url":  "https://dashboard.heroku.com",
			"key":  "container_registry",
		},

		"harbor": {
			"patterns": []string{
				"harbor/", "goharbor.io/", "harbor-core",
				"docker push.*harbor", "HARBOR_USERNAME", "HARBOR_PASSWORD",
			},
			"name": "Harbor Registry",
			"url":  "https://goharbor.io",
			"key":  "container_registry",
		},

		"jfrog_artifactory": {
			"patterns": []string{
				"artifactory/", "jfrog.io/", "jfrog.com/",
				"docker push.*artifactory", "JFROG_", "RT_",
				"jfrog rt docker-push",
			},
			"name": "JFrog Artifactory",
			"url":  "https://jfrog.com/artifactory",
			"key":  "container_registry",
		},
	}

	// Check for specific registries in order of popularity
	registryOrder := []string{
		"docker_hub", "ghcr", "gitlab_registry", "gcr", "artifact_registry",
		"ecr", "acr", "quay", "digitalocean_registry", "heroku_registry",
		"harbor", "jfrog_artifactory",
	}

	for _, registryKey := range registryOrder {
		registryInfo := registries[registryKey]
		patterns := registryInfo["patterns"].([]string)

		for _, pattern := range patterns {
			if strings.Contains(content, pattern) {
				results = append(results, &DetectionResult{
					Key:         registryInfo["key"].(string),
					Value:       registryInfo["url"].(string),
					Description: registryInfo["name"].(string) + " detected in container configuration",
					Confidence:  0.90,
				})
				break // Only add each registry once
			}
		}
	}

	return results, nil
}