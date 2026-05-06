# alpha-kube-rest-gateway (Go version)

A lightweight Go-based REST gateway that acts as a secure proxy between an AI agent and a Kubernetes cluster.

## Features

- **Lightweight & Fast**: Rewritten in Go for better performance and lower footprint.
- **Secure**: Authentication via static bearer token.
- **AI Friendly**: Provides a `/api/help` endpoint that gives the AI agent instructions on how to use the cluster.
- **Read-Only**: Designed to allow the agent to inspect the cluster without full administrative access.

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
| `DEFAULT_LOG_TAIL_LINES` | Default pod log tail size | `100` |

## API Endpoints

- `GET /healthz` - Liveness probe
- `GET /api/help` - AI Agent instructions
- `GET /api/pods` - List pods
- `GET /api/pods/:name` - Describe a pod
- `GET /api/pods/:name/status` - Detailed pod/container lifecycle status
- `GET /api/logs/:name` - Fetch pod logs
- `GET /api/services` - List services
- `GET /api/deployments` - List deployments
- `GET /api/replicasets` - List ReplicaSets
- `GET /api/statefulsets` - List StatefulSets
- `GET /api/daemonsets` - List DaemonSets
- `GET /api/jobs` - List Jobs
- `GET /api/cronjobs` - List CronJobs
- `GET /api/nodes` - List nodes
- `GET /api/namespaces` - List namespaces
- `GET /api/events` - List cluster events
- `GET /api/endpoints` - List Endpoints
- `GET /api/endpointslices` - List EndpointSlices
- `GET /api/ingresses` - List Ingresses
- `GET /api/pvcs` - List PersistentVolumeClaims
- `GET /api/pvs` - List PersistentVolumes
- `GET /api/storageclasses` - List StorageClasses
- `GET /api/networkpolicies` - List NetworkPolicies
- `GET /api/resourcequotas` - List ResourceQuotas
- `GET /api/limitranges` - List LimitRanges
- `GET /api/configmaps` - List ConfigMap metadata and keys only
- `GET /api/metrics/pods` - List pod metrics if metrics-server is installed
- `GET /api/metrics/nodes` - List node metrics if metrics-server is installed

## API reference

All endpoints (except `/healthz`) require the header:

```
Authorization: Bearer <GATEWAY_API_TOKEN>
```

### `GET /healthz`
Liveness probe — no auth required.

### `GET /api/pods`
List pods in a namespace.

| Query param | Default | Description |
|---|---|---|
| `namespace` | `default` | Kubernetes namespace |
| `label_selector` | — | e.g. `app=nginx` |
| `field_selector` | — | e.g. `status.phase=Running` |

### `GET /api/pods/{pod_name}`
Describe a single pod.

### `GET /api/pods/{pod_name}/status`
Detailed pod status for troubleshooting `Pending`, `CrashLoopBackOff`, `ImagePullBackOff`, init-container failures, and readiness issues.

Returns pod phase/reason/message, node placement, restart policy, service account, conditions, init container states, container states, last terminated states, restart counts, waiting reasons, messages, exit codes, and timestamps.

### `GET /api/logs/{pod_name}`
Tail pod logs.

| Query param | Default | Description |
|---|---|---|
| `namespace` | `default` | Kubernetes namespace |
| `container` | — | Container name (required for multi-container pods) |
| `tail_lines` | `100` | Lines to return (max 10 000) |
| `previous` | `false` | Return logs of the previous (crashed) instance |

### `GET /api/services`
List services in a namespace.

### `GET /api/deployments`
List deployments in a namespace.

Supports `namespace` and `label_selector`.

### Workload endpoints

The following endpoints support `namespace` and `label_selector`:

- `GET /api/replicasets`
- `GET /api/statefulsets`
- `GET /api/daemonsets`
- `GET /api/jobs`
- `GET /api/cronjobs`

Use these to troubleshoot rollout ownership, desired vs ready replicas, DaemonSet scheduling, Job failures, and CronJob schedules.

### `GET /api/nodes`
List cluster nodes (cluster-scoped, no namespace param).

### `GET /api/namespaces`
List all namespaces (cluster-scoped).

### `GET /api/events`
List recent events — great for diagnosing pod failures.

| Query param | Default | Description |
|---|---|---|
| `namespace` | `default` | Kubernetes namespace |
| `field_selector` | — | e.g. `involvedObject.name=my-pod` |

### Service routing endpoints

The following endpoints support `namespace` and `label_selector`:

- `GET /api/endpoints`
- `GET /api/endpointslices`
- `GET /api/ingresses`

Use these to verify service backends, ready/not-ready endpoint addresses, EndpointSlice conditions, and ingress host/path/TLS/backend mappings.

### Storage endpoints

- `GET /api/pvcs` supports `namespace` and `label_selector`.
- `GET /api/pvs` is cluster-scoped and supports `label_selector`.
- `GET /api/storageclasses` is cluster-scoped and supports `label_selector`.

Use these to troubleshoot PVC `Pending`, PV binding, reclaim policies, storage classes, and volume expansion settings.

### Policy and quota endpoints

The following endpoints support `namespace` and `label_selector`:

- `GET /api/networkpolicies`
- `GET /api/resourcequotas`
- `GET /api/limitranges`

Use these to troubleshoot pod admission failures, quota exhaustion, default limits, and network isolation.

### `GET /api/configmaps`
List ConfigMap metadata and key names only. Values are intentionally omitted because ConfigMaps may contain sensitive data.

Supports `namespace` and `label_selector`.

### Metrics endpoints

- `GET /api/metrics/pods` supports `namespace` and `label_selector`.
- `GET /api/metrics/nodes` is cluster-scoped and supports `label_selector`.

These require Kubernetes metrics-server. If metrics-server is not installed or the gateway service account lacks access, the gateway returns the Kubernetes API error.

---

## curl test examples

> Replace `YOUR_TOKEN` with your actual `GATEWAY_API_TOKEN` value.

```bash
# ── Liveness probe (no auth needed) ──────────────────────────────────────────
curl http://YOUR_SERVER_IP:8081/healthz

# ── List pods in the default namespace ───────────────────────────────────────
curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://YOUR_SERVER_IP:8081/api/pods

# ── List pods in a specific namespace ────────────────────────────────────────
curl -H "Authorization: Bearer YOUR_TOKEN" \
     "http://YOUR_SERVER_IP:8081/api/pods?namespace=kube-system"

# ── Filter pods by label ─────────────────────────────────────────────────────
curl -H "Authorization: Bearer YOUR_TOKEN" \
     "http://YOUR_SERVER_IP:8081/api/pods?label_selector=app%3Dnginx"

# ── Describe a single pod ─────────────────────────────────────────────────────
curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://YOUR_SERVER_IP:8081/api/pods/my-pod-7d9f8b-xkqzp

# ── Fetch the last 50 log lines from a pod ───────────────────────────────────
curl -H "Authorization: Bearer YOUR_TOKEN" \
     "http://YOUR_SERVER_IP:8081/api/logs/my-pod-7d9f8b-xkqzp?tail_lines=50"

# ── Logs for a specific container in a multi-container pod ───────────────────
curl -H "Authorization: Bearer YOUR_TOKEN" \
     "http://YOUR_SERVER_IP:8081/api/logs/my-pod-7d9f8b-xkqzp?container=sidecar&tail_lines=200"

# ── List services ─────────────────────────────────────────────────────────────
curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://YOUR_SERVER_IP:8081/api/services

# ── List deployments ──────────────────────────────────────────────────────────
curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://YOUR_SERVER_IP:8081/api/deployments

# ── List nodes ────────────────────────────────────────────────────────────────
curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://YOUR_SERVER_IP:8081/api/nodes

# ── List namespaces ───────────────────────────────────────────────────────────
curl -H "Authorization: Bearer YOUR_TOKEN" \
     http://YOUR_SERVER_IP:8081/api/namespaces

# ── Events for a specific pod ─────────────────────────────────────────────────
curl -H "Authorization: Bearer YOUR_TOKEN" \
     "http://YOUR_SERVER_IP:8081/api/events?field_selector=involvedObject.name%3Dmy-pod-7d9f8b-xkqzp"

# ── Verify auth rejection (expect 401) ───────────────────────────────────────
curl -H "Authorization: Bearer wrong-token" \
     http://YOUR_SERVER_IP:8081/api/pods
```

---
