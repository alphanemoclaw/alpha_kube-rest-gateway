# alpha-kube-rest-gateway (Go version)

A lightweight Go-based REST gateway that acts as a secure proxy between an AI agent and a Kubernetes cluster.

## Features

- **Lightweight & Fast**: Rewritten in Go for better performance and lower footprint.
- **Secure**: Authentication via static bearer token.
- **AI Friendly**: Provides a `/api/help` endpoint that gives the AI agent instructions on how to use the cluster.
- **Read-Only (mostly)**: Designed to allow the agent to inspect the cluster without full administrative access.

## Quick Start

### Build and Run locally

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/your-org/alpha-kube-rest-gateway.git
    cd alpha-kube-rest-gateway
    ```

2.  **Install dependencies**:
    ```bash
    go mod download
    ```

3.  **Run the server**:
    ```bash
    go run main.go
    ```
    The server will start on `http://0.0.0.0:8081`.

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `GATEWAY_API_TOKEN` | Bearer token for authentication | (empty, auth disabled) |
| `KUBECONFIG_PATH` | Path to kubeconfig file | `~/.kube/config` |
| `GATEWAY_PORT` | Port to listen on | `8081` |
| `GATEWAY_HOST` | Host to listen on | `0.0.0.0` |
| `DEFAULT_NAMESPACE` | Default K8s namespace | `default` |

## API Endpoints

- `GET /healthz` - Liveness probe
- `GET /api/help` - AI Agent instructions
- `GET /api/pods` - List pods
- `GET /api/pods/:name` - Describe a pod
- `GET /api/logs/:name` - Fetch pod logs
- `GET /api/services` - List services
- `GET /api/deployments` - List deployments
- `GET /api/nodes` - List nodes
- `GET /api/namespaces` - List namespaces
- `GET /api/events` - List cluster events

## Docker

Build the image:
```bash
docker build -t kube-rest-gateway .
```

Run the container:
```bash
docker run -p 8081:8081 -v ~/.kube/config:/home/gatewayuser/.kube/config kube-rest-gateway
```
