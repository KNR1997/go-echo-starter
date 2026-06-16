#!/bin/sh

# Build connection string
DB_STRING="${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}?parseTime=true"

echo "Running migrations on ${DB_HOST}:${DB_PORT}/${DB_NAME}"
goose -dir /migrations mysql "$DB_STRING" up

# Exit with goose's exit code
exit $?