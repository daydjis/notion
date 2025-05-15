#!/bin/sh

echo "🔄 Waiting for Postgres to be ready..."

until migrate -path ./migrations -database "$DATABASE_URL" up; do
  echo "⏳ Database not ready, waiting..."
  sleep 2
done

echo "✅ Migrations applied. Starting server..."
exec ./server
