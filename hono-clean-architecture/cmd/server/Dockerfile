FROM oven/bun:latest

WORKDIR /app

# Copy package files
COPY package.json bun.lock* ./

# Install dependencies
RUN bun install --frozen-lockfile

# Copy source code
COPY . .

# Set environment
ENV ENV=prd

# Expose port
EXPOSE 8080

# Start the server
CMD ["bun", "run", "cmd/server/main.ts"]
