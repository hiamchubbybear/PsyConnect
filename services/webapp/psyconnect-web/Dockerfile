FROM node:22-alpine AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
ARG CLOUDINARY_API_KEY=PLACEHOLDER_API_KEY
RUN sed -i "s|PLACEHOLDER_API_KEY|${CLOUDINARY_API_KEY}|g" src/environments/environment.ts
RUN npm run build -- --configuration production
FROM nginx:alpine
RUN rm -rf /usr/share/nginx/html/*
COPY --from=build /app/dist/psyconnect-web/browser /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
