# --- Build stage ---
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

ARG TARGETOS
ARG TARGETARCH

COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o server .

# --- Runtime stage ---
FROM alpine:3.20

ARG TARGETARCH

RUN apk add --no-cache \
      ca-certificates \
      ffmpeg \
      curl

# Use MUSL builds of yt-dlp (important for Alpine!)
RUN set -eux; \
    case "$TARGETARCH" in \
      "amd64")  YT_BIN="yt-dlp_musllinux" ;; \
      "arm64")  YT_BIN="yt-dlp_musllinux_aarch64" ;; \
      *) echo "Unsupported arch: $TARGETARCH"; exit 1 ;; \
    esac; \
    curl -L "https://github.com/yt-dlp/yt-dlp/releases/latest/download/${YT_BIN}" \
      -o /usr/local/bin/yt-dlp; \
    chmod +x /usr/local/bin/yt-dlp

WORKDIR /app
COPY --from=builder /app/server /app/server

ENV AUDIO_OUTPUT_DIR=/downloads
RUN mkdir -p /downloads && chmod 775 /downloads

EXPOSE 3333
CMD ["/app/server"]
