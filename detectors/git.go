package detectors

import (
	"os/exec"
	"regexp"
	"strings"
)

// GitInfo holds information about a git repository
type GitInfo struct {
	Host      string // e.g., "gitlab.com", "github.com"
	Owner     string // e.g., "owner-name"
	Repo      string // e.g., "repo-name"
	FullURL   string // e.g., "https://gitlab.com/owner-name/repo-name"
	IsGitLab  bool
	IsGitHub  bool
}

// GitDetector detects git repository information
type GitDetector struct{}

func (g *GitDetector) Name() string {
	return "git"
}

func (g *GitDetector) Description() string {
	return "Git repository detector"
}

func (g *GitDetector) ShouldRun() bool {
	return isGitRepository()
}

func (g *GitDetector) Detect() ([]*DetectionResult, error) {
	originURL, err := getGitOriginURL()
	if err != nil {
		return nil, err
	}

	if originURL == "" {
		return nil, nil
	}

	repoURL := convertToHTTPSURL(originURL)
	return []*DetectionResult{{
		Key:         "repo",
		Value:       repoURL,
		Description: "Git repository URL",
		Confidence:  1.0,
		DebugInfo:   "Found git remote origin: " + originURL,
	}}, nil
}

// Helper functions for git operations
func isGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

func getGitOriginURL() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func convertToHTTPSURL(gitURL string) string {
	// Pattern for SSH URLs like git@github.com:user/repo.git
	sshPattern := regexp.MustCompile(`^git@([^:]+):(.+)\.git$`)
	if matches := sshPattern.FindStringSubmatch(gitURL); len(matches) == 3 {
		return "https://" + matches[1] + "/" + matches[2]
	}

	// Pattern for SSH URLs like git@github.com:user/repo (without .git)
	sshPatternNoGit := regexp.MustCompile(`^git@([^:]+):(.+)$`)
	if matches := sshPatternNoGit.FindStringSubmatch(gitURL); len(matches) == 3 {
		return "https://" + matches[1] + "/" + matches[2]
	}

	// If it's already HTTPS or HTTP, remove .git suffix if present
	if strings.HasPrefix(gitURL, "http://") || strings.HasPrefix(gitURL, "https://") {
		return strings.TrimSuffix(gitURL, ".git")
	}

	return gitURL
}

// parseGitInfo extracts detailed repository information from git remote origin
func parseGitInfo() (*GitInfo, error) {
	originURL, err := getGitOriginURL()
	if err != nil {
		return nil, err
	}

	if originURL == "" {
		return nil, nil
	}

	return parseGitURL(originURL)
}

// parseGitURL parses various git URL formats and extracts repository information
func parseGitURL(url string) (*GitInfo, error) {
	// Remove .git suffix if present
	url = strings.TrimSuffix(url, ".git")

	var host, owner, repo string

	// Handle SSH format: git@host:owner/repo
	sshPattern := regexp.MustCompile(`^git@([^:]+):([^/]+)/(.+)$`)
	if matches := sshPattern.FindStringSubmatch(url); len(matches) == 4 {
		host = matches[1]
		owner = matches[2]
		repo = matches[3]
	} else {
		// Handle HTTPS format: https://host/owner/repo
		httpsPattern := regexp.MustCompile(`^https://([^/]+)/([^/]+)/(.+)$`)
		if matches := httpsPattern.FindStringSubmatch(url); len(matches) == 4 {
			host = matches[1]
			owner = matches[2]
			repo = matches[3]
		} else {
			// Handle HTTP format: http://host/owner/repo
			httpPattern := regexp.MustCompile(`^http://([^/]+)/([^/]+)/(.+)$`)
			if matches := httpPattern.FindStringSubmatch(url); len(matches) == 4 {
				host = matches[1]
				owner = matches[2]
				repo = matches[3]
			} else {
				// Unable to parse - return basic info
				return &GitInfo{
					FullURL: convertToHTTPSURL(url),
				}, nil
			}
		}
	}

	// Build full HTTPS URL
	fullURL := "https://" + host + "/" + owner + "/" + repo

	return &GitInfo{
		Host:     host,
		Owner:    owner,
		Repo:     repo,
		FullURL:  fullURL,
		IsGitLab: strings.Contains(host, "gitlab"),
		IsGitHub: strings.Contains(host, "github"),
	}, nil
}

// buildGitLabContainerRegistryURL builds the correct GitLab container registry URL
func buildGitLabContainerRegistryURL(gitInfo *GitInfo) string {
	if gitInfo == nil || !gitInfo.IsGitLab {
		return "https://gitlab.com/-/packages/container_registry" // fallback
	}

	// Build the correct GitLab container registry URL
	// Format: https://host/owner/repo/container_registry
	return "https://" + gitInfo.Host + "/" + gitInfo.Owner + "/" + gitInfo.Repo + "/container_registry"
}

// buildGitHubContainerRegistryURL builds the correct GitHub container registry URL
func buildGitHubContainerRegistryURL(gitInfo *GitInfo) string {
	if gitInfo == nil || !gitInfo.IsGitHub {
		return "https://github.com/settings/packages" // fallback
	}

	// Build the correct GitHub packages URL
	// Format: https://github.com/owner/repo/pkgs/container/repo
	return "https://github.com/" + gitInfo.Owner + "/" + gitInfo.Repo + "/pkgs/container/" + gitInfo.Repo
}