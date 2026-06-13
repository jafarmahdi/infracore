# InfraCore

InfraCore is an enterprise infrastructure operations platform that brings data
center management, network inventory, monitoring, assets, alerts, agents,
licenses, and contracts into one secure control plane.

It is designed for NOC, IT operations, infrastructure, and data center teams
that need a unified view of physical, virtual, and network resources.

## Platform Modules

### Authentication and Access Control

- JWT access and refresh token authentication
- Protected API and frontend routes
- Multi-tenant user isolation
- Role-based access control and granular permissions
- Account lockout, audit logs, and secure password hashing

### Dashboard

- Total asset and device health summaries
- Online and offline device counts
- Active alerts and license expiration indicators
- SLA and availability widgets
- CPU, memory, disk, and network traffic charts
- Recent infrastructure events

### DCIM

Data Center Infrastructure Management models the physical environment:

- Data centers, rooms, and racks
- Rack units and equipment placement
- PDUs and power feeds
- Patch panels, cable paths, and terminations
- Capacity, power, and environmental information

### Asset Management

Provides a unified inventory for:

- Servers, switches, routers, and firewalls
- UPS units, printers, access points, and storage
- Virtual machines and hypervisor relationships
- Manufacturers and device types
- Network interfaces, warranties, lifecycle, and ownership

### IPAM

IP Address Management covers:

- VRFs and routing domains
- VLANs
- IPv4 and IPv6 prefixes
- IP address allocation
- DNS zones and records
- DHCP leases

### Monitoring

Monitoring supports infrastructure and service health:

- Monitored hosts and reusable monitoring profiles
- Agent, SNMP, WMI, ICMP, SSH, and API checks
- CPU, memory, disk, network, process, and service metrics
- Availability records and uptime calculations
- Time-series metric storage with TimescaleDB

### Alerts

- Active alerts and alert history
- Severity and status tracking
- Alert acknowledgement and assignment
- Notification delivery records
- Escalation rules and maintenance silences

### Agents

- Agent registration and authentication
- Agent groups and registration tokens
- Heartbeat and health tracking
- Version and compliance information
- Remote task delivery

### Licenses and Contracts

- Software licenses, vendors, seats, and assignments
- Expiration and renewal tracking
- Support and maintenance contracts
- Warranty records and covered assets

### Reports and Administration

- Infrastructure and availability reporting
- License and contract reports
- User, role, and permission management
- Companies, branches, departments, and sites
- Notification channels and API token administration

## Architecture

```text
React + TypeScript Frontend
            |
       REST API / JWT
            |
        Go + Gin API
            |
    Application Services
            |
       Domain Models
            |
 PostgreSQL/TimescaleDB + Redis
```

### Backend Layers

- **Domain:** Core entities, value objects, repository contracts, and business
  rules.
- **Application:** Use cases such as authentication and monitored-host
  management.
- **Infrastructure:** PostgreSQL repositories, Redis integration, logging, and
  cryptography.
- **Interfaces:** Gin handlers, middleware, routing, and HTTP error responses.

### Frontend Layers

- Route-level pages for each InfraCore module
- Reusable layout, dashboard, form, and UI components
- Typed Axios API services
- TanStack Query for server state
- Zustand for authentication and interface preferences
- React Hook Form and Zod for validated forms
- Recharts for operational visualizations

## Technology Stack

### Backend

- Go
- Gin
- PostgreSQL and TimescaleDB
- Redis
- JWT
- sqlx
- Zap

### Frontend

- React and TypeScript
- Vite
- Tailwind CSS and shadcn/ui patterns
- React Router
- TanStack Query
- Zustand
- Axios
- Recharts
- React Hook Form and Zod

## Quick Start

The easiest way to start the complete stack is Docker Compose:

```bash
docker compose up --build
```

Open:

- Frontend: http://localhost
- API health: http://localhost:8080/health

The migration container initializes the database before the API starts.

### Development Login

```text
Organization: demo
Email: admin@demo.com
Password: Admin@123456
```

These credentials are for local development only. Replace all database, JWT,
encryption, and seeded credentials before a production deployment.

## Local Development

Start PostgreSQL/TimescaleDB and Redis, then apply the migrations under
`migrations/postgres`.

Run the API:

```bash
go run ./cmd/api
```

Run the frontend in another terminal:

```bash
cd web
cp .env.example .env
npm install
npm run dev
```

Local services:

- Frontend: http://localhost:5173
- API: http://localhost:8080
- PostgreSQL: localhost:5432
- Redis: localhost:6379

## Repository Structure

```text
cmd/api/                     API entry point
configs/                     Application configuration
internal/application/        Application services and use cases
internal/domain/             Domain models and repository contracts
internal/infrastructure/     PostgreSQL and Redis implementations
internal/interfaces/http/    HTTP handlers, middleware, and routing
migrations/postgres/         Database migrations and development seed
pkg/                         Shared configuration, crypto, errors, and logging
scripts/                     Local development helpers
web/                         React frontend
```

## Current Implementation Status

Available foundations include:

- Real database-backed authentication
- JWT middleware and protected routes
- Tenant-aware identity repositories
- Monitoring host APIs and repositories
- Full PostgreSQL/TimescaleDB schema
- Enterprise frontend shell and dashboard
- Module pages for DCIM, assets, IPAM, monitoring, alerts, agents, licenses,
  contracts, reports, security, and settings
- Dark mode, responsive navigation, and RTL-ready UI structure

Some modules currently provide complete domain models and frontend workflows
while their remaining CRUD API handlers are still under development.

## Validation

```bash
go test ./...
go vet ./...

cd web
npm run build
npm run lint
```
