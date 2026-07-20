# Production image for the Go backend. Builds three static binaries — api,
# worker, and admin-seed — into one small runtime image. Which one runs is
# chosen at deploy time via the container's start command (CMD below runs
# the API; Render's "Docker Command" override picks worker/admin-seed for
# those services), so one image serves every Go process this project runs.
#
# This is distinct from Dockerfile.dev, which installs the full toolchain
# plus the `air` hot-reload watcher for local development and is never meant
# to be deployed.

# ---- build stage -----------------------------------------------------------
FROM golang:1.26-alpine AS builder

WORKDIR /src

# Cached separately from the source copy so `go mod download` only reruns
# when go.mod/go.sum actually change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 produces a static binary (no libc dependency), which is what
# lets the runtime stage below be a minimal Alpine image with no Go/gcc
# toolchain. -trimpath/-ldflags strip local paths and debug symbols to keep
# the binaries smaller.
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api && \
    go build -trimpath -ldflags="-s -w" -o /out/worker ./cmd/worker && \
    go build -trimpath -ldflags="-s -w" -o /out/admin-seed ./cmd/admin-seed

# ---- runtime stage ----------------------------------------------------------
FROM alpine:3.20

# ca-certificates is required for outbound HTTPS calls this app makes
# (Gemini, YouTube Data API, Razorpay/Stripe, MailerSend, S3). tzdata lets
# time.LoadLocation work if ever needed; both are tiny.
RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -u 10001 appuser

WORKDIR /app
COPY --from=builder /out/api /out/worker /out/admin-seed ./

USER appuser

# Only the API serves HTTP; the worker and admin-seed binaries don't listen
# on a port. Render still only routes traffic for whichever service is
# configured to accept it.
EXPOSE 8080

# Default process is the API server. Override via Render's "Docker Command"
# field (or `docker run <image> ./worker` / `./admin-seed` locally) to run
# the worker or the one-off admin bootstrap instead.
CMD ["./api"]
