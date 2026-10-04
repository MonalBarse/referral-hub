# Referral Hub

A small real-time job referral platform. Users post job listings, refer candidates to them, and discuss each referral in its own comment thread. New jobs, referrals, and comments are pushed to every connected user over WebSockets.

## Tech Stack

- **Backend:** Go, chi, pgx, gorilla/websocket
- **Frontend:** Next.js (App Router), React, TailwindCSS, shadcn/ui
- **Database:** PostgreSQL

## Prerequisites

- Go 1.26
- Node.js 20 or newer
- Docker and Docker Compose
- Git

## Quick Start

Run each service in a separate terminal from the repository root.

### 1. Start PostgreSQL

```bash
docker compose up -d
```

### 2. Start the Backend

```bash
cd backend
go run ./cmd/api
```

The API runs on port 8090. Database migrations run automatically at startup.

### 3. Start the Frontend

Open another terminal:

```bash
cd frontend
npm install
npm run dev
```

The frontend runs on port 3001.

Open http://localhost:3001 in your browser.

## Health Check

```bash
curl -s http://localhost:8090/healthz
```

Expected response:

```json
{ "status": "ok" }
```

## Authentication

Authentication is mocked for development. Enter any email address and display name to sign in. The first login creates the user; subsequent logins with the same email identify the same user.

Suggested demo accounts:

| Email                                         | Display Name |
| --------------------------------------------- | ------------ |
| [alice@example.com](mailto:alice@example.com) | Alice Chen   |
| [bob@example.com](mailto:bob@example.com)     | Bob Iyer     |

Use a normal browser window for Alice and an incognito window for Bob to demonstrate real-time updates between two users.

## Features

- Create and browse job listings.
- Submit candidate referrals for jobs.
- Add comments to referral threads.
- Receive live updates through WebSockets.
- Validate input on the server.
- Prevent duplicate referrals for the same job and candidate email.
- Track job and referral statuses.
- Check API and database health.

## API Endpoints

Base URL: `http://localhost:8090/api/v1`

| Method | Endpoint                           | Description                   |
| ------ | ---------------------------------- | ----------------------------- |
| POST   | `/auth/login`                      | Sign in or create a user      |
| GET    | `/auth/me`                         | Get the current user          |
| POST   | `/jobs`                            | Create a job                  |
| GET    | `/jobs`                            | List jobs                     |
| GET    | `/jobs/{jobID}`                    | Get a job and its referrals   |
| POST   | `/jobs/{jobID}/referrals`          | Submit a referral             |
| GET    | `/jobs/{jobID}/referrals`          | List referrals for a job      |
| POST   | `/referrals/{referralID}/comments` | Create a comment              |
| GET    | `/referrals/{referralID}/comments` | List comments                 |
| GET    | `/healthz`                         | Check API and database health |

Read endpoints are public. Write endpoints require a valid bearer token.

## WebSocket Events

WebSocket URL:

`ws://localhost:8090/ws?token=<jwt>`

Supported topics:

- `jobs` — new jobs and dashboard referral-count updates.
- `job:<jobID>` — referrals for a specific job.
- `referral:<referralID>` — comments for a specific referral.

Supported event types:

- `job.created`
- `referral.created`
- `comment.created`

## Running Tests

Backend tests:

```bash
cd backend
go test ./...
```

Frontend checks:

```bash
cd frontend
npm run typecheck
npm run lint
```

Run these commands only after installing the frontend dependencies.

## Configuration

Backend environment variables:

| Variable       | Default                                                                |
| -------------- | ---------------------------------------------------------------------- |
| `PORT`         | `8090`                                                                 |
| `DATABASE_URL` | `postgres://referhub:referhub@localhost:5434/referhub?sslmode=disable` |
| `JWT_SECRET`   | `dev-only-insecure-secret`                                             |
| `CORS_ORIGIN`  | `http://localhost:3001`                                                |

Frontend environment variables:

| Variable              | Default                        |
| --------------------- | ------------------------------ |
| `NEXT_PUBLIC_API_URL` | `http://localhost:8090/api/v1` |
| `NEXT_PUBLIC_WS_URL`  | `ws://localhost:8090/ws`       |

These defaults must match your actual Docker Compose and application configuration.

## Project Structure

```text
backend/
  cmd/
    api/
    wsprobe/
  migrations/
  internal/
    auth/
    config/
    db/
    handler/
    middleware/
    models/
    router/
    store/
    ws/

frontend/
  app/
  components/
  lib/
```

The exact structure may vary depending on the files implemented in your repository.

## Important Notes

- Authentication is mocked and is not suitable for production.
- The WebSocket hub stores subscriptions in memory, so this setup assumes a single backend instance.
- WebSocket authentication tokens are passed in the query string.
- Authorization rules and pagination may need to be expanded for production use.
- Change the development JWT secret before deploying.
- Never commit production credentials or secrets.

## Troubleshooting

### Backend cannot connect to PostgreSQL

Check that PostgreSQL is running:

```bash
docker compose ps
```

Inspect its logs:

```bash
docker compose logs postgres
```

### API is not responding

```bash
curl -i http://localhost:8090/healthz
```

### Frontend cannot reach the backend

Check the API URL and CORS origin. The default frontend port is 3001 and the backend port is 8090.

### Port 5434 is already in use

Check which service is using the port, or change the host-side PostgreSQL port in `docker-compose.yml` and update `DATABASE_URL`.

```

```
