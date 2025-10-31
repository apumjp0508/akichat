# 開発用 Dockerfile（ホットリロード対応）
FROM golang:1.25 AS dev

WORKDIR /app

# 依存関係だけ先にコピーしてキャッシュ効率UP
COPY ./go.mod ./go.sum ./
RUN go mod download

# ソースをコピー
COPY . .

# air（ホットリロードツール）を導入
RUN go install github.com/air-verse/air@latest

EXPOSE 8080

# ホットリロード実行
CMD ["air", "-c", ".air.toml"]
