.PHONY: dev build start migrate clean

# Development
dev:
	ENV=dev bun run --hot cmd/server/main.ts 

# Build
build:
	bun build cmd/server/main.ts --outdir=dist --target=bun

# Production start
start:
	ENV=prd bun run dist/main.js
# Run Drizzle migrations
migrate:
	bun run drizzle-kit generate
	bun run drizzle-kit push

# Generate Drizzle migrations
migrate-gen:
	bun run drizzle-kit generate

# Push Drizzle migrations
migrate-push:
	bun run drizzle-kit push

# Clean build artifacts
clean:
	rm -rf dist node_modules

# Install dependencies
install:
	bun install

# Type check
typecheck:
	bun run tsc --noEmit

# Lint
lint:
	bun run eslint . --ext .ts

# Format
format:
	bun run prettier --write "**/*.ts"
