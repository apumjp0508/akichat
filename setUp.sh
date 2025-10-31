#!/bin/bash
set -e
# スクリプトの位置からプロジェクトルートを解決
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$SCRIPT_DIR"

echo "SCRIPT_DIR = $SCRIPT_DIR"
echo "ROOT_DIR   = $ROOT_DIR"

ENV_FILE="$ROOT_DIR/environment/dev/.env"
COMPOSE_FILE="$ROOT_DIR/environment/dev/docker-compose.yml"

echo "📄 環境ファイルを確認中..."
if [ ! -f "$ENV_FILE" ]; then
  echo "❌ $ENV_FILE が見つかりません。"
  exit 1
fi

echo "🧹 既存の Docker 環境を削除中..."
docker compose -f "$COMPOSE_FILE" down -v || true

echo "🔨 ビルド & 起動..."
docker compose -f "$COMPOSE_FILE" up --build -d

echo "✅ 起動完了"
echo "Frontend: http://localhost:3000"
echo "Backend:  http://localhost:8080/health"
