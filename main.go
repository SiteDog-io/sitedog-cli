package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/go-shiori/obelisk"
	"gopkg.in/yaml.v2"
)

const (
	defaultConfigPath  = "./sitedog.yml"
	defaultTemplate    = "demo.html.tpl"
	defaultPort        = 8081
	globalTemplatePath = ".sitedog/demo.html.tpl"
	authFilePath       = ".sitedog/auth"
	apiBaseURL         = "https://app.sitedog.io" // Change to your actual API URL
	Version            = "v0.2.1"
	exampleConfig      = `# Describe your project with a free key-value format, think simple.
#
# Random sample:
registrar: gandi # registrar service
dns: Route 53 # dns service
hosting: https://carrd.com # hosting service
mail: zoho # mail service
`
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}
	switch os.Args[1] {
	case "init":
		handleInit()
	case "live":
		handleLive()
	case "push":
		handlePush()
	case "render":
		handleRender()
	case "scan":
		handleScan()
	case "version":
		fmt.Println("sitedog version", Version)
	case "help":
		showHelp()
	default:
		fmt.Println("Unknown command:", os.Args[1])
		showHelp()
	}
}

func showHelp() {
	fmt.Println(`Usage: sitedog <command>

Commands:
  init    Create sitedog.yml configuration file
  live    Start live server with preview
  push    Push configuration to cloud
  render  Render template to HTML
  scan    Scan git repository and suggest adding repo URL to config
  version Print version
  help    Show this help message

Options for init:
  --config PATH    Path to config file (default: ./sitedog.yml)

Options for live:
  --config PATH    Path to config file (default: ./sitedog.yml)
  --port PORT      Port to run server on (default: 8081)

Options for push:
  --config PATH    Path to config file (default: ./sitedog.yml)
  --title TITLE    Configuration title (default: current directory name)
  --remote URL     Custom API base URL (e.g., localhost:3000, api.example.com)
  SITEDOG_TOKEN    Environment variable for authentication token

Options for render:
  --config PATH    Path to config file (default: ./sitedog.yml)
  --output PATH    Path to output HTML file (default: sitedog.html)

Options for scan:
  --config PATH    Path to config file (default: ./sitedog.yml)

Examples:
  sitedog init --config my-config.yml
  sitedog live --port 3030
  sitedog push --title my-project
  sitedog push --remote localhost:3000 --title my-project
  sitedog push --remote api.example.com --title my-project
  sitedog push --remote https://api.example2.com --title my-project
  SITEDOG_TOKEN=your_token sitedog push --title my-project
  sitedog render --output index.html
  sitedog scan --config my-config.yml`)
}

func handleInit() {
	configPath := flag.NewFlagSet("init", flag.ExitOnError)
	configFile := configPath.String("config", defaultConfigPath, "Path to config file")
	configPath.Parse(os.Args[2:])
	if _, err := os.Stat(*configFile); err == nil {
		fmt.Println("Error:", *configFile, "already exists")
		os.Exit(1)
	}
	if err := ioutil.WriteFile(*configFile, []byte(exampleConfig), 0644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created", *configFile, "configuration file")
}

func startServer(configFile *string, port int) (*http.Server, string) {
	// Handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		config, err := ioutil.ReadFile(*configFile)
		if err != nil {
			http.Error(w, "Error reading config", http.StatusInternalServerError)
			return
		}

		faviconCache := getFaviconCache(config)
		tmpl, _ := ioutil.ReadFile(findTemplate())
		tmpl = bytes.Replace(tmpl, []byte("{{CONFIG}}"), config, -1)
		tmpl = bytes.Replace(tmpl, []byte("{{FAVICON_CACHE}}"), faviconCache, -1)
		w.Header().Set("Content-Type", "text/html")
		w.Write(tmpl)
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		config, err := ioutil.ReadFile(*configFile)
		if err != nil {
			http.Error(w, "Error reading config", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/yaml")
		w.Write(config)
	})

	// Start the server
	addr := fmt.Sprintf(":%d", port)
	server := &http.Server{
		Addr: addr,
	}

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Wait for server to start
	time.Sleep(1 * time.Second)

	return server, addr
}

func handleLive() {
	liveFlags := flag.NewFlagSet("live", flag.ExitOnError)
	configFile := liveFlags.String("config", defaultConfigPath, "Path to config file")
	port := liveFlags.Int("port", defaultPort, "Port to run server on")
	liveFlags.Parse(os.Args[2:])

	if _, err := os.Stat(*configFile); err != nil {
		fmt.Println("Error:", *configFile, "not found. Run 'sitedog init' first.")
		os.Exit(1)
	}

	server, addr := startServer(configFile, *port)
	url := "http://localhost" + addr

	go func() {
		time.Sleep(500 * time.Millisecond)
		openBrowser(url)
	}()
	fmt.Println("Starting live server at", url)
	fmt.Println("Press Ctrl+C to stop")

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	// Gracefully shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

func findTemplate() string {
	local := filepath.Join(".", defaultTemplate)
	if _, err := os.Stat(local); err == nil {
		return local
	}
	usr, _ := user.Current()
	global := filepath.Join(usr.HomeDir, globalTemplatePath)
	if _, err := os.Stat(global); err == nil {
		return global
	}
	fmt.Println("Template not found.")
	os.Exit(1)
	return ""
}

func openBrowser(url string) {
	var cmd string
	var args []string
	switch {
	case strings.Contains(strings.ToLower(os.Getenv("OS")), "windows"):
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case strings.Contains(strings.ToLower(os.Getenv("OSTYPE")), "darwin"):
		cmd = "open"
		args = []string{url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}
	exec.Command(cmd, args...).Start()
}

func handlePush() {
	pushFlags := flag.NewFlagSet("push", flag.ExitOnError)
	configFile := pushFlags.String("config", defaultConfigPath, "Path to config file")
	configName := pushFlags.String("title", "", "Configuration title")
	remoteURL := pushFlags.String("remote", "", "Custom API base URL (e.g., localhost:3000, api.example.com)")
	pushFlags.Parse(os.Args[2:])

	if _, err := os.Stat(*configFile); err != nil {
		fmt.Println("Error:", *configFile, "not found. Run 'sitedog init' first.")
		os.Exit(1)
	}

	// Get authorization token
	token, err := getAuthToken()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Read configuration
	config, err := ioutil.ReadFile(*configFile)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}

	// Get configuration name from directory name if not specified
	if *configName == "" {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			os.Exit(1)
		}
		*configName = filepath.Base(dir)
	}

	// Determine API base URL
	apiURL := apiBaseURL
	if *remoteURL != "" {
		// Add protocol if not specified
		if !strings.HasPrefix(*remoteURL, "http://") && !strings.HasPrefix(*remoteURL, "https://") {
			apiURL = "http://" + *remoteURL
		} else {
			apiURL = *remoteURL
		}
	}

	// Send configuration to server
	err = pushConfig(token, *configName, string(config), apiURL)
	if err != nil {
		fmt.Println("Error pushing config:", err)
		os.Exit(1)
	}

	fmt.Printf("Configuration '%s' pushed successfully to %s!\n", *configName, apiURL)
}

func getAuthToken() (string, error) {
	// First check for environment variable
	if token := os.Getenv("SITEDOG_TOKEN"); token != "" {
		return strings.TrimSpace(token), nil
	}

	// Fall back to file-based authentication
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("error getting current user: %v", err)
	}

	authFile := filepath.Join(usr.HomeDir, authFilePath)
	if _, err := os.Stat(authFile); err == nil {
		// If file exists, read the token
		token, err := ioutil.ReadFile(authFile)
		if err != nil {
			return "", fmt.Errorf("error reading auth file: %v", err)
		}
		return strings.TrimSpace(string(token)), nil
	}

	// If file doesn't exist, request authorization
	fmt.Print("Email: ")
	var email string
	fmt.Scanln(&email)

	fmt.Print("Password: ")
	var password string
	fmt.Scanln(&password)

	// Create authorization request
	reqBody, err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	resp, err := http.Post(apiBaseURL+"/cli/auth", "application/json", strings.NewReader(string(reqBody)))
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authentication failed: %s", resp.Status)
	}

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	// Create .sitedog directory if it doesn't exist
	authDir := filepath.Dir(authFile)
	if err := os.MkdirAll(authDir, 0700); err != nil {
		return "", fmt.Errorf("error creating auth directory: %v", err)
	}

	// Save the token
	if err := ioutil.WriteFile(authFile, []byte(result.Token), 0600); err != nil {
		return "", fmt.Errorf("error saving token: %v", err)
	}

	return result.Token, nil
}

func pushConfig(token, name, content, apiURL string) error {
	reqBody, err := json.Marshal(map[string]string{
		"name":    name,
		"content": content,
	})
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL+"/cli/push", strings.NewReader(string(reqBody)))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("push failed: %s", resp.Status)
	}

	return nil
}

func spinner(stopSpinner chan bool, message string) {
	spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-stopSpinner:
			fmt.Print("\r")
			return
		default:
			fmt.Printf("\r%s %s", spinner[i], message)
			i = (i + 1) % len(spinner)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func handleRender() {
	// Start loading indicator
	stopSpinner := make(chan bool)
	go spinner(stopSpinner, "Rendering...")

	renderFlags := flag.NewFlagSet("render", flag.ExitOnError)
	configFile := renderFlags.String("config", defaultConfigPath, "Path to config file")
	outputFile := renderFlags.String("output", "sitedog.html", "Path to output HTML file")
	renderFlags.Parse(os.Args[2:])

	if _, err := os.Stat(*configFile); err != nil {
		stopSpinner <- true
		fmt.Println("Error:", *configFile, "not found. Run 'sitedog init' first.")
		os.Exit(1)
	}

	port := 34324
	server, addr := startServer(configFile, port)
	url := "http://localhost" + addr

	// Check server availability
	resp, err := http.Get(url)
	if err != nil {
		stopSpinner <- true
		fmt.Println("Error checking server:", err)
		server.Close()
		os.Exit(1)
	}
	resp.Body.Close()

	// Use Obelisk to save the page
	archiver := &obelisk.Archiver{
		EnableLog:             false,
		MaxConcurrentDownload: 10,
	}

	// Validate archiver
	archiver.Validate()

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Create request
	req := obelisk.Request{
		URL: url,
	}

	// Save the page
	html, _, err := archiver.Archive(ctx, req)
	if err != nil {
		stopSpinner <- true
		fmt.Println("\nError archiving page:", err)
		server.Close()
		os.Exit(1)
	}

	// Save result to file
	if err := ioutil.WriteFile(*outputFile, html, 0644); err != nil {
		stopSpinner <- true
		fmt.Println("Error saving file:", err)
		server.Close()
		os.Exit(1)
	}

	// Stop loading indicator
	stopSpinner <- true

	// Close server
	server.Close()

	fmt.Printf("Rendered cards saved to %s\n", *outputFile)
}

func getFaviconCache(config []byte) []byte {
	// Parse YAML config
	var configMap map[string]interface{}
	if err := yaml.Unmarshal(config, &configMap); err != nil {
		return []byte("{}")
	}

	// Create map for storing favicon cache
	faviconCache := make(map[string]string)

	// Function to extract domain from URL
	extractDomain := func(urlStr string) string {
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			return ""
		}
		return parsedURL.Hostname()
	}

	// Function for recursive value traversal
	var traverseValue func(value interface{})
	traverseValue = func(value interface{}) {
		switch v := value.(type) {
		case string:
			// Check if string is a URL
			if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
				domain := extractDomain(v)
				if domain != "" {
					// Get favicon
					faviconURL := fmt.Sprintf("https://www.google.com/s2/favicons?domain=%s&sz=64", url.QueryEscape(domain))
					resp, err := http.Get(faviconURL)
					if err != nil {
						return
					}
					defer resp.Body.Close()

					if resp.StatusCode == http.StatusOK {
						// Read favicon
						faviconData, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							return
						}

						// Convert to base64
						base64Data := base64.StdEncoding.EncodeToString(faviconData)
						// Get content type from response headers
						contentType := resp.Header.Get("Content-Type")
						if contentType == "" {
							contentType = "image/png" // fallback to png if no content type specified
						}
						dataURL := fmt.Sprintf("data:%s;base64,%s", contentType, base64Data)
						faviconCache[v] = dataURL
					}
				}
			}
		case map[interface{}]interface{}:
			for _, val := range v {
				traverseValue(val)
			}
		case []interface{}:
			for _, val := range v {
				traverseValue(val)
			}
		case map[string]interface{}:
			for _, val := range v {
				traverseValue(val)
			}
		}
	}

	// Traverse all values in config
	traverseValue(configMap)

	// Convert map to JSON
	jsonData, err := json.Marshal(faviconCache)
	if err != nil {
		return []byte("{}")
	}

	return jsonData
}

// handleScan scans the current git repository and suggests adding repo URL to config
func handleScan() {
	scanFlags := flag.NewFlagSet("scan", flag.ExitOnError)
	configFile := scanFlags.String("config", defaultConfigPath, "Path to config file")
	scanFlags.Parse(os.Args[2:])

	// Check if we're in a git repository
	if !isGitRepository() {
		fmt.Println("Error: Not in a git repository. Please run this command from within a git repository.")
		os.Exit(1)
	}

	// Get git remote origin URL
	originURL, err := getGitOriginURL()
	if err != nil {
		fmt.Println("Error getting git origin URL:", err)
		os.Exit(1)
	}

	if originURL == "" {
		fmt.Println("No git remote origin found. Please add a remote origin to your repository.")
		os.Exit(1)
	}

	// Convert SSH URL to HTTPS if needed
	repoURL := convertToHTTPSURL(originURL)
	fmt.Printf("Found git repository: %s\n", repoURL)

	// Check if config file exists
	if _, err := os.Stat(*configFile); err != nil {
		fmt.Printf("Config file %s not found. Run 'sitedog init' first.\n", *configFile)
		os.Exit(1)
	}

	// Read existing config
	configData, err := ioutil.ReadFile(*configFile)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}

	// Parse YAML config
	var config map[string]interface{}
	if err := yaml.Unmarshal(configData, &config); err != nil {
		fmt.Println("Error parsing config file:", err)
		os.Exit(1)
	}

	// Check if repo key already exists
	if checkRepoExists(config, repoURL) {
		fmt.Println("Repository URL already exists in config file.")
		return
	}

	// Ask user if they want to add the repo
	fmt.Printf("Do you want to add 'repo: %s' to your config? (y/N): ", repoURL)
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(strings.TrimSpace(response)) != "y" {
		fmt.Println("Cancelled.")
		return
	}

	// Add repo to config
	if err := addRepoToConfig(*configFile, repoURL); err != nil {
		fmt.Println("Error updating config file:", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully added 'repo: %s' to %s\n", repoURL, *configFile)
}

// isGitRepository checks if current directory is a git repository
func isGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// getGitOriginURL gets the origin URL from git remote
func getGitOriginURL() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// convertToHTTPSURL converts SSH git URLs to HTTPS URLs
func convertToHTTPSURL(gitURL string) string {
	// Pattern for SSH URLs like git@github.com:user/repo.git
	sshPattern := regexp.MustCompile(`^git@([^:]+):(.+)\.git$`)
	if matches := sshPattern.FindStringSubmatch(gitURL); len(matches) == 3 {
		return fmt.Sprintf("https://%s/%s", matches[1], matches[2])
	}

	// Pattern for SSH URLs like git@github.com:user/repo (without .git)
	sshPatternNoGit := regexp.MustCompile(`^git@([^:]+):(.+)$`)
	if matches := sshPatternNoGit.FindStringSubmatch(gitURL); len(matches) == 3 {
		return fmt.Sprintf("https://%s/%s", matches[1], matches[2])
	}

	// If it's already HTTPS or HTTP, remove .git suffix if present
	if strings.HasPrefix(gitURL, "http://") || strings.HasPrefix(gitURL, "https://") {
		return strings.TrimSuffix(gitURL, ".git")
	}

	// Return as is if we can't parse it
	return gitURL
}

// checkRepoExists checks if repo URL already exists in config
func checkRepoExists(config map[string]interface{}, repoURL string) bool {
	// Check if there's a direct repo key
	if repo, exists := config["repo"]; exists {
		if repoStr, ok := repo.(string); ok && repoStr == repoURL {
			return true
		}
	}

	// Check nested objects for repo key
	for _, value := range config {
		if nestedMap, ok := value.(map[interface{}]interface{}); ok {
			if repo, exists := nestedMap["repo"]; exists {
				if repoStr, ok := repo.(string); ok && repoStr == repoURL {
					return true
				}
			}
		}
		if nestedMap, ok := value.(map[string]interface{}); ok {
			if repo, exists := nestedMap["repo"]; exists {
				if repoStr, ok := repo.(string); ok && repoStr == repoURL {
					return true
				}
			}
		}
	}

	return false
}

// addRepoToConfig adds repo URL to the config file while preserving order
func addRepoToConfig(configFile, repoURL string) error {
	// Read existing config as text
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	configText := string(configData)
	lines := strings.Split(configText, "\n")

	// Parse YAML to understand structure
	var config map[string]interface{}
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return err
	}

	// Find the first root key (project name)
	var projectKey string
	if len(config) > 0 {
		// Find the first key in the original text order
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if strings.Contains(trimmedLine, ":") && !strings.HasPrefix(trimmedLine, "#") && !strings.HasPrefix(trimmedLine, " ") && !strings.HasPrefix(trimmedLine, "\t") {
				projectKey = strings.Split(trimmedLine, ":")[0]
				break
			}
		}
	}

	if projectKey == "" {
		// Create a project key based on current directory
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		projectKey = filepath.Base(dir)
		// Add new project section
		repoLine := fmt.Sprintf("  repo: %s", repoURL)
		newConfig := fmt.Sprintf("%s:\n%s\n", projectKey, repoLine)
		if configText != "" {
			newConfig = newConfig + "\n" + configText
		}
		return ioutil.WriteFile(configFile, []byte(newConfig), 0644)
	}

	// Find where to insert the repo line
	repoLine := fmt.Sprintf("  repo: %s", repoURL)
	insertIndex := -1
	inProjectSection := false
	indentLevel := 0

	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Check if this is the start of our project section
		if strings.HasPrefix(line, projectKey+":") {
			inProjectSection = true
			indentLevel = getIndentLevel(line) + 2 // Next level indent
			continue
		}

		// If we're in the project section
		if inProjectSection {
			currentIndent := getIndentLevel(line)

			// If we hit a line with same or less indent than project key, we've left the section
			if trimmedLine != "" && currentIndent <= indentLevel-2 {
				insertIndex = i
				break
			}

			// If this is the last line and we haven't found a place to insert
			if i == len(lines)-1 {
				insertIndex = len(lines)
				break
			}
		}
	}

	// Insert the repo line
	if insertIndex >= 0 {
		newLines := make([]string, 0, len(lines)+1)
		newLines = append(newLines, lines[:insertIndex]...)
		newLines = append(newLines, repoLine)
		newLines = append(newLines, lines[insertIndex:]...)
		return ioutil.WriteFile(configFile, []byte(strings.Join(newLines, "\n")), 0644)
	}

	// Fallback: add at the end of the project section
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if strings.TrimSpace(line) != "" {
			newLines := make([]string, 0, len(lines)+1)
			newLines = append(newLines, lines[:i+1]...)
			newLines = append(newLines, repoLine)
			newLines = append(newLines, lines[i+1:]...)
			return ioutil.WriteFile(configFile, []byte(strings.Join(newLines, "\n")), 0644)
		}
	}

	return fmt.Errorf("could not find appropriate place to insert repo")
}

// getIndentLevel returns the number of spaces at the beginning of a line
func getIndentLevel(line string) int {
	count := 0
	for _, char := range line {
		if char == ' ' {
			count++
		} else if char == '\t' {
			count += 2 // Count tab as 2 spaces
		} else {
			break
		}
	}
	return count
}
