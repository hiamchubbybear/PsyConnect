# Stage 1: Build the Angular application
FROM node:22-alpine AS build

WORKDIR /app

# Copy package files and install dependencies
COPY package.json package-lock.json ./
RUN npm ci

# Copy the rest of the application code
COPY . .

# ===== Build Args: Key bí mật được truyền vào lúc build =====
ARG CLOUDINARY_API_KEY=PLACEHOLDER_API_KEY
ARG CLOUDINARY_API_SECRET=PLACEHOLDER_API_SECRET

# Thay thế placeholder trong environment.ts bằng giá trị thật
RUN sed -i "s|PLACEHOLDER_API_KEY|${CLOUDINARY_API_KEY}|g" src/environments/environment.ts && \
    sed -i "s|PLACEHOLDER_API_SECRET|${CLOUDINARY_API_SECRET}|g" src/environments/environment.ts

# Build the application for production
RUN npm run build -- --configuration production

# Stage 2: Serve the application with Nginx
FROM nginx:alpine

# Remove default nginx website
RUN rm -rf /usr/share/nginx/html/*

# Copy built application from build stage
# Note: The output path depends on angular.json, usually dist/psyconnect-web/browser or dist/psyconnect-web
COPY --from=build /app/dist/psyconnect-web/browser /usr/share/nginx/html

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Exposure port 80 (Container Port)
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
