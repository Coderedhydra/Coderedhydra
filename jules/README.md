# Jules Web Security Framework

Jules is an open-source web security framework inspired by Burp Suite, offering both CLI and web-UI components for web application security testing.

## Features (Planned)

*   **Crawler**: Discovers URLs and site structure.
*   **HTTP(S) Proxy**: Intercepts and allows modification of web traffic.
*   **Active Scanner**: Probes for vulnerabilities like XSS, SQLi, SSRF.
*   **Repeater**: Manually send and replay HTTP requests.
*   **Web UI**: Provides a graphical interface for sitemap, repeater, and scan results.
*   **Reporting**: Generates reports in JSON (HTML planned).

## Project Structure

*   `/cli/`: Python-based command-line interface.
*   `/ui/`: React (Vite) web front-end.
*   `/core/`: Shared logic in Go (HTTP client, plugin system, report generator).
*   `/config/`: Configuration files.
*   `/docs/`: (Optional, for more detailed documentation later)

## Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-org/jules.git 
    # Replace your-org with the actual organization or username
    cd jules
    ```

2.  **Prerequisites:**
    *   Go (version 1.21+)
    *   Node.js (version 20+) and npm
    *   Python (version 3.9+) and pip

    The `run.sh` script will attempt to check for these but not install them system-wide. Please ensure they are available in your PATH.

3.  **Run the setup and launch script:**
    The `run.sh` script handles dependency installation (local to the project where possible), building components, and launching the application.

## One-Line Run (Example)

This command will (eventually) initialize dependencies, build components, and start the CLI with a target, potentially running a default set of actions like crawling and scanning.

```bash
./run.sh --target https://example.com
```

**Note:** The `run.sh` script is currently under development and primarily sets up the project and shows placeholder commands. Full functionality will be implemented incrementally.

## Modules Explained

*   **Proxy (`./run.sh --action proxy --proxy-port 8080`)**:
    Starts an HTTP/S proxy on the specified port (default 8080). Allows you to route browser or tool traffic through Jules to intercept, view, and modify requests and responses.
    *(Functionality to be implemented in `jules/cli/proxy.py` and `jules/core/http_client.go`)*

*   **Crawler (`./run.sh --action crawl --target https://example.com`)**:
    Spiders the target website to discover accessible pages and resources. The discovered sitemap will be available in the UI and can be used as a basis for other modules.
    *(Functionality to be implemented in `jules/cli/crawler.py`)*

*   **Scanner (`./run.sh --action scan --target https://example.com --scan-type xss`)**:
    Actively scans the target for common web vulnerabilities.
    Available scan types (planned): `xss`, `sqli`, `ssrf`.
    *(Functionality to be implemented in `jules/cli/scanner.py` and `jules/core/plugin_system.go`)*

*   **Repeater (Web UI)**:
    A tool within the web UI that allows you to take any HTTP request, modify it, resend it, and view the response. Useful for manual testing and vulnerability verification.
    *(Functionality to be implemented in `jules/ui/src/components/RepeaterPanel.jsx` and backend)*

## Configuration

An example configuration file can be found at `jules/config/default.yaml`. This file will be used to set options for the proxy, scanner, and other components. (Actual loading and usage of this config is TBD).

## Development

*   **CLI**: `python3 jules/cli/main.py --help`
*   **UI**: `cd jules/ui && npm run dev`
*   **Core**: `cd jules/core && go build .`

## Contributing

Contributions are welcome! Please fork the repository, create a new branch for your feature or fix, and submit a pull request.
(Further contribution guidelines to be added)

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
