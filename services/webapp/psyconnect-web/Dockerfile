# Stage 1: Build the Angular application
FROM node:22-alpine AS build

WORKDIR /app

# Copy package files and install dependencies
COPY package.json package-lock.json ./
RUN npm ci

# Copy the rest of the application code
COPY . .

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
