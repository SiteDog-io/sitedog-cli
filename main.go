package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	defaultConfigPath  = "./sitedog.yml"
	defaultTemplate    = "demo.html.erb"
	defaultPort        = 8081
	globalTemplatePath = ".sitedog/demo.html.erb"
	exampleConfig = `# Describe your project with a free key-value format, think simple.
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
  help    Show this help message

Options for init:
  --config PATH    Path to config file (default: ./sitedog.yml)

Options for live:
  --config PATH    Path to config file (default: ./sitedog.yml)
  --port PORT      Port to run server on (default: 8081)
`)
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

func handleLive() {
	liveFlags := flag.NewFlagSet("live", flag.ExitOnError)
	configFile := liveFlags.String("config", defaultConfigPath, "Path to config file")
	port := liveFlags.Int("port", defaultPort, "Port to run server on")
	liveFlags.Parse(os.Args[2:])

	if _, err := os.Stat(*configFile); err != nil {
		fmt.Println("Error:", *configFile, "not found. Run 'sitedog init' first.")
		os.Exit(1)
	}

	templatePath := findTemplate()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		config, _ := ioutil.ReadFile(*configFile)
		tmpl, _ := ioutil.ReadFile(templatePath)
		page := strings.Replace(string(tmpl), "{{CONFIG}}", string(config), -1)
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(page))
	})

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		config, _ := ioutil.ReadFile(*configFile)
		var data map[string]interface{}
		yaml.Unmarshal(config, &data)
		
		// Извлекаем только верхнеуровневые ключи в порядке их появления в YAML файле
		re := regexp.MustCompile(`(?m)^([a-zA-Z0-9_.]+):`)
		matches := re.FindAllStringSubmatch(string(config), -1)
		orderedKeys := make([]string, 0, len(matches))
		for _, match := range matches {
			if len(match) > 1 {
				orderedKeys = append(orderedKeys, match[1])
			}
		}

		// Создаем структуру ответа
		response := struct {
			Config      map[string]interface{} `json:"config"`
			OrderedKeys []string               `json:"orderedKeys"`
		}{
			Config:      data,
			OrderedKeys: orderedKeys,
		}

		w.Header().Set("Content-Type", "application/json")
		jsonData, _ := json.Marshal(response)
		w.Write(jsonData)
	})

	addr := fmt.Sprintf(":%d", *port)
	go func() {
		time.Sleep(500 * time.Millisecond)
		openBrowser("http://localhost" + addr)
	}()
	fmt.Println("Starting live server at http://localhost" + addr)
	fmt.Println("Press Ctrl+C to stop")
	log.Fatal(http.ListenAndServe(addr, nil))
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