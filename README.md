# Pack Calculator

A small Go service plus web UI that calculates how to fulfil customer orders using **whole packs only**, following the rules from the task:

1. Only whole packs can be sent. Packs cannot be broken open.  
2. Within that constraint, send out the **fewest total items** that still fulfil the order.  
3. If multiple combinations ship the same number of items, prefer the one with **fewer packs**.  

> Rule #2 takes precedence over rule #3.

The service exposes a JSON API and serves a minimal single-page UI for easy manual testing.

---

## Prerequisites

- Go (any recent version)
- `make`
- Docker (for running the containerized version)

---

## Running the project

### 1. Run tests

```bash
make test
```

Runs unit tests

---

### 2. Run the app locally

```bash
make run
```

Starts the server at:

- UI → http://localhost:8080  
- API → http://localhost:8080/api/pack-sizes

---

### 3. Build binary

```bash
make build
```

---

### 4. Build and run via Docker

```bash
make docker-build
make docker-run
```

Runs the service in a container on http://localhost:8080.

---

## API Endpoints

### GET /api/pack-sizes

Returns current pack sizes.

```bash
curl -s http://localhost:8080/api/pack-sizes
```

---

### PUT /api/pack-sizes

Replaces the current pack size configuration with a new list.

```bash
curl -s -X PUT http://localhost:8080/api/pack-sizes -H "Content-Type: application/json" -d '{"packSizes":[23,31,53]}'
```

---

### POST /api/calculate

Response includes pack breakdown, total shipped, and extra items.

```bash
curl -s -X POST http://localhost:8080/api/calculate   -H "Content-Type: application/json"   -d '{"items":263}'
```

---

## Architecture Overview

```
cmd/server/main.go   – application wiring  
internal/api         – API handlers, DTOs  
internal/calculator  – domain logic
web/index.html       – frontend UI
```

### Flow

```
UI → API → Calculator → PackResult → JSON → UI
```

---

## Persistence

For simplicity, pack sizes are stored **in-memory**.  
Restart resets pack sizes to defaults.  
Architecture allows adding DB/file storage easily.

## Deployment

This project is fully containerized and can run on any platform that supports Docker.
For demonstration purposes, it is temporarily deployed on **AWS** using:

- **Docker image** stored in **Amazon ECR**
- **AWS App Runner** (ECS-backed managed container service)

## Live Demo (Temporary)

👉 https://pack-calculator.aleksandre.net/
