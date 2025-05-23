#!/bin/bash

# Exit on any error
set -e

# --- Configuration ---
TARGET_URL=""
CLI_ACTION="default_action" # e.g., crawl, scan, proxy, or a default combined action
PROXY_PORT="8080"
SCAN_TYPE="" # xss, sqli, ssrf

# --- Helper Functions ---
print_usage() {
  echo "Usage: $0 --target <URL> [options]"
  echo "Options:"
  echo "  --target <URL>         Specify the target URL (required for most operations)"
  echo "  --action <action>      CLI action: crawl, proxy, scan (default: a combined action)"
  echo "  --proxy-port <port>    Port for the intercepting proxy (default: 8080)"
  echo "  --scan-type <type>     Type of scan to run (e.g., xss, sqli, ssrf)"
  echo "  --help                 Display this help message"
  exit 1
}

# --- Parse Command Line Arguments ---
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --target) TARGET_URL="$2"; shift ;;
        --action) CLI_ACTION="$2"; shift ;;
        --proxy-port) PROXY_PORT="$2"; shift ;;
        --scan-type) SCAN_TYPE="$2"; shift ;;
        --help) print_usage ;;
        *) echo "Unknown parameter passed: $1"; print_usage ;;
    esac
    shift
done

# Validate target URL for relevant actions
if [[ "$CLI_ACTION" != "proxy" && -z "$TARGET_URL" && "$CLI_ACTION" != "default_action_no_target_needed" ]]; then # Adjust condition as needed
    echo "Error: --target <URL> is required for the action '$CLI_ACTION'."
    print_usage
fi


# --- Dependency Installation ---
echo "[+] Checking and installing dependencies..."

# Check for Go
if ! command -v go &> /dev/null; then
    echo "Go not found. Please install Go (https://golang.org/doc/install)."
    # Add instructions or attempt to install if feasible and desired
    exit 1
fi

# Check for Node.js/npm
if ! command -v npm &> /dev/null; then
    echo "npm not found. Please install Node.js and npm (https://nodejs.org/)."
    # Add instructions or attempt to install if feasible and desired
    exit 1
fi

# Check for Python/pip
if ! command -v python3 &> /dev/null || ! command -v pip3 &> /dev/null; then
    echo "Python 3 and pip3 not found. Please install them."
    # Add instructions or attempt to install if feasible and desired
    exit 1
fi

echo "[+] Installing Python dependencies..."
# Create a virtual environment if desired
# python3 -m venv .venv
# source .venv/bin/activate
# pip3 install -r jules/cli/requirements.txt # Assuming a requirements.txt will be added

echo "[+] Installing Go dependencies for core..."
(cd jules/core && go mod tidy) # Fetches dependencies listed in go.mod

echo "[+] Installing Node.js dependencies for UI..."
(cd jules/ui && npm install)


# --- Build Steps ---
echo "[+] Building components..."

echo "[+] Building Go core..."
(cd jules/core && go build -o ../../jules-core .) # Output to top-level for easier access or to a bin/ dir

echo "[+] Building React UI..."
(cd jules/ui && npm run build) # Assumes 'build' script in ui/package.json creates a dist/ folder


# --- Launch Application ---
echo "[+] Launching Jules Framework..."

# Example: Launch UI server in the background and then CLI
# This is a simplified example. Proper process management (e.g., using a process manager like systemd or supervisord for production,
# or just backgrounding with '&' and managing PIDs for a script) is needed for robust operation.

echo "[+] Starting UI server (placeholder)..."
# (cd jules/ui && npm run dev -- --port 3000) & # Example: Run Vite dev server
# UI_PID=$!
# Or serve the build directory:
# (cd jules/ui/dist && npx http-server -p 3000) &
# UI_PID=$!
# echo "UI server started with PID $UI_PID on port 3000 (placeholder)."
echo "UI launch is placeholder. To run UI: cd jules/ui && npm run dev"


echo "[+] Executing CLI command..."
# Construct the CLI command based on parsed arguments
CLI_COMMAND_ARGS=()
if [[ -n "$TARGET_URL" ]]; then
    CLI_COMMAND_ARGS+=(--target "$TARGET_URL")
fi

if [[ "$CLI_ACTION" == "default_action" && -n "$TARGET_URL" ]]; then
    echo "Running default action: Crawl and Scan (XSS, SQLi) on $TARGET_URL"
    # python3 jules/cli/main.py crawl "$TARGET_URL"
    # python3 jules/cli/main.py scan "$TARGET_URL" --type xss
    # python3 jules/cli/main.py scan "$TARGET_URL" --type sqli
    echo "CLI: (Placeholder) Would crawl and scan $TARGET_URL"
elif [[ "$CLI_ACTION" == "crawl" && -n "$TARGET_URL" ]]; then
    # python3 jules/cli/main.py crawl "$TARGET_URL"
    echo "CLI: (Placeholder) Would crawl $TARGET_URL"
elif [[ "$CLI_ACTION" == "proxy" ]]; then
    # python3 jules/cli/main.py proxy --port "$PROXY_PORT"
    echo "CLI: (Placeholder) Would start proxy on port $PROXY_PORT"
elif [[ "$CLI_ACTION" == "scan" && -n "$TARGET_URL" && -n "$SCAN_TYPE" ]]; then
    # python3 jules/cli/main.py scan "$TARGET_URL" --type "$SCAN_TYPE"
    echo "CLI: (Placeholder) Would scan $TARGET_URL for $SCAN_TYPE"
elif [[ "$CLI_ACTION" == "scan" && -n "$TARGET_URL" && -z "$SCAN_TYPE" ]]; then
    echo "CLI: (Placeholder) Would run all scans on $TARGET_URL"
    # python3 jules/cli/main.py scan "$TARGET_URL" --type xss
    # python3 jules/cli/main.py scan "$TARGET_URL" --type sqli
    # python3 jules/cli/main.py scan "$TARGET_URL" --type ssrf
else
    echo "No specific CLI action specified or missing parameters. See --help."
    # python3 jules/cli/main.py "${CLI_COMMAND_ARGS[@]}" # Or default help
fi

# --- Cleanup (optional, for background processes) ---
# trap 'kill $UI_PID; echo "UI server stopped."' EXIT
# wait $UI_PID # Wait for UI server if it's critical for the script's foreground lifetime

echo "[+] Jules run script finished."
# Deactivate virtual environment if used
# deactivate
