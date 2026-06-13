# InfraCore

InfraCore is an enterprise infrastructure operations platform combining DCIM,
asset management, IPAM, monitoring, alerting, agent management, licensing, and
contracts in one control plane.

## Stack

- Go, Gin, PostgreSQL/TimescaleDB, Redis, JWT
- React, TypeScript, Vite, Tailwind CSS
- TanStack Query, Zustand, Axios, Recharts
- React Hook Form and Zod

## Quick Start

The easiest way to run the full stack is Docker Compose:

```bash
docker compose up --build
```

Open:

- Frontend: http://localhost
- API health: http://localhost:8080/health

Seeded development login:

```text
Organization: demo
Email: admin@demo.com
Password: Admin@123456
```

Change all database, JWT, encryption, and seeded credentials before any
production deployment.

## Local Development

Start PostgreSQL/TimescaleDB and Redis, apply the migrations under
`migrations/postgres`, then run:

```bash
go run ./cmd/api
```

In another terminal:

```bash
cd web
cp .env.example .env
npm install
npm run dev
```

The frontend is available at http://localhost:5173 and the API at
http://localhost:8080.

## Current Status

Authentication, JWT middleware, tenant-aware identity repositories, database
migrations, the enterprise application shell, and dashboard views are present.
Several operational modules currently provide UI workflows or domain contracts
while their API handlers are still under development.
