# 開発用 Dockerfile
FROM node:20-alpine

WORKDIR /app

# 依存関係をインストール
COPY ./package.json ./pnpm-lock.yaml* ./
RUN corepack enable && pnpm install

# Tailwind設定ファイルも含めてソースをコピー
COPY . .

# Tailwind関連ファイルを確実に認識させる
RUN ls -la && cat tailwind.config.js || echo "⚠️ tailwind.config.js not found"

EXPOSE 3000

CMD ["pnpm", "dev", "-H", "0.0.0.0", "-p", "3000"]
