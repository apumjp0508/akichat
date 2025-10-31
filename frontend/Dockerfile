# 開発用 Dockerfile
FROM node:20-alpine

WORKDIR /app

# 依存ファイルをコピー
COPY ./package.json ./pnpm-lock.yaml* ./
RUN corepack enable && pnpm install

# ソースをコピー
COPY . .

EXPOSE 3000

# 開発サーバーを起動（ホットリロード）
CMD ["pnpm", "dev", "-H", "0.0.0.0", "-p", "3000"]

