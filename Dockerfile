# ── Stage 1: Build the JS bundle ─────────────────────────────────────────────
FROM node:22-alpine AS js-builder

WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm ci

COPY tsconfig.json vite.config.ts ./
COPY src/ src/
RUN npm run build
# Output: dist/plugin.js

# ── Stage 2: Build the Go gRPC server ────────────────────────────────────────
FROM golang:1.23-alpine AS go-builder

WORKDIR /app/server
COPY server/go.mod server/go.sum* ./
RUN go mod download

COPY server/ .
RUN CGO_ENABLED=0 go build -o /plugin .

# ── Stage 3: Minimal runtime image ───────────────────────────────────────────
FROM alpine:3.20

COPY --from=go-builder /plugin /plugin
COPY --from=js-builder /app/dist/plugin.js /ui/plugin.js

ENV PLUGIN_PORT=50051
EXPOSE 50051

ENTRYPOINT ["/plugin"]
