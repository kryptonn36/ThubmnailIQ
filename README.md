# ThumbnailIQ — Getting Started (Testing & Debugging)

ThumbnailIQ is a multi-service application that analyzes YouTube thumbnails and predicts their performance. It provides a dashboard UI, backend API, background workers, and computer vision microservices to score thumbnails based on various factors like color, clarity, face detection, competitor analysis, and more.

This repository contains a runnable, zero‑API‑key version of the system (using synthetic data and heuristics) that you can run locally to explore the codebase and understand how the pieces fit together.

--- 

## Table of Contents
- [Project Overview](#project-overview)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [One-time Setup](#one-time-setup)
- [Running with Docker (Recommended)](#running-with-docker-recommended)
- [Running without Docker](#running-without-docker)
- [Verifying the Setup](#verifying-the-setup)
- [Debugging Tips](#debugging-tips)
- [Running Checks](#running-checks)
- [Known Simplifications](#known-simplifications)
- [Project Layout](#project-layout)

---

## Project Overview

ThumbnailIQ consists of five core services:

| Service   | Technology | Port | Purpose |
|-----------|------------|------|---------|
| `web`     | Next.js 14 | 3000 | Dashboard UI for registering users, uploading thumbnails, and viewing analysis results. |
| `api`     | Go / Gin   | 8080 | REST API handling authentication, analysis queuing, result retrieval, and billing. |
| `worker`  | Go / Asynq | –    | Background processor that pulls analysis jobs from Redis and runs the scoring pipeline. |
| `cv-service` | Python / FastAPI | 8001 | Computer‑vision microservice that extracts OCR, face, color, and clutter features from thumbnail images. |
| `postgres`| PostgreSQL | 5432 | Primary datastore for users, workspaces, analyses, and billing records. |
| `redis`   | Redis      | 6379 | Job queue (via Asynq) and caching layer. |
| `minio`   | MinIO      | 9000/9001 | S3‑compatible object storage for uploaded thumbnails and generated previews. |

All services communicate over HTTP or via the Redis broker. The system is designed so that you can run each piece independently, making it easy to debug and extend.

---

## Architecture

```
+----------------+       +----------------+       +-----------------+
|   Web UI       | <---> |    API (Gin)   | <---> |   Postgres DB   |
+----------------+       +----------------+       +-----------------+
        ^                         |
        |                         v
        |               +-----------------+
        |               |   Redis Queue   |
        |               +-----------------+
        |                         ^
        |                         |
        v                         |
+----------------+       +-----------------+
|  Worker (Asynq)|------>| CV Service (Py)|
+----------------+       +-----------------+
        ^                         |
        |                         v
        |                 +-----------------+
        |                 |    MinIO S3     |
        |                 +-----------------+
        |
        v
+-----------------+
|   Docker Compose|
+-----------------+
```

*The arrows indicate request/data flow. The web UI talks to the API; the API enqueues jobs in Redis; the worker consumes jobs and calls the CV service; the CV service reads/writes thumbnails from MinIO and returns results to the worker.*

---

## Prerequisites

- **Go** 1.22+ (for `api` and `worker`)
- **Node.js** 20+ (for `web`)
- **Python** 3.12+ (for `cv-service`)
- **Docker & Docker Compose** (recommended for running all infra services)
- Optional CLI tools (only needed if you plan to modify migrations or queries):
  - `goose` (DB migration tool)
  - `sqlc` (SQL query code generator)

Install them with:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

---

## One-time Setup

```bash
# From the repository root:
cp .env.example .env            # Adjust if you want real API keys (YouTube, Gemini, Razorpay, Stripe)
cd web && cp .env.local.example .env.local && cd ..
```

Edit `.env` and `.env.local` as needed. The defaults work for a fully local, zero‑API‑key run.

---

## Running with Docker (Recommended)

Docker Compose brings up the dependent infra services (Postgres, Redis, MinIO, and the CV service) in one step.

```bash
# Start Postgres, Redis, MinIO, and cv-service
make infra-up

# Apply any pending database migrations
make migrate-up
```

Then, in **three separate terminal tabs/windows**, run the three main processes:

```bash
# Terminal 1: API server
make api      # -> http://localhost:8080

# Terminal 2: Background worker
make worker   # (no HTTP port; processes the analysis queue)

# Terminal 3: Web UI
make web      # cd web && npm run dev -> http://localhost:3000
```

Open **http://localhost:3000** in your browser, register an account, and upload a thumbnail under **“New Analysis.”** A workspace is auto‑created on registration, so no further configuration is needed.

To stop everything:
```bash
make infra-down   # stops Docker containers
# Press Ctrl+C in each of the three terminal tabs to stop api, worker, and web.
```

---

## Running without Docker

If you don’t have Docker, you can run each infra service manually as local processes. This is how the project was originally verified.

### 1. Postgres
```bash
# If you already have a local Postgres server, just create the role and DB:
psql -U <youruser> -d postgres -c "CREATE ROLE thumbnailiq LOGIN PASSWORD 'thumbnailiq' CREATEDB;"
psql -U <youruser> -d postgres -c "CREATE DATABASE thumbnailiq OWNER thumbnailiq;"
```

### 2. Redis
```bash
# Install redis-server if missing (e.g., apt-get install redis-server) or use a static binary:
/usr/bin/redis-server --daemonize yes --port 6379
```

### 3. MinIO (S3‑compatible storage)
```bash
curl -sSL -o minio https://dl.min.io/server/minio/release/linux-amd64/minio
chmod +x minio
MINIO_ROOT_USER=minioadmin MINIO_ROOT_PASSWORD=minioadmin \
  ./minio server /tmp/minio-data --address ":9000" --console-address ":9001" &
```

### 4. CV Service (Python/FastAPI)
```bash
cd services/cv
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
uvicorn app.main:app --host 0.0.0.0 --port 8001 &
```

### 5. Apply DB migrations and run the apps
```bash
make migrate-up   # requires the Postgres DB from step 1

# Then, in three terminals:
make api
make worker
make web
```

---

## Verifying the Setup (cURL walkthrough)

You can verify the backend works without touching the UI by using curl:

```bash
# 1️⃣ Register a user (this also creates a workspace)
curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"you@example.com","password":"password123","full_name":"Your Name"}'

# Save the returned access_token; you’ll need it for the next calls.
TOKEN="<paste the access_token here>"

# 2️⃣ Upload a thumbnail for analysis
curl -s -X POST http://localhost:8080/api/v1/analyses \
  -H "Authorization: Bearer $TOKEN" \
  -F "keyword=how to lose weight fast" \
  -F "thumbnail=@/path/to/some.jpg;type=image/jpeg"

# Response includes an analysis ID, e.g. {"id":"...","status":"pending",...}
ANALYSIS_ID="<id from previous response>"

# 3️⃣ Poll until the worker finishes processing (usually 1–2 seconds)
while true; do
  RESULT=$(curl -s http://localhost:8080/api/v1/analyses/$ANALYSIS_ID \
    -H "Authorization: Bearer $TOKEN")
  echo "$RESULT" | grep -q '"status":"complete"' && break
  sleep 1
done

echo "Analysis complete!"
echo "$RESULT" | jq .
```

A `"complete"` response contains the overall `score`, six sub‑scores, `cv_results`, competitor averages, and improvement suggestions.

---

## Debugging Tips

- **Logs**: `api` and `worker` output structured logs to stdout (zerolog). Redirect to a file if running detached:  
  `go run ./cmd/api > /tmp/api.log 2>&1 &`
- **Worker not processing?** Ensure Redis is reachable (`redis-cli ping` → `PONG`) and that `make worker` is actually running.
- **CV service connection refused?** Verify `curl http://localhost:8001/health` returns `200 OK`. The API calls it synchronously; if it’s down, analyses will fail with an `error_message` you can inspect via `GET /analyses/{id}`.
- **OCR returns `text_detected: false`?** Make sure `tesseract-ocr` is installed and on the `PATH` for the CV service process. The Docker image includes it; for a bare‑metal venv run, install it via your package manager (e.g., `apt-get install tesseract-ocr`).
- **Next.js build issues?** See the “Known simplifications” section below for recurring gotchas (e.g., `next.config.ts` vs `.mjs`, SWC binary download hiccups).
- **401 errors from the web app?** The API client clears invalid tokens and redirects to `/login`. If you restarted the API with a new `JWT_ACCESS_SECRET`, old tokens in `localStorage` become invalid—just log in again.
- **Workspace‑related 400s?** Endpoints that need a `workspace_id` will derive it from the logged‑in user’s first workspace if omitted. If you manually inserted a user into the DB without creating a workspace, you’ll need to create one via the UI or API.

---

## Running Checks

```bash
# Backend
go build ./...          # compile all Go code
go vet ./...
go test ./... -race     # run existing tests (none yet, but add your own)

# Frontend
cd web && npm run build  # production build
cd web && npm run lint   # ESLint check
```

---

## Known Simplifications vs. the Full Blueprint

This runnable preview deliberately omits or replaces some heavyweight components to keep the setup lightweight and dependency‑free:

- **Computer Vision**: Uses `pytesseract` + OpenCV Haar cascades + scikit‑learn color clustering instead of EasyOCR/InsightFace/DeepFace/YOLOv8 (which would require large model downloads). Emotion detection is a smile‑cascade heuristic; object detection/clutter relies on Sobel edge density alone. See `services/cv/README.md` for details.
- **YouTube Competitor Data**: Synthetic (deterministic per keyword) by default. Set a real `YOUTUBE_API_KEY` in `.env` to fetch actual competitor thumbnails and run them through the CV service.
- **Curiosity Score**: Heuristic‑based unless `GEMINI_API_KEY` is provided, in which case it calls the Gemini API.
- **Billing**: Real Razorpay test‑mode integration (creates an order, opens the Razorpay Checkout widget, verifies HMAC). Stripe remains a mock behind the same `payment.Gateway` interface—switch `PAYMENT_PROVIDER=stripe` once you have a Stripe account.
- **Branding Sub‑score**: Fixed at 50 (the blueprint’s “3+ historical thumbnails from the same channel” signal is not collected in this build).
- **Eye‑contact / Gaze and Arrow / Object Detection**: Not available without heavier models; those scoring terms are zeroed out rather than fabricated.
- **Real‑time Updates**: The analysis detail page polls every 2 seconds instead of using WebSocket/SSE push.
- **Storage / Auth / CORS**: Dev‑grade defaults (MinIO public‑read bucket, HS256 JWT tokens, `Access-Control-Allow-Origin: *`). Tighten these before any production deployment.

---

## Project Layout

```
cmd/api, cmd/worker          # Go entrypoints for API and background worker
internal/domain              # Entities + repository interfaces
internal/usecase             # Business logic
internal/handler             # Gin HTTP handlers
internal/infra               # Clients for Postgres, Redis, S3, YouTube, Gemini, Razorpay, Stripe, CV
internal/scoring             # ThumbnailIQ score engine + suggestion catalog
internal/worker              # Asynq task handlers (the actual analysis pipeline)
db/migrations, db/queries    # Goose migrations + SQLC queries
services/cv                  # Python FastAPI computer‑vision microservice
web/                         # Next.js 14 dashboard application
```

---

*Happy hacking! If you run into trouble, check the logs, verify each service is reachable, and feel free to open an issue.*