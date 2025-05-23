# Dockerfile for Jules Web Security Framework

# --- Stage 1: Build Go Core ---
FROM golang:1.21-alpine AS go-builder
WORKDIR /app
# Copy Go module files and source code
COPY jules/core/go.mod jules/core/go.sum ./core/
COPY jules/core/ ./core/
# Download dependencies and build
RUN cd core && go mod download
RUN cd core && go build -o /app/jules-core .

# --- Stage 2: Build React UI ---
FROM node:20-alpine AS ui-builder
WORKDIR /app
# Copy UI source code and package files
COPY jules/ui/package.json jules/ui/package-lock.json ./ui/
COPY jules/ui/ ./ui/
# Install dependencies and build
RUN cd ui && npm install && npm run build

# --- Stage 3: Final Application Image ---
FROM python:3.11-slim
WORKDIR /app

# Install runtime dependencies (e.g., for serving UI, running Python)
# RUN apt-get update && apt-get install -y --no-install-recommends some-package && rm -rf /var/lib/apt/lists/*

# Copy built Go application from builder stage
COPY --from=go-builder /app/jules-core /app/jules-core

# Copy built UI assets from builder stage
COPY --from=ui-builder /app/ui/dist /app/jules/ui/dist

# Copy Python CLI application
COPY jules/cli/ /app/jules/cli/
# COPY jules/cli/requirements.txt /app/jules/cli/
# RUN pip3 install --no-cache-dir -r /app/jules/cli/requirements.txt

# Copy run script and example config
COPY jules/run.sh /app/run.sh
RUN chmod +x /app/run.sh
COPY jules/config/ /app/jules/config/

# Set environment variables (optional)
# ENV PATH="/app:${PATH}"

# Expose ports (e.g., for proxy and UI)
EXPOSE 8080 # Default proxy port
EXPOSE 3000 # Default UI port (if served by http-server or similar)

# Entrypoint or CMD
# The run.sh script is complex and might need adjustments for container environment
# For example, it might try to install system dependencies which isn't ideal for a container CMD.
# A container-specific run script might be better.
# For now, we can make it the entrypoint and let it handle calls to the CLI.
ENTRYPOINT ["/app/run.sh"]
# Example: to run proxy by default, or provide --help
CMD ["--help"]
