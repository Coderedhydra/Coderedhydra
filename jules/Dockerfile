# Dockerfile for Jules Web Security Framework

# --- Stage 1: Build Go Components (Core, LLM, KB) ---
FROM golang:1.21-alpine AS go-builder
WORKDIR /app

# Copy Go module files first for dependency caching
COPY jules/core/go.mod jules/core/go.sum ./core/
COPY jules/llm/go.mod jules/llm/go.sum ./llm/
COPY jules/kb/go.mod jules/kb/go.sum ./kb/

# Download Go dependencies
RUN cd core && go mod download
RUN cd llm && go mod download
RUN cd kb && go mod download

# Copy rest of the Go source code
COPY jules/core/ ./core/
COPY jules/llm/ ./llm/
COPY jules/kb/ ./kb/

# Build Go components (as modules/libraries, not necessarily final executables yet)
# These RUN commands are more for ensuring they *can* be built.
# A final application structure might have a main Go program that imports these.
RUN cd core && go build -v -o /app/build/jules-core-module .
RUN cd llm && go build -v -o /app/build/jules-llm-module .
RUN cd kb && go build -v -o /app/build/jules-kb-module .
# For a real application, you'd likely build a single main Go executable if Python calls Go.
# Or, if Go calls Python for LLM, then Python is the main entry.

# --- Stage 2: Build React UI ---
FROM node:20-alpine AS ui-builder
WORKDIR /app
COPY jules/ui/package.json jules/ui/package-lock.json ./ui/
COPY jules/ui/ ./ui/
RUN cd ui && npm install && npm run build

# --- Stage 3: Final Application Image ---
FROM python:3.11-slim
WORKDIR /app

# Install runtime OS dependencies if any (e.g., git for some Python packages)
# RUN apt-get update && apt-get install -y --no-install-recommends git && rm -rf /var/lib/apt/lists/*

# Copy Python application code and LLM/KB module code (as Python might need to know their structure)
COPY jules/cli/ /app/jules/cli/
COPY jules/llm/ /app/jules/llm/ # Python LLM handler and Go files (Go files for reference if needed)
COPY jules/kb/ /app/jules/kb/   # Go KB files (for reference if needed)
COPY jules/core/ /app/jules/core/ # Go Core files (for reference if needed)


# Install Python dependencies
COPY jules/cli/requirements.txt /app/jules/cli/requirements.txt
RUN pip3 install --no-cache-dir -r /app/jules/cli/requirements.txt

COPY jules/llm/requirements.txt /app/jules/llm/requirements.txt
# Only install if not empty, as it's mostly a placeholder
RUN if [ -s /app/jules/llm/requirements.txt ]; then         pip3 install --no-cache-dir -r /app/jules/llm/requirements.txt;     fi

# --- LLM Model & Knowledge Base Data ---
# IMPORTANT: Managing LLM models (often several GBs) and extensive KB data in Docker images is complex.
# Option 1: Bake into the image (leads to large images, potential licensing issues for redistribution).
#   COPY models/qwen3-8b-awq /app/models/qwen3-8b-awq
#   COPY jules/kb/data /app/jules/kb/data
# Option 2: Mount as volumes at runtime (more flexible, keeps image smaller).
#   This Dockerfile assumes Option 2 for LLM & KB data.
#   Users will need to mount them to /app/models/ and /app/jules/kb/data respectively.
RUN mkdir -p /app/models/qwen3-8b-awq
RUN mkdir -p /app/jules/kb/data
ENV JULES_MODEL_PATH="/app/models/qwen3-8b-awq"
ENV JULES_KB_DATA_PATH="/app/jules/kb/data"
ENV JULES_LLM_SCRIPT_DIR="/app/jules/llm"
ENV JULES_PYTHON_EXEC="python3"


# Copy built Go applications/modules from go-builder stage (if they were executables)
# For now, Go components are mainly libraries called by Python or vice-versa via scripts.
# If Go parts become standalone servers, copy them here.
# COPY --from=go-builder /app/build/jules-main-app /app/jules-main-app

# Copy built UI assets
COPY --from=ui-builder /app/ui/dist /app/jules/ui/dist

# Copy run script and example config
COPY jules/run.sh /app/run.sh
RUN chmod +x /app/run.sh
COPY jules/config/ /app/jules/config/

# Expose ports (e.g., for proxy and UI)
EXPOSE 8080 # Default proxy port (if run via Python/Go)
EXPOSE 3000 # Default UI port (if UI is served by a simple server in container)

# Entrypoint or CMD
# The run.sh script is complex for a container entrypoint if it tries to install system deps.
# For now, it will be used as the entrypoint, but a container-specific startup script is often better.
# The Python CLI (main.py) will be the primary driver.
ENTRYPOINT ["/app/jules/cli/main.py"] 
# Default command if no args are provided to `docker run`
CMD ["--help"] 
# Example: docker run jules-image --target https://example.com --action scan
