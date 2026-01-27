# Build stage
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install --legacy-peer-deps
COPY . .
ARG VITE_API_BASE_URL=http://localhost:8080/api/v1
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL
RUN npx vite build
# 复制vditor资源到dist目录（匹配本地开发环境结构）
RUN mkdir -p /app/dist/vditor/dist && \
    cp /app/node_modules/vditor/dist/index.* /app/dist/vditor/ && \
    cp -r /app/node_modules/vditor/dist/js /app/node_modules/vditor/dist/css /app/node_modules/vditor/dist/images /app/dist/vditor/dist/

# Run stage
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
