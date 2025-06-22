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
	"strings"
	"syscall"
	"time"

	"github.com/go-shiori/obelisk"
	"gopkg.in/yaml.v2"
	"sitedog-cli/detectors"
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
  --detector NAME  Run specific detector only (git, gitlab-ci, github-actions, gemfile)

Examples:
  sitedog init --config my-config.yml
  sitedog live --port 3030
  sitedog push --title my-project
  sitedog push --remote localhost:3000 --title my-project
  sitedog push --remote api.example.com --title my-project
  sitedog push --remote https://api.example2.com --title my-project
  SITEDOG_TOKEN=your_token sitedog push --title my-project
  sitedog render --output index.html
  sitedog scan --config my-config.yml
  sitedog scan --detector git
  sitedog scan --detector gitlab-ci
  sitedog scan --detector github-actions
  sitedog scan --detector gemfile`)
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



// handleScan runs all available detectors and suggests config additions
func handleScan() {
	scanFlags := flag.NewFlagSet("scan", flag.ExitOnError)
	configFile := scanFlags.String("config", defaultConfigPath, "Path to config file")
	detectorName := scanFlags.String("detector", "", "Run specific detector only (git, package, docker, etc.)")
	scanFlags.Parse(os.Args[2:])

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

	// Initialize all detectors
	allDetectors := detectors.GetAllDetectors()

	// Filter detectors if specific one requested
	if *detectorName != "" {
		detector := detectors.FindDetectorByName(strings.ToLower(*detectorName))
		if detector == nil {
			fmt.Printf("Unknown detector: %s\n", *detectorName)
			fmt.Printf("Available detectors: %s\n", strings.Join(detectors.GetDetectorNames(), ", "))
			os.Exit(1)
		}
		allDetectors = []detectors.Detector{detector}
	}

	// Run detectors
	var results []*detectors.DetectionResult
	for _, detector := range allDetectors {
		if !detector.ShouldRun() {
			continue
		}

		fmt.Printf("Running %s detector...\n", detector.Name())
		result, err := detector.Detect()
		if err != nil {
			fmt.Printf("Warning: %s detector failed: %v\n", detector.Name(), err)
			continue
		}

		if result != nil {
			// Check if this key already exists
			if !keyExistsInConfig(config, result.Key, result.Value) {
				results = append(results, result)
			} else {
				fmt.Printf("Skipping %s: already exists in config\n", result.Key)
			}
		}
	}

	if len(results) == 0 {
		fmt.Println("No new suggestions found.")
		return
	}

	// Ask for each result individually
	fmt.Printf("\nFound %d suggestion(s):\n", len(results))
	addedCount := 0
	for i, result := range results {
		fmt.Printf("\n%d. %s\n", i+1, result.Description)
		fmt.Printf("   %s: %v\n", result.Key, result.Value)
		if result.Confidence < 1.0 {
			fmt.Printf("   (Confidence: %.0f%%)\n", result.Confidence*100)
		}

		fmt.Print("Add this to config? (y/N): ")
		var response string
		fmt.Scanln(&response)

		if strings.ToLower(strings.TrimSpace(response)) == "y" {
			if err := addKeyToConfig(*configFile, result.Key, result.Value); err != nil {
				fmt.Printf("Error adding %s: %v\n", result.Key, err)
			} else {
				fmt.Printf("✓ Added %s: %v\n", result.Key, result.Value)
				addedCount++
			}
		} else {
			fmt.Println("Skipped.")
		}
	}

	fmt.Printf("\nSuccessfully added %d item(s) to %s\n", addedCount, *configFile)
}



func keyExistsInConfig(config map[string]interface{}, key string, value interface{}) bool {
	// Check if there's a direct key
	if existing, exists := config[key]; exists {
		if fmt.Sprintf("%v", existing) == fmt.Sprintf("%v", value) {
			return true
		}
	}

	// Check nested objects for key
	for _, configValue := range config {
		if nestedMap, ok := configValue.(map[interface{}]interface{}); ok {
			if existing, exists := nestedMap[key]; exists {
				if fmt.Sprintf("%v", existing) == fmt.Sprintf("%v", value) {
					return true
				}
			}
		}
		if nestedMap, ok := configValue.(map[string]interface{}); ok {
			if existing, exists := nestedMap[key]; exists {
				if fmt.Sprintf("%v", existing) == fmt.Sprintf("%v", value) {
					return true
				}
			}
		}
	}

	return false
}

// addKeyToConfig adds any key-value pair to the config file while preserving order
func addKeyToConfig(configFile, key string, value interface{}) error {
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
		keyLine := formatKeyValue(key, value, 2)
		newConfig := fmt.Sprintf("%s:\n%s\n", projectKey, keyLine)
		if configText != "" {
			newConfig = newConfig + "\n" + configText
		}
		return ioutil.WriteFile(configFile, []byte(newConfig), 0644)
	}

	// Find where to insert the key line
	keyLine := formatKeyValue(key, value, 2)
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

	// Insert the key line
	if insertIndex >= 0 {
		newLines := make([]string, 0, len(lines)+1)
		newLines = append(newLines, lines[:insertIndex]...)
		newLines = append(newLines, keyLine)
		newLines = append(newLines, lines[insertIndex:]...)
		return ioutil.WriteFile(configFile, []byte(strings.Join(newLines, "\n")), 0644)
	}

	// Fallback: add at the end of the project section
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if strings.TrimSpace(line) != "" {
			newLines := make([]string, 0, len(lines)+1)
			newLines = append(newLines, lines[:i+1]...)
			newLines = append(newLines, keyLine)
			newLines = append(newLines, lines[i+1:]...)
			return ioutil.WriteFile(configFile, []byte(strings.Join(newLines, "\n")), 0644)
		}
	}

	return fmt.Errorf("could not find appropriate place to insert key")
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

// formatKeyValue formats a key-value pair into YAML format with specified indentation
func formatKeyValue(key string, value interface{}, indent int) string {
	indentStr := strings.Repeat(" ", indent)

	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%s%s: %s", indentStr, key, v)
	case bool:
		return fmt.Sprintf("%s%s: %t", indentStr, key, v)
	case int, int32, int64, float32, float64:
		return fmt.Sprintf("%s%s: %v", indentStr, key, v)
	case map[string]string:
		// For commands and similar map structures
		lines := []string{fmt.Sprintf("%s%s:", indentStr, key)}
		for k, val := range v {
			lines = append(lines, fmt.Sprintf("%s  %s: %s", indentStr, k, val))
		}
		return strings.Join(lines, "\n")
	case map[string]interface{}:
		lines := []string{fmt.Sprintf("%s%s:", indentStr, key)}
		for k, val := range v {
			lines = append(lines, formatKeyValue(k, val, indent+2))
		}
		return strings.Join(lines, "\n")
	case []string:
		// For arrays
		lines := []string{fmt.Sprintf("%s%s:", indentStr, key)}
		for _, item := range v {
			lines = append(lines, fmt.Sprintf("%s- %s", indentStr, item))
		}
		return strings.Join(lines, "\n")
	default:
		// Fallback to string representation
		return fmt.Sprintf("%s%s: %v", indentStr, key, v)
	}
}
