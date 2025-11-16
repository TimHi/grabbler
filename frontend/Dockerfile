# 1️⃣ Builder stage: install deps and build the app
FROM node:20-alpine AS builder
WORKDIR /app

# Copy dependency files first for caching
COPY package*.json ./
# If using Yarn or pnpm:
# COPY yarn.lock ./
# COPY pnpm-lock.yaml ./

RUN npm ci

# Copy rest of the app
COPY . .

# Build the Next.js app
RUN npm run build

# 2️⃣ Runner stage: use a lightweight Node image
FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3000

# Copy only necessary files from builder
COPY --from=builder /app/package*.json ./
RUN npm ci --omit=dev

COPY --from=builder /app/.next ./.next
COPY --from=builder /public ./public
COPY --from=builder /app/next.config.ts ./next.config.ts

EXPOSE 3000

# Start the Next.js server
CMD ["npm", "run", "start"]
