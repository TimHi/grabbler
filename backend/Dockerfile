# --- Build stage ---
FROM golang:1.22-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# --- Runtime stage ---
FROM debian:bookworm-slim
# yt-dlp + ffmpeg are required for many formats/merges
RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates curl python3 python3-pip ffmpeg \
    && pip3 install --no-cache-dir yt-dlp \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/server /app/server

# Where we'll write downloads INSIDE the container (you will mount over this)
ENV AUDIO_OUTPUT_DIR=/downloads

# Create the directory and make it writable
RUN mkdir -p /downloads && chmod 775 /downloads

EXPOSE 3333
CMD ["/app/server"]
