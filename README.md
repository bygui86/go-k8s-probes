
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

Root URL: `localhost:8080/`

| Method | URL | Description
| --- | --- | --- |
| GET | /products | Fetch list of products |
| GET | /products/{id} | Fetch a product by ID |
| POST | /products | Create a new product |
| PUT | /products/{id} | Update an existing product retrieved by ID |
| DELETE | /products/{id} | Delete a product by ID |

### Prometheus metrics

URL: `localhost:9090/metrics`

### Kubernetes probes

Liveness URL: `localhost:9091/live`

Readiness URL: `localhost:9091/ready`
