# ThumbnailIQ — Admin Panel

A back-office panel for managing customers, moderating uploads, watching
product analytics, tuning app-wide settings, and auditing admin activity —
fully separate from the customer-facing product (`web/`) in both identity
and deployment.

---

## 1. Architecture at a glance

- **Backend**: lives inside the existing Go/Gin API (`cmd/api`), following
  the same Clean Architecture layering as every other feature
  (`domain` → `usecase` → `handler` → `infra/postgres`, sqlc-only queries,
  goose migrations). No parallel architecture was introduced.
- **Frontend**: a brand-new, standalone Vite + React + TypeScript app at
  `admin-web/`, sibling to the customer `web/` app — separate dev server,
  separate deployable, separate auth token storage.
- **Identity**: `admin_users` is a completely separate table/login from
  customer `users`. Admin JWTs are signed with their own secret
  (`ADMIN_JWT_*`), so an admin token and a customer token can never be
  confused with or forged from one another, even though both use the same
  `pkg/jwt.Service` code.
- **No self-service admin registration.** The only way to create an admin
  account is `make admin-seed` (see §5). This is intentional.

---

## 2. Database changes

### New tables (migration `20240101000013_create_admin_tables.sql`)
- `admin_users` — id, email, password_hash (bcrypt), full_name, role,
  is_active, last_login_at, timestamps.
- `admin_refresh_tokens` — same shape as the existing customer
  `refresh_tokens` table (SHA-256-hashed, rotated on every refresh).
- `admin_audit_logs` — admin_id, action, target_type, target_id, metadata
  (jsonb), created_at. Every mutating admin action writes one row here.

### Altered tables
- `20240101000014_add_user_status.sql` — adds `users.status`
  (`active`/`suspended`, CHECK-constrained). "Delete user" reuses the
  **existing** `users.deleted_at` soft-delete column.
- `20240101000015_add_upload_file_size.sql` — adds nullable
  `file_size_bytes` to `analyses` and `thumbnail_versions`, captured on
  every new upload/version going forward (old rows are `NULL`). "Delete/
  restore upload" reuses the **existing** `analyses.deleted_at` column.

### New table
- `20240101000016_create_app_settings.sql` — singleton `app_settings` row
  (`id` fixed at 1) holding upload limits, allowed extensions, feature
  flags (jsonb), storage/email provider config. Seeded with defaults.

All new sqlc queries live in `db/queries/admin.sql`; the corresponding
`analyses.sql`/`users.sql` queries were extended in place, not duplicated.

---

## 3. New backend code

```
internal/domain/admin/admin.go        Admin, AuditLog, UserSummary, UserDetail,
                                       UploadSummary, DashboardStats, Analytics,
                                       Settings, SystemHealth types + the
                                       single admin.Repository interface
internal/usecase/admin/admin.go       Login/Refresh (mirrors usecase/user),
                                       every user/upload/settings mutation,
                                       audit logging, health aggregation
internal/infra/postgres/admin_repo.go sqlc-backed Repository implementation
internal/infra/health/health.go       Thin DB/Redis/CV-service liveness probe
                                       for the dashboard's system-health widget
internal/middleware/admin_auth.go     AdminAuth — structurally identical to
                                       the existing Auth middleware, but
                                       keyed off the separate admin JWT secret
internal/handler/admin_*.go           One flat handler file per resource,
                                       matching the existing handler/ layout
                                       (auth, dashboard, users, uploads,
                                       analytics, settings, logs, profile)
cmd/admin-seed/main.go                CLI to bootstrap the first admin account
```

No existing repository interface, domain type, or route changed shape —
the admin panel is additive.

---

## 4. Route list

All under the existing `/api/v1` prefix, versioned identically to every
other endpoint.

```
POST   /api/v1/admin/auth/login              rate-limited (10/min/IP)
POST   /api/v1/admin/auth/refresh

GET    /api/v1/admin/dashboard               stats + system_health

GET    /api/v1/admin/users                   ?page&per_page&search&status
GET    /api/v1/admin/users/:id
PATCH  /api/v1/admin/users/:id/suspend
PATCH  /api/v1/admin/users/:id/activate
DELETE /api/v1/admin/users/:id
POST   /api/v1/admin/users/:id/reset-password
PATCH  /api/v1/admin/users/:id/role          { workspace_id, role }
GET    /api/v1/admin/users/:id/uploads

GET    /api/v1/admin/uploads                 ?page&per_page&search&status&include_deleted
GET    /api/v1/admin/uploads/:id
DELETE /api/v1/admin/uploads/:id
POST   /api/v1/admin/uploads/:id/restore
GET    /api/v1/admin/uploads/:id/download    302 → CDN URL

GET    /api/v1/admin/analytics

GET    /api/v1/admin/settings
PATCH  /api/v1/admin/settings

GET    /api/v1/admin/logs/audit              ?page&per_page

GET    /api/v1/admin/profile
PATCH  /api/v1/admin/profile/password        { current_password, new_password }
```

Every route except `auth/login` and `auth/refresh` requires
`Authorization: Bearer <admin access token>` via the `AdminAuth` middleware.

---

## 5. Environment variables

Already added to `.env.example`:

```
ADMIN_JWT_ACCESS_SECRET=dev-admin-access-secret-change-me
ADMIN_JWT_REFRESH_SECRET=dev-admin-refresh-secret-change-me

# Used only by `make admin-seed` — there is no self-service admin registration.
ADMIN_SEED_EMAIL=
ADMIN_SEED_PASSWORD=
```

New for `admin-web/.env.local` (see `admin-web/.env.example`):

```
VITE_API_URL=http://localhost:8080/api/v1
```

---

## 6. Frontend structure (`admin-web/`)

Vite + React 19 + TypeScript + Tailwind v4 + React Router v6 + TanStack
Query + Axios + React Hook Form + Zod + recharts.

```
src/
  api/          One thin module per resource (auth, dashboard, users,
                 uploads, analytics, settings, logs, profile) — typed axios calls
  hooks/         TanStack Query hooks per resource + useAuth (context-based
                 session state)
  lib/
    axios.ts     Centralized client: request interceptor attaches the bearer
                 token; response interceptor handles (a) single-flight 401 →
                 refresh → retry-once → redirect-to-login on failure, and
                 (b) backed-off retry (2 attempts) for GET requests that hit
                 a network error or 5xx
    tokenStorage.ts   localStorage helpers, namespaced separately from the
                      customer web app
    queryClient.ts
  components/
    layout/      Sidebar, MobileSidebar (slide-in overlay), Topbar, AppLayout,
                 ProtectedRoute
    table/       DataTable<T> — one generic component reused by Users,
                 Uploads, and Logs: server-side pagination, search, client-
                 side per-page sort, filters slot, bulk-select + bulk action
                 bar, loading/empty/error states
    charts/      TrendChart — recharts line-chart wrapper, reused by both
                 the Dashboard growth charts and the Analytics page
    ui/          Card, StatCard, Badge, Button
    forms/       FormField (RHF error display + label wrapper)
  pages/         Login, Dashboard, Users, UserDetail, Uploads, UploadDetail,
                 Analytics, Settings, Logs, Profile, ChangePassword
  router.tsx     Route tree, /login public, everything else behind
                 ProtectedRoute + AppLayout
  types/         Shared TypeScript interfaces matching the backend's JSON
                 shapes exactly (snake_case)
```

---

## 7. Setup instructions

```bash
# Backend (from repo root)
cp .env.example .env    # already includes ADMIN_JWT_*/ADMIN_SEED_* — edit as needed
make infra-up
make migrate-up
ADMIN_SEED_EMAIL=admin@example.com ADMIN_SEED_PASSWORD='ChangeMe123!' make admin-seed
make api        # terminal 1
make worker     # terminal 2

# Frontend (new terminal)
cd admin-web
cp .env.example .env.local
npm install
make admin-web  # or: npm run dev   (from within admin-web/)
```

Visit `http://localhost:5174/login` and sign in with the seeded admin
credentials.

## 8. Build instructions

```bash
go build ./...                       # backend
cd admin-web && npm run build         # frontend (tsc -b && vite build)
cd admin-web && npm run lint          # oxlint
```

## 9. Deployment notes

- The admin panel is **not** exposed by anything today beyond the same
  `/api/v1` origin as the customer API — no new load-balancer route or
  CORS change is required beyond what already exists (`middleware.CORS()`
  already allows all origins in this dev-grade setup; tighten this the
  same way you'd tighten it for the customer app before any real deploy).
- `admin-web` is a static Vite build (`dist/`) — deploy it behind whatever
  static hosting/CDN you use for `web/`, pointed at `VITE_API_URL` for
  production.
- **Rotate `ADMIN_JWT_ACCESS_SECRET`/`ADMIN_JWT_REFRESH_SECRET`** away from
  the dev defaults before any real deployment, exactly like the existing
  `JWT_ACCESS_SECRET`/`JWT_REFRESH_SECRET` warning in `GETTING_STARTED.md`.
- Run `make admin-seed` once per environment to create the first admin;
  it's idempotent (no-op if the email already exists), so it's safe to
  include in a deploy script.

---

## 10. Known simplifications

- **Sorting** in tables is client-side, over whatever page is currently
  loaded — the backend list endpoints always order by `created_at DESC`.
  Full-dataset sorting would need dynamic `ORDER BY` per sqlc query, which
  was intentionally not added to keep every query static/sqlc-native.
- **File-type analytics** will show entirely "jpg" today, because the
  upload pipeline (`usecase/analysis.Create`) always writes `.jpg`
  regardless of the original content type — the query groups by the real
  stored extension, so this becomes accurate automatically if that ever
  changes.
- **Bulk actions** in the frontend are implemented as parallel calls to the
  existing single-item endpoints (`Promise.allSettled`), not new bulk API
  endpoints.
- **Of the four log types** requested (login history, admin actions, error
  logs, request logs), only **admin action audit log** got real
  infrastructure — the other three don't exist anywhere in this codebase
  today and were out of scope for this pass.
