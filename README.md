
# go-k8s-probes

Flexible solution to fastly implement Kubernetes probes in Golang

## Prerequisites

- PostgreSQL

## Run

1. Start PostgreSQL

    ```bash
    make start-postgres
    ```

1. Start application

    ```bash
    make run
    ```

## Endpoints

### Application

Root URL: `localhost:8080`

| Method | URL | Description |
| --- | --- | --- |
| GET | /api/v1/products | Fetch list of products |
| GET | /api/v1/products/{id} | Fetch a product by ID |
| POST | /api/v1/products | Create a new product |
| PUT | /api/v1/products/{id} | Update an existing product retrieved by ID |
| DELETE | /api/v1/products/{id} | Delete a product by ID |

### Prometheus metrics

Root URL: `localhost:9090`

| Method | URL | Description |
| --- | --- | --- |
| GET | /metrics | Fetch Prometheus metrics |

### Kubernetes probes

Root URL: `localhost:9091`

| Method | URL | Description |
| --- | --- | --- |
| GET | /live | Fetch liveness info |
| GET | /ready | Fetch readiness info |
