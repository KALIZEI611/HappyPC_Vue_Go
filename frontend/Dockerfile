# Этап сборки статики
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# Этап раздачи статики через nginx
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

# Если нужен кастомный конфиг nginx (например, для прокси), можно добавить:
# COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]