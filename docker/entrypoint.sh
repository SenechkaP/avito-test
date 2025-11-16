#!/bin/sh
set -e

make migrate-up

echo "Migrations applied. Starting app..."
exec ./pr_service