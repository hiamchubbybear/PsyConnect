FROM node:18-alpine AS builder
WORKDIR /app
COPY notificationservice/package*.json ./
RUN npm ci --only=production
COPY notificationservice .
FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app .
EXPOSE 8082
CMD ["node", "main.js"]
