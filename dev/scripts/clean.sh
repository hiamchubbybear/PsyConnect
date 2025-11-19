#!/bin/bash
set -e

echo "🧹 Cleaning all development data..."
echo "⚠️  This will remove all volumes and data!"
read -p "Are you sure? [y/N]: " -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
    docker-compose down -v
    echo "✅ All data cleaned!"
else
    echo "❌ Cancelled"
fi
