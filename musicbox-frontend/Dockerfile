# 构建阶段
FROM node:22-alpine AS builder

WORKDIR /app

RUN npm install -g pnpm

ENV CI=true
ARG VITE_BACKEND_URL=
ENV VITE_BACKEND_URL=${VITE_BACKEND_URL}

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile

COPY index.html tsconfig.app.json tsconfig.json tsconfig.node.json uno.config.ts vite.config.ts ./
COPY public ./public
COPY src ./src

RUN pnpm build


# 运行阶段
FROM nginx:alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
