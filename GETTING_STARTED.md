# ThumbnailIQ — Getting Started (Testing & Debugging)

This is a working build of the ThumbnailIQ blueprint (`thumbnailiq_blueprint.md`), scoped down to a
runnable multi-service stack:

| Service | Tech | Port | Purpose |
|---|---|---|---|
| `web` | Next.js 14 | 3000 | Dashboard UI |
| `api` | Go / Gin | 8080 | REST API (`/api/v1/...`) |
| `worker` | Go / Asynq | — | Background analysis pipeline |
| `cv-service` | Python / FastAPI | 8001 | OCR, face, color, clutter analysis |
| `postgres` | PostgreSQL 16 | 5432 | Primary database |
| `redis` | Redis 7 | 6379 | Queue + cache |
| `minio` | MinIO | 9000/9001 | S3-compatible thumbnail storage |

Everything runs with **zero API keys**. YouTube competitor data is synthetic (seeded by keyword,
so it's deterministic), and the "curiosity" sub-score uses a heuristic instead of a real Gemini
call. Set `YOUTUBE_API_KEY` / `GEMINI_API_KEY` / `STRIPE_SECRET_KEY` in `.env` to switch any of
those to the real thing — see `.env.example`.

Out of scope for this build (see blueprint for the full vision): browser extension, Kubernetes/
Terraform manifests, CI/CD pipelines, team comments/PDF export, viral thumbnail DB seeding.

---

## 1. Prerequisites

- Go 1.22+
- Node.js 20+
- Python 3.12+
- Docker + Docker Compose (recommended — see §3). If Docker isn't available, see §4 for a
  no-Docker fallback that was actually used to verify this build.
- `goose` and `sqlc` CLIs only if you plan to change migrations/queries:
  ```bash
  go install github.com/pressly/goose/v3/cmd/goose@latest
  go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
  ```

## 2. One-time setup

```bash
cd cmd/api && cd ../.. # (just to confirm you're at repo root)
cp .env.example .env            # edit if you want real API keys
cd web && cp .env.local.example .env.local && cd ..
```

## 3. Quick start (Docker)

```bash
make infra-up        # starts postgres, redis, minio, cv-service via docker-compose
make migrate-up       # applies db/migrations
```

Then in three separate terminals:

```bash
make api      # go run ./cmd/api      -> http://localhost:8080
make worker   # go run ./cmd/worker   (no HTTP port, processes the analysis queue)
make web      # cd web && npm run dev -> http://localhost:3000
```

Visit **http://localhost:3000**, register an account, and upload a thumbnail under
"New Analysis". A workspace is auto-created on registration, so nothing else to configure.

To tear down: `make infra-down` (stops the docker-compose services; `Ctrl+C` the three `make`
processes above).

## 4. Quick start (no Docker)

If you don't have a working Docker daemon, every infra piece can run as a plain local process.
This is exactly how this build was verified in its own development sandbox (which had no Docker
daemon at all):

```bash
# Postgres — if you already have a local postgres server running, just:
psql -U <youruser> -d postgres -c "CREATE ROLE thumbnailiq LOGIN PASSWORD 'thumbnailiq' CREATEDB;"
psql -U <youruser> -d postgres -c "CREATE DATABASE thumbnailiq OWNER thumbnailiq;"

# Redis — download a static build if redis-server isn't installed:
#   apt-get download redis-server && dpkg-deb -x redis-server*.deb /tmp/redis-local
/usr/bin/redis-server --daemonize yes --port 6379

# MinIO — single static binary, no install needed:
curl -sSL -o minio https://dl.min.io/server/minio/release/linux-amd64/minio
chmod +x minio
MINIO_ROOT_USER=minioadmin MINIO_ROOT_PASSWORD=minioadmin \
  ./minio server /tmp/minio-data --address ":9000" --console-address ":9001" &

# CV service — see services/cv/README.md for the venv setup. If tesseract-ocr
# isn't installed system-wide, OCR degrades gracefully (text_detected: false)
# rather than crashing; everything else (faces, colors, clutter) still works.
cd services/cv && python3 -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt
uvicorn app.main:app --host 0.0.0.0 --port 8001 &
```

Then run `make migrate-up`, `make api`, `make worker`, `make web` as in §3.

## 5. Verifying it works (curl walkthrough)

```bash
# Register (auto-creates a workspace)
curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"you@example.com","password":"password123","full_name":"Your Name"}'
# -> save the access_token from the response

TOKEN="<paste access_token here>"

# Upload a thumbnail (workspace_id is optional — defaults to your first workspace)
curl -s -X POST http://localhost:8080/api/v1/analyses \
  -H "Authorization: Bearer $TOKEN" \
  -F "keyword=how to lose weight fast" \
  -F "thumbnail=@/path/to/some.jpg;type=image/jpeg"
# -> {"id": "...", "status": "pending", ...}

# Poll until status flips to "complete" (worker processes it within ~1-2s)
curl -s http://localhost:8080/api/v1/analyses/<id> -H "Authorization: Bearer $TOKEN"
```

A `"complete"` response includes `score`, the six sub-scores, `cv_results`, `competitor_avg`,
`competitors[]`, and a ranked `suggestions[]` list.

## 6. Debugging

**Logs**: `api` and `worker` log to stdout (zerolog, pretty-printed in dev). If running detached,
redirect to a file: `go run ./cmd/api > /tmp/api.log 2>&1 &`.

**Worker not processing jobs**: check Redis is reachable (`redis-cli ping` should return `PONG`)
and that `make worker` is actually running — the API only enqueues; nothing scores without it.

**`cv-service` returns connection refused**: check `curl http://localhost:8001/health`. The Go
side calls it synchronously per-thumbnail; if it's down, analyses will land in `status: "failed"`
with `error_message` populated — check that field via `GET /analyses/{id}`.

**OCR always returns `text_detected: false`**: tesseract-ocr isn't installed/on `PATH` for the
`cv-service` process. The Docker image installs it automatically; for a bare venv run, install
`tesseract-ocr` (`apt-get install tesseract-ocr` or download a `.deb` and extract it without root —
see services/cv/README.md). Everything else in the pipeline still scores normally without it.

**Next.js build fails with "Configuring Next.js via 'next.config.ts' is not supported"**: this
project's Next.js version (14.2.5) only supports `next.config.mjs`/`.js`, not `.ts` — already
fixed in this repo (`web/next.config.mjs`). If you see this again after pulling changes, make sure
no stray `next.config.ts` exists alongside it.

**`npm install` produces a working build but `next build` segfaults with "Bus error"**: this means
`node_modules/@next/swc-linux-x64-gnu/*.node` got truncated mid-download (seen on flaky/sandboxed
networks where npm's tarball stream cuts off but still reports success). Fix:
```bash
curl -sSL -o /tmp/swc.tgz https://registry.npmjs.org/@next/swc-linux-x64-gnu/-/swc-linux-x64-gnu-14.2.5.tgz
tar -xzf /tmp/swc.tgz -C /tmp/swc-extract
cp /tmp/swc-extract/package/next-swc.linux-x64-gnu.node web/node_modules/@next/swc-linux-x64-gnu/
```
Or just delete `node_modules` and `package-lock.json` and re-run `npm install` — usually resolves
on retry.

**401s from the web app**: the API client (`web/lib/api.ts`) clears the token and redirects to
`/login` on any 401. If you registered against a different `api` instance (e.g. restarted with a
new `JWT_ACCESS_SECRET`), old tokens in `localStorage` will look "expired" — just log in again.

**Workspace-related 400s when calling the API directly with curl**: every endpoint that needs a
workspace accepts an optional `workspace_id`; if omitted, it resolves to the caller's first
workspace. If you get `"could not resolve workspace_id: not found"`, the user truly has no
workspace — shouldn't happen via normal registration, but can happen if you create a user directly
in the DB.

## 7. Running checks

```bash
go build ./...          # backend compiles
go vet ./...
go test ./... -race     # no test suite exists yet beyond what you add
cd web && npm run build  # frontend production build
cd web && npm run lint
```

## 8. Known simplifications vs. the full blueprint

- **Computer vision**: uses `pytesseract` + OpenCV Haar cascades + scikit-learn color clustering
  instead of EasyOCR/InsightFace/DeepFace/YOLOv8 (those need GB-scale model downloads). Emotion
  detection is a smile-cascade heuristic, not a real classifier. Object detection/clutter from
  YOLO is replaced by Sobel edge density alone. See `services/cv/README.md` for details.
- **YouTube competitor data**: synthetic by default (deterministic per keyword), used only when
  `YOUTUBE_API_KEY` is unset. With a real key, it calls the actual YouTube Data API v3 and
  analyzes real competitor thumbnails through `cv-service`.
- **Curiosity score**: heuristic (keyword/number/emotion presence) unless `GEMINI_API_KEY` is set,
  in which case it calls the real Gemini API.
- **Billing**: Stripe is mocked — `Subscribe` mints a fake subscription ID and immediately marks
  it active. No real payment flow.
- **Branding sub-score**: fixed at a neutral 50; the blueprint's "3+ historical thumbnails from
  the same channel" signal isn't collected in this build.
- **Eye-contact/gaze and arrow/object detection**: not available without heavier models; those
  scoring terms are zeroed out rather than fabricated (see comments in
  `internal/scoring/engine.go`).
- **Real-time updates**: the analysis detail page polls every 2s instead of using the blueprint's
  WebSocket/SSE push.
- **Storage/Auth/CORS**: dev-grade only — MinIO with a public-read bucket policy, HS256 JWT, and
  `Access-Control-Allow-Origin: *`. Tighten all of this before any real deployment.

## 9. Project layout

```
cmd/api, cmd/worker        Go entrypoints
internal/domain            Entities + repository interfaces
internal/usecase           Business logic
internal/handler           Gin HTTP handlers
internal/infra             Postgres/Redis/S3/YouTube/Gemini/Stripe/CV clients
internal/scoring           ThumbnailIQ score engine + suggestion catalog
internal/worker            Asynq task handlers (the actual analysis pipeline)
db/migrations, db/queries  goose migrations + sqlc queries
services/cv                Python FastAPI computer-vision microservice
web/                       Next.js dashboard
```
