# InfraCore Frontend

Enterprise NOC/DCIM frontend built with React, TypeScript, Vite, Tailwind CSS, TanStack Query, Zustand, Axios, Recharts, React Hook Form, and Zod.

## Run locally

```bash
cp .env.example .env
npm install
npm run dev
```

Authentication uses the local InfraCore API. Dashboard data remains mocked until its backend endpoint is implemented.

## API contract

- `POST /auth/login` returns a JWT access token and sets an httpOnly refresh cookie
- `GET /dashboard` returns the typed shape in `src/types/dashboard.ts`
- Requests automatically include `Authorization: Bearer <accessToken>`
- HTTP 401 responses clear the session and redirect to login

## Structure

```text
src/
  app/             providers and router
  components/      reusable UI, auth, layout, and dashboard components
  config/          navigation and application configuration
  lib/             shared utilities
  pages/           route-level screens
  services/api/    Axios client and endpoint services
  stores/          persisted Zustand state
  styles/          global theme and Tailwind styles
  types/           domain interfaces
```

Use tenant `demo`, email `admin@demo.com`, and password `Admin@123456` after applying the seed migration.
