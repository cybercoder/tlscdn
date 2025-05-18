# TLSCDN - Kubernetes-Native CDN Controller for TLS Edge Routing

TLSCDN controller is a Kubernetes controller that enables edge routing for http traffic in a Content Delivery Network (CDN) setup. It watches for CdnGateway and CdnHTTPRoute custom resources in your Kubernetes cluster and dynamically configures the underlying CDN infrastructure through a Redis database.

## Overview

TLSCDN bridges the gap between Kubernetes and CDN infrastructure by:

1. Watching for custom resource definitions (CRDs) in your Kubernetes cluster
2. Processing Gateway and HTTPRoute resources.
3. Storing the routing configuration in Redis.
4. Supporting various load-balancing methods and upstream configurations.

## Installation

### Prerequisites

- Kubernetes cluster (v1.19+)
- Redis server
- `kubectl` configured to connect to your cluster

### Deploying the Controller

1. Clone this repository:
   ```
   git clone https://github.com/cybercoder/tlscdn.git
   cd tlscdn
   ```

2. Install the CRDs:
   ```
   kubectl apply -f config/crds/
   ```

3. Deploy the controller:
   ```
   kubectl apply -f config/deployment/
   ```

## Usage

TLSCDN uses two custom resources to define your CDN routing:

1. **CdnGateway**: Defines upstream servers and configurations
2. **CdnHTTPRoute**: Maps paths to specific upstreams defined in a gateway

### Example: Creating a Gateway

```yaml
apiVersion: cdn.ik8s.ir/v1alpha1
kind: CdnGateway
metadata:
  name: example
spec:
  upstreams:
    - name: default
      hostHeader: example.com
      servers:
        - protocol: "https"
          port: 443
          address: "example.com"
          weight: 1
```

### Example: Creating an HTTP Route

```yaml
apiVersion: cdn.ik8s.ir/v1alpha1
kind: CdnHTTPRoute
metadata:
  name: example
spec:
  gateway:
    name: example
  upstreamName: default
  lb_method: ip_hash
  path:
    type: prefix
    path: "/"
```

## Configuration

The controller can be configured with the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `REDIS_HOST` | Redis server host | `localhost` |
| `REDIS_PORT` | Redis server port | `6379` |
| `REDIS_PASSWORD` | Redis server password | `` |
| `REDIS_DB` | Redis database index | `0` |
| `CDN_HOSTNAME` | Default CDN hostname | `cdntls.ir` |

## Architecture

TLSCDN follows a Kubernetes controller pattern:

1. The controller watches for CdnGateway and CdnHTTPRoute resources
2. When new resources are added, appropriate handlers are triggered
3. Configuration is stored in Redis with keys derived from resource UIDs
4. The CDN's edge servers can query Redis to determine routing rules

## Development

### Prerequisites

- Go 1.24+
- kubectl
- Access to a Kubernetes cluster (or minikube/kind for local development)
- Redis instance

### Building from Source

```
git clone https://github.com/cybercoder/tlscdn.git
cd tlscdn
go build -o tlscdn .
```

### Running Locally

```
export REDIS_HOST=your-redis-host
export CDN_HOSTNAME=your-cdn-hostname
./tlscdn
```

## License

[Include your license information here]

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
