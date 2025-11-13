# --- Build stage ---
FROM golang:1.25-alpine AS builder
WORKDIR /app

# If your go.mod pulls from git repos, you may need git:
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# --- Runtime stage ---
FROM alpine:3.20

RUN apk add --no-cache \
      ca-certificates \
      ffmpeg \
      curl \
    && curl -L "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_musl" \
         -o /usr/local/bin/yt-dlp \
    && chmod +x /usr/local/bin/yt-dlp

WORKDIR /app
COPY --from=builder /app/server /app/server

ENV AUDIO_OUTPUT_DIR=/downloads
RUN mkdir -p /downloads && chmod 775 /downloads

EXPOSE 3333
CMD ["/app/server"]
