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
)

const (
	defaultConfigPath  = "./sitedog.yml"
	defaultTemplate    = "demo.html.erb"
	defaultPort        = 8081
	globalTemplatePath = ".sitedog/demo.html.erb"
	authFilePath      = ".sitedog/auth"
	apiBaseURL        = "http://localhost:4567" // Измените на реальный URL вашего API
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
	case "push":
		handlePush()
	case "render":
		handleRender()
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
  help    Show this help message

Options for init:
  --config PATH    Path to config file (default: ./sitedog.yml)

Options for live:
  --config PATH    Path to config file (default: ./sitedog.yml)
  --port PORT      Port to run server on (default: 8081)

Options for push:
  --config PATH    Path to config file (default: ./sitedog.yml)
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

func startServer(configFile *string, port int) (*http.Server, string) {
	// Обработчики
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
		w.Header().Set("Content-Type", "application/json")
		w.Write(config)
	})

	// Запускаем сервер
	addr := fmt.Sprintf(":%d", port)
	server := &http.Server{
		Addr: addr,
	}

	// Запускаем сервер в горутине
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Ждем запуска сервера
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

	// Ждем сигнала завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	// Корректно завершаем сервер
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
	pushFlags.Parse(os.Args[2:])

	if _, err := os.Stat(*configFile); err != nil {
		fmt.Println("Error:", *configFile, "not found. Run 'sitedog init' first.")
		os.Exit(1)
	}

	// Получаем токен авторизации
	token, err := getAuthToken()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Читаем конфигурацию
	config, err := ioutil.ReadFile(*configFile)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}

	// Получаем имя конфигурации из имени директории
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		os.Exit(1)
	}
	configName := filepath.Base(dir)

	// Отправляем конфигурацию на сервер
	err = pushConfig(token, configName, string(config))
	if err != nil {
		fmt.Println("Error pushing config:", err)
		os.Exit(1)
	}

	fmt.Println("Configuration pushed successfully!")
}

func getAuthToken() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("error getting current user: %v", err)
	}

	authFile := filepath.Join(usr.HomeDir, authFilePath)
	if _, err := os.Stat(authFile); err == nil {
		// Если файл существует, читаем токен
		token, err := ioutil.ReadFile(authFile)
		if err != nil {
			return "", fmt.Errorf("error reading auth file: %v", err)
		}
		return strings.TrimSpace(string(token)), nil
	}

	// Если файл не существует, запрашиваем авторизацию
	fmt.Print("Email: ")
	var email string
	fmt.Scanln(&email)

	fmt.Print("Password: ")
	var password string
	fmt.Scanln(&password)

	// Создаем запрос на авторизацию
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

	// Создаем директорию .sitedog если её нет
	authDir := filepath.Dir(authFile)
	if err := os.MkdirAll(authDir, 0700); err != nil {
		return "", fmt.Errorf("error creating auth directory: %v", err)
	}

	// Сохраняем токен
	if err := ioutil.WriteFile(authFile, []byte(result.Token), 0600); err != nil {
		return "", fmt.Errorf("error saving token: %v", err)
	}

	return result.Token, nil
}

func pushConfig(token, name, content string) error {
	reqBody, err := json.Marshal(map[string]string{
		"name":    name,
		"content": content,
	})
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req, err := http.NewRequest("POST", apiBaseURL+"/cli/push", strings.NewReader(string(reqBody)))
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

func handleRender() {
	renderFlags := flag.NewFlagSet("render", flag.ExitOnError)
	configFile := renderFlags.String("config", defaultConfigPath, "Path to config file")
	outputFile := renderFlags.String("output", "sitedog.html", "Path to output HTML file")
	renderFlags.Parse(os.Args[2:])

	if _, err := os.Stat(*configFile); err != nil {
		fmt.Println("Error:", *configFile, "not found. Run 'sitedog init' first.")
		os.Exit(1)
	}
  
	port := 34324
	server, addr := startServer(configFile, port)
	url := "http://localhost" + addr

	// Проверяем доступность сервера
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error checking server:", err)
		server.Close()
		os.Exit(1)
	}
	resp.Body.Close()

	// Используем Obelisk для сохранения страницы
	archiver := &obelisk.Archiver{
		EnableLog:             false,
		MaxConcurrentDownload: 10,
	}

	// Валидируем архиватор
	archiver.Validate()

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Создаем запрос
	req := obelisk.Request{
		URL: url,
	}

	// Сохраняем страницу
	html, _, err := archiver.Archive(ctx, req)
	if err != nil {
		fmt.Println("Error archiving page:", err)
		server.Close()
		os.Exit(1)
	}

	// Сохраняем результат в файл
	if err := ioutil.WriteFile(*outputFile, html, 0644); err != nil {
		fmt.Println("Error saving file:", err)
		server.Close()
		os.Exit(1)
	}

	// Закрываем сервер
	server.Close()

	fmt.Printf("Page saved to %s\n", *outputFile)
}

func getFaviconCache(config []byte) []byte {
	// Парсим YAML конфиг
	var configMap map[string]interface{}
	if err := yaml.Unmarshal(config, &configMap); err != nil {
		log.Printf("Error parsing YAML: %v", err)
		return []byte("{}")
	}

	log.Printf("Parsed config: %+v", configMap)

	// Создаем карту для хранения favicon кэша
	faviconCache := make(map[string]string)

	// Функция для извлечения домена из URL
	extractDomain := func(urlStr string) string {
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			log.Printf("Error parsing URL %s: %v", urlStr, err)
			return ""
		}
		domain := parsedURL.Hostname()
		log.Printf("Extracted domain from %s: %s", urlStr, domain)
		return domain
	}

	// Функция для рекурсивного обхода значений
	var traverseValue func(value interface{})
	traverseValue = func(value interface{}) {
		switch v := value.(type) {
		case string:
			// Проверяем, является ли строка URL
			if strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
				log.Printf("Found URL: %s", v)
				domain := extractDomain(v)
				if domain != "" {
					// Получаем favicon
					faviconURL := fmt.Sprintf("https://www.google.com/s2/favicons?domain=%s&sz=64", url.QueryEscape(domain))
					log.Printf("Fetching favicon from: %s", faviconURL)
					resp, err := http.Get(faviconURL)
					if err != nil {
						log.Printf("Error fetching favicon for %s: %v", v, err)
						return
					}
					defer resp.Body.Close()

					if resp.StatusCode == http.StatusOK {
						// Читаем favicon
						faviconData, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							log.Printf("Error reading favicon for %s: %v", v, err)
							return
						}

						// Конвертируем в base64
						base64Data := base64.StdEncoding.EncodeToString(faviconData)
						// Get content type from response headers
						contentType := resp.Header.Get("Content-Type")
						if contentType == "" {
							contentType = "image/png" // fallback to png if no content type specified
						}
						dataURL := fmt.Sprintf("data:%s;base64,%s", contentType, base64Data)
						faviconCache[v] = dataURL
						log.Printf("Successfully cached favicon for %s", v)
					} else {
						log.Printf("Failed to fetch favicon for %s: status code %d", v, resp.StatusCode)
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

	// Обходим все значения в конфиге
	traverseValue(configMap)

	// Конвертируем карту в JSON
	jsonData, err := json.Marshal(faviconCache)
	if err != nil {
		log.Printf("Error marshaling favicon cache: %v", err)
		return []byte("{}")
	}

	log.Printf("Final favicon cache: %s", string(jsonData))
	return jsonData
}