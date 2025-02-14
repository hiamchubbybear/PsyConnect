# Stage 1: Build
FROM node:18-alpine AS builder

WORKDIR /app

# Chỉ copy package.json trước để tận dụng cache
COPY package.json package-lock.json ./

# Cài dependencies, nhưng không copy toàn bộ source code để tránh cache bị invalid
RUN npm ci --only=production

# Copy toàn bộ code vào container
COPY . .

# Stage 2: Chạy ứng dụng
FROM node:18-alpine

WORKDIR /app

# Chỉ copy node_modules từ stage build trước
COPY --from=builder /app/node_modules ./node_modules
# Copy source code
COPY --from=builder /app .

# Expose port
EXPOSE 8082

# Start the app
CMD ["node", "main.js"]
