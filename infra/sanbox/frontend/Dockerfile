# Stage 1: Build
FROM node:20.17.0-alpine AS builder

WORKDIR /app

# Copy dependencies
COPY package.json yarn.lock ./
RUN yarn --frozen-lockfile
ARG NEXT_PUBLIC_API_BASE_URL=http://localhost:3010/api/v1
ENV NEXT_PUBLIC_API_BASE_URL=${NEXT_PUBLIC_API_BASE_URL}
ARG ENV_TYPE=production
ENV ENV_TYPE=${ENV_TYPE}
COPY . .
RUN yarn build

# Stage 2: Serve with nginx
FROM nginx:stable-alpine AS runner

RUN rm -rf /usr/share/nginx/html/*

COPY --from=builder /app/out /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
