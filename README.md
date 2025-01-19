
# 5G Full-Stack Virtualization

This repository provides a **containerized microservices** architecture to demonstrate a conceptual **5G NR** environment with **network slicing**. It includes services for **Radio Resource Management (RRM)**, **MAC scheduling**, **PDCP** data forwarding, a dedicated **Slice Manager**, and an **API Gateway** with **JWT-based authentication**.  
Everything is orchestrated in **Kubernetes**, with **Prometheus**-based monitoring and an optional **Horizontal Pod Autoscaler** setup.

## Table of Contents

- [Key Features](#key-features)  
- [Architecture Overview](#architecture-overview)  
- [Repository Structure](#repository-structure)  
- [Prerequisites](#prerequisites)  
- [Setup & Deployment](#setup--deployment)  
  - [Local Environment (Docker Compose)](#local-environment-docker-compose-optional)  
  - [Kubernetes Deployment](#kubernetes-deployment)  
- [Usage](#usage)  
  - [Creating & Deleting Slices](#creating--deleting-slices)  
  - [Monitoring & Metrics](#monitoring--metrics)  
  - [Authentication (JWT)](#authentication-jwt)  
- [Extending the Project](#extending-the-project)  
- [License](#license)

---

## Key Features

1. **Network Slicing**  
   - Create, modify, and delete slices with custom bandwidth and priority attributes.  
2. **Microservices-Based Architecture**  
   - **RRM Service** for allocating resources, **MAC Service** for scheduling, **PDCP Service** for data forwarding, plus a **Slice Manager** and **API Gateway**.  
3. **Containerization & Orchestration**  
   - Docker images for each service; Kubernetes manifests for deployments, services, and an optional Horizontal Pod Autoscaler.  
4. **Observability**  
   - **Prometheus metrics** exported from each microservice; can integrate with Grafana for real-time monitoring.  
5. **Security & Authentication**  
   - **JWT-based auth** at the API Gateway.  
6. **Example Scheduling Logic**  
   - Simulated scheduling in the MAC layer; an RRM layer that “allocates” bandwidth resources.  
7. **Slice Manager**  
   - Dedicated microservice orchestrating slice creation/deletion across the RRM, MAC, and PDCP services.

---

## Architecture Overview

```
                       +----------------------+
                       |     API Gateway     | <-- JWT Authentication
                       +---------+------------+
                                 |
                                 v
    +--------------+     +--------------+     +---------------+
    |  RRM Service | <-> |  MAC Service | <-> | PDCP Service  |
    +--------------+     +--------------+     +---------------+
                                 ^
                                 |
                       +----------------------+
                       |   Slice Manager     |
                       +----------------------+

    - All services are containerized and communicate via internal Service URLs in Kubernetes
    - Prometheus scrapes /metrics endpoints for real-time stats
    - Horizontal Pod Autoscaler can scale selected microservices based on CPU or custom metrics
```

- **RRM Service**: Manages “allocated” bandwidth and resources per slice.  
- **MAC Service**: Simulates scheduling decisions, using data from the RRM.  
- **PDCP Service**: Demonstrates a simplified data-forwarding path with minimal encryption logic.  
- **Slice Manager**: Central point for orchestrating slice lifecycle.  
- **API Gateway**: Provides external endpoints for slice operations and enforces JWT-based authentication.

---

## Repository Structure

```
5g-fullstack-virtualization/
├── README.md
├── go.mod
├── go.sum
├── api-gateway/
│   ├── main.go
│   ├── gateway.go
│   ├── security.go
│   └── Dockerfile
├── rrm-service/
│   ├── main.go
│   ├── rrm.go
│   ├── metrics.go
│   └── Dockerfile
├── mac-service/
│   ├── main.go
│   ├── mac.go
│   ├── metrics.go
│   └── Dockerfile
├── pdcp-service/
│   ├── main.go
│   ├── pdcp.go
│   ├── metrics.go
│   └── Dockerfile
├── slice-manager/
│   ├── main.go
│   ├── manager.go
│   ├── metrics.go
│   └── Dockerfile
└── k8s-manifests/
    ├── namespace.yaml
    ├── api-gateway-deployment.yaml
    ├── rrm-deployment.yaml
    ├── mac-deployment.yaml
    ├── pdcp-deployment.yaml
    ├── slice-manager-deployment.yaml
    ├── hpa-example.yaml
```

---

## Prerequisites

1. **Go 1.19+**  
2. **Docker** (for building images)  
3. **Kubernetes** (Minikube, kind, or a real cluster)  
4. **kubectl** and (optionally) **Helm**  
5. **Prometheus** (for metrics, if desired)  

*(You can also adapt this for Docker Compose if you prefer local container orchestration.)*

---

## Setup & Deployment

### Local Environment (Docker Compose) (Optional)
While the code and structure is primarily for Kubernetes, you could adapt it as follows:
1. Write a `docker-compose.yml` referencing each Dockerfile and binding ports.
2. Build images via `docker compose build`.
3. Start everything with `docker compose up`.

*(This is not provided in detail here, since Kubernetes is the default approach.)*

### Kubernetes Deployment

1. **Build & Push Images**  
   From each service folder:
   ```bash
   docker build -t your-registry/rrm-service:v1 .
   docker push your-registry/rrm-service:v1
   ```
   Repeat for **mac-service**, **pdcp-service**, **slice-manager**, and **api-gateway**.

2. **Apply K8s Manifests**  
   ```bash
   kubectl apply -f k8s-manifests/namespace.yaml
   kubectl apply -f k8s-manifests/rrm-deployment.yaml
   kubectl apply -f k8s-manifests/mac-deployment.yaml
   kubectl apply -f k8s-manifests/pdcp-deployment.yaml
   kubectl apply -f k8s-manifests/slice-manager-deployment.yaml
   kubectl apply -f k8s-manifests/api-gateway-deployment.yaml
   kubectl apply -f k8s-manifests/hpa-example.yaml  # optional
   ```

3. **Validate**  
   ```bash
   kubectl get pods -n fiveg-virtualization
   kubectl get svc -n fiveg-virtualization
   ```

4. **Access the API Gateway**  
   - If `type: NodePort`, find the node port and then do:
     ```bash
     curl http://<NodeIP>:<NodePort>/health
     ```
   - Or, if `type: LoadBalancer`, use the provided external IP/host.

---

## Usage

### Creating & Deleting Slices

1. **Obtain (or generate) a JWT** (assuming you have a test token).
2. **Create a Slice**:
   ```bash
   curl -X POST \
     -H "Authorization: Bearer <your-jwt>" \
     -H "Content-Type: application/json" \
     -d '{"sliceID":"sliceA","requiredBW":"10MHz","priority":1}' \
     http://<NodeIP>:<NodePort>/createSlice
   ```
3. **Delete a Slice**:
   ```bash
   curl -X DELETE \
     -H "Authorization: Bearer <your-jwt>" \
     -H "Content-Type: application/json" \
     -d '{"sliceID":"sliceA"}' \
     http://<NodeIP>:<NodePort>/deleteSlice
   ```

### Monitoring & Metrics

- Each service (`RRM`, `MAC`, `PDCP`, `Slice Manager`, `API Gateway`) exposes Prometheus metrics at `http://<service>:<port>/metrics`.  
- Configure your Prometheus to scrape these endpoints (or use port forwarding to view them locally).

### Authentication (JWT)

- The **API Gateway** enforces **JWT-based** auth for slice management endpoints.  
- For demo purposes, the environment variable `JWT_SECRET` can be set (default: `"my-secret-key"`).  
- In production, integrate a real identity provider (Keycloak, Okta, etc.) and handle token expiration, claims, etc.
