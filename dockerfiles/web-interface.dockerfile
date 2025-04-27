FROM node:20-alpine AS builder

ARG VITE_DOMAIN
ENV VITE_DOMAIN=${VITE_DOMAIN}

WORKDIR /app

COPY package.json /app
# COPY package-lock.json /app

RUN npm install
RUN npm rebuild esbuild

COPY . /app
RUN npm run build

FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html
COPY ./nginx/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]