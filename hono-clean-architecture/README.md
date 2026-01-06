# SkillSync - Hono Clean Architecture

A TypeScript backend service built with [Hono](https://hono.dev/) framework, [Drizzle ORM](https://orm.drizzle.team/), and PostgreSQL following Clean Architecture principles.

## Tech Stack

- **Framework**: Hono (TypeScript)
- **ORM**: Drizzle ORM
- **Database**: PostgreSQL
- **Cache**: Redis
- **Runtime**: Bun
- **Authentication**: JWT (EC/RSA/EdDSA keys supported)

## Project Structure

```
hono-clean-architecture/
├── cmd/
│   └── server/
│       ├── main.ts          # Application entry point
│       └── Dockerfile        # Docker configuration
├── configs/
│   ├── dev.yaml             # Development configuration
│   ├── uat.yaml             # UAT configuration
│   └── prd.yaml             # Production configuration
├── internal/
│   ├── assets/              # Static assets
│   ├── datasources/         # Database, Redis, JWT connections
│   ├── handlers/            # HTTP middleware and route handlers
│   └── infrastructure/      # Server, router, error handling
├── pkg/
│   ├── auth/                # Authentication module
│   ├── config/              # Configuration loader
│   ├── domain/              # Domain interfaces
│   ├── logs/                # Logging utilities
│   ├── models/              # Database schemas and types
│   ├── permission/          # Permission management
│   ├── role/                # Role management
│   ├── role_permission/     # Role-permission associations
│   ├── tx_context/          # Transaction management
│   ├── user/                # User management
│   ├── user_role/           # User-role associations
│   └── utils/               # Utility functions
├── uploads/                 # File uploads directory
├── package.json
├── tsconfig.json
├── drizzle.config.ts
└── Makefile
```

## Prerequisites

- [Bun](https://bun.sh/) >= 1.0.0
- PostgreSQL >= 14
- Redis >= 6

## Setup

### 1. Install Dependencies

```bash
bun install
```

### 2. Configure Environment

Copy and edit the appropriate config file in `configs/`:

```yaml
# configs/dev.yaml
api:
  host: localhost
  port: 8080

db:
  postgres:
    host: localhost
    port: 15432
    username: postgres
    password: your_password
    db_name: skillsync
    conn:
      min: 2
      max: 10
  redis:
    host: localhost
    port: 16379
    username: ""
    password: ""
    db_name: 0

jwt:
  private_key_path: "./internal/assets/dev/private_key.pem"
```

### 3. Generate JWT Keys

Generate EC keys (recommended):

```bash
# Generate private key
openssl ecparam -genkey -name prime256v1 -noout -out internal/assets/dev/private_key.pem

# Extract public key (optional, derived automatically)
openssl ec -in internal/assets/dev/private_key.pem -pubout -out internal/assets/dev/public_key.pem
```

### 4. Run Database Migrations

```bash
make migrate
```

### 5. Start Development Server

```bash
make dev
```

Or directly:

```bash
ENV=dev bun run --hot cmd/server/main.ts
```

## API Endpoints

### Health

- `GET /` - API info
- `GET /health` - Health check

### Authentication

- `POST /auth/register` - Register new user
- `POST /auth/login` - Login
- `POST /auth/refresh` - Refresh token
- `GET /auth/me` - Get current user

### Users

- `GET /users` - List users
- `GET /users/:id` - Get user
- `POST /users` - Create user
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user

### Roles

- `GET /roles` - List roles
- `GET /roles/:id` - Get role
- `POST /roles` - Create role
- `PUT /roles/:id` - Update role
- `DELETE /roles/:id` - Delete role

### Permissions

- `GET /permissions` - List permissions
- `GET /permissions/:id` - Get permission
- `POST /permissions` - Create permission
- `PUT /permissions/:id` - Update permission
- `DELETE /permissions/:id` - Delete permission

## Docker

### Build

```bash
docker build -t skillsync-api -f cmd/server/Dockerfile .
```

### Run

```bash
docker run -p 8080:8080 -e ENV=prd skillsync-api
```

## Makefile Commands

| Command             | Description                              |
| ------------------- | ---------------------------------------- |
| `make dev`          | Start development server with hot reload |
| `make build`        | Build for production                     |
| `make start`        | Start production server                  |
| `make migrate`      | Run database migrations                  |
| `make migrate-gen`  | Generate migration files                 |
| `make migrate-push` | Push migrations to database              |
| `make typecheck`    | Run TypeScript type checking             |
| `make lint`         | Run ESLint                               |
| `make format`       | Format code with Prettier                |
| `make clean`        | Remove build artifacts                   |
| `make install`      | Install dependencies                     |

## License

MIT
