# ThumbnailIQ — Complete Production-Grade SaaS Blueprint

> A YouTube Thumbnail Intelligence Platform that helps creators compare, analyze, score, improve, and optimize thumbnails against real competitors.

---

## Table of Contents

1. [Product Vision & Problem Statement](#1-product-vision--problem-statement)
2. [Market Research & Competitive Analysis](#2-market-research--competitive-analysis)
3. [Core Features — Deep Specification](#3-core-features--deep-specification)
4. [System Architecture](#4-system-architecture)
5. [Tech Stack — Mandatory Specification](#5-tech-stack--mandatory-specification)
6. [Database Design](#6-database-design)
7. [API Design & OpenAPI Specification](#7-api-design--openapi-specification)
8. [AI & Computer Vision Pipeline](#8-ai--computer-vision-pipeline)
9. [YouTube Data Acquisition](#9-youtube-data-acquisition)
10. [Background Job System (Asynq)](#10-background-job-system-asynq)
11. [Authentication & Security](#11-authentication--security)
12. [Billing & Monetization (Stripe)](#12-billing--monetization-stripe)
13. [Frontend Architecture](#13-frontend-architecture)
14. [Browser Extension](#14-browser-extension)
15. [Observability & Monitoring](#15-observability--monitoring)
16. [Deployment & Infrastructure](#16-deployment--infrastructure)
17. [Analytics](#17-analytics)
18. [Product Roadmap](#18-product-roadmap)
19. [Repository Structure & Claude Code Execution Plan](#19-repository-structure--claude-code-execution-plan)

---

## 1. Product Vision & Problem Statement

### The Core Problem

YouTube is the world's second-largest search engine with over 800 million videos. Click-Through Rate (CTR) is the single most important metric separating successful creators from struggling ones. Thumbnails account for 70–90% of the CTR decision a viewer makes in under 300 milliseconds.

Yet creators make thumbnail decisions in a vacuum:

- No data on what competitors are doing
- No objective scoring of visual quality
- No way to A/B test before publishing
- No insight into what design patterns win in their niche
- No historical tracking of what changes top channels make

**ThumbnailIQ** closes this gap by building an intelligence layer on top of YouTube search results, applying computer vision and AI to extract actionable insight, and giving creators a data-driven advantage before they hit publish.

### Strategic Pillars

| Pillar | Description |
|--------|-------------|
| Competitive Intelligence | Show exactly how a thumbnail performs vs. what's ranking |
| Objective Scoring | Replace gut feeling with a repeatable 0–100 score |
| Actionable Guidance | Convert analysis into specific, ranked improvement steps |
| Historical Insight | Track what changes winners make over time |
| Team Workflow | Make thumbnail decisions collaborative, not solo |

### Target Users

- **Solo Creators** (100K–10M subs): Need an edge, have budget, technically capable
- **Creator Agencies**: Manage 10–200 channels, need efficiency and white-label tools
- **Brand YouTube Teams**: Corporates running YouTube as a growth channel
- **Freelance Thumbnail Designers**: Need to prove ROI to clients
- **MCNs (Multi-Channel Networks)**: Need analytics across portfolios

---

## 2. Market Research & Competitive Analysis

### Existing Competitors

| Tool | Core Feature | Pricing | Weakness |
|------|-------------|---------|----------|
| TubeBuddy | SEO + Tags + A/B Test | $9–$49/mo | No visual AI analysis |
| VidIQ | Analytics + Keywords | $10–$99/mo | CTR focused, not thumbnail-focused |
| Canva | Design tool | Free–$15/mo | No competitor context |
| Adobe Express | Design tool | $10/mo | No YouTube-specific scoring |
| Thumblytics | Basic thumbnail test | $29/mo | No AI, limited competitors |
| Checkmate | A/B testing post-publish | $15/mo | No pre-publish analysis |
| Rapidtags | Tag research only | Free | Not thumbnail-related |

### Competitor Strengths

- **TubeBuddy**: Deep YouTube API integration, large user base, established brand trust
- **VidIQ**: Strong analytics dashboard, Chrome extension widely used
- **Canva**: Excellent UX, templates, collaborative design

### Competitor Weaknesses

- None offer **pre-publish visual AI scoring**
- None show **your thumbnail in the actual SERP context**
- None perform **face/emotion/clutter analysis** on competitors
- None track **competitor thumbnail history and changes**
- None provide a **public API for agencies**

### Feature Comparison Matrix

| Feature | ThumbnailIQ | TubeBuddy | VidIQ | Canva |
|---------|-------------|-----------|-------|-------|
| SERP preview (desktop/mobile) | ✅ | ❌ | ❌ | ❌ |
| AI visual scoring | ✅ | ❌ | ❌ | ❌ |
| Competitor face/emotion analysis | ✅ | ❌ | ❌ | ❌ |
| Contrast/clutter scoring | ✅ | ❌ | ❌ | ❌ |
| Multi-version A/B comparison | ✅ | ✅ | ❌ | ❌ |
| Historical competitor tracking | ✅ | ❌ | ❌ | ❌ |
| Browser extension overlay | ✅ | ✅ | ✅ | ❌ |
| Viral thumbnail database | ✅ | ❌ | ❌ | ❌ |
| Public REST API | ✅ | ❌ | ❌ | ❌ |
| Team collaboration | ✅ | ❌ | ❌ | ✅ |

### Competitive Moat Opportunities

1. **Proprietary scoring model** — Train on CTR data from high-performing channels to build a scoring model that predicts CTR probability
2. **Historical database** — First-mover in storing competitor thumbnail snapshots over time
3. **Niche benchmarks** — The only tool showing "how does a gaming thumbnail compare to average gaming thumbnails"
4. **Agency white-label** — White-label reports for agencies managing client channels
5. **API marketplace** — Let other tools integrate ThumbnailIQ scoring into their platforms

### Monetization Opportunities

- Direct SaaS subscriptions (primary)
- Public API usage billing
- White-label reseller program
- Agency bulk licensing
- Affiliate partnerships with Canva / Adobe Express
- Enterprise custom model training

---

## 3. Core Features — Deep Specification

### 3.1 Thumbnail Search Preview

**User Flow:**
1. User uploads a PNG/JPG thumbnail (max 2MB, 1280×720 recommended)
2. User enters a search keyword (e.g., "how to lose weight fast")
3. System calls YouTube Data API v3 search endpoint
4. Fetches top 20 results with thumbnails
5. Injects user's thumbnail at positions 1, 5, and 10 of the results
6. Renders three responsive preview frames: desktop (1920px), tablet (768px), mobile (375px)
7. User can toggle between positions and devices
8. User can see their thumbnail in dark mode and light mode contexts

**UI Components:**
- Drag-and-drop upload zone with thumbnail preview
- Keyword input with autocomplete based on YouTube Suggest API
- Device switcher (Desktop / Tablet / Mobile tabs)
- Scrollable SERP mockup with realistic YouTube Chrome (logo, search bar, sidebar, titles, channel info)
- Position selector slider (1–20)
- Dark/light mode toggle

**Backend Logic:**
- YouTube search results cached in Redis for 1 hour per keyword
- Thumbnail stored in S3, resized to YouTube dimensions (480×270 for SERP)
- Results rendered server-side as an HTML snapshot OR pixel-perfect React component

### 3.2 Competitor Thumbnail Analysis

**Data Collected Per Competitor Thumbnail:**

```
Dimensions:           Width × Height
Face Count:           0, 1, 2, 3+
Dominant Faces:       Primary face bounding box coordinates
Facial Emotion:       Happy / Surprised / Neutral / Angry / Sad / Fear
Text Detected:        Boolean
Text Strings:         Array of detected text strings
Text Density:         Percentage of image area covered by text (0–100%)
Word Count:           Number of distinct words in thumbnail
Text Size:            Average character height relative to image height
Dominant Colors:      Top 5 hex values + percentage coverage
Color Palette Type:   Monochromatic / Analogous / Complementary / Triadic
Contrast Score:       WCAG-based contrast ratio (0–21)
Brightness Score:     Average luminance (0–255)
Saturation Score:     Average saturation (0–100%)
Edge Density:         Sobel edge detection count normalized
Object Count:         Detected objects via YOLO or similar
Clutter Score:        Derived from object count + edge density + text density
Visual Complexity:    Combined composite score (0–100)
Background Type:      Solid / Gradient / Scene / Studio
Has Arrow/Pointer:    Boolean
Has Circle Callout:   Boolean
Has Price/Number:     Boolean
```

**Industry Average Generation:**
After collecting data from 20–50 competitors, compute mean, median, and standard deviation for each numeric metric. These become the "niche benchmark" for that keyword at that timestamp.

**Comparison Output:**
Each metric shows:
- User's value
- Competitor average
- Percentile rank
- Visual indicator (better / average / worse)

### 3.3 Thumbnail Ranking Engine (Proprietary Score)

The **ThumbnailIQ Score** is a composite 0–100 score built from six sub-dimensions:

#### Score Architecture

```
FINAL SCORE (0–100) = Weighted Sum of:

  Visibility Score     (25%)   — Will it stand out at thumbnail size?
  Contrast Score       (20%)   — Is text/subject legible against background?
  Attention Score      (20%)   — Does it draw the eye within 300ms?
  Mobile Score         (15%)   — Is it readable on a 375px screen?
  Branding Score       (10%)   — Is it consistent with channel style?
  Curiosity Score      (10%)   — Does it trigger a psychological need to click?
```

#### Sub-Score Algorithms

**Visibility Score:**
```
visibility = (saturation_score × 0.3) + (contrast_vs_competitors × 0.4) + (unique_color_distance × 0.3)
```
Where `unique_color_distance` measures how different the dominant color is from the average competitor dominant color using ΔE (CIE76 color difference).

**Contrast Score:**
```
contrast = wcag_contrast_ratio normalized to 0–100
         + bonus for text_contrast_ratio > 4.5
         - penalty for text overlapping busy regions
```

**Attention Score:**
```
attention = (face_present × 0.4) + (eye_contact_detected × 0.3) + (large_text_present × 0.2) + (arrow_present × 0.1)
```
Eye contact is detected by checking if detected face gaze direction points toward viewer.

**Mobile Score:**
```
mobile = 100 - text_density_penalty - clutter_penalty - small_text_penalty
text_density_penalty = max(0, (text_coverage_pct - 30) × 2)
clutter_penalty = clutter_score × 0.5
small_text_penalty = if avg_text_height < 8% of image height: 20 else 0
```

**Branding Score:**
Requires at least 3 historical thumbnails from the same channel to compute color palette consistency, font consistency, and compositional consistency.

**Curiosity Score:**
AI-driven via prompt to Claude API:
```
Score this thumbnail's curiosity factor 0–100:
- Does it create an information gap?
- Does it promise a reward/transformation?
- Does it feature an unusual or emotional situation?
- Does it use numbers or specific claims?
Return JSON: {"curiosity_score": N, "reasoning": "..."}
```

#### Score Bands

| Band | Range | Label |
|------|-------|-------|
| 🔴 Poor | 0–39 | Needs Major Work |
| 🟡 Fair | 40–59 | Room to Improve |
| 🟢 Good | 60–79 | Above Average |
| 💎 Excellent | 80–100 | Optimized |

### 3.4 Thumbnail Improvement Suggestions

After scoring, generate a ranked list of improvement actions sorted by estimated impact on CTR:

**Suggestion Types:**

```
HIGH_IMPACT:
  - add_human_face          "Add a human face showing strong emotion"
  - increase_text_size      "Your text is too small for mobile viewers"
  - improve_contrast        "Low contrast makes text hard to read"
  - reduce_clutter          "Too many elements compete for attention"
  - use_complementary_color "Your palette blends with competitor averages"

MEDIUM_IMPACT:
  - add_curiosity_gap       "No clear promise or payoff visible"
  - improve_composition     "Subject is too centered — use rule of thirds"
  - add_directional_element "Add an arrow pointing to key subject"
  - increase_saturation     "Increase color vibrancy to stand out"

LOW_IMPACT:
  - add_channel_branding    "Consistent branding builds recognition"
  - optimize_for_dark_mode  "Check thumbnail in YouTube dark mode"
```

Each suggestion includes:
- A 1-sentence headline
- A 2-3 sentence explanation with data backing it
- A before/after visual example (from the viral database or AI-generated placeholder)
- An estimated CTR impact range (e.g., "+0.5% to +1.2% CTR")

### 3.5 Multi-Version Comparison

- Upload up to 5 thumbnail variants
- System scores all variants using the same pipeline
- Side-by-side grid view with score overlays
- Highlight differences using diff visualization:
  - Heat map overlay showing where attention would land (simulated saliency)
  - Color palette comparison
  - Text density comparison
  - Face emotion comparison
- Predict the likely winner with confidence percentage
- Allow team members to vote on versions (collaborative mode)

### 3.6 Historical Competitor Tracking

**How it works:**
- User adds a YouTube channel or keyword to "Track"
- Asynq background job runs nightly
- Fetches current top-20 results for keyword OR current thumbnails for channel
- Compares to stored previous snapshot
- If thumbnail URL changes → download, analyze, store new version
- Generate diff report: which metrics changed, score delta

**UI:**
- Timeline view of a competitor channel's thumbnail changes
- Side-by-side old vs. new
- Score change chart
- Alert system: "Competitor [Channel Name] updated their thumbnail for [Video Title]"

**Retention:** 12 months of history on Pro plan, 24 months on Agency plan.

### 3.7 Viral Thumbnail Database

A searchable gallery of high-performing thumbnails collected over time:

**Indexing fields:**
- Niche / Category (Gaming, Finance, Fitness, Beauty, Tech, etc.)
- Keyword tags
- Score range
- Color palette
- Face presence
- Emotion detected
- View count range
- Channel size range

**Use cases:**
- Find inspiration for "gaming thumbnails with surprised face, red background"
- See what the top 10 thumbnails for "make money online" look like
- Export a collection of reference thumbnails

**Data sourcing:** All thumbnails collected through analysis runs across all users (with appropriate data policies).

### 3.8 Team Collaboration

**Structure:**
```
Organization
  └── Workspace (per team/client)
        └── Projects
              └── Thumbnail Sets
                    └── Versions
                          └── Comments & Annotations
```

**Roles:**
- Owner: Full control, billing, member management
- Admin: Manage projects, members
- Editor: Upload, analyze, comment
- Viewer: View results only

**Workflow features:**
- Comment threads on any analysis or comparison
- @mentions with email notifications
- Approval workflow: Draft → In Review → Approved → Published
- Version history for uploaded thumbnails
- Shared keyword tracking lists
- Export reports as PDF (branded, white-label for Agency plan)

### 3.9 Browser Extension

**Platform:** Chrome (Manifest V3), Firefox (Manifest V2 compatible)

**Features on youtube.com:**
- Overlay ThumbnailIQ score badge on every thumbnail in search results
- Click any thumbnail to open a sidebar with full analysis
- "Analyze My Thumbnail" button in YouTube Studio sidebar
- Quick compare: right-click any thumbnail → "Compare with mine"

**Technical approach:**
- Content script injects overlay DOM elements
- Extension background service worker calls ThumbnailIQ API with JWT
- Results cached in chrome.storage.local for 24 hours
- Only triggers on youtube.com domains

### 3.10 Public API

**Endpoints exposed:**
```
POST /api/v1/analyze       — Analyze a thumbnail URL or upload
GET  /api/v1/score/{id}    — Get score for a previous analysis
POST /api/v1/compare       — Compare multiple thumbnails
GET  /api/v1/competitors   — Get competitors for a keyword
POST /api/v1/track         — Start tracking a channel or keyword
GET  /api/v1/trends        — Get trending thumbnail patterns for a niche
```

**Authentication:** API keys (Bearer token), issued per organization.

**Rate limits:**
```
Free:       100 requests/month
Starter:    1,000 requests/month
Pro:        10,000 requests/month
Agency:     100,000 requests/month
Enterprise: Unlimited (SLA)
```

**Documentation:** Auto-generated Swagger UI at `/api/docs`, also published to developers.thumbnailiq.com.

---

## 4. System Architecture

### High-Level Architecture

```mermaid
graph TB
    subgraph Clients
        WEB[Next.js Web App]
        EXT[Browser Extension]
        API_CLIENT[3rd Party API Clients]
    end

    subgraph CDN
        CF[CloudFront CDN]
    end

    subgraph API Layer
        GW[API Gateway / Load Balancer]
        GO_API[Go API Server - Gin]
    end

    subgraph Background Jobs
        ASYNQ[Asynq Worker Pool]
        REDIS_Q[Redis Queue]
    end

    subgraph AI Services
        CV[Computer Vision Service]
        OCR[OCR Service - Tesseract]
        FACE[Face Detection - InsightFace]
        OBJ[Object Detection - YOLOv8]
        EMOTION[Emotion Analysis]
        LLM[Claude API - Curiosity Score]
    end

    subgraph Data Layer
        PG[(PostgreSQL)]
        REDIS_C[(Redis Cache)]
        S3[(AWS S3)]
    end

    subgraph External APIs
        YT[YouTube Data API v3]
        STRIPE[Stripe API]
        SMTP[Email - Resend]
    end

    subgraph Observability
        PROM[Prometheus]
        GRAF[Grafana]
        ZEROLOG[Zerolog JSON Logs]
    end

    WEB --> CF
    CF --> GW
    EXT --> GW
    API_CLIENT --> GW
    GW --> GO_API
    GO_API --> REDIS_C
    GO_API --> PG
    GO_API --> S3
    GO_API --> REDIS_Q
    GO_API --> YT
    GO_API --> STRIPE
    REDIS_Q --> ASYNQ
    ASYNQ --> CV
    ASYNQ --> OCR
    ASYNQ --> FACE
    ASYNQ --> OBJ
    ASYNQ --> EMOTION
    ASYNQ --> LLM
    ASYNQ --> PG
    ASYNQ --> S3
    ASYNQ --> REDIS_C
    GO_API --> PROM
    ASYNQ --> PROM
    PROM --> GRAF
```

### Service Decomposition

| Service | Language | Responsibility |
|---------|----------|----------------|
| `api-server` | Go (Gin) | REST API, auth, request routing |
| `worker` | Go (Asynq) | Background jobs, CV pipeline coordination |
| `cv-service` | Python (FastAPI) | Computer vision, OCR, face/object detection |
| `web` | Next.js | Frontend SPA |
| `extension` | TypeScript | Browser extension |

### Data Flow — Thumbnail Analysis

```mermaid
sequenceDiagram
    participant U as User
    participant API as Go API
    participant S3 as AWS S3
    participant Q as Redis Queue
    participant W as Asynq Worker
    participant CV as CV Service
    participant LLM as Claude API
    participant DB as PostgreSQL
    participant WS as WebSocket

    U->>API: POST /analyses {thumbnail, keyword}
    API->>S3: Upload thumbnail
    S3-->>API: thumbnail_url
    API->>DB: Create analysis record (status=pending)
    API->>Q: Enqueue analysis job
    API-->>U: 202 {analysis_id, status: "pending"}
    
    W->>Q: Dequeue job
    W->>S3: Download thumbnail
    W->>CV: Analyze image (OCR, face, objects, color)
    CV-->>W: cv_results
    W->>LLM: Score curiosity
    LLM-->>W: curiosity_score
    W->>API: Fetch competitor thumbnails (YouTube API)
    W->>CV: Batch analyze competitors
    CV-->>W: competitor_results[]
    W->>W: Compute ThumbnailIQ Score
    W->>DB: Update analysis record (status=complete, results)
    W->>WS: Push completion event to user session
    
    U->>API: GET /analyses/{id}
    API->>DB: Fetch complete results
    API-->>U: Full analysis response
```

---

## 5. Tech Stack — Mandatory Specification

### 5.1 Frontend — Next.js

**Version:** Next.js 14+ (App Router)

**Key packages:**
```json
{
  "dependencies": {
    "next": "^14.2.0",
    "react": "^18.3.0",
    "react-dom": "^18.3.0",
    "typescript": "^5.4.0",
    "tailwindcss": "^3.4.0",
    "@shadcn/ui": "latest",
    "axios": "^1.7.0",
    "react-query": "^5.0.0",
    "zustand": "^4.5.0",
    "react-dropzone": "^14.2.0",
    "recharts": "^2.12.0",
    "framer-motion": "^11.0.0",
    "react-hot-toast": "^2.4.0",
    "date-fns": "^3.6.0",
    "zod": "^3.23.0",
    "react-hook-form": "^7.51.0",
    "@hookform/resolvers": "^3.3.0"
  }
}
```

**App Router structure:**
```
app/
  (auth)/
    login/page.tsx
    register/page.tsx
    forgot-password/page.tsx
  (dashboard)/
    layout.tsx
    dashboard/page.tsx
    analyses/
      page.tsx
      [id]/page.tsx
    compare/page.tsx
    competitors/page.tsx
    tracking/page.tsx
    database/page.tsx
    team/page.tsx
    settings/page.tsx
    billing/page.tsx
  api/
    auth/[...nextauth]/route.ts
  layout.tsx
  page.tsx (landing)
```

### 5.2 Backend — Go with Gin

**Version:** Go 1.22+

**Project structure (Clean Architecture):**
```
cmd/
  api/main.go          — API server entrypoint
  worker/main.go       — Asynq worker entrypoint

internal/
  domain/              — Entities, value objects
    analysis/
    competitor/
    user/
    workspace/
    billing/
  
  repository/          — Repository interfaces (ports)
    analysis_repo.go
    user_repo.go
    workspace_repo.go
    competitor_repo.go

  usecase/             — Business logic
    analysis_usecase.go
    competitor_usecase.go
    user_usecase.go
    billing_usecase.go

  handler/             — HTTP handlers (Gin)
    analysis_handler.go
    user_handler.go
    workspace_handler.go
    auth_handler.go

  worker/              — Asynq task handlers
    analysis_worker.go
    competitor_worker.go
    tracking_worker.go

  infra/               — Repository implementations
    postgres/
      analysis_postgres.go
      user_postgres.go
    redis/
      cache.go
      queue.go
    s3/
      storage.go
    youtube/
      client.go
    cv/
      client.go
    stripe/
      client.go

  middleware/
    auth.go
    rate_limit.go
    tenant.go
    metrics.go

  config/
    config.go

pkg/
  jwt/
  validator/
  pagination/
  errors/
  logger/
  hash/
```

**Dependency Injection (Wire):**
```go
// wire.go
//go:build wireinject

func InitializeAPIServer(cfg *config.Config) (*gin.Engine, error) {
    wire.Build(
        postgres.NewDB,
        redis.NewClient,
        s3.NewClient,
        youtube.NewClient,
        cv.NewClient,
        stripe.NewClient,
        repository.NewAnalysisRepo,
        repository.NewUserRepo,
        repository.NewWorkspaceRepo,
        usecase.NewAnalysisUsecase,
        usecase.NewUserUsecase,
        usecase.NewBillingUsecase,
        handler.NewAnalysisHandler,
        handler.NewUserHandler,
        handler.NewAuthHandler,
        server.NewGinRouter,
    )
    return nil, nil
}
```

### 5.3 Database — PostgreSQL + SQLC

**SQLC configuration (`sqlc.yaml`):**
```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "db/queries/"
    schema: "db/migrations/"
    gen:
      go:
        package: "db"
        out: "internal/infra/postgres/db"
        emit_json_tags: true
        emit_prepared_queries: true
        emit_interface: true
        emit_exact_table_names: false
        emit_empty_slices: true
        emit_enum_valid_method: true
        emit_all_enum_values: true
```

**Query file example (`db/queries/analyses.sql`):**
```sql
-- name: CreateAnalysis :one
INSERT INTO analyses (
    id, workspace_id, user_id, keyword, thumbnail_url, status, created_at
) VALUES (
    $1, $2, $3, $4, $5, 'pending', NOW()
) RETURNING *;

-- name: GetAnalysisByID :one
SELECT * FROM analyses
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListAnalysesByWorkspace :many
SELECT * FROM analyses
WHERE workspace_id = $1 AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateAnalysisStatus :one
UPDATE analyses
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateAnalysisResults :one
UPDATE analyses
SET 
    status = 'complete',
    score = $2,
    visibility_score = $3,
    contrast_score = $4,
    attention_score = $5,
    mobile_score = $6,
    branding_score = $7,
    curiosity_score = $8,
    cv_results = $9,
    competitor_avg = $10,
    suggestions = $11,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
```

### 5.4 Migrations — Goose

**Configuration (`db/migrations/`):**

Files follow naming: `YYYYMMDDHHMMSS_description.sql`

```sql
-- 20240101000001_create_users.sql
-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    full_name VARCHAR(255),
    avatar_url TEXT,
    google_id VARCHAR(255),
    email_verified BOOLEAN DEFAULT FALSE,
    plan VARCHAR(50) DEFAULT 'free',
    stripe_customer_id VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_stripe_id ON users(stripe_customer_id);

-- +goose Down
DROP TABLE IF EXISTS users;
```

**Rollback strategy:**
```bash
# Apply all pending
goose -dir db/migrations postgres "$DATABASE_URL" up

# Roll back one
goose -dir db/migrations postgres "$DATABASE_URL" down

# Roll back to specific version
goose -dir db/migrations postgres "$DATABASE_URL" down-to 20240101000001

# Status check
goose -dir db/migrations postgres "$DATABASE_URL" status
```

**Versioning strategy:** Sequential timestamp-based versioning. Each migration is a paired Up/Down. Breaking changes always get a new migration rather than editing existing ones. Migrations run automatically at service startup in `main.go` using `goose.Up()`.

### 5.5 Cache — Redis

**Redis usage patterns:**

```go
// Cache keys (using namespaced prefixes)
const (
    CacheKeyYouTubeSearch   = "yt:search:%s:%d"    // keyword, page
    CacheKeyAnalysis        = "analysis:%s"          // analysis_id
    CacheKeyCompetitors     = "competitors:%s"        // keyword
    CacheKeyUserSession     = "session:%s"            // session_id
    CacheKeyRateLimit       = "ratelimit:%s:%s"       // user_id, endpoint
    CacheKeyAPIKey          = "apikey:%s"             // hashed_key
    CacheKeyNicheAvg        = "niche_avg:%s"          // keyword
)

// TTLs
const (
    TTLYouTubeSearch   = 1 * time.Hour
    TTLAnalysis        = 24 * time.Hour
    TTLCompetitors     = 6 * time.Hour
    TTLUserSession     = 7 * 24 * time.Hour
    TTLNicheAvg        = 12 * time.Hour
)
```

**Rate limiting with Redis:**
```go
func (r *RedisRateLimiter) Allow(ctx context.Context, userID, endpoint string) (bool, error) {
    key := fmt.Sprintf(CacheKeyRateLimit, userID, endpoint)
    pipe := r.client.Pipeline()
    incr := pipe.Incr(ctx, key)
    pipe.Expire(ctx, key, time.Minute)
    _, err := pipe.Exec(ctx)
    if err != nil {
        return false, err
    }
    count := incr.Val()
    limit := r.getLimitForPlan(userID)
    return count <= int64(limit), nil
}
```

### 5.6 Background Jobs — Asynq

**Task type definitions:**
```go
const (
    TypeAnalyzeThumbnail    = "thumbnail:analyze"
    TypeAnalyzeCompetitors  = "competitor:analyze"
    TypeTrackChannel        = "tracking:channel"
    TypeTrackKeyword        = "tracking:keyword"
    TypeGenerateReport      = "report:generate"
    TypeSendNotification    = "notification:send"
    TypeCleanupExpired      = "cleanup:expired"
)

// Asynq server config
srv := asynq.NewServer(
    asynq.RedisClientOpt{Addr: cfg.RedisAddr},
    asynq.Config{
        Concurrency: 20,
        Queues: map[string]int{
            "critical": 6,  // user-initiated analyses
            "default":  3,  // tracking jobs
            "low":      1,  // cleanup, reports
        },
        RetryDelayFunc: asynq.DefaultRetryDelayFunc,
        MaxRetry:       3,
        ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
            log.Error().Err(err).Str("task", task.Type()).Msg("task failed")
        }),
    },
)

// Register handlers
mux := asynq.NewServeMux()
mux.HandleFunc(TypeAnalyzeThumbnail, worker.HandleAnalyzeThumbnail)
mux.HandleFunc(TypeAnalyzeCompetitors, worker.HandleAnalyzeCompetitors)
mux.HandleFunc(TypeTrackChannel, worker.HandleTrackChannel)
mux.HandleFunc(TypeTrackKeyword, worker.HandleTrackKeyword)
```

### 5.7 Storage — AWS S3

**Bucket structure:**
```
thumbnailiq-uploads/
  {workspace_id}/{analysis_id}/original.{ext}
  {workspace_id}/{analysis_id}/resized_480x270.webp
  {workspace_id}/{analysis_id}/resized_120x90.webp

thumbnailiq-competitors/
  {keyword_hash}/{video_id}/{snapshot_date}.webp

thumbnailiq-reports/
  {workspace_id}/{report_id}.pdf

thumbnailiq-assets/
  suggestions/before-after/{suggestion_type}.webp
  viral-db/{thumbnail_id}.webp
```

**Pre-signed URL generation:**
```go
func (s *S3Storage) GeneratePresignedUpload(ctx context.Context, key string, ttl time.Duration) (string, error) {
    client := s3.NewPresignClient(s.client)
    req, err := client.PresignPutObject(ctx, &s3.PutObjectInput{
        Bucket:      aws.String(s.bucket),
        Key:         aws.String(key),
        ContentType: aws.String("image/jpeg"),
    }, s3.WithPresignExpires(ttl))
    if err != nil {
        return "", err
    }
    return req.URL, nil
}
```

### 5.8 Configuration — Viper

```go
// config/config.go
type Config struct {
    Server struct {
        Port    int    `mapstructure:"port"`
        Host    string `mapstructure:"host"`
        Env     string `mapstructure:"env"`
    } `mapstructure:"server"`
    
    Database struct {
        URL         string `mapstructure:"url"`
        MaxConns    int    `mapstructure:"max_conns"`
        MinConns    int    `mapstructure:"min_conns"`
    } `mapstructure:"database"`
    
    Redis struct {
        Addr     string `mapstructure:"addr"`
        Password string `mapstructure:"password"`
        DB       int    `mapstructure:"db"`
    } `mapstructure:"redis"`
    
    JWT struct {
        AccessSecret  string        `mapstructure:"access_secret"`
        RefreshSecret string        `mapstructure:"refresh_secret"`
        AccessTTL     time.Duration `mapstructure:"access_ttl"`
        RefreshTTL    time.Duration `mapstructure:"refresh_ttl"`
    } `mapstructure:"jwt"`
    
    AWS struct {
        Region          string `mapstructure:"region"`
        UploadBucket    string `mapstructure:"upload_bucket"`
        CompetitorBucket string `mapstructure:"competitor_bucket"`
        ReportBucket    string `mapstructure:"report_bucket"`
    } `mapstructure:"aws"`
    
    YouTube struct {
        APIKey          string `mapstructure:"api_key"`
        DailyQuota      int    `mapstructure:"daily_quota"`
    } `mapstructure:"youtube"`
    
    Stripe struct {
        SecretKey       string `mapstructure:"secret_key"`
        WebhookSecret   string `mapstructure:"webhook_secret"`
        PriceIDStarter  string `mapstructure:"price_id_starter"`
        PriceIDPro      string `mapstructure:"price_id_pro"`
        PriceIDAgency   string `mapstructure:"price_id_agency"`
    } `mapstructure:"stripe"`
    
    CVService struct {
        URL     string `mapstructure:"url"`
        APIKey  string `mapstructure:"api_key"`
    } `mapstructure:"cv_service"`
    
    Claude struct {
        APIKey string `mapstructure:"api_key"`
        Model  string `mapstructure:"model"`
    } `mapstructure:"claude"`
}

func Load() (*Config, error) {
    v := viper.New()
    v.SetConfigName("config")
    v.SetConfigType("yaml")
    v.AddConfigPath(".")
    v.AddConfigPath("./config")
    v.AutomaticEnv()
    v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    
    if err := v.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("reading config: %w", err)
    }
    
    var cfg Config
    if err := v.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("unmarshaling config: %w", err)
    }
    return &cfg, nil
}
```

---

## 6. Database Design

### 6.1 ER Diagram

```mermaid
erDiagram
    USERS ||--o{ WORKSPACE_MEMBERS : "belongs to"
    WORKSPACES ||--o{ WORKSPACE_MEMBERS : "has"
    WORKSPACES ||--o{ PROJECTS : "contains"
    WORKSPACES ||--o{ API_KEYS : "issues"
    PROJECTS ||--o{ ANALYSES : "contains"
    ANALYSES ||--o{ THUMBNAIL_VERSIONS : "has"
    ANALYSES ||--o{ COMPETITOR_SNAPSHOTS : "generates"
    ANALYSES ||--o{ SUGGESTIONS : "produces"
    ANALYSES ||--o{ COMMENTS : "receives"
    USERS ||--o{ COMMENTS : "writes"
    WORKSPACES ||--o{ TRACKING_JOBS : "runs"
    TRACKING_JOBS ||--o{ COMPETITOR_SNAPSHOTS : "captures"
    WORKSPACES ||--o{ SUBSCRIPTIONS : "has"
    USERS ||--o{ REFRESH_TOKENS : "holds"

    USERS {
        uuid id PK
        varchar email
        varchar password_hash
        varchar full_name
        text avatar_url
        varchar google_id
        bool email_verified
        varchar stripe_customer_id
        timestamptz created_at
        timestamptz updated_at
        timestamptz deleted_at
    }

    WORKSPACES {
        uuid id PK
        varchar name
        varchar slug
        text logo_url
        varchar plan
        uuid owner_id FK
        int analyses_this_month
        int analyses_limit
        timestamptz created_at
        timestamptz updated_at
        timestamptz deleted_at
    }

    WORKSPACE_MEMBERS {
        uuid id PK
        uuid workspace_id FK
        uuid user_id FK
        varchar role
        timestamptz joined_at
    }

    PROJECTS {
        uuid id PK
        uuid workspace_id FK
        varchar name
        text description
        varchar color
        timestamptz created_at
        timestamptz updated_at
        timestamptz deleted_at
    }

    ANALYSES {
        uuid id PK
        uuid workspace_id FK
        uuid project_id FK
        uuid user_id FK
        varchar keyword
        text thumbnail_url
        text thumbnail_s3_key
        varchar status
        int score
        int visibility_score
        int contrast_score
        int attention_score
        int mobile_score
        int branding_score
        int curiosity_score
        jsonb cv_results
        jsonb competitor_avg
        jsonb suggestions
        int competitor_count
        int rank_in_competitors
        timestamptz created_at
        timestamptz updated_at
        timestamptz deleted_at
    }

    THUMBNAIL_VERSIONS {
        uuid id PK
        uuid analysis_id FK
        int version_number
        text s3_key
        text thumbnail_url
        int score
        jsonb cv_results
        bool is_selected_winner
        timestamptz created_at
    }

    COMPETITOR_SNAPSHOTS {
        uuid id PK
        uuid analysis_id FK
        varchar video_id
        text video_title
        varchar channel_id
        varchar channel_name
        text thumbnail_url
        text thumbnail_s3_key
        int view_count
        int subscriber_count
        int rank_position
        jsonb cv_results
        int score
        date snapshot_date
        timestamptz created_at
    }

    TRACKING_JOBS {
        uuid id PK
        uuid workspace_id FK
        varchar type
        varchar channel_id
        varchar keyword
        varchar status
        int check_interval_hours
        timestamptz last_checked_at
        timestamptz next_check_at
        timestamptz created_at
        timestamptz updated_at
    }

    SUGGESTIONS {
        uuid id PK
        uuid analysis_id FK
        varchar type
        varchar impact_level
        text headline
        text explanation
        text before_image_url
        text after_image_url
        float estimated_ctr_min
        float estimated_ctr_max
        int display_order
        timestamptz created_at
    }

    COMMENTS {
        uuid id PK
        uuid analysis_id FK
        uuid user_id FK
        uuid parent_id FK
        text body
        varchar status
        timestamptz created_at
        timestamptz updated_at
        timestamptz deleted_at
    }

    API_KEYS {
        uuid id PK
        uuid workspace_id FK
        varchar name
        varchar key_hash
        varchar key_prefix
        int requests_this_month
        int requests_limit
        timestamptz last_used_at
        timestamptz expires_at
        timestamptz created_at
    }

    SUBSCRIPTIONS {
        uuid id PK
        uuid workspace_id FK
        varchar stripe_subscription_id
        varchar stripe_price_id
        varchar plan
        varchar status
        timestamptz current_period_start
        timestamptz current_period_end
        timestamptz cancel_at
        timestamptz created_at
        timestamptz updated_at
    }

    REFRESH_TOKENS {
        uuid id PK
        uuid user_id FK
        varchar token_hash
        varchar device_info
        bool is_revoked
        timestamptz expires_at
        timestamptz created_at
    }
```

### 6.2 Full SQL Schema

```sql
-- db/migrations/20240101000001_create_users.sql
-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    full_name VARCHAR(255) NOT NULL DEFAULT '',
    avatar_url TEXT,
    google_id VARCHAR(255) UNIQUE,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    stripe_customer_id VARCHAR(255) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_google_id ON users(google_id) WHERE google_id IS NOT NULL;

-- +goose Down
DROP TABLE IF EXISTS users;

-- db/migrations/20240101000002_create_workspaces.sql
-- +goose Up
CREATE TABLE workspaces (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    logo_url TEXT,
    plan VARCHAR(50) NOT NULL DEFAULT 'free',
    owner_id UUID NOT NULL REFERENCES users(id),
    analyses_this_month INT NOT NULL DEFAULT 0,
    analyses_limit INT NOT NULL DEFAULT 5,
    api_requests_this_month INT NOT NULL DEFAULT 0,
    api_requests_limit INT NOT NULL DEFAULT 100,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_workspaces_owner ON workspaces(owner_id);
CREATE INDEX idx_workspaces_slug ON workspaces(slug) WHERE deleted_at IS NULL;

CREATE TABLE workspace_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'editor',
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(workspace_id, user_id)
);

CREATE INDEX idx_workspace_members_workspace ON workspace_members(workspace_id);
CREATE INDEX idx_workspace_members_user ON workspace_members(user_id);

-- +goose Down
DROP TABLE IF EXISTS workspace_members;
DROP TABLE IF EXISTS workspaces;

-- db/migrations/20240101000003_create_projects.sql
-- +goose Up
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    color VARCHAR(7) DEFAULT '#6366F1',
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_projects_workspace ON projects(workspace_id) WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE IF EXISTS projects;

-- db/migrations/20240101000004_create_analyses.sql
-- +goose Up
CREATE TABLE analyses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    project_id UUID REFERENCES projects(id) ON DELETE SET NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    keyword VARCHAR(500) NOT NULL,
    keyword_normalized VARCHAR(500) NOT NULL,
    thumbnail_url TEXT NOT NULL,
    thumbnail_s3_key TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    -- Scores
    score INT,
    visibility_score INT,
    contrast_score INT,
    attention_score INT,
    mobile_score INT,
    branding_score INT,
    curiosity_score INT,
    -- Analysis Results
    cv_results JSONB,
    competitor_avg JSONB,
    suggestions JSONB,
    competitor_count INT DEFAULT 0,
    rank_in_competitors INT,
    -- Error handling
    error_message TEXT,
    retry_count INT DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_analyses_workspace ON analyses(workspace_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_analyses_project ON analyses(project_id) WHERE project_id IS NOT NULL;
CREATE INDEX idx_analyses_keyword ON analyses USING gin(to_tsvector('english', keyword));
CREATE INDEX idx_analyses_status ON analyses(status) WHERE status != 'complete';
CREATE INDEX idx_analyses_created ON analyses(created_at DESC);

CREATE TABLE thumbnail_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    analysis_id UUID NOT NULL REFERENCES analyses(id) ON DELETE CASCADE,
    version_number INT NOT NULL DEFAULT 1,
    s3_key TEXT NOT NULL,
    thumbnail_url TEXT NOT NULL,
    score INT,
    cv_results JSONB,
    is_selected_winner BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(analysis_id, version_number)
);

-- +goose Down
DROP TABLE IF EXISTS thumbnail_versions;
DROP TABLE IF EXISTS analyses;

-- db/migrations/20240101000005_create_competitors.sql
-- +goose Up
CREATE TABLE competitor_snapshots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    analysis_id UUID REFERENCES analyses(id) ON DELETE SET NULL,
    tracking_job_id UUID,
    video_id VARCHAR(20) NOT NULL,
    video_title TEXT NOT NULL,
    channel_id VARCHAR(50) NOT NULL,
    channel_name VARCHAR(255) NOT NULL,
    thumbnail_url TEXT NOT NULL,
    thumbnail_s3_key TEXT,
    view_count BIGINT,
    like_count BIGINT,
    subscriber_count BIGINT,
    rank_position INT,
    keyword VARCHAR(500),
    cv_results JSONB,
    score INT,
    snapshot_date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_competitor_snapshots_analysis ON competitor_snapshots(analysis_id);
CREATE INDEX idx_competitor_snapshots_video ON competitor_snapshots(video_id);
CREATE INDEX idx_competitor_snapshots_channel ON competitor_snapshots(channel_id);
CREATE INDEX idx_competitor_snapshots_keyword ON competitor_snapshots(keyword, snapshot_date);
CREATE INDEX idx_competitor_snapshots_date ON competitor_snapshots(snapshot_date DESC);

-- Partition by month for large datasets
-- CREATE TABLE competitor_snapshots_2024_01 PARTITION OF competitor_snapshots
-- FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

-- +goose Down
DROP TABLE IF EXISTS competitor_snapshots;

-- db/migrations/20240101000006_create_tracking.sql
-- +goose Up
CREATE TABLE tracking_jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, -- 'channel' or 'keyword'
    channel_id VARCHAR(50),
    keyword VARCHAR(500),
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    check_interval_hours INT NOT NULL DEFAULT 24,
    last_checked_at TIMESTAMPTZ,
    next_check_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tracking_jobs_workspace ON tracking_jobs(workspace_id);
CREATE INDEX idx_tracking_jobs_next_check ON tracking_jobs(next_check_at) WHERE status = 'active';

-- +goose Down
DROP TABLE IF EXISTS tracking_jobs;

-- db/migrations/20240101000007_create_comments_api_keys.sql
-- +goose Up
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    analysis_id UUID NOT NULL REFERENCES analyses(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    parent_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'open',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_comments_analysis ON comments(analysis_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_comments_parent ON comments(parent_id) WHERE parent_id IS NOT NULL;

CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    key_hash VARCHAR(64) NOT NULL UNIQUE,
    key_prefix VARCHAR(12) NOT NULL,
    requests_this_month INT NOT NULL DEFAULT 0,
    requests_limit INT NOT NULL DEFAULT 100,
    last_used_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ
);

CREATE INDEX idx_api_keys_workspace ON api_keys(workspace_id);
CREATE INDEX idx_api_keys_hash ON api_keys(key_hash);

CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    stripe_subscription_id VARCHAR(255) UNIQUE NOT NULL,
    stripe_price_id VARCHAR(255) NOT NULL,
    plan VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    current_period_start TIMESTAMPTZ,
    current_period_end TIMESTAMPTZ,
    cancel_at TIMESTAMPTZ,
    canceled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_subscriptions_workspace ON subscriptions(workspace_id);
CREATE INDEX idx_subscriptions_stripe ON subscriptions(stripe_subscription_id);

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    device_info TEXT,
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_hash ON refresh_tokens(token_hash) WHERE NOT is_revoked;

-- +goose Down
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS comments;

-- db/migrations/20240101000008_create_viral_db.sql
-- +goose Up
CREATE TABLE viral_thumbnails (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    video_id VARCHAR(20) UNIQUE NOT NULL,
    channel_id VARCHAR(50) NOT NULL,
    channel_name VARCHAR(255) NOT NULL,
    video_title TEXT NOT NULL,
    thumbnail_url TEXT NOT NULL,
    thumbnail_s3_key TEXT,
    niche VARCHAR(100),
    tags TEXT[],
    view_count BIGINT,
    view_count_when_captured BIGINT,
    score INT,
    cv_results JSONB,
    captured_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_viral_thumbnails_niche ON viral_thumbnails(niche);
CREATE INDEX idx_viral_thumbnails_score ON viral_thumbnails(score DESC);
CREATE INDEX idx_viral_thumbnails_tags ON viral_thumbnails USING gin(tags);
CREATE INDEX idx_viral_thumbnails_cv ON viral_thumbnails USING gin(cv_results);

-- +goose Down
DROP TABLE IF EXISTS viral_thumbnails;
```

### 6.3 Index Strategy

```sql
-- Full-text search on keywords
CREATE INDEX idx_analyses_keyword_fts ON analyses USING gin(to_tsvector('english', keyword));

-- Partial indexes for active/pending items (dramatically reduces index size)
CREATE INDEX idx_analyses_pending ON analyses(created_at) WHERE status = 'pending';
CREATE INDEX idx_tracking_active ON tracking_jobs(next_check_at) WHERE status = 'active';

-- JSONB indexes for CV results queries
CREATE INDEX idx_analyses_cv_face ON analyses USING gin((cv_results->'face_count'));
CREATE INDEX idx_competitor_cv_dominant_color ON competitor_snapshots USING gin((cv_results->'dominant_colors'));

-- Composite indexes for common query patterns
CREATE INDEX idx_analyses_workspace_created ON analyses(workspace_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_competitor_snapshots_keyword_date ON competitor_snapshots(keyword, snapshot_date DESC);
```

### 6.4 Partitioning Strategy

For `competitor_snapshots` (will be the largest table):

```sql
-- Convert to range partitioning by snapshot_date
CREATE TABLE competitor_snapshots (
    -- ... columns ...
) PARTITION BY RANGE (snapshot_date);

-- Create monthly partitions
CREATE TABLE competitor_snapshots_y2024m01
    PARTITION OF competitor_snapshots
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

-- Automate partition creation with pg_partman
SELECT partman.create_parent(
    p_parent_table => 'public.competitor_snapshots',
    p_control => 'snapshot_date',
    p_type => 'range',
    p_interval => 'monthly',
    p_premake => 3
);
```

---

## 7. API Design & OpenAPI Specification

### 7.1 API Versioning Strategy

All endpoints are prefixed with `/api/v1/`. Breaking changes trigger a new version prefix `/api/v2/`. Non-breaking additions (new fields) are backward-compatible within a version.

Deprecated endpoints return `Deprecation: true` and `Sunset: {date}` headers.

### 7.2 Core Endpoints

```yaml
openapi: 3.1.0
info:
  title: ThumbnailIQ API
  version: 1.0.0
  description: YouTube Thumbnail Intelligence Platform API

servers:
  - url: https://api.thumbnailiq.com/api/v1

security:
  - BearerAuth: []
  - ApiKeyAuth: []

paths:
  /auth/register:
    post:
      summary: Register new user
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [email, password, full_name]
              properties:
                email: { type: string, format: email }
                password: { type: string, minLength: 8 }
                full_name: { type: string }
      responses:
        201:
          description: User created
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AuthResponse'

  /auth/login:
    post:
      summary: Login
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required: [email, password]
              properties:
                email: { type: string }
                password: { type: string }
      responses:
        200:
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AuthResponse'

  /auth/refresh:
    post:
      summary: Refresh access token
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required: [refresh_token]
              properties:
                refresh_token: { type: string }

  /analyses:
    post:
      summary: Create new thumbnail analysis
      requestBody:
        content:
          multipart/form-data:
            schema:
              type: object
              required: [thumbnail, keyword]
              properties:
                thumbnail:
                  type: string
                  format: binary
                keyword:
                  type: string
                project_id:
                  type: string
                  format: uuid
      responses:
        202:
          description: Analysis queued
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AnalysisResponse'
    
    get:
      summary: List analyses
      parameters:
        - name: workspace_id
          in: query
          schema: { type: string, format: uuid }
        - name: project_id
          in: query
          schema: { type: string, format: uuid }
        - name: page
          in: query
          schema: { type: integer, default: 1 }
        - name: per_page
          in: query
          schema: { type: integer, default: 20, maximum: 100 }
        - name: status
          in: query
          schema: { type: string, enum: [pending, processing, complete, failed] }

  /analyses/{id}:
    get:
      summary: Get analysis by ID
      parameters:
        - name: id
          in: path
          required: true
          schema: { type: string, format: uuid }
      responses:
        200:
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AnalysisFullResponse'

  /analyses/{id}/compare:
    post:
      summary: Add version for comparison
      requestBody:
        content:
          multipart/form-data:
            schema:
              type: object
              required: [thumbnail]
              properties:
                thumbnail:
                  type: string
                  format: binary

  /keywords/{keyword}/competitors:
    get:
      summary: Get competitor thumbnails for keyword
      parameters:
        - name: keyword
          in: path
          required: true
          schema: { type: string }
        - name: count
          in: query
          schema: { type: integer, default: 20, maximum: 50 }
      responses:
        200:
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/CompetitorListResponse'

  /tracking:
    post:
      summary: Create tracking job
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required: [type]
              properties:
                type: { type: string, enum: [channel, keyword] }
                channel_id: { type: string }
                keyword: { type: string }
                interval_hours: { type: integer, default: 24 }
    get:
      summary: List tracking jobs

  /workspaces:
    post:
      summary: Create workspace
    get:
      summary: List user workspaces

  /workspaces/{id}/members:
    get:
      summary: List workspace members
    post:
      summary: Invite member
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required: [email, role]
              properties:
                email: { type: string, format: email }
                role: { type: string, enum: [admin, editor, viewer] }

  /api-keys:
    post:
      summary: Create API key
    get:
      summary: List API keys
  
  /api-keys/{id}:
    delete:
      summary: Revoke API key

  /billing/plans:
    get:
      summary: Get available plans

  /billing/subscribe:
    post:
      summary: Create or update subscription
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required: [plan, payment_method_id]
              properties:
                plan: { type: string, enum: [starter, pro, agency] }
                payment_method_id: { type: string }

  /viral-db:
    get:
      summary: Search viral thumbnail database
      parameters:
        - name: niche
          in: query
          schema: { type: string }
        - name: tags
          in: query
          schema: { type: array, items: { type: string } }
        - name: min_score
          in: query
          schema: { type: integer }
        - name: has_face
          in: query
          schema: { type: boolean }

components:
  schemas:
    AuthResponse:
      type: object
      properties:
        access_token: { type: string }
        refresh_token: { type: string }
        expires_in: { type: integer }
        user:
          $ref: '#/components/schemas/User'

    User:
      type: object
      properties:
        id: { type: string, format: uuid }
        email: { type: string }
        full_name: { type: string }
        avatar_url: { type: string }

    AnalysisResponse:
      type: object
      properties:
        id: { type: string, format: uuid }
        status: { type: string }
        keyword: { type: string }
        thumbnail_url: { type: string }
        created_at: { type: string, format: date-time }

    AnalysisFullResponse:
      allOf:
        - $ref: '#/components/schemas/AnalysisResponse'
        - type: object
          properties:
            score: { type: integer }
            visibility_score: { type: integer }
            contrast_score: { type: integer }
            attention_score: { type: integer }
            mobile_score: { type: integer }
            branding_score: { type: integer }
            curiosity_score: { type: integer }
            cv_results:
              $ref: '#/components/schemas/CVResults'
            competitor_avg:
              $ref: '#/components/schemas/CompetitorAvg'
            suggestions:
              type: array
              items:
                $ref: '#/components/schemas/Suggestion'
            competitor_count: { type: integer }
            rank_in_competitors: { type: integer }

    CVResults:
      type: object
      properties:
        face_count: { type: integer }
        emotions: { type: array, items: { type: string } }
        text_detected: { type: boolean }
        text_strings: { type: array, items: { type: string } }
        text_density_pct: { type: number }
        word_count: { type: integer }
        dominant_colors:
          type: array
          items:
            type: object
            properties:
              hex: { type: string }
              percentage: { type: number }
        contrast_score: { type: number }
        brightness_score: { type: number }
        saturation_score: { type: number }
        object_count: { type: integer }
        clutter_score: { type: number }
        visual_complexity: { type: number }

    CompetitorAvg:
      type: object
      additionalProperties: true

    Suggestion:
      type: object
      properties:
        type: { type: string }
        impact_level: { type: string, enum: [high, medium, low] }
        headline: { type: string }
        explanation: { type: string }
        before_image_url: { type: string }
        after_image_url: { type: string }
        estimated_ctr_min: { type: number }
        estimated_ctr_max: { type: number }

  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
    ApiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
```

---

## 8. AI & Computer Vision Pipeline

### 8.1 Architecture Overview

The CV pipeline runs as a Python microservice (FastAPI) separate from the Go API, allowing independent scaling of compute-intensive tasks.

```
CV Service (Python 3.12 + FastAPI)
├── /analyze           — Full analysis endpoint
├── /batch-analyze     — Analyze multiple images
└── /health

Models:
├── OCR:        Tesseract 5.x / EasyOCR
├── Faces:      InsightFace (Buffalo_L model)
├── Emotions:   FER+ / DeepFace
├── Objects:    YOLOv8n (lightweight)
└── Colors:     Custom Python (PIL/scikit-learn)
```

### 8.2 OCR (Text Detection)

**Model:** EasyOCR (better than Tesseract for thumbnail-style text with varied fonts/backgrounds)

```python
import easyocr
import numpy as np
from PIL import Image

reader = easyocr.Reader(['en'], gpu=False)  # GPU when available

def extract_text(image_path: str) -> dict:
    results = reader.readtext(image_path, detail=1)
    
    texts = []
    total_area = 0
    img = Image.open(image_path)
    img_area = img.width * img.height
    
    for (bbox, text, confidence) in results:
        if confidence < 0.5:
            continue
        # Calculate bounding box area
        x_coords = [p[0] for p in bbox]
        y_coords = [p[1] for p in bbox]
        width = max(x_coords) - min(x_coords)
        height = max(y_coords) - min(y_coords)
        area = width * height
        total_area += area
        
        # Calculate text height relative to image
        relative_height = height / img.height
        
        texts.append({
            "text": text,
            "confidence": confidence,
            "bbox": bbox,
            "relative_height": relative_height,
        })
    
    text_density_pct = (total_area / img_area) * 100
    word_count = sum(len(t["text"].split()) for t in texts)
    avg_text_height = np.mean([t["relative_height"] for t in texts]) if texts else 0
    
    return {
        "text_detected": len(texts) > 0,
        "text_strings": [t["text"] for t in texts],
        "text_density_pct": round(text_density_pct, 2),
        "word_count": word_count,
        "avg_text_height_pct": round(avg_text_height * 100, 2),
        "texts_detail": texts,
    }
```

### 8.3 Face Detection & Emotion Analysis

**Face Detection:** InsightFace (accurate, fast)
**Emotion Analysis:** DeepFace with FER2013 backend

```python
import insightface
from deepface import DeepFace
import cv2

# Initialize InsightFace
app = insightface.app.FaceAnalysis(name='buffalo_l')
app.prepare(ctx_id=0, det_size=(640, 640))  # ctx_id=-1 for CPU

def analyze_faces(image_path: str) -> dict:
    img = cv2.imread(image_path)
    faces = app.get(img)
    
    face_results = []
    for face in faces:
        bbox = face.bbox.tolist()
        
        # Get emotion from DeepFace on face crop
        x1, y1, x2, y2 = [int(c) for c in bbox]
        face_crop = img[y1:y2, x1:x2]
        
        try:
            emotion_result = DeepFace.analyze(
                face_crop,
                actions=['emotion'],
                enforce_detection=False,
                silent=True,
            )
            dominant_emotion = emotion_result[0]['dominant_emotion']
            emotions = emotion_result[0]['emotion']
        except Exception:
            dominant_emotion = "unknown"
            emotions = {}
        
        # Check if face is looking at camera (gaze estimation)
        gaze = estimate_gaze(face)
        
        face_results.append({
            "bbox": bbox,
            "dominant_emotion": dominant_emotion,
            "emotions": emotions,
            "eye_contact": gaze["looking_at_camera"],
            "confidence": float(face.det_score),
        })
    
    return {
        "face_count": len(faces),
        "faces": face_results,
        "has_eye_contact": any(f["eye_contact"] for f in face_results),
        "primary_emotion": face_results[0]["dominant_emotion"] if face_results else None,
    }

def estimate_gaze(face) -> dict:
    # Simple gaze estimation using face pose
    # face.pose returns [pitch, yaw, roll]
    if hasattr(face, 'pose'):
        yaw = abs(face.pose[1])
        pitch = abs(face.pose[0])
        looking_at_camera = yaw < 25 and pitch < 25
        return {"looking_at_camera": looking_at_camera}
    return {"looking_at_camera": False}
```

### 8.4 Color Analysis

```python
from sklearn.cluster import KMeans
import numpy as np
from PIL import Image
import colorsys

def analyze_colors(image_path: str, n_colors: int = 5) -> dict:
    img = Image.open(image_path).convert('RGB')
    img = img.resize((150, 84))  # Resize for speed
    data = np.array(img).reshape(-1, 3)
    
    # K-means clustering for dominant colors
    kmeans = KMeans(n_clusters=n_colors, random_state=42, n_init=10)
    kmeans.fit(data)
    
    colors = []
    total_pixels = len(data)
    
    for i, center in enumerate(kmeans.cluster_centers_):
        r, g, b = [int(c) for c in center]
        hex_color = f"#{r:02x}{g:02x}{b:02x}"
        count = np.sum(kmeans.labels_ == i)
        percentage = (count / total_pixels) * 100
        
        # Calculate luminance
        lum = 0.299 * r + 0.587 * g + 0.114 * b
        
        # Calculate saturation
        h, s, v = colorsys.rgb_to_hsv(r/255, g/255, b/255)
        
        colors.append({
            "hex": hex_color,
            "rgb": [r, g, b],
            "percentage": round(percentage, 2),
            "luminance": round(lum, 2),
            "saturation": round(s * 100, 2),
        })
    
    colors.sort(key=lambda x: x["percentage"], reverse=True)
    
    # Calculate contrast between top 2 colors
    if len(colors) >= 2:
        contrast = calculate_wcag_contrast(colors[0]["rgb"], colors[1]["rgb"])
    else:
        contrast = 0
    
    # Overall brightness (average luminance)
    avg_brightness = np.mean([c["luminance"] for c in colors])
    avg_saturation = np.mean([c["saturation"] for c in colors])
    
    return {
        "dominant_colors": colors,
        "contrast_score": round(contrast, 2),
        "brightness_score": round(avg_brightness, 2),
        "saturation_score": round(avg_saturation, 2),
    }

def calculate_wcag_contrast(rgb1: list, rgb2: list) -> float:
    def relative_luminance(rgb):
        r, g, b = [x / 255 for x in rgb]
        r = r / 12.92 if r <= 0.03928 else ((r + 0.055) / 1.055) ** 2.4
        g = g / 12.92 if g <= 0.03928 else ((g + 0.055) / 1.055) ** 2.4
        b = b / 12.92 if b <= 0.03928 else ((b + 0.055) / 1.055) ** 2.4
        return 0.2126 * r + 0.7152 * g + 0.0722 * b
    
    l1 = relative_luminance(rgb1) + 0.05
    l2 = relative_luminance(rgb2) + 0.05
    return max(l1, l2) / min(l1, l2)
```

### 8.5 Object Detection & Clutter Scoring

**Model:** YOLOv8n (nano variant, fast, ~6MB)

```python
from ultralytics import YOLO
import numpy as np

model = YOLO('yolov8n.pt')

def detect_objects(image_path: str) -> dict:
    results = model(image_path, verbose=False)[0]
    
    objects = []
    for box in results.boxes:
        objects.append({
            "label": results.names[int(box.cls)],
            "confidence": float(box.conf),
            "bbox": box.xyxy[0].tolist(),
        })
    
    object_count = len(objects)
    
    # Calculate edge density (Sobel)
    img = cv2.imread(image_path, cv2.IMREAD_GRAYSCALE)
    img_resized = cv2.resize(img, (480, 270))
    sobelx = cv2.Sobel(img_resized, cv2.CV_64F, 1, 0, ksize=3)
    sobely = cv2.Sobel(img_resized, cv2.CV_64F, 0, 1, ksize=3)
    magnitude = np.sqrt(sobelx**2 + sobely**2)
    edge_density = float(np.mean(magnitude)) / 255.0
    
    # Clutter score: 0–100 (higher = more cluttered)
    # Normalized: >10 objects = very cluttered
    object_clutter = min(object_count / 10.0, 1.0) * 50
    edge_clutter = edge_density * 50
    clutter_score = object_clutter + edge_clutter
    
    return {
        "object_count": object_count,
        "objects": objects,
        "edge_density": round(edge_density, 4),
        "clutter_score": round(clutter_score, 2),
    }
```

### 8.6 Curiosity Score (Claude API)

```go
// internal/infra/claude/client.go
type CuriosityRequest struct {
    ThumbnailURL  string   `json:"thumbnail_url"`
    TextContent   []string `json:"text_content"`
    ObjectLabels  []string `json:"object_labels"`
    PrimaryEmotion string  `json:"primary_emotion"`
}

type CuriosityResponse struct {
    Score     int    `json:"curiosity_score"`
    Reasoning string `json:"reasoning"`
    Factors   struct {
        InformationGap   bool `json:"information_gap"`
        PromisesReward   bool `json:"promises_reward"`
        UnusualSituation bool `json:"unusual_situation"`
        HasNumbers       bool `json:"has_numbers"`
    } `json:"factors"`
}

func (c *ClaudeClient) ScoreCuriosity(ctx context.Context, req CuriosityRequest) (*CuriosityResponse, error) {
    prompt := fmt.Sprintf(`You are a YouTube CTR optimization expert.

Analyze this thumbnail and score its curiosity factor from 0-100.

Thumbnail text: %v
Detected objects: %v  
Primary emotion shown: %s

Score based on:
1. Information gap (does it create a "need to know"?)
2. Promised reward/transformation
3. Unusual or unexpected element
4. Specific numbers or claims
5. Emotional resonance

Return ONLY valid JSON with this structure:
{
  "curiosity_score": <0-100>,
  "reasoning": "<2 sentence explanation>",
  "factors": {
    "information_gap": <true/false>,
    "promises_reward": <true/false>,
    "unusual_situation": <true/false>,
    "has_numbers": <true/false>
  }
}`,
        req.TextContent, req.ObjectLabels, req.PrimaryEmotion)

    response, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
        Model:     anthropic.F(anthropic.ModelClaude3HaikuLatest),
        MaxTokens: anthropic.F(int64(500)),
        Messages: anthropic.F([]anthropic.MessageParam{
            anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
        }),
    })
    if err != nil {
        return nil, fmt.Errorf("claude api: %w", err)
    }
    
    var result CuriosityResponse
    text := response.Content[0].Text
    if err := json.Unmarshal([]byte(text), &result); err != nil {
        return nil, fmt.Errorf("parsing claude response: %w", err)
    }
    return &result, nil
}
```

### 8.7 AI Model Alternatives

| Function | Open Source (Self-hosted) | Hosted API Alternative |
|----------|--------------------------|----------------------|
| OCR | EasyOCR, Tesseract 5 | AWS Textract ($1.50/1K) |
| Face Detection | InsightFace, MediaPipe | AWS Rekognition ($1/1K) |
| Emotion Analysis | DeepFace, FER+ | Clarifai Emotion Model |
| Object Detection | YOLOv8n | Google Vision API ($1.5/1K) |
| Curiosity Score | Self-hosted Llama 3 | Claude API (Haiku) |

**Recommendation for MVP:** Use self-hosted models for OCR, faces, objects, and colors. Use Claude API for curiosity score (high value, low volume).

---

## 9. YouTube Data Acquisition

### 9.1 YouTube Data API v3 Integration

```go
// internal/infra/youtube/client.go
type YouTubeClient struct {
    apiKey     string
    httpClient *http.Client
    cache      *redis.Client
    logger     zerolog.Logger
}

type SearchResult struct {
    VideoID      string
    Title        string
    ChannelID    string
    ChannelName  string
    ThumbnailURL string
    PublishedAt  time.Time
    ViewCount    int64
}

func (c *YouTubeClient) SearchVideos(ctx context.Context, keyword string, maxResults int) ([]SearchResult, error) {
    // Check cache first
    cacheKey := fmt.Sprintf("yt:search:%s:%d", normalizeKeyword(keyword), maxResults)
    cached, err := c.cache.Get(ctx, cacheKey).Result()
    if err == nil {
        var results []SearchResult
        json.Unmarshal([]byte(cached), &results)
        return results, nil
    }
    
    // API call
    url := fmt.Sprintf(
        "https://www.googleapis.com/youtube/v3/search?part=snippet&q=%s&type=video&maxResults=%d&key=%s&videoCategoryId=0",
        url.QueryEscape(keyword), maxResults, c.apiKey,
    )
    
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, fmt.Errorf("youtube search: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode == 429 {
        return nil, ErrQuotaExceeded
    }
    
    var searchResp YouTubeSearchResponse
    json.NewDecoder(resp.Body).Decode(&searchResp)
    
    // Fetch video statistics (view counts) in a second call
    videoIDs := make([]string, len(searchResp.Items))
    for i, item := range searchResp.Items {
        videoIDs[i] = item.ID.VideoID
    }
    
    stats, err := c.GetVideoStatistics(ctx, videoIDs)
    if err != nil {
        c.logger.Warn().Err(err).Msg("failed to fetch video stats, continuing without")
    }
    
    results := make([]SearchResult, 0, len(searchResp.Items))
    for _, item := range searchResp.Items {
        result := SearchResult{
            VideoID:      item.ID.VideoID,
            Title:        item.Snippet.Title,
            ChannelID:    item.Snippet.ChannelID,
            ChannelName:  item.Snippet.ChannelTitle,
            ThumbnailURL: item.Snippet.Thumbnails.High.URL,
            PublishedAt:  item.Snippet.PublishedAt,
        }
        if s, ok := stats[item.ID.VideoID]; ok {
            result.ViewCount = s.ViewCount
        }
        results = append(results, result)
    }
    
    // Cache results
    data, _ := json.Marshal(results)
    c.cache.Set(ctx, cacheKey, string(data), 1*time.Hour)
    
    return results, nil
}
```

### 9.2 Quota Management

YouTube Data API v3 costs:
- **Search.list:** 100 units per call
- **Videos.list (statistics):** 1 unit per call
- **Default daily quota:** 10,000 units/day (free tier)
- **Cost per additional quota:** $0.04/unit above 10K

**Quota budget per analysis:**
- 1 search.list = 100 units
- 1 videos.list (50 results) = 1 unit
- Total: ~101 units per analysis

**At 10,000 daily quota:** ~99 analyses/day using search
**For scale:** Apply for quota increase (common for SaaS, approved within 2-3 days)

### 9.3 Caching Strategy

```
Search results:     Cache 1 hour  (keyword unlikely to change rapidly)
Video statistics:   Cache 6 hours
Channel data:       Cache 24 hours
Thumbnail URLs:     Cache 12 hours (thumbnails are semi-permanent)
Competitor snapshots: Store in S3 permanently for historical tracking
```

### 9.4 Thumbnail Download Pipeline

```go
func (w *AnalysisWorker) downloadAndStoreThumbnail(ctx context.Context, videoID, thumbnailURL, s3Key string) error {
    // Download from YouTube CDN
    resp, err := http.Get(thumbnailURL)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    // Read to buffer
    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    
    // Upload to S3
    _, err = w.s3.PutObject(ctx, &s3.PutObjectInput{
        Bucket:      aws.String(w.cfg.AWS.CompetitorBucket),
        Key:         aws.String(s3Key),
        Body:        bytes.NewReader(data),
        ContentType: aws.String("image/jpeg"),
    })
    return err
}
```

---

## 10. Background Job System (Asynq)

### 10.1 Job Definitions

```go
// internal/worker/analyze_thumbnail.go
type AnalyzeThumbnailPayload struct {
    AnalysisID   string `json:"analysis_id"`
    WorkspaceID  string `json:"workspace_id"`
    UserID       string `json:"user_id"`
    ThumbnailURL string `json:"thumbnail_url"`
    S3Key        string `json:"s3_key"`
    Keyword      string `json:"keyword"`
}

func (w *Worker) HandleAnalyzeThumbnail(ctx context.Context, task *asynq.Task) error {
    var payload AnalyzeThumbnailPayload
    if err := json.Unmarshal(task.Payload(), &payload); err != nil {
        return fmt.Errorf("unmarshal payload: %w", err)
    }
    
    log := w.logger.With().Str("analysis_id", payload.AnalysisID).Logger()
    
    // Update status to processing
    if err := w.analysisRepo.UpdateStatus(ctx, payload.AnalysisID, "processing"); err != nil {
        return err
    }
    
    // Step 1: Download thumbnail from S3
    thumbnailData, err := w.storage.Download(ctx, payload.S3Key)
    if err != nil {
        w.markFailed(ctx, payload.AnalysisID, "failed to download thumbnail")
        return err
    }
    
    // Step 2: Run CV analysis
    cvResults, err := w.cvClient.Analyze(ctx, thumbnailData)
    if err != nil {
        w.markFailed(ctx, payload.AnalysisID, "CV analysis failed")
        return err
    }
    log.Info().Interface("cv_results", cvResults).Msg("CV analysis complete")
    
    // Step 3: Get curiosity score from Claude
    curiosityResp, err := w.claudeClient.ScoreCuriosity(ctx, claude.CuriosityRequest{
        TextContent:    cvResults.TextStrings,
        ObjectLabels:   cvResults.ObjectLabels,
        PrimaryEmotion: cvResults.PrimaryEmotion,
    })
    if err != nil {
        log.Warn().Err(err).Msg("curiosity scoring failed, defaulting to 50")
        curiosityResp = &claude.CuriosityResponse{Score: 50}
    }
    
    // Step 4: Fetch and analyze competitors
    competitors, err := w.youtubeClient.SearchVideos(ctx, payload.Keyword, 20)
    if err != nil {
        log.Warn().Err(err).Msg("failed to fetch competitors, continuing with limited data")
    }
    
    competitorResults := make([]CVResults, 0, len(competitors))
    for _, comp := range competitors {
        compData, err := w.downloadImage(ctx, comp.ThumbnailURL)
        if err != nil {
            continue
        }
        compCV, err := w.cvClient.Analyze(ctx, compData)
        if err != nil {
            continue
        }
        // Store competitor snapshot
        w.storeCompetitorSnapshot(ctx, payload.AnalysisID, comp, compCV)
        competitorResults = append(competitorResults, *compCV)
    }
    
    // Step 5: Compute industry averages
    competitorAvg := computeCompetitorAverages(competitorResults)
    
    // Step 6: Calculate ThumbnailIQ Score
    score := w.scoreEngine.Calculate(ScoreInput{
        CVResults:      cvResults,
        CuriosityScore: curiosityResp.Score,
        CompetitorAvg:  competitorAvg,
        Keyword:        payload.Keyword,
    })
    
    // Step 7: Generate suggestions
    suggestions := w.suggestionEngine.Generate(SuggestionInput{
        CVResults:     cvResults,
        Score:         score,
        CompetitorAvg: competitorAvg,
    })
    
    // Step 8: Calculate rank among competitors
    rank := calculateRank(score.Total, competitorResults)
    
    // Step 9: Update database with complete results
    if err := w.analysisRepo.UpdateResults(ctx, UpdateResultsParams{
        ID:              payload.AnalysisID,
        Score:           score.Total,
        VisibilityScore: score.Visibility,
        ContrastScore:   score.Contrast,
        AttentionScore:  score.Attention,
        MobileScore:     score.Mobile,
        BrandingScore:   score.Branding,
        CuriosityScore:  curiosityResp.Score,
        CVResults:       cvResults,
        CompetitorAvg:   competitorAvg,
        Suggestions:     suggestions,
        CompetitorCount: len(competitorResults),
        RankInCompetitors: rank,
    }); err != nil {
        return err
    }
    
    // Step 10: Notify user via WebSocket/SSE
    w.notifier.NotifyUser(ctx, payload.UserID, NotifyEvent{
        Type:       "analysis_complete",
        AnalysisID: payload.AnalysisID,
        Score:      score.Total,
    })
    
    log.Info().Int("score", score.Total).Msg("analysis complete")
    return nil
}
```

### 10.2 Tracking Worker

```go
// internal/worker/tracking.go
func (w *Worker) HandleTrackKeyword(ctx context.Context, task *asynq.Task) error {
    var payload TrackKeywordPayload
    json.Unmarshal(task.Payload(), &payload)
    
    // Fetch current top results
    results, err := w.youtubeClient.SearchVideos(ctx, payload.Keyword, 20)
    if err != nil {
        return err
    }
    
    // Load previous snapshot from DB
    previousSnapshots, _ := w.competitorRepo.GetSnapshotsForKeyword(ctx, payload.Keyword, time.Now().AddDate(0, 0, -1))
    
    previousIDs := make(map[string]bool)
    for _, s := range previousSnapshots {
        previousIDs[s.VideoID] = true
    }
    
    // Detect new or changed thumbnails
    for _, result := range results {
        if !previousIDs[result.VideoID] {
            // New video appeared - analyze and store
            w.analyzeAndStoreCompetitor(ctx, result, payload.TrackingJobID, payload.Keyword)
        } else {
            // Check if thumbnail URL changed
            prev := findPrevious(previousSnapshots, result.VideoID)
            if prev != nil && prev.ThumbnailURL != result.ThumbnailURL {
                // Thumbnail was updated - analyze new version
                w.analyzeAndStoreCompetitor(ctx, result, payload.TrackingJobID, payload.Keyword)
                w.sendChangeAlert(ctx, payload.WorkspaceID, result, prev)
            }
        }
    }
    
    // Update tracking job's next_check_at
    w.trackingRepo.UpdateNextCheck(ctx, payload.TrackingJobID, time.Now().Add(
        time.Duration(payload.IntervalHours) * time.Hour,
    ))
    
    return nil
}
```

---

## 11. Authentication & Security

### 11.1 JWT Authentication Flow

```go
// Access token: 15 minutes TTL
// Refresh token: 7 days TTL, stored hashed in DB

type AccessTokenClaims struct {
    UserID      string `json:"uid"`
    WorkspaceID string `json:"wid"`
    Plan        string `json:"plan"`
    Role        string `json:"role"`
    jwt.RegisteredClaims
}

func (j *JWTService) GenerateTokenPair(ctx context.Context, user *domain.User, workspaceID string) (*TokenPair, error) {
    // Access token
    accessClaims := AccessTokenClaims{
        UserID:      user.ID,
        WorkspaceID: workspaceID,
        Plan:        user.Plan,
        Role:        "member",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "thumbnailiq",
        },
    }
    
    accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(j.accessSecret))
    if err != nil {
        return nil, err
    }
    
    // Refresh token (random 32 bytes)
    rawRefresh := make([]byte, 32)
    rand.Read(rawRefresh)
    refreshTokenStr := base64.URLEncoding.EncodeToString(rawRefresh)
    
    // Store hashed refresh token
    hash := sha256.Sum256([]byte(refreshTokenStr))
    hashStr := hex.EncodeToString(hash[:])
    
    err = j.refreshRepo.Create(ctx, db.CreateRefreshTokenParams{
        UserID:    user.ID,
        TokenHash: hashStr,
        ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
    })
    
    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshTokenStr,
        ExpiresIn:    900,
    }, err
}
```

### 11.2 Multi-tenancy

Every database query filters by `workspace_id`. The workspace ID is embedded in the JWT claims and validated by middleware on every request:

```go
func TenantMiddleware(db *postgres.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        claims := GetClaims(c)
        workspaceID := c.Param("workspace_id")
        
        // Validate user has access to this workspace
        member, err := db.GetWorkspaceMember(c, claims.UserID, workspaceID)
        if err != nil || member == nil {
            c.AbortWithStatusJSON(403, gin.H{"error": "access denied"})
            return
        }
        
        c.Set("workspace_id", workspaceID)
        c.Set("workspace_role", member.Role)
        c.Next()
    }
}
```

### 11.3 OWASP Protections

```go
// middleware/security.go
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        c.Header("Content-Security-Policy", "default-src 'self'")
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        c.Next()
    }
}

// File upload validation
func ValidateThumbnailUpload(file multipart.File, header *multipart.FileHeader) error {
    // Check file size (max 5MB)
    if header.Size > 5*1024*1024 {
        return ErrFileTooLarge
    }
    
    // Read first 512 bytes for MIME sniffing
    buf := make([]byte, 512)
    file.Read(buf)
    file.Seek(0, 0)
    
    mimeType := http.DetectContentType(buf)
    allowed := map[string]bool{
        "image/jpeg": true,
        "image/png":  true,
        "image/webp": true,
    }
    
    if !allowed[mimeType] {
        return ErrInvalidFileType
    }
    return nil
}
```

### 11.4 Rate Limiting

```go
// 3-layer rate limiting:
// 1. Global IP-based (nginx level): 100 req/min per IP
// 2. API-key based: per plan limits
// 3. Expensive operations (analysis): per workspace monthly quota

func AnalysisQuotaMiddleware(repo WorkspaceRepo) gin.HandlerFunc {
    return func(c *gin.Context) {
        workspaceID := c.GetString("workspace_id")
        
        workspace, err := repo.GetByID(c, workspaceID)
        if err != nil {
            c.AbortWithStatusJSON(500, gin.H{"error": "internal error"})
            return
        }
        
        if workspace.AnalysesThisMonth >= workspace.AnalysesLimit {
            c.AbortWithStatusJSON(429, gin.H{
                "error": "monthly analysis limit reached",
                "limit": workspace.AnalysesLimit,
                "used": workspace.AnalysesThisMonth,
                "upgrade_url": "https://thumbnailiq.com/billing",
            })
            return
        }
        
        // Increment counter atomically
        repo.IncrementAnalysisCount(c, workspaceID)
        c.Next()
    }
}
```

---

## 12. Billing & Monetization (Stripe)

### 12.1 Plan Structure

| Feature | Free | Starter ($19/mo) | Pro ($49/mo) | Agency ($149/mo) | Enterprise |
|---------|------|-----------------|-------------|-----------------|------------|
| Analyses/month | 5 | 50 | 250 | 1,000 | Unlimited |
| Keyword tracking | 0 | 5 | 25 | 100 | Unlimited |
| Team members | 1 | 3 | 10 | 50 | Unlimited |
| API access | ❌ | ❌ | ✅ (1K/mo) | ✅ (100K/mo) | ✅ |
| Viral database | ❌ | Limited | ✅ | ✅ | ✅ |
| History retention | 7 days | 30 days | 12 months | 24 months | Unlimited |
| White-label reports | ❌ | ❌ | ❌ | ✅ | ✅ |
| Priority support | ❌ | ❌ | Email | Chat | Dedicated CSM |
| Custom model training | ❌ | ❌ | ❌ | ❌ | ✅ |

### 12.2 Stripe Integration

```go
// internal/infra/stripe/client.go

func (s *StripeClient) CreateSubscription(ctx context.Context, params CreateSubscriptionParams) (*stripe.Subscription, error) {
    // Create or retrieve customer
    customer, err := s.ensureCustomer(ctx, params.Email, params.WorkspaceID)
    if err != nil {
        return nil, err
    }
    
    // Attach payment method
    _, err = paymentmethod.Attach(params.PaymentMethodID, &stripe.PaymentMethodAttachParams{
        Customer: stripe.String(customer.ID),
    })
    if err != nil {
        return nil, fmt.Errorf("attach payment method: %w", err)
    }
    
    // Set default payment method
    _, err = customer2.Update(customer.ID, &stripe.CustomerParams{
        InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
            DefaultPaymentMethod: stripe.String(params.PaymentMethodID),
        },
    })
    
    // Create subscription
    sub, err := subscription.New(&stripe.SubscriptionParams{
        Customer: stripe.String(customer.ID),
        Items: []*stripe.SubscriptionItemsParams{
            {Price: stripe.String(params.PriceID)},
        },
        PaymentBehavior:       stripe.String("default_incomplete"),
        PaymentSettings: &stripe.SubscriptionPaymentSettingsParams{
            SaveDefaultPaymentMethod: stripe.String("on_subscription"),
        },
        Expand: []*string{stripe.String("latest_invoice.payment_intent")},
    })
    return sub, err
}

// Webhook handler
func (s *StripeClient) HandleWebhook(ctx context.Context, payload []byte, sig string) error {
    event, err := webhook.ConstructEvent(payload, sig, s.webhookSecret)
    if err != nil {
        return fmt.Errorf("invalid webhook: %w", err)
    }
    
    switch event.Type {
    case "customer.subscription.updated":
        var sub stripe.Subscription
        json.Unmarshal(event.Data.Raw, &sub)
        return s.handleSubscriptionUpdated(ctx, &sub)
    
    case "customer.subscription.deleted":
        var sub stripe.Subscription
        json.Unmarshal(event.Data.Raw, &sub)
        return s.handleSubscriptionDeleted(ctx, &sub)
    
    case "invoice.payment_succeeded":
        var invoice stripe.Invoice
        json.Unmarshal(event.Data.Raw, &invoice)
        return s.handlePaymentSucceeded(ctx, &invoice)
    
    case "invoice.payment_failed":
        var invoice stripe.Invoice
        json.Unmarshal(event.Data.Raw, &invoice)
        return s.handlePaymentFailed(ctx, &invoice)
    }
    return nil
}
```

### 12.3 Usage Metering

At the start of each calendar month, reset `analyses_this_month` for all workspaces using a Asynq scheduled job:

```go
// Cron: 0 0 1 * * (first of each month)
func (w *Worker) HandleResetMonthlyUsage(ctx context.Context, task *asynq.Task) error {
    return w.workspaceRepo.ResetMonthlyCounters(ctx)
}
```

---

## 13. Frontend Architecture

### 13.1 Key Pages

**Dashboard (`/dashboard`):**
- Recent analyses grid
- Usage quota progress bar
- Quick "New Analysis" CTA
- Recent tracking alerts
- Score trend chart (last 30 days)

**New Analysis (`/analyses/new`):**
- Drag-and-drop thumbnail upload
- Keyword input with suggestions
- Project selector
- "Analyze" button → shows progress animation → redirects to results

**Analysis Results (`/analyses/[id]`):**
- Score gauge (animated 0→score)
- 6 sub-score cards
- Device preview tabs
- Competitor grid
- Suggestions list with before/after
- CV data breakdown accordion

**Compare (`/compare`):**
- Multi-thumbnail upload
- Side-by-side grid
- Score comparison bar chart
- Winner prediction card

**Competitor Tracking (`/tracking`):**
- List of active tracking jobs
- Timeline view of changes
- Alert history

**Viral Database (`/database`):**
- Filterable/searchable gallery
- Filter panel: niche, emotion, colors, score range
- Masonry grid of thumbnails

### 13.2 State Management

```typescript
// stores/analysisStore.ts
import { create } from 'zustand'
import { devtools, persist } from 'zustand/middleware'

interface AnalysisState {
  currentAnalysis: Analysis | null
  analyses: Analysis[]
  isLoading: boolean
  error: string | null
  
  setCurrentAnalysis: (analysis: Analysis | null) => void
  addAnalysis: (analysis: Analysis) => void
  setLoading: (loading: boolean) => void
  clearError: () => void
}

export const useAnalysisStore = create<AnalysisState>()(
  devtools(
    (set) => ({
      currentAnalysis: null,
      analyses: [],
      isLoading: false,
      error: null,
      
      setCurrentAnalysis: (analysis) => set({ currentAnalysis: analysis }),
      addAnalysis: (analysis) => set((state) => ({
        analyses: [analysis, ...state.analyses]
      })),
      setLoading: (loading) => set({ isLoading: loading }),
      clearError: () => set({ error: null }),
    }),
    { name: 'analysis-store' }
  )
)
```

### 13.3 Real-time Updates (SSE)

While an analysis is processing, the frontend polls for updates via Server-Sent Events:

```typescript
// hooks/useAnalysisSSE.ts
export function useAnalysisSSE(analysisId: string, onComplete: (analysis: Analysis) => void) {
  useEffect(() => {
    if (!analysisId) return
    
    const eventSource = new EventSource(`/api/v1/analyses/${analysisId}/events`, {
      withCredentials: true,
    })
    
    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data)
      
      if (data.status === 'complete') {
        onComplete(data.analysis)
        eventSource.close()
      } else if (data.status === 'failed') {
        toast.error('Analysis failed. Please try again.')
        eventSource.close()
      }
    }
    
    eventSource.onerror = () => {
      eventSource.close()
    }
    
    return () => eventSource.close()
  }, [analysisId])
}
```

### 13.4 Device Preview Component

```tsx
// components/SERPPreview.tsx
interface SERPPreviewProps {
  userThumbnail: string
  competitors: CompetitorSnapshot[]
  userPosition: number
  device: 'desktop' | 'tablet' | 'mobile'
}

export function SERPPreview({ userThumbnail, competitors, userPosition, device }: SERPPreviewProps) {
  const deviceConfig = {
    desktop: { width: '100%', thumbnailWidth: 168, thumbnailHeight: 94 },
    tablet: { width: '768px', thumbnailWidth: 130, thumbnailHeight: 73 },
    mobile: { width: '375px', thumbnailWidth: 80, thumbnailHeight: 45 },
  }
  
  const config = deviceConfig[device]
  
  // Inject user's thumbnail at the specified position
  const displayResults = [
    ...competitors.slice(0, userPosition - 1),
    { isUser: true, thumbnailUrl: userThumbnail },
    ...competitors.slice(userPosition - 1),
  ]
  
  return (
    <div className="border rounded-lg bg-white overflow-hidden" style={{ maxWidth: config.width }}>
      {/* YouTube Header */}
      <YouTubeHeader device={device} />
      
      {/* Results List */}
      <div className="p-4 space-y-3">
        {displayResults.slice(0, 10).map((result, i) => (
          <SERPResultRow
            key={i}
            result={result}
            thumbnailWidth={config.thumbnailWidth}
            thumbnailHeight={config.thumbnailHeight}
            isHighlighted={result.isUser}
          />
        ))}
      </div>
    </div>
  )
}
```

---

## 14. Browser Extension

### 14.1 Manifest V3 Structure

```json
{
  "manifest_version": 3,
  "name": "ThumbnailIQ",
  "version": "1.0.0",
  "description": "Real-time thumbnail scoring and analysis on YouTube",
  "permissions": ["storage", "activeTab"],
  "host_permissions": ["https://www.youtube.com/*", "https://api.thumbnailiq.com/*"],
  "action": {
    "default_popup": "popup.html",
    "default_icon": { "32": "icons/icon32.png" }
  },
  "content_scripts": [{
    "matches": ["https://www.youtube.com/*"],
    "js": ["content.js"],
    "css": ["content.css"],
    "run_at": "document_idle"
  }],
  "background": {
    "service_worker": "background.js",
    "type": "module"
  }
}
```

### 14.2 Content Script

```typescript
// extension/src/content.ts
const SCORE_OVERLAY_CLASS = 'tiq-score-overlay'

function injectScoreOverlays() {
  const thumbnails = document.querySelectorAll('ytd-video-renderer, ytd-compact-video-renderer')
  
  thumbnails.forEach(async (element) => {
    const thumbnail = element.querySelector('img#img') as HTMLImageElement
    const videoLink = element.querySelector('a#thumbnail') as HTMLAnchorElement
    
    if (!thumbnail || !videoLink || element.querySelector(`.${SCORE_OVERLAY_CLASS}`)) {
      return
    }
    
    const videoId = extractVideoId(videoLink.href)
    if (!videoId) return
    
    // Get score from extension background (cached)
    const score = await chrome.runtime.sendMessage({
      type: 'GET_THUMBNAIL_SCORE',
      videoId,
      thumbnailUrl: thumbnail.src,
    })
    
    // Inject score badge
    const overlay = document.createElement('div')
    overlay.className = SCORE_OVERLAY_CLASS
    overlay.innerHTML = `
      <span class="tiq-score tiq-score-${getScoreBand(score)}">${score}</span>
    `
    
    const container = thumbnail.parentElement
    container.style.position = 'relative'
    container.appendChild(overlay)
  })
}

// Watch for DOM changes (YouTube is a SPA)
const observer = new MutationObserver(injectScoreOverlays)
observer.observe(document.body, { childList: true, subtree: true })
injectScoreOverlays()
```

---

## 15. Observability & Monitoring

### 15.1 Prometheus Metrics

```go
// pkg/metrics/metrics.go
var (
    HTTPRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration",
            Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
        },
        []string{"method", "path", "status"},
    )
    
    AnalysisJobDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "analysis_job_duration_seconds",
            Help:    "Analysis job processing time",
            Buckets: []float64{1, 5, 10, 30, 60, 120, 300},
        },
        []string{"status"},
    )
    
    CVServiceErrors = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "cv_service_errors_total",
            Help: "Computer vision service errors",
        },
        []string{"operation"},
    )
    
    YouTubeAPIQuotaUsed = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "youtube_api_quota_used",
        Help: "Current YouTube API quota usage",
    })
    
    ActiveAnalysisJobs = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "analysis_jobs_active",
        Help: "Currently processing analysis jobs",
    })
)
```

### 15.2 Zerolog Setup

```go
// pkg/logger/logger.go
func New(env string) zerolog.Logger {
    if env == "production" {
        return zerolog.New(os.Stdout).
            With().
            Timestamp().
            Str("service", "thumbnailiq-api").
            Logger()
    }
    
    return zerolog.New(zerolog.ConsoleWriter{
        Out:        os.Stdout,
        TimeFormat: time.RFC3339,
    }).With().Timestamp().Logger().Level(zerolog.DebugLevel)
}
```

### 15.3 Grafana Dashboards

Key dashboards to build:

1. **API Health:** request rate, error rate, p50/p95/p99 latency
2. **Analysis Pipeline:** job queue depth, processing time, success rate
3. **CV Service:** analysis duration per operation, error rate
4. **Business Metrics:** new signups, active users, analyses created, plan distribution
5. **YouTube Quota:** daily quota used vs. limit, projections
6. **Infrastructure:** CPU/memory per container, PostgreSQL connections, Redis memory

---

## 16. Deployment & Infrastructure

### 16.1 Docker Compose (Development)

```yaml
# docker-compose.yml
version: '3.9'

services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: thumbnailiq
      POSTGRES_USER: thumbnailiq
      POSTGRES_PASSWORD: secret
    ports: ["5432:5432"]
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports: ["6379:6379"]
    command: redis-server --save 60 1 --loglevel warning

  api:
    build:
      context: .
      dockerfile: cmd/api/Dockerfile
    ports: ["8080:8080"]
    environment:
      DATABASE_URL: postgres://thumbnailiq:secret@postgres:5432/thumbnailiq?sslmode=disable
      REDIS_ADDR: redis:6379
    depends_on: [postgres, redis]
    volumes:
      - ./config.yaml:/app/config.yaml

  worker:
    build:
      context: .
      dockerfile: cmd/worker/Dockerfile
    environment:
      DATABASE_URL: postgres://thumbnailiq:secret@postgres:5432/thumbnailiq?sslmode=disable
      REDIS_ADDR: redis:6379
    depends_on: [postgres, redis, api]

  cv-service:
    build: ./services/cv
    ports: ["8001:8001"]
    environment:
      MODEL_DIR: /models

  web:
    build: ./web
    ports: ["3000:3000"]
    environment:
      NEXT_PUBLIC_API_URL: http://localhost:8080

  prometheus:
    image: prom/prometheus:latest
    ports: ["9090:9090"]
    volumes:
      - ./deploy/prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana:latest
    ports: ["3001:3000"]
    environment:
      GF_SECURITY_ADMIN_PASSWORD: admin

volumes:
  postgres_data:
```

### 16.2 Kubernetes (Production)

```yaml
# deploy/k8s/api-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: thumbnailiq-api
  namespace: thumbnailiq
spec:
  replicas: 3
  selector:
    matchLabels:
      app: thumbnailiq-api
  template:
    metadata:
      labels:
        app: thumbnailiq-api
    spec:
      containers:
        - name: api
          image: ghcr.io/thumbnailiq/api:latest
          ports:
            - containerPort: 8080
          env:
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: thumbnailiq-secrets
                  key: database-url
            - name: REDIS_ADDR
              value: "redis-service:6379"
          resources:
            requests:
              cpu: "250m"
              memory: "256Mi"
            limits:
              cpu: "500m"
              memory: "512Mi"
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 10
            periodSeconds: 5
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
---
apiVersion: v1
kind: Service
metadata:
  name: thumbnailiq-api-service
  namespace: thumbnailiq
spec:
  selector:
    app: thumbnailiq-api
  ports:
    - port: 80
      targetPort: 8080
  type: ClusterIP
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: thumbnailiq-api-hpa
  namespace: thumbnailiq
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: thumbnailiq-api
  minReplicas: 2
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
```

### 16.3 Terraform (AWS Infrastructure)

```hcl
# infrastructure/terraform/main.tf

terraform {
  required_providers {
    aws = { source = "hashicorp/aws", version = "~> 5.0" }
  }
  backend "s3" {
    bucket = "thumbnailiq-terraform-state"
    key    = "production/terraform.tfstate"
    region = "us-east-1"
  }
}

# VPC
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.0.0"
  name    = "thumbnailiq-vpc"
  cidr    = "10.0.0.0/16"
  azs             = ["us-east-1a", "us-east-1b", "us-east-1c"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]
  enable_nat_gateway = true
}

# EKS Cluster
module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "20.0.0"
  cluster_name    = "thumbnailiq-prod"
  cluster_version = "1.29"
  vpc_id          = module.vpc.vpc_id
  subnet_ids      = module.vpc.private_subnets
  
  eks_managed_node_groups = {
    general = {
      min_size     = 2
      max_size     = 10
      desired_size = 3
      instance_types = ["t3.medium"]
    }
    cv_service = {
      min_size     = 1
      max_size     = 5
      desired_size = 2
      instance_types = ["c5.xlarge"]  # CPU-optimized for CV
      labels = { "workload" = "cv" }
    }
  }
}

# RDS PostgreSQL
resource "aws_db_instance" "postgres" {
  identifier             = "thumbnailiq-postgres"
  engine                 = "postgres"
  engine_version         = "16.2"
  instance_class         = "db.t3.medium"
  allocated_storage      = 100
  storage_type           = "gp3"
  storage_encrypted      = true
  db_name                = "thumbnailiq"
  username               = "thumbnailiq"
  password               = var.db_password
  vpc_security_group_ids = [aws_security_group.rds.id]
  db_subnet_group_name   = aws_db_subnet_group.main.name
  backup_retention_period = 7
  deletion_protection     = true
  skip_final_snapshot     = false
  multi_az               = true
}

# ElastiCache Redis
resource "aws_elasticache_cluster" "redis" {
  cluster_id           = "thumbnailiq-redis"
  engine               = "redis"
  node_type            = "cache.t3.medium"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis7"
  port                 = 6379
  subnet_group_name    = aws_elasticache_subnet_group.main.name
  security_group_ids   = [aws_security_group.redis.id]
}

# S3 Buckets
resource "aws_s3_bucket" "uploads" {
  bucket = "thumbnailiq-uploads-${var.environment}"
}

resource "aws_s3_bucket_server_side_encryption_configuration" "uploads" {
  bucket = aws_s3_bucket.uploads.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# CloudFront
resource "aws_cloudfront_distribution" "main" {
  origin {
    domain_name = aws_s3_bucket.uploads.bucket_regional_domain_name
    origin_id   = "S3-uploads"
    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.main.cloudfront_access_identity_path
    }
  }
  enabled = true
  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "S3-uploads"
    forwarded_values {
      query_string = false
      cookies { forward = "none" }
    }
    viewer_protocol_policy = "redirect-to-https"
    min_ttl     = 0
    default_ttl = 86400
    max_ttl     = 31536000
  }
  restrictions {
    geo_restriction { restriction_type = "none" }
  }
  viewer_certificate {
    cloudfront_default_certificate = true
  }
}
```

### 16.4 GitHub Actions CI/CD

```yaml
# .github/workflows/deploy.yml
name: Deploy Production

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_DB: thumbnailiq_test
          POSTGRES_USER: thumbnailiq
          POSTGRES_PASSWORD: secret
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
      redis:
        image: redis:7
        options: --health-cmd "redis-cli ping" --health-interval 10s
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.22' }
      - name: Run tests
        run: go test ./... -race -coverprofile=coverage.txt
        env:
          DATABASE_URL: postgres://thumbnailiq:secret@localhost:5432/thumbnailiq_test
          REDIS_ADDR: localhost:6379
      - name: Upload coverage
        uses: codecov/codecov-action@v4

  build-and-push:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Log in to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - name: Build and push API
        uses: docker/build-push-action@v5
        with:
          context: .
          file: cmd/api/Dockerfile
          push: true
          tags: ghcr.io/thumbnailiq/api:${{ github.sha }},ghcr.io/thumbnailiq/api:latest
      - name: Build and push Worker
        uses: docker/build-push-action@v5
        with:
          context: .
          file: cmd/worker/Dockerfile
          push: true
          tags: ghcr.io/thumbnailiq/worker:${{ github.sha }},ghcr.io/thumbnailiq/worker:latest

  deploy:
    needs: build-and-push
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: aws-actions/configure-aws-credentials@v4
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: us-east-1
      - name: Update kubeconfig
        run: aws eks update-kubeconfig --name thumbnailiq-prod
      - name: Deploy API
        run: |
          kubectl set image deployment/thumbnailiq-api \
            api=ghcr.io/thumbnailiq/api:${{ github.sha }} \
            -n thumbnailiq
          kubectl rollout status deployment/thumbnailiq-api -n thumbnailiq
      - name: Deploy Worker
        run: |
          kubectl set image deployment/thumbnailiq-worker \
            worker=ghcr.io/thumbnailiq/worker:${{ github.sha }} \
            -n thumbnailiq
          kubectl rollout status deployment/thumbnailiq-worker -n thumbnailiq
      - name: Run DB Migrations
        run: |
          kubectl run migrations --image=ghcr.io/thumbnailiq/api:${{ github.sha }} \
            --restart=Never --rm -it \
            -n thumbnailiq \
            -- /app/api migrate
```

### 16.5 Cost Estimates (AWS, Production)

| Component | Spec | Monthly Cost |
|-----------|------|-------------|
| EKS Control Plane | Managed | $73 |
| EKS Nodes (general, 3× t3.medium) | On-demand | $90 |
| EKS Nodes (cv, 2× c5.xlarge) | On-demand | $275 |
| RDS PostgreSQL (db.t3.medium, Multi-AZ) | | $98 |
| ElastiCache Redis (cache.t3.medium) | | $37 |
| S3 (500GB storage + transfer) | | $15 |
| CloudFront (1TB transfer) | | $85 |
| YouTube Data API (quota increases) | | $50 |
| Claude API (curiosity scoring) | | $30 |
| Route53 + ACM | | $5 |
| ALB | | $20 |
| Data transfer (egress) | | $40 |
| **Total** | | **~$818/month** |

At $19 average revenue per user, breakeven is ~43 paying customers.

---

## 17. Analytics

### 17.1 Event Tracking (PostHog / Mixpanel)

```typescript
// lib/analytics.ts
type EventName = 
  | 'analysis_created'
  | 'analysis_completed'
  | 'comparison_started'
  | 'tracking_job_created'
  | 'viral_db_searched'
  | 'upgrade_clicked'
  | 'api_key_created'
  | 'team_member_invited'

export function trackEvent(event: EventName, properties?: Record<string, unknown>) {
  if (typeof window === 'undefined') return
  
  // PostHog
  posthog.capture(event, {
    ...properties,
    timestamp: new Date().toISOString(),
  })
}

// Key funnels to track:
// 1. Signup → First Analysis → Score Viewed → Plan Upgraded
// 2. Free User → Analysis Limit Hit → Upgrade Shown → Converted
// 3. User → Tracking Setup → Alert Received → Engaged
```

### 17.2 Server-Side Product Analytics

```go
// Track key business events server-side for accuracy
type AnalyticsEvent struct {
    UserID     string
    Event      string
    Properties map[string]interface{}
    Timestamp  time.Time
}

func (a *AnalyticsService) Track(ctx context.Context, event AnalyticsEvent) {
    // Write to analytics DB table for retention
    // Also forward to PostHog via server-side SDK
}
```

### 17.3 Revenue Tracking

Track via Stripe webhooks → write to `revenue_events` table:

```sql
CREATE TABLE revenue_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID REFERENCES workspaces(id),
    event_type VARCHAR(100), -- 'subscription_created', 'payment_succeeded', 'churned'
    plan VARCHAR(50),
    amount_cents INT,
    currency VARCHAR(3),
    stripe_event_id VARCHAR(255) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

MRR calculation:
```sql
SELECT 
    date_trunc('month', created_at) AS month,
    SUM(amount_cents) / 100.0 AS mrr
FROM revenue_events
WHERE event_type = 'payment_succeeded'
GROUP BY 1
ORDER BY 1 DESC;
```

---

## 18. Product Roadmap

### Phase 1 — MVP (Weeks 1–4)

**Goal:** Validate core value proposition with paying users

**Features:**
- User auth (email + Google OAuth)
- Single thumbnail upload + keyword analysis
- SERP preview (desktop + mobile)
- Basic ThumbnailIQ score (simplified, without historical competitors)
- Top 5 improvement suggestions
- Stripe billing (Free + Starter)

**Team:** 1 backend engineer, 1 frontend engineer, 1 designer

**Estimated Costs:** $2,000–$4,000/month (infrastructure + tools)

**Expected Outcomes:**
- 100 beta signups
- 10 paying Starter plan users
- Validate "aha moment" (seeing your thumbnail in SERP preview)

---

### Phase 2 — Paid Product (Weeks 5–12)

**Goal:** Reach $5,000 MRR

**Features:**
- Full 6-dimension scoring
- Competitor analysis (20 videos)
- Multi-version comparison
- Pro plan launch
- Team collaboration (comments, roles)
- Browser extension v1
- Email notifications

**Team:** +1 backend (CV specialist), +1 growth marketer

**Estimated Costs:** $5,000–$8,000/month

**Expected Outcomes:**
- 50 Starter + 20 Pro subscribers
- $1,500 MRR → $5,000 MRR
- NPS > 50

---

### Phase 3 — Growth (Months 4–6)

**Goal:** Reach $20,000 MRR

**Features:**
- Historical competitor tracking
- Viral thumbnail database (initial 10,000 thumbnails)
- Agency plan launch
- White-label reports (PDF)
- Public API v1
- YouTube Studio integration (unofficial)

**Team:** +1 frontend, +1 data engineer

**Estimated Costs:** $10,000–$15,000/month

**Expected Outcomes:**
- 500+ paying users
- $20,000 MRR
- First 5 agency customers ($149/mo each)

---

### Phase 4 — Scale (Months 7–12)

**Goal:** Reach $80,000 MRR

**Features:**
- Advanced AI improvement suggestions with visual examples
- A/B test framework integration (post-publish tracking)
- Custom branding score (channel-specific)
- Niche benchmarks (Gaming, Finance, Fitness, etc.)
- Mobile app (React Native)
- SEO integration (keyword CTR data from Search Console)

**Team:** Hire product manager, 2 more engineers

**Estimated Costs:** $25,000–$40,000/month

**Expected Outcomes:**
- 2,000+ paying users
- $80,000 MRR
- Enterprise pipeline (5+ contracts)

---

### Phase 5 — Enterprise (Year 2)

**Goal:** $500,000 ARR

**Features:**
- Enterprise SSO (SAML/OKTA)
- Custom model training (niche-specific scoring)
- Dedicated API SLA
- MCN portfolio management
- Data export / data warehouse integration
- SOC2 Type II compliance

**Team:** Enterprise sales, customer success, compliance

**Estimated Costs:** $100,000+/month

**Expected Outcomes:**
- 50+ enterprise contracts
- $500K ARR
- Acquisition conversations or Series A

---

## 19. Repository Structure & Claude Code Execution Plan

### 19.1 Complete Repository Structure

```
thumbnailiq/
├── .github/
│   ├── workflows/
│   │   ├── ci.yml
│   │   ├── deploy.yml
│   │   └── release.yml
│   └── CODEOWNERS
│
├── cmd/
│   ├── api/
│   │   ├── main.go
│   │   └── Dockerfile
│   ├── worker/
│   │   ├── main.go
│   │   └── Dockerfile
│   └── migrate/
│       └── main.go
│
├── internal/
│   ├── domain/
│   │   ├── analysis/
│   │   │   ├── analysis.go         (entity)
│   │   │   ├── repository.go       (interface)
│   │   │   └── errors.go
│   │   ├── competitor/
│   │   │   ├── competitor.go
│   │   │   └── repository.go
│   │   ├── user/
│   │   │   ├── user.go
│   │   │   └── repository.go
│   │   ├── workspace/
│   │   │   ├── workspace.go
│   │   │   └── repository.go
│   │   └── billing/
│   │       ├── subscription.go
│   │       └── repository.go
│   │
│   ├── usecase/
│   │   ├── analysis/
│   │   │   ├── create_analysis.go
│   │   │   ├── get_analysis.go
│   │   │   └── list_analyses.go
│   │   ├── competitor/
│   │   │   ├── fetch_competitors.go
│   │   │   └── track_competitors.go
│   │   ├── user/
│   │   │   ├── register.go
│   │   │   ├── login.go
│   │   │   └── refresh_token.go
│   │   ├── workspace/
│   │   │   ├── create_workspace.go
│   │   │   └── manage_members.go
│   │   └── billing/
│   │       ├── subscribe.go
│   │       └── handle_webhook.go
│   │
│   ├── handler/
│   │   ├── analysis_handler.go
│   │   ├── auth_handler.go
│   │   ├── competitor_handler.go
│   │   ├── workspace_handler.go
│   │   ├── billing_handler.go
│   │   ├── tracking_handler.go
│   │   ├── viral_db_handler.go
│   │   └── apikey_handler.go
│   │
│   ├── worker/
│   │   ├── server.go
│   │   ├── analyze_thumbnail.go
│   │   ├── analyze_competitors.go
│   │   ├── track_keyword.go
│   │   ├── track_channel.go
│   │   ├── generate_report.go
│   │   ├── send_notification.go
│   │   └── reset_monthly_usage.go
│   │
│   ├── scoring/
│   │   ├── engine.go
│   │   ├── visibility.go
│   │   ├── contrast.go
│   │   ├── attention.go
│   │   ├── mobile.go
│   │   ├── branding.go
│   │   └── suggestions.go
│   │
│   ├── infra/
│   │   ├── postgres/
│   │   │   ├── db/                  (SQLC generated)
│   │   │   │   ├── db.go
│   │   │   │   ├── models.go
│   │   │   │   ├── analyses.sql.go
│   │   │   │   ├── users.sql.go
│   │   │   │   ├── workspaces.sql.go
│   │   │   │   └── competitors.sql.go
│   │   │   ├── analysis_repo.go     (implements domain.AnalysisRepository)
│   │   │   ├── user_repo.go
│   │   │   └── workspace_repo.go
│   │   ├── redis/
│   │   │   ├── cache.go
│   │   │   └── rate_limiter.go
│   │   ├── s3/
│   │   │   └── storage.go
│   │   ├── youtube/
│   │   │   └── client.go
│   │   ├── cv/
│   │   │   └── client.go
│   │   ├── stripe/
│   │   │   └── client.go
│   │   └── claude/
│   │       └── client.go
│   │
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── rate_limit.go
│   │   ├── tenant.go
│   │   ├── metrics.go
│   │   ├── cors.go
│   │   └── security.go
│   │
│   ├── server/
│   │   ├── router.go
│   │   └── server.go
│   │
│   └── config/
│       └── config.go
│
├── pkg/
│   ├── jwt/
│   │   └── jwt.go
│   ├── validator/
│   │   └── validator.go
│   ├── pagination/
│   │   └── pagination.go
│   ├── errors/
│   │   └── errors.go
│   ├── logger/
│   │   └── logger.go
│   └── hash/
│       └── hash.go
│
├── db/
│   ├── migrations/
│   │   ├── 20240101000001_create_users.sql
│   │   ├── 20240101000002_create_workspaces.sql
│   │   ├── 20240101000003_create_projects.sql
│   │   ├── 20240101000004_create_analyses.sql
│   │   ├── 20240101000005_create_competitors.sql
│   │   ├── 20240101000006_create_tracking.sql
│   │   ├── 20240101000007_create_comments_api_keys.sql
│   │   └── 20240101000008_create_viral_db.sql
│   ├── queries/
│   │   ├── analyses.sql
│   │   ├── users.sql
│   │   ├── workspaces.sql
│   │   ├── competitors.sql
│   │   ├── tracking.sql
│   │   └── viral_thumbnails.sql
│   └── sqlc.yaml
│
├── services/
│   └── cv/
│       ├── main.py
│       ├── app/
│       │   ├── routes.py
│       │   ├── models.py
│       │   ├── ocr.py
│       │   ├── face.py
│       │   ├── colors.py
│       │   └── objects.py
│       ├── requirements.txt
│       └── Dockerfile
│
├── web/                             (Next.js app)
│   ├── app/
│   │   ├── (auth)/
│   │   ├── (dashboard)/
│   │   └── layout.tsx
│   ├── components/
│   │   ├── ui/                      (shadcn components)
│   │   ├── analysis/
│   │   │   ├── ScoreGauge.tsx
│   │   │   ├── SERPPreview.tsx
│   │   │   ├── CompetitorGrid.tsx
│   │   │   ├── SuggestionList.tsx
│   │   │   └── CVBreakdown.tsx
│   │   ├── compare/
│   │   │   └── CompareGrid.tsx
│   │   ├── tracking/
│   │   │   └── TrackingTimeline.tsx
│   │   └── shared/
│   │       ├── Layout.tsx
│   │       ├── Navbar.tsx
│   │       └── Sidebar.tsx
│   ├── hooks/
│   │   ├── useAnalysis.ts
│   │   ├── useAnalysisSSE.ts
│   │   └── useWorkspace.ts
│   ├── lib/
│   │   ├── api.ts
│   │   ├── auth.ts
│   │   └── analytics.ts
│   ├── stores/
│   │   ├── analysisStore.ts
│   │   └── workspaceStore.ts
│   ├── types/
│   │   └── index.ts
│   ├── next.config.ts
│   ├── tailwind.config.ts
│   └── package.json
│
├── extension/                       (Browser extension)
│   ├── src/
│   │   ├── background.ts
│   │   ├── content.ts
│   │   ├── popup/
│   │   │   ├── popup.html
│   │   │   └── popup.ts
│   │   └── types.ts
│   ├── public/
│   │   ├── manifest.json
│   │   └── icons/
│   ├── tsconfig.json
│   └── package.json
│
├── deploy/
│   ├── k8s/
│   │   ├── namespace.yaml
│   │   ├── api-deployment.yaml
│   │   ├── worker-deployment.yaml
│   │   ├── cv-deployment.yaml
│   │   ├── web-deployment.yaml
│   │   ├── ingress.yaml
│   │   └── secrets.yaml
│   └── prometheus.yml
│
├── infrastructure/
│   └── terraform/
│       ├── main.tf
│       ├── variables.tf
│       ├── outputs.tf
│       ├── vpc.tf
│       ├── eks.tf
│       ├── rds.tf
│       ├── redis.tf
│       ├── s3.tf
│       └── cloudfront.tf
│
├── docs/
│   ├── api.yaml                     (OpenAPI spec)
│   ├── architecture.md
│   └── runbook.md
│
├── config.yaml.example
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### 19.2 Claude Code Execution Sequence

**Execute in this exact order to build the product from scratch:**

#### Step 1: Project Scaffolding
```bash
# 1.1 Initialize Go module
go mod init github.com/thumbnailiq/thumbnailiq

# 1.2 Install Go dependencies
go get github.com/gin-gonic/gin
go get github.com/hibiken/asynq
go get github.com/redis/go-redis/v9
go get github.com/aws/aws-sdk-go-v2
go get github.com/stripe/stripe-go/v78
go get github.com/golang-jwt/jwt/v5
go get github.com/rs/zerolog
go get github.com/spf13/viper
go get github.com/jackc/pgx/v5
go get github.com/pressly/goose/v3
go get github.com/google/wire/cmd/wire

# 1.3 Initialize Next.js web app
npx create-next-app@latest web --typescript --tailwind --app --src-dir=false

# 1.4 Install shadcn/ui
cd web && npx shadcn-ui@latest init

# 1.5 Initialize Python CV service
mkdir -p services/cv && cd services/cv
python -m venv .venv && source .venv/bin/activate
pip install fastapi uvicorn easyocr insightface deepface ultralytics scikit-learn pillow opencv-python
```

#### Step 2: Core Infrastructure
```
Build in order:
2.1  pkg/logger/logger.go         — Zerolog setup
2.2  pkg/errors/errors.go         — Standard error types
2.3  internal/config/config.go    — Viper config loading
2.4  pkg/jwt/jwt.go               — JWT service
2.5  pkg/hash/hash.go             — bcrypt + SHA256 helpers
2.6  pkg/validator/validator.go   — Input validation
2.7  pkg/pagination/pagination.go — Cursor/offset pagination
```

#### Step 3: Database Layer
```
Build in order:
3.1  db/migrations/               — All 8 migration files
3.2  db/queries/                  — All SQLC query files
3.3  db/sqlc.yaml                 — SQLC configuration
3.4  Run: sqlc generate           — Generate typed Go code
3.5  internal/infra/postgres/     — Repository implementations
```

#### Step 4: Domain + Use Cases
```
Build in order:
4.1  internal/domain/user/        — User entity + repository interface
4.2  internal/domain/workspace/   — Workspace entity
4.3  internal/domain/analysis/    — Analysis entity
4.4  internal/domain/competitor/  — Competitor entity
4.5  internal/domain/billing/     — Subscription entity
4.6  internal/usecase/user/       — Register, Login, RefreshToken
4.7  internal/usecase/workspace/  — Create workspace, Manage members
4.8  internal/usecase/analysis/   — Create analysis, Get analysis
4.9  internal/usecase/billing/    — Subscribe, Handle webhook
```

#### Step 5: External Service Clients
```
Build in order:
5.1  internal/infra/redis/        — Cache + rate limiter
5.2  internal/infra/s3/           — S3 storage
5.3  internal/infra/youtube/      — YouTube API client
5.4  internal/infra/cv/           — CV service HTTP client
5.5  internal/infra/claude/       — Claude API client
5.6  internal/infra/stripe/       — Stripe client
```

#### Step 6: Scoring Engine
```
Build in order:
6.1  internal/scoring/engine.go       — Main scoring coordinator
6.2  internal/scoring/visibility.go   — Visibility sub-score
6.3  internal/scoring/contrast.go     — Contrast sub-score
6.4  internal/scoring/attention.go    — Attention sub-score
6.5  internal/scoring/mobile.go       — Mobile sub-score
6.6  internal/scoring/branding.go     — Branding sub-score
6.7  internal/scoring/suggestions.go  — Suggestion generator
```

#### Step 7: HTTP API Layer
```
Build in order:
7.1  internal/middleware/auth.go       — JWT auth middleware
7.2  internal/middleware/rate_limit.go — Rate limiting
7.3  internal/middleware/tenant.go     — Workspace validation
7.4  internal/middleware/metrics.go    — Prometheus middleware
7.5  internal/middleware/cors.go       — CORS config
7.6  internal/middleware/security.go   — Security headers
7.7  internal/handler/auth_handler.go  — /auth routes
7.8  internal/handler/analysis_handler.go
7.9  internal/handler/workspace_handler.go
7.10 internal/handler/billing_handler.go
7.11 internal/handler/tracking_handler.go
7.12 internal/handler/viral_db_handler.go
7.13 internal/handler/apikey_handler.go
7.14 internal/server/router.go         — Route registration
7.15 cmd/api/main.go                   — API entrypoint
```

#### Step 8: Background Workers
```
Build in order:
8.1  internal/worker/analyze_thumbnail.go   — Main analysis job
8.2  internal/worker/analyze_competitors.go — Competitor batch analysis
8.3  internal/worker/track_keyword.go       — Keyword tracking
8.4  internal/worker/track_channel.go       — Channel tracking
8.5  internal/worker/reset_monthly_usage.go — Monthly reset cron
8.6  internal/worker/send_notification.go   — Notification job
8.7  internal/worker/server.go              — Asynq server setup
8.8  cmd/worker/main.go                     — Worker entrypoint
```

#### Step 9: Python CV Service
```
Build in order:
9.1  services/cv/app/models.py         — Pydantic request/response models
9.2  services/cv/app/ocr.py            — EasyOCR integration
9.3  services/cv/app/face.py           — InsightFace + DeepFace
9.4  services/cv/app/colors.py         — Color analysis
9.5  services/cv/app/objects.py        — YOLOv8 object detection
9.6  services/cv/app/routes.py         — FastAPI routes
9.7  services/cv/main.py               — FastAPI app entrypoint
9.8  services/cv/Dockerfile            — CV service container
```

#### Step 10: Next.js Frontend
```
Build in order:
10.1  web/types/index.ts                        — TypeScript types
10.2  web/lib/api.ts                            — API client (axios)
10.3  web/lib/auth.ts                           — Auth helpers
10.4  web/stores/workspaceStore.ts              — Workspace store
10.5  web/stores/analysisStore.ts               — Analysis store
10.6  web/hooks/useWorkspace.ts                 — Workspace hook
10.7  web/hooks/useAnalysis.ts                  — Analysis hook
10.8  web/hooks/useAnalysisSSE.ts               — SSE hook
10.9  web/components/shared/Layout.tsx          — App layout
10.10 web/components/shared/Sidebar.tsx         — Navigation
10.11 web/components/analysis/ScoreGauge.tsx    — Score display
10.12 web/components/analysis/SERPPreview.tsx   — SERP mockup
10.13 web/components/analysis/CompetitorGrid.tsx
10.14 web/components/analysis/SuggestionList.tsx
10.15 web/components/analysis/CVBreakdown.tsx
10.16 web/components/compare/CompareGrid.tsx
10.17 web/app/(auth)/login/page.tsx
10.18 web/app/(auth)/register/page.tsx
10.19 web/app/(dashboard)/layout.tsx
10.20 web/app/(dashboard)/dashboard/page.tsx
10.21 web/app/(dashboard)/analyses/new/page.tsx
10.22 web/app/(dashboard)/analyses/[id]/page.tsx
10.23 web/app/(dashboard)/compare/page.tsx
10.24 web/app/(dashboard)/tracking/page.tsx
10.25 web/app/(dashboard)/database/page.tsx
10.26 web/app/(dashboard)/billing/page.tsx
10.27 web/app/(dashboard)/settings/page.tsx
```

#### Step 11: Browser Extension
```
Build in order:
11.1  extension/public/manifest.json
11.2  extension/src/types.ts
11.3  extension/src/background.ts     — Service worker
11.4  extension/src/content.ts        — YouTube DOM injection
11.5  extension/src/popup/popup.html
11.6  extension/src/popup/popup.ts
```

#### Step 12: Infrastructure & Deployment
```
Build in order:
12.1  cmd/api/Dockerfile
12.2  cmd/worker/Dockerfile
12.3  services/cv/Dockerfile
12.4  web/Dockerfile
12.5  docker-compose.yml               — Local development
12.6  deploy/prometheus.yml
12.7  deploy/k8s/*.yaml                — Kubernetes manifests
12.8  infrastructure/terraform/*.tf    — AWS infrastructure
12.9  .github/workflows/ci.yml
12.10 .github/workflows/deploy.yml
12.11 Makefile                         — Common commands
```

### 19.3 Makefile

```makefile
.PHONY: dev build test migrate generate lint

# Development
dev:
	docker-compose up -d postgres redis
	go run cmd/api/main.go &
	go run cmd/worker/main.go &
	cd web && npm run dev

# Build
build:
	go build -o bin/api cmd/api/main.go
	go build -o bin/worker cmd/worker/main.go
	cd web && npm run build

# Test
test:
	go test ./... -race -count=1

test-cover:
	go test ./... -race -coverprofile=coverage.txt
	go tool cover -html=coverage.txt

# Database
migrate-up:
	goose -dir db/migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir db/migrations postgres "$(DATABASE_URL)" down

migrate-status:
	goose -dir db/migrations postgres "$(DATABASE_URL)" status

# Code generation
generate:
	sqlc generate
	go generate ./...
	wire ./cmd/api/... ./cmd/worker/...

# Lint
lint:
	golangci-lint run ./...
	cd web && npm run lint

# Docker
docker-build:
	docker build -f cmd/api/Dockerfile -t thumbnailiq-api .
	docker build -f cmd/worker/Dockerfile -t thumbnailiq-worker .
	docker build -f services/cv/Dockerfile -t thumbnailiq-cv services/cv/
	docker build -f web/Dockerfile -t thumbnailiq-web web/

# Infrastructure
tf-plan:
	cd infrastructure/terraform && terraform plan

tf-apply:
	cd infrastructure/terraform && terraform apply

# Secrets (using AWS Secrets Manager)
rotate-secrets:
	./scripts/rotate_secrets.sh
```

---

## Appendix A: Environment Variables Reference

```bash
# Server
SERVER_PORT=8080
SERVER_ENV=production

# Database
DATABASE_URL=postgres://user:pass@host:5432/thumbnailiq?sslmode=require
DATABASE_MAX_CONNS=25
DATABASE_MIN_CONNS=5

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

# JWT
JWT_ACCESS_SECRET=<32-char-random>
JWT_REFRESH_SECRET=<32-char-random>
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h

# AWS
AWS_REGION=us-east-1
AWS_UPLOAD_BUCKET=thumbnailiq-uploads-prod
AWS_COMPETITOR_BUCKET=thumbnailiq-competitors-prod
AWS_REPORT_BUCKET=thumbnailiq-reports-prod

# YouTube
YOUTUBE_API_KEY=<api-key>
YOUTUBE_DAILY_QUOTA=10000

# Stripe
STRIPE_SECRET_KEY=sk_live_...
STRIPE_WEBHOOK_SECRET=whsec_...
STRIPE_PRICE_ID_STARTER=price_...
STRIPE_PRICE_ID_PRO=price_...
STRIPE_PRICE_ID_AGENCY=price_...

# CV Service
CV_SERVICE_URL=http://cv-service:8001
CV_SERVICE_API_KEY=<internal-key>

# Claude API
CLAUDE_API_KEY=sk-ant-...
CLAUDE_MODEL=claude-3-haiku-20240307

# Email
RESEND_API_KEY=re_...
EMAIL_FROM=noreply@thumbnailiq.com
```

---

## Appendix B: Security Checklist

- [ ] All secrets in AWS Secrets Manager (not in environment or code)
- [ ] Database SSL enabled (`sslmode=require`)
- [ ] S3 buckets private (no public read), accessed via CloudFront + pre-signed URLs
- [ ] JWT RS256 for production (not HS256)
- [ ] Refresh token rotation on every use
- [ ] Rate limiting at nginx, API, and application level
- [ ] File upload validation: MIME type check, max size, virus scanning (ClamAV)
- [ ] OWASP headers on all responses
- [ ] SQL injection prevention: SQLC only, no raw string interpolation
- [ ] XSS prevention: React escaping + CSP headers
- [ ] CSRF protection on state-changing operations
- [ ] Audit log table for all sensitive operations
- [ ] PII handling: user emails encrypted at rest
- [ ] GDPR: data export + deletion endpoints implemented
- [ ] Dependency scanning: `govulncheck`, `npm audit` in CI

---

*End of ThumbnailIQ Production Blueprint — Version 1.0*
*Total estimated development time: 16–20 weeks for full feature set*
*Minimum viable product: 4 weeks*
