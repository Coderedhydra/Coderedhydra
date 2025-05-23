# Jules: Next-Generation Offline Web Security Toolkit

Jules is a cutting-edge, fully offline web-security toolkit built around the **Qwen/Qwen3-8B-AWQ Large Language Model (LLM)**. It empowers security teams to perform deep, automated assessments without reliance on external services by leveraging advanced AI for intelligent vulnerability discovery and payload generation.

## Core Features

*   **Form Harvesting & Context Analysis**: Jules meticulously crawls target sites, automatically discovering all HTML forms and their parameters. The integrated LLM then analyzes and classifies each input’s precise context (e.g., within HTML body, specific attribute, URL parameter, JSON data, JavaScript block).
*   **LLM-Powered Advanced Payload Generation**: Leveraging the Qwen3-8B-AWQ model, Jules crafts bespoke, context-aware security payloads. These aren't just generic attack strings; they combine sophisticated encoding techniques, known bypass methods, and novel fuzz-vectors, all tailored to the specific structure, context, and potential risk profile of each discovered input field.
*   **Automated Submission & Validation**: Generated payloads are seamlessly submitted to the target application through either a built-in intercepting proxy or a direct HTTP client. Jules then parses application responses, actively looking for proof-of-concept evidence of vulnerabilities (e.g., reflected XSS, SQL error messages, SSRF callback indicators).
*   **Self-Evolving Intelligence via Offline Knowledge Base**: When new error messages, parameter names, or potential bypass clues are identified during scanning, Jules consults its comprehensive, locally mirrored vulnerability knowledge base. This knowledge base contains CVE archives, OWASP guidelines, Exploit-DB dumps, and other security resources. The LLM uses this "local Google" to understand new patterns and refine its payload generation strategies in real-time, effectively expanding its exploit library.
*   **End-to-End Reporting**: All findings—including successfully exploited vulnerabilities, the exact payloads used, contextual information, crawl statistics, and discovered parameters—are collated into a detailed JSON report. HTML reporting is also planned for easier human consumption and integration into CI/CD pipelines.
*   **Fully Offline Operation**: Designed for sensitive environments, Jules and its core LLM capabilities operate entirely offline, ensuring no data leaves your local environment.

## Project Structure

*   `/cli/`: Python-based command-line interface (enhanced for LLM-driven scans).
*   `/ui/`: React (Vite) web front-end (for visualizing complex findings and interactions).
*   `/core/`: Shared logic in Go, including HTTP handling and plugin systems.
*   `/llm/`: (New) Go/Python components for Qwen3-8B-AWQ model interaction, prompt engineering, and context processing.
*   `/kb/`: (New) Local knowledge base resources and management scripts.
*   `/config/`: Configuration files for Jules, including LLM and KB settings.

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
    *   **Qwen/Qwen3-8B-AWQ Model**: Specific setup instructions for the LLM will be provided (e.g., model download, local API endpoint if applicable). This is essential for core functionality.
    *   **Knowledge Base**: Instructions for populating the local knowledge base.

3.  **Run the setup and launch script:**
    The `run.sh` script will handle dependency installation, building components, and launching the application.
    ```bash
    ./run.sh --target https://example.com 
    ```
    *(This command will eventually trigger an LLM-enhanced scan)*

## Modules Explained (High-Level)

*   **Crawler (`jules/cli/crawler.py`)**: Discovers site structure, focusing on harvesting all forms and input parameters for LLM analysis.
*   **LLM Orchestrator (`jules/llm/`)**: Manages interaction with the Qwen model for context analysis, payload generation, and knowledge base queries.
*   **Scanner (`jules/cli/scanner.py`, `jules/core/plugin_system.go`)**: Drives the automated testing process, submitting LLM-generated payloads and validating responses. Leverages plugins for different vulnerability types, all enhanced by the LLM.
*   **Knowledge Base (`jules/kb/`)**: Stores and provides access to offline security data (CVEs, etc.) for the LLM's self-improvement cycle.
*   **Proxy & HTTP Client (`jules/core/http_client.go`)**: Handles all HTTP(S) traffic for crawling, scanning, and payload submission.
*   **Web UI (`jules/ui/`)**: Visualizes the sitemap, scan progress, detailed findings (including LLM insights), and facilitates manual request crafting/replaying.

## Development Roadmap (Illustrative)

*   Integrate Qwen3-8B-AWQ model for basic context classification.
*   Develop initial payload generation prompts based on context.
*   Implement local knowledge base querying.
*   Refine LLM feedback loops for self-improvement.
*   Expand UI to showcase LLM-specific data.

## Contributing

Contributions are highly welcome, especially in areas of LLM prompt engineering, knowledge base curation, and advanced exploit techniques. Please fork the repository and submit a pull request.

## Licensing Considerations

While the Jules Framework source code itself is released under the MIT License (see below), its full functionality relies on external components and data sources that users must procure and manage, each potentially under its own license:

*   **Qwen/Qwen3-8B-AWQ LLM Model**: Jules is designed to integrate with the Qwen/Qwen3-8B-AWQ language model. Users are responsible for downloading this model and adhering to its specific license terms. The Jules project does not distribute the LLM model. Please consult the official Qwen model sources for licensing details.

*   **Knowledge Base Data**: The `jules/kb/data/` directory is intended to be populated by users with information from various sources such as CVE databases (e.g., NVD), OWASP documentation, Exploit-DB, and other security resources. Users must ensure that their use of any such data complies with the original licensing terms of those sources. The Jules project does not provide this data directly, other than potentially example placeholder structures.

*   **Third-Party Libraries**: Jules utilizes various third-party Go and Python libraries, each with its own license (typically permissive, like MIT, Apache 2.0, BSD). These are managed via Go modules and Python requirements files.

It is the user's responsibility to ensure compliance with all applicable licenses for any components or data they integrate with or use in conjunction with the Jules framework.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
