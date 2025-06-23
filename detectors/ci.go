package detectors

import (
	"io/ioutil"
	"os"
	"strings"
)

// GitLabCIDetector detects GitLab CI configuration
type GitLabCIDetector struct{}

func (g *GitLabCIDetector) Name() string {
	return "gitlab-ci"
}

func (g *GitLabCIDetector) Description() string {
	return "GitLab CI detector"
}

func (g *GitLabCIDetector) ShouldRun() bool {
	_, err := os.Stat(".gitlab-ci.yml")
	return err == nil
}

func (g *GitLabCIDetector) Detect() ([]*DetectionResult, error) {
	// Get git origin URL to construct CI link
	originURL, err := getGitOriginURL()
	if err == nil && originURL != "" {
		repoURL := convertToHTTPSURL(originURL)
		if strings.Contains(repoURL, "gitlab.com") {
			ciURL := strings.TrimSuffix(repoURL, "/") + "/-/pipelines"
			return []*DetectionResult{{
				Key:         "ci",
				Value:       ciURL,
				Description: "GitLab CI pipelines URL",
				Confidence:  1.0,
			}}, nil
		}
	}

	// Fallback: just indicate CI is used
	return []*DetectionResult{{
		Key:         "ci",
		Value:       "gitlab-ci",
		Description: "GitLab CI configuration detected",
		Confidence:  0.8,
	}}, nil
}

// GitHubActionsDetector detects GitHub Actions configuration
type GitHubActionsDetector struct{}

func (g *GitHubActionsDetector) Name() string {
	return "github-actions"
}

func (g *GitHubActionsDetector) Description() string {
	return "GitHub Actions detector"
}

func (g *GitHubActionsDetector) ShouldRun() bool {
	// Check if .github/workflows directory exists and has YAML files
	workflowsDir := ".github/workflows"
	if _, err := os.Stat(workflowsDir); err != nil {
		return false
	}

	// Check for YAML files in workflows directory
	files, err := ioutil.ReadDir(workflowsDir)
	if err != nil {
		return false
	}

	for _, file := range files {
		if !file.IsDir() {
			name := strings.ToLower(file.Name())
			if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
				return true
			}
		}
	}

	return false
}

func (g *GitHubActionsDetector) Detect() ([]*DetectionResult, error) {
	// Get git origin URL to construct CI link
	originURL, err := getGitOriginURL()
	if err == nil && originURL != "" {
		repoURL := convertToHTTPSURL(originURL)
		if strings.Contains(repoURL, "github.com") {
			ciURL := strings.TrimSuffix(repoURL, "/") + "/actions"
			return []*DetectionResult{{
				Key:         "ci",
				Value:       ciURL,
				Description: "GitHub Actions URL",
				Confidence:  1.0,
			}}, nil
		}
	}

	// Fallback: just indicate CI is used
	return []*DetectionResult{{
		Key:         "ci",
		Value:       "github-actions",
		Description: "GitHub Actions configuration detected",
		Confidence:  0.8,
	}}, nil
}