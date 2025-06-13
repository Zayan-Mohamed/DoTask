# Multi-stage build for SvelteKit frontend
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm ci

# Copy source code
COPY . .

# Create production environment file
RUN echo 'PUBLIC_GRAPHQL_ENDPOINT=http://localhost/api/query' > .env && \
    echo 'PUBLIC_WS_ENDPOINT=ws://localhost/api/query' >> .env && \
    echo 'PUBLIC_APP_TITLE=DoTask' >> .env && \
    echo 'PUBLIC_APP_VERSION=1.0.0' >> .env

# Build the application
RUN npm run build

# Production stage with NGINX
FROM nginx:alpine

# Copy custom nginx configuration (use production optimized version)
COPY nginx.prod.conf /etc/nginx/nginx.conf

# Create cache directory
RUN mkdir -p /var/cache/nginx/api && \
    chown -R nginx:nginx /var/cache/nginx

# Copy built application
COPY --from=builder /app/build /usr/share/nginx/html

# Create static directory and copy if exists
RUN mkdir -p /usr/share/nginx/html/static

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
