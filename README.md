# ML Monitoring

## Table of Contents

- [Prerequisites](#prerequisites)  
- [Getting Started](#getting-started)  
  - [Clone the Repo](#clone-the-repo)  
  - [Clean Slate (Optional)](#clean-slate-optional)  
  - [Run with Docker Compose](#run-with-docker-compose)  
- [Database & Seed Data](#database--seed-data)  
- [API Endpoints](#api-endpoints)  
- [Running Tests](#running-tests)  
- [Project Structure](#project-structure)  
- [Further Development](#further-development)  

---

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) (20.10+)  
- [Docker Compose](https://docs.docker.com/compose/) (1.29+)  
- [Go](https://golang.org/dl/) (1.20+) **only for running unit tests locally**  

---

## Getting Started

### Clone the Repo

```bash
git clone https://github.com/Olt-Kondirolli91/ml-monitoring.git
cd ml-monitoring
```

### Clean Slate (Optional)

If you have run this before and want to pick up new seed data, remove the old Postgres volume:

```bash
docker-compose down
docker volume rm ml-monitoring_db_data
```

> This ensures Postgres re‑initializes and runs your migrations & seed inserts.

### Run with Docker Compose

```bash
docker-compose up --build
```

- **Postgres** listens on `localhost:5432` (user: `postgres`, password: `postgres`, db: `postgres`).  
- **Go app** starts on port **8080** inside the container.

You should see logs:

```
ml_monitoring_db   | database system is ready to accept connections
ml_monitoring_app  | Database migrated successfully!
ml_monitoring_app  | Starting HTTP server on port 8080
```

---

## Database & Seed Data

Migrations live in `migrations/20250408001_init_schema.up.sql`. They:

1. Create `inferences` and `feedback` tables  
2. Define indexes and FK constraints  
3. **Seed** two inferences and one feedback row:

| Inference ID                             | model_name | model_version | has_feedback |
|------------------------------------------|------------|---------------|--------------|
| `c3c4c350-7d8a-4b02-816c-245ced77ff01`   | seed_model | 1.0.0         | false        |
| `11111111-2222-3333-4444-555555555555`   | seed_model | 1.0.0         | true         |

One feedback row:

| Feedback ID                              | Inference ID                          |
|------------------------------------------|---------------------------------------|
| `66666666-7777-8888-9999-000000000000`   | `11111111-2222-3333-4444-555555555555`|

To inspect:

```bash
psql -h localhost -p 5432 -U postgres -W
\c postgres
SELECT * FROM inferences;
SELECT * FROM feedback;
```

---

## API Endpoints

All endpoints are JSON over HTTP on port **8080** (container) → **localhost:8080** (host).

### Health Check

```
GET /health
```

Response:
```json
{"status":"ok"}
```

### Create Inference

```
POST /inferences
Content-Type: application/json

{
  "model_name":    "my_model",
  "model_version": "1.2.3",
  "input_data":    { ... },
  "output_data":   { ... }
}
```

Response `201 Created`:
```json
{"inference_id":"<uuid>"}
```

### Get Inference by ID

```
GET /inferences/{inference_id}
```

Response `200 OK`:
```json
{
  "id":"<uuid>",
  "model_name":"my_model",
  "model_version":"1.2.3",
  "input_data":"{...}",
  "output_data":"{...}",
  "created_at":"2025-04-10T...",
  "has_feedback":false
}
```

`404 Not Found` if ID doesn’t exist.

### Create Feedback

```
POST /inferences/{inference_id}/feedback
Content-Type: application/json

{
  "feedback_data": { ... }
}
```

Response `201 Created`:
```json
{"feedback_id":"<uuid>"}
```

Also sets `has_feedback=true` on the inference.

### Get Feedback for Inference

```
GET /inferences/{inference_id}/feedback
```

Response `200 OK`:
```json
[
  {
    "id":"<uuid>",
    "inference_id":"<uuid>",
    "feedback_data":"{...}",
    "created_at":"2025-04-10T..."
  },
  ...
]
```

Empty array `[]` if no feedback.

---

## Running Tests

All unit tests live under `tests/` and mock out DB or HTTP repos. No containers needed.

```bash
go test ./tests/... -v
```

You should see all tests **PASS**:

```
PASS
ok   github.com/Olt-Kondirolli91/ml-monitoring/tests
```

To run **all** tests (including any future integration tests):

```bash
go test ./... -v
```

---

## Project Structure

```
ml-monitoring/
├── cmd/
│   └── main.go               # Entry point: loads config, runs migrations, starts server
├── internal/
│   ├── config/
│   │   └── config.go         # Env var loader
│   ├── db/
│   │   ├── db.go             # DB connect
│   │   └── migrations.go     # golang-migrate runner
│   ├── models/
│   │   ├── inference.go      # Inference struct
│   │   └── feedback.go       # Feedback struct
│   ├── repository/
│   │   ├── inference_repo.go # SQL CRUD for inferences
│   │   └── feedback_repo.go  # SQL CRUD for feedback
│   └── server/
│       ├── server.go         # HTTP router & startup
│       └── handlers.go       # HTTP handler implementations
├── migrations/
│   ├── 20250408001_init_schema.up.sql   # Schema + seed data
│   └── 20250408001_init_schema.down.sql # Drop tables
├── tests/                     # Unit tests (sqlmock & in‑memory mocks)
│   ├── repo_inference_test.go
│   ├── repo_feedback_test.go
│   ├── mock_repository.go
│   └── server_test.go
├── Dockerfile                 # Multi‑stage build for Go app
├── docker-compose.yml         # Postgres + app services
├── go.mod & go.sum            # Go module
└── README.md                  # This file
```

---
