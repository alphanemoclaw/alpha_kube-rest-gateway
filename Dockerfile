# Comment added by OpenClaw assistant
# Use official Python image as base
FROM python:3.11-slim

# Set working directory
WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements and install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Create non-root user for security
RUN useradd --create-home --shell /bin/bash appuser \
    && chown -R appuser:appuser /app
USER appuser

# Expose the port the gateway runs on
EXPOSE 8081

# Environment variables with defaults (can be overridden at runtime)
ENV KUBECONFIG_PATH=/app/.kube/config \
    GATEWAY_API_TOKEN= \
    GATEWAY_HOST=0.0.0.0 \
    GATEWAY_PORT=8081 \
    DEFAULT_NAMESPACE=default \
    DEFAULT_LOG_TAIL_LINES=100 \
    LOG_LEVEL=INFO

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8081/healthz || exit 1

# Run the application
CMD ["python", "main.py"]
