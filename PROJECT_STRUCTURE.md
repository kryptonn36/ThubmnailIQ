# Project Structure — What's In Each Directory

A map of every directory in this repo, what lives in it, and where to look when something breaks.
Pairs with [GETTING_STARTED.md](GETTING_STARTED.md) (run/test/debug instructions) — this file is
about *where the code is*, that one is about *how to run and diagnose it*.

The three independently-runnable pieces are `cmd/` + `internal/` + `pkg/` (Go backend),
`services/cv/` (Python CV microservice), and `web/` (Next.js frontend). `db/` is shared
schema/queries consumed by the Go backend at build time (sqlc) and migration time (goose).

---

## `cmd/` — Go binaries (entrypoints only)

Each subdirectory is a `main.go` that wires dependencies together and starts a process. No
business logic lives here — if you're debugging *behavior*, go to `internal/`; if you're debugging
*startup* (config loading, DB connection, which implementation got wired in), look here.

- **`cmd/api/main.go`** — starts the HTTP server (port from `SERVER_PORT`, default 8080). Loads
  config, connects to Postgres/S3, decides real-vs-mock YouTube client based on whether
  `YOUTUBE_API_KEY` is set, builds every usecase and handler, and calls `router.Run()`. If the API
  won't start, the error is almost always logged right here (bad `DATABASE_URL`, S3 endpoint
  unreachable, etc.) before it ever reaches `internal/`.
- **`cmd/worker/main.go`** — starts the Asynq worker pool (no HTTP port). Same config/wiring
  pattern as the API, but builds `internal/worker.AnalysisHandler` and registers it against task
  type `thumbnail:analyze`. If analyses get stuck on `"pending"` forever, check this process is
  actually running and connected to the same Redis as the API.

## `internal/` — all backend business logic (not importable outside this module)

Organized as a fairly standard clean-architecture layering: `domain` → `usecase` → `handler`/
`worker`, with `infra` providing concrete implementations of the interfaces `domain` declares.

### `internal/domain/` — entities + repository interfaces, no behavior

One subpackage per aggregate: `user`, `workspace`, `analysis`, `competitor`, `billing`, `viraldb`.
Each file defines plain structs (with `json:"..."` tags — these tags are what shapes the JSON the
API actually returns, so if a field is missing/misnamed in an API response, check the tag here
first) and a `Repository` interface. **No SQL, no HTTP, no business rules** — if you find logic
here, it's misplaced. This is also where `billing.Plans` (the hardcoded Free/Starter/Pro/Agency
plan catalog) lives.

### `internal/usecase/` — orchestration logic

One subpackage per use-case area: `user` (register/login/refresh), `workspace` (create/invite/
list), `analysis` (create + enqueue, list, add compare version), `billing` (subscribe, API keys),
`tracking`, `viraldb`. Each `Usecase` struct takes `domain.Repository` interfaces (not concrete
infra types) as dependencies — this is the layer to read when you want to understand *what happens*
on a given API call, before any HTTP or SQL specifics. Notably:
- `usecase/user/user.go` — registration also auto-creates a workspace and an "owner" membership
  row. If a fresh user has no workspace, this is the first place to check.
- `usecase/analysis/analysis.go` — `Create()` uploads to S3 and enqueues the Asynq job but does
  **not** score anything itself (that's `internal/worker`). `AddCompareVersion()` *does* score
  synchronously (single CV call, reuses the parent analysis's stored competitor average).

### `internal/handler/` — HTTP layer (Gin)

One file per resource (`auth_handler.go`, `analysis_handler.go`, `billing_handler.go`, etc.).
Handlers parse the request, call a usecase, shape the JSON response — they should have basically
no business logic. Two files are cross-cutting:
- **`respond.go`** — maps domain sentinel errors (`pkg/errors`) to HTTP status codes. If an
  endpoint returns the wrong status code for an error case, this is the single place to fix it.
- **`workspace_resolve.go`** — every handler that needs a `workspace_id` calls
  `resolveWorkspaceID()` instead of requiring the client to send one explicitly. If omitted, it
  falls back to the caller's first workspace. This is why the web app never sends `workspace_id`
  anywhere and things still work — **if you add a new workspace-scoped endpoint, use this helper**
  rather than `uuid.Parse(c.Query("workspace_id"))` directly, or multi-workspace-naive clients will
  break against it.

### `internal/worker/` — the actual analysis pipeline

- **`tasks.go`** — Asynq task-type constants and payload structs/constructors
  (`NewAnalyzeThumbnailTask`, etc.). The API (`usecase/analysis`) and worker both import this file
  so they agree on task names/payloads without depending on each other.
- **`analyze_thumbnail.go`** — `AnalysisHandler.HandleAnalyzeThumbnail` is the real pipeline:
  load the analysis row → mark `processing` → CV-analyze the user's thumbnail → fetch N
  competitors (real or mock, via the `youtube.Fetcher` interface) → CV-analyze/score each →
  compute the competitor average → call Gemini (or heuristic) for curiosity → run
  `internal/scoring` for all six sub-scores → build suggestions → persist everything via
  `UpdateResults`. **This is the first file to read when an analysis comes back wrong, scored
  oddly, or with the wrong competitor count.** Failures here set `status: "failed"` with
  `error_message` populated on the analysis row — always check that field before guessing.

### `internal/scoring/` — pure functions, no I/O

- **`engine.go`** — `Visibility`, `Contrast`, `Attention`, `Mobile`, `Branding`, `FinalScore`. Each
  sub-score's formula has a comment explaining any deviation from the blueprint's original spec
  (e.g. why `Branding` is a fixed 50, why `Attention` doesn't have a real arrow-detection term).
  If a score looks wrong, the bug is either here (formula) or in the CV result feeding it
  (`internal/infra/cv`).
- **`suggestions.go`** — threshold-based rules that turn CV metrics + sub-scores into the ranked
  `suggestions[]` list. Add a new suggestion type here.
- **`types.go`** — `CompetitorAvg`, `SubScores`, `Suggestion`, and the `clamp()` helper.

### `internal/infra/` — concrete implementations of domain interfaces

- **`postgres/`** — one `*_repo.go` per domain aggregate, implementing the corresponding
  `domain.Repository` interface using the sqlc-generated code in `postgres/db/` (gitignored from
  this listing, regenerate with `make sqlc-generate` — **never hand-edit those files**).
  `convert.go` holds the `pgtype.X` ↔ Go-native conversion helpers (`textVal`, `int4Ptr`, etc.)
  used by every repo — if a field comes back zero/null when it shouldn't, check (a) the repo's
  `toDomainX()` mapping function actually assigns that field, and (b) the right convert helper is
  used. Two real bugs of exactly this shape were caught during initial verification: a missing
  `CreatedBy` mapping in `competitor_repo.go`, and a workspace-membership row that was queried
  before it was created in the user usecase.
- **`redis/`** — `cache.go` (generic get/set/del) and `rate_limiter.go` (fixed-window limiter,
  not currently wired into any route but ready to attach via `internal/middleware`).
- **`s3/storage.go`** — wraps `aws-sdk-go-v2`'s S3 client, pointed at MinIO by default
  (`S3_ENDPOINT`). `EnsurePublicReadBucket()` sets a public-read bucket policy at API startup so
  uploaded thumbnails are fetchable by URL without presigning — fine for dev, not for production.
- **`youtube/`** — `client.go` (real YouTube Data API v3 + per-thumbnail CV analysis) and
  `mock.go` (deterministic synthetic competitors, seeded by keyword hash, used whenever
  `YOUTUBE_API_KEY` is unset). Both implement the `Fetcher` interface in `types.go`. **If
  competitor data looks suspiciously repeatable/fake, you're on the mock path — check
  `YOUTUBE_API_KEY`.**
- **`cv/`** — `client.go` is the HTTP client the Go side uses to call `services/cv`'s `/analyze`
  endpoint; `types.go` mirrors that service's JSON response shape exactly. If you change the
  Python service's response shape, update this struct too or fields will silently zero-out.
- **`gemini/client.go`** — real Gemini API call if `GEMINI_API_KEY` is set, otherwise
  `heuristicCuriosity()` (keyword/number/emotion presence scoring). Any Gemini API failure also
  falls back to the heuristic rather than failing the whole analysis.
- **`payment/`** — billing runs against the `payment.Gateway` interface
  (`internal/domain/payment`), not a concrete provider. `payment/razorpay/client.go` is the live
  implementation (real Orders API + HMAC signature verification) used while Stripe access is
  invite-gated. `payment/stripe/client.go` mints a fake order ID and always verifies — swap its
  body for the real `stripe-go` SDK once your Stripe account is approved, then flip
  `PAYMENT_PROVIDER=stripe`; no usecase/handler changes needed.

### `internal/middleware/`, `internal/server/`, `internal/config/`

- **`middleware/auth.go`** — JWT bearer parsing; sets `user_id`/`email` in the Gin context.
- **`middleware/cors.go`** — wide-open CORS (`*`) for dev convenience.
- **`middleware/rate_limit.go`** — Redis-backed limiter, implemented but not yet attached to any
  route in `server/router.go`.
- **`server/router.go`** — the single source of truth for every route → handler mapping. **Start
  here when a request 404s** — if the route isn't listed in this file, it doesn't exist.
- **`config/config.go`** — Viper config with defaults for every setting (so the app runs with zero
  env vars). Env var names are the dotted key with dots replaced by underscores, e.g.
  `youtube.api_key` → `YOUTUBE_API_KEY`. If an env var doesn't seem to take effect, check the
  struct tag (`mapstructure:"..."`) matches what you're exporting.

## `pkg/` — small standalone helpers (no domain knowledge)

`jwt/` (token issuing/parsing), `hash/` (bcrypt + SHA256 + API key generation), `errors/`
(sentinel errors like `ErrNotFound`, mapped to HTTP codes in `internal/handler/respond.go`),
`validator/` (email/password/slug checks), `pagination/`, `logger/` (zerolog setup). These have no
dependency on `internal/`, so they're safe to reuse as-is if you ever split this into more
services.

## `db/` — schema and queries (Go-only; sqlc/goose consume this)

- **`migrations/`** — goose-managed, sequential, paired up/down SQL files. Run `make migrate-up`/
  `make migrate-down`/`make migrate-status`. **Never edit an already-applied migration** — add a
  new one instead, even for a one-line fix.
- **`queries/`** — sqlc query files (`-- name: X :one/:many/:exec` annotations). Editing these
  requires `make sqlc-generate` (or `cd db && sqlc generate`) to regenerate
  `internal/infra/postgres/db/*.sql.go` before the Go code will compile against your changes.
- **`sqlc.yaml`** — sqlc config; notably overrides Postgres `uuid` columns to Go's
  `google/uuid.UUID` instead of `pgtype.UUID` for ergonomics. `jsonb` columns come through as
  `[]byte` (sqlc's pgx/v5 default) — that's why `analysis.CVResults` etc. are `[]byte` end to end.

## `services/cv/` — Python computer-vision microservice (independent of Go)

- **`app/main.py`** — FastAPI app instantiation, mounts the router.
- **`app/routes.py`** — `/health` and `/analyze` (accepts multipart upload *or* JSON
  `{"image_url": ...}`, the latter accepting a local path or http(s) URL).
- **`app/ocr.py`** — pytesseract wrapper; catches a missing-tesseract-binary error and degrades to
  `text_detected: false` instead of crashing the request.
- **`app/face.py`** — OpenCV Haar-cascade face detection + a smile-cascade heuristic standing in
  for real emotion detection.
- **`app/colors.py`** — scikit-learn KMeans dominant-color extraction + WCAG contrast calculation.
- **`app/clutter.py`** — Sobel edge-density based clutter score (no object detector — YOLO was
  skipped, see `README.md`'s "Known simplifications" section for the full list).
- **`app/models.py`** — Pydantic response models; this is the formal contract that
  `internal/infra/cv/types.go` on the Go side must match field-for-field.
- **`requirements.txt`** / **`Dockerfile`** — the Dockerfile installs `tesseract-ocr`, `libgl1`,
  `libglib2.0-0` as system deps before `pip install`; a bare venv won't have OCR working unless you
  install `tesseract-ocr` yourself (see GETTING_STARTED.md §6).

## `web/` — Next.js 14 (App Router) dashboard

- **`app/`** — one folder per route. `app/page.tsx` is the marketing landing page;
  `app/login`, `app/register` are unauthenticated; everything under `app/dashboard/` shares
  `dashboard/layout.tsx`, which is the auth guard (redirects to `/login` if no token in
  `localStorage`) and renders `Sidebar`/`Navbar`. **If a dashboard page is blank/redirecting
  unexpectedly, check `dashboard/layout.tsx` first.**
  - `dashboard/analyses/new/page.tsx` — the upload form; posts multipart to `/analyses` with no
    `workspace_id` (relies on backend auto-resolution, see `workspace_resolve.go` above).
  - `dashboard/analyses/[id]/page.tsx` — the results page; polls via `hooks/useAnalysis.ts` every
    2s until `status` leaves `pending`/`processing`.
  - `dashboard/compare/page.tsx` — multi-upload comparison UI hitting `/analyses/{id}/compare`.
- **`components/`** — `ScoreGauge`, `CVBreakdown`, `CompetitorGrid`, `SuggestionList`,
  `CompareGrid` (presentational, take parsed API data as props), `Sidebar`/`Navbar` (layout
  chrome).
- **`lib/api.ts`** — the single fetch wrapper everything goes through: attaches the bearer token,
  clears auth + redirects to `/login` on any 401, throws `ApiError` on non-2xx. **If API calls
  silently fail or redirect unexpectedly, this file's `request()` function is where to look.**
- **`lib/auth.ts`** — localStorage read/write for tokens + user object.
- **`hooks/useAuth.ts`**, **`hooks/useAnalysis.ts`** — thin wrappers around `lib/api.ts` for the
  auth flow and the polling-for-completion flow respectively.
- **`types/index.ts`** — TypeScript interfaces mirroring the Go API's JSON shapes. If you change a
  field name/shape in a Go handler's response, update this file or TypeScript won't catch the
  mismatch (it's not generated from the Go types — there's no shared schema).
- **`next.config.mjs`** — deliberately `.mjs`, not `.ts`: Next.js 14.2.5 doesn't support
  `next.config.ts`. Don't reintroduce a `.ts` version of this file.

## Top-level files

- **`docker-compose.yml`** — `postgres`, `redis`, `minio`, `cv-service` only. The Go `api`/`worker`
  and the Next.js `web` app are deliberately run via `go run`/`npm run dev` (see the `Makefile`)
  rather than containerized, for faster local iteration.
- **`Makefile`** — `infra-up`/`infra-down` (docker-compose), `migrate-*` (goose),
  `sqlc-generate`, `api`/`worker`/`web` (run each service), `build`, `test`.
- **`.env.example`** — every env var the app reads, each with an explanation of what using vs.
  omitting it changes (real vs. mocked YouTube/Gemini; live Razorpay vs. mocked Stripe billing).
- **`thumbnailiq_blueprint.md`** — the original full product spec this build was scoped down
  from. Useful for understanding the *intended* full-scale design (Kubernetes, Terraform, browser
  extension, real ML models, etc.) beyond what's actually implemented here.
- **`GETTING_STARTED.md`** — run instructions, a curl walkthrough, and a debugging FAQ for the
  specific issues hit while building this (npm binary truncation, missing tesseract, etc.).
