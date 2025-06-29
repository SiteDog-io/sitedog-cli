<!DOCTYPE html>
<html>
  <head>
    <title>SiteDog Demo</title>
    <link rel="stylesheet" href="https://sitedog.io/css/preview.css" />
    <style>
      /* Import Inter font for modern look */
      @import url("https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600&display=swap");

      /* CSS variables for centralized color management */
      :root {
        --bg-body: #f4f6f8;
        --bg-card: #ffffff;
        --border: #e0e3e7;
        --text-main: #111827;
        --text-muted: #6b7280;
        --accent: #f4b760;
        --editor-bg: #1e1e1e;
        --editor-text: #d4d4d4;
        --error-color: #e53935;
      }

      /* Reset and base styles */
      html, body {
        width: 100vw;
        min-height: 100vh;
        overflow-x: hidden;
      }

      * {
        margin: 0;
        padding: 0;
        box-sizing: border-box;
      }

      body {
        font-family: 'Inter', sans-serif;
        background-color: var(--bg-body);
        color: var(--text-main);
        min-height: 100vh;
        display: flex;
        flex-direction: column;
        width: 100vw;
        overflow-x: hidden;
      }

      /* Header styles */
      header {
        background-color: var(--bg-card);
        border-bottom: 1px solid var(--border);
        padding: 1rem 2rem;
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        z-index: 100;
        display: flex;
        align-items: center;
        justify-content: space-between;
        width: 100vw;
        box-sizing: border-box;

      }

      header h1 {
        font-size: 1.5rem;
        font-weight: 600;
        color: var(--text-main);
      }

      /* Main content styles */
      main {
        flex: 1;
        margin-top: 4rem; /* Height of header */
        padding: 0;
        min-height: calc(100vh - 4rem);
        width: 100vw;
        box-sizing: border-box;
        display: flex;
        flex-direction: column;
        align-items: stretch;
      }

      .preview-section {
        background-color: var(--bg-card);
        border-radius: 0;
        box-shadow: none;
        width: 100%;
        min-height: calc(100vh - 4rem);
        height: 100%;
        box-sizing: border-box;
        overflow-x: auto;
        padding: 0;
      }

      #card-container {
        padding: 2rem;
        height: 100%;
        width: 100%;
        box-sizing: border-box;
      }

      .error-message {
        text-align: center;
        color: var(--error-color);
        padding: 1rem;
        font-size: 0.9rem;
      }

      #config {
        display: none;
      }
    </style>
  </head>
  <body>
    <header>
      <h1>SITEDOG Preview</h1>
    </header>
    <main>
      <div class="preview-section">
        <div id="card-container"></div>
      </div>
    </main>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/js-yaml/4.1.0/js-yaml.min.js"></script>
    <script src="https://sitedog.io/js/renderCards.js"></script>
    <script>
      const cardContainer = document.getElementById("card-container");
      let lastConfig = null;

      function updateCards(yamlText, faviconCache = null) {
        if (yamlText !== lastConfig) {
          lastConfig = yamlText;
          renderCards(yamlText, cardContainer, (config, result, error) => {
            if (!result) {
              cardContainer.innerHTML = `<div class="error-message">YAML Error: ${error}</div>`;
            }
          },
          {
            faviconCache: faviconCache,
          });
        }
      }

      const staticConfig = `{{CONFIG}}`;
      if (staticConfig !== (`{{` + `CONFIG` + `}}`) && location.hostname !== "localhost") {
        const faviconCache = JSON.parse(`{{FAVICON_CACHE}}`);
        updateCards(staticConfig, faviconCache);
      } else {
        // Long polling for checking updates
        function checkForUpdates() {
          fetch('/config')
            .then(response => response.text())
            .then(config => {
              updateCards(config);
              // Continue polling
              setTimeout(checkForUpdates, 100);
            })
            .catch(error => {
              console.error('Error checking for updates:', error);
              // On error, try again after 5 seconds
              setTimeout(checkForUpdates, 5000);
            });
        }

        // Start checking for updates
        checkForUpdates();        
      }
    </script>
  </body>
</html>

