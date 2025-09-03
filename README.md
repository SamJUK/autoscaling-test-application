# Autoscaling Test Application

Web server to test scaling functionality. 

## Configuration

Configuration is managed via Environment Variables

ENV_NAME | DEFAULT | Purpose
--- | --- | ---
CONNECTION_COUNT | `10` | Limit for HTTP requests
REQUEST_TIME | `3` | Time to sleep each request for
LISTEN_ADDRESS | `0.0.0.0` | Address to listen on
LISTEN_PORT | `80` | Port to listen on

## Example Usage

### K8s

### Docker

### Other
