# Referral Hub

A small real-time job referral platform. Users post job listings, refer candidates to them, and discuss each referral in its own comment thread. New jobs, referrals, and comments are pushed to every connected user over WebSockets, so nothing needs a refresh.

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

The API runs on port 8090. Database migrations run automatically at startup, so there is no separate migrate step. You should see `api: listening on :8090`.

### 3. Start the Frontend

Open another terminal:

```bash
cd frontend
npm install
npm run dev
```

The frontend runs on port 3001. Open http://localhost:3001 in your browser.

## Health Check

```bash
curl -s http://localhost:8090/healthz
```

Expected response, returned only if the API can also reach PostgreSQL:

```json
{ "status": "ok" }
```

## Authentication and Test Credentials

Authentication is mocked. There is no signup and no password. Enter any email address and display name to sign in. The first login creates the user; later logins with the same email identify the same user.

Suggested demo accounts:

| Email               | Display Name |
| ------------------- | ------------ |
| `alice@example.com` | Alice Chen   |
| `bob@example.com`   | Bob Iyer     |

## Demo Walkthrough

Use **one normal window and one incognito window**, placed side by side. Two ordinary tabs share `localStorage`, so they would both be signed in as the same person.

1. Normal window: sign in as `alice@example.com`.
2. Incognito window: sign in as `bob@example.com`.
3. **A new job appears live.** Alice posts a job from the dashboard. It appears at the top of Bob's feed with no refresh.
4. **A new referral appears live.** Bob opens that job and submits a referral with a candidate name and email. Alice, with the same job open, sees the referral appear.
5. **A new comment appears live.** Both users type into that referral's comment thread. Messages land in the other window as they are sent.
6. **Server-side validation.** Submit the referral form with empty fields. The inputs show per-field messages. These come from the server's 422 response, not from the browser: the form sets `noValidate` deliberately so the backend is what is actually exercised.
7. **Duplicate protection.** Refer the same candidate email to the same job a second time. It is rejected with a 409 and a clear message.

### Demonstrating Live Events Without a Second Browser

`cmd/wsprobe` is a terminal WebSocket client that prints events as they arrive.

```bash
cd backend
TOKEN=$(curl -s -X POST http://localhost:8090/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"bob@example.com","name":"Bob Iyer"}' | python3 -c "import json,sys;print(json.load(sys.stdin)['token'])")
go run ./cmd/wsprobe "$TOKEN" jobs
```

Post a job in the browser and the event prints in the terminal.

## Features

- Create and browse job listings.
- Submit candidate referrals for jobs.
- Discuss each referral in its own comment thread.
- Receive live updates through WebSockets, with automatic reconnection.
- Validate every write on the server and report errors per field.
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
| GET    | `/jobs`                            | List jobs, newest first       |
| GET    | `/jobs/{jobID}`                    | Get a job and its referrals   |
| POST   | `/jobs/{jobID}/referrals`          | Submit a referral             |
| GET    | `/jobs/{jobID}/referrals`          | List referrals for a job      |
| POST   | `/referrals/{referralID}/comments` | Create a comment              |
| GET    | `/referrals/{referralID}/comments` | List comments                 |
| GET    | `/healthz`                         | Check API and database health |

`/healthz` is served at the server root, outside `/api/v1`. Read endpoints are public, so the dashboard renders before sign in. Write endpoints require a valid bearer token.

### Error Responses

Every failure uses the same envelope:

```json
{ "error": "this candidate has already been referred for this job" }
```

Validation failures add a per-field map:

```json
{
  "error": "validation failed",
  "fields": {
    "title": "is required",
    "candidateEmail": "must be a valid email address"
  }
}
```

| Status | Meaning                                           |
| ------ | ------------------------------------------------- |
| 400    | Body is not valid JSON, or has an unknown field   |
| 401    | Missing, malformed, or expired token              |
| 404    | Unknown job or referral                           |
| 409    | Candidate already referred to this job            |
| 422    | Validation failed, with the offending field names |

## WebSocket Events

WebSocket URL:

`ws://localhost:8090/ws?token=<jwt>`

The token is a query parameter because the browser WebSocket API cannot set headers on the handshake.

Clients subscribe to topics rather than receiving everything:

| Topic                   | Carries                                            |
| ----------------------- | -------------------------------------------------- |
| `jobs`                  | New jobs, and referral counts changing on the feed |
| `job:<jobID>`           | Referrals arriving on one job                      |
| `referral:<referralID>` | Comments arriving on one referral                  |

Messages sent by the client:

```json
{ "action": "subscribe", "topic": "jobs" }
```

`unsubscribe` takes the same shape. Messages received from the server:

```json
{
  "type": "job.created",
  "topic": "jobs",
  "payload": { "id": "...", "title": "...", "company": "..." },
  "timestamp": "2026-10-04T14:22:31Z"
}
```

Event types are `job.created`, `referral.created`, and `comment.created`. Each event carries the whole created object, so the client appends it to state directly instead of refetching.

## Data Model

Four normalized tables, with cascading deletes down the ownership chain.

```text
users --< jobs --< referrals --< comments
```

| Table       | Notable columns                                                                |
| ----------- | ------------------------------------------------------------------------------ |
| `users`     | `email` unique                                                                 |
| `jobs`      | `status` checked against open, paused, closed; `posted_by`                     |
| `referrals` | `job_id`, `referrer_id`, candidate name and email, `status` (submitted, reviewed, accepted, rejected) |
| `comments`  | `referral_id`, `author_id`, `body`                                             |

Two constraints do real work rather than living only in Go:

- `UNIQUE (job_id, candidate_email)` on `referrals` prevents the same candidate being referred twice for one role. It surfaces as a 409.
- `CHECK` constraints keep status values valid even if a write bypasses the API.

## Running Tests

Backend tests cover request validation, JWT issue and verify, and WebSocket hub routing. They need no database.

```bash
cd backend
go test ./...
```

Frontend checks, after installing dependencies:

```bash
cd frontend
npm run typecheck
npm run lint
```

## Configuration

Both sides have working defaults, so the app boots with no `.env` file.

Backend environment variables, see `backend/.env.example`:

| Variable       | Default                                                                |
| -------------- | ---------------------------------------------------------------------- |
| `PORT`         | `8090`                                                                 |
| `DATABASE_URL` | `postgres://referhub:referhub@localhost:5434/referhub?sslmode=disable` |
| `JWT_SECRET`   | `dev-only-insecure-secret`                                             |
| `CORS_ORIGIN`  | `http://localhost:3001`                                                |

Frontend environment variables, see `frontend/.env.local.example`:

| Variable              | Default                        |
| --------------------- | ------------------------------ |
| `NEXT_PUBLIC_API_URL` | `http://localhost:8090/api/v1` |
| `NEXT_PUBLIC_WS_URL`  | `ws://localhost:8090/ws`       |

PostgreSQL is published on 5434 and the web app on 3001 to avoid clashing with anything already running on the usual 5432 and 3000. If you change the web port, change `CORS_ORIGIN` to match, or the browser will block every request.

## Project Structure

```text
backend/
  cmd/
    api/              server entrypoint, graceful shutdown
    wsprobe/          terminal WebSocket client for demos
  migrations/         SQL, applied automatically at boot
  internal/
    auth/             JWT issue and verify
    config/           environment settings
    db/               connection pool and migration runner
    handler/          decode, validate, store, publish
    middleware/       bearer token to authenticated user
    models/           API wire types
    respond/          JSON and error response helpers
    router/           the whole URL surface in one file
    store/            all SQL lives here
    ws/               topic hub and socket plumbing

frontend/
  app/                dashboard and job detail pages
  components/         auth and socket providers, forms, cards, thread
  lib/                API client and shared types
```

## Design Notes

**Events are published after the write succeeds.** A client is never told about a row that failed to save.

**The author is not special-cased.** When you post a job, your own tab learns about it from the same broadcast everyone else receives, with no optimistic insert. The live path is therefore exercised on every write.

**One socket per browser tab.** A React context owns the connection and multiplexes topic subscriptions over it, reconnecting and re-subscribing by itself if the connection drops. Components never open their own sockets.

**A slow client cannot stall the others.** The hub does a non-blocking send and drops a frame if a client's buffer is full, rather than blocking the fan-out.

## Assumptions and Simplifications

These were deliberate.

1. **Authentication is mocked.** Email plus a display name returns a normal JWT. No passwords, no email verification, no refresh tokens. Swapping in a real identity provider would touch only `internal/auth` and the login handler, since the other handlers work off a verified user id.
2. **A single backend instance.** The WebSocket hub keeps subscriptions in memory. Running several instances would mean publishing through Redis pub/sub and having each instance fan out to its own sockets. `Hub.Publish` is the only seam that changes.
3. **The socket token is a query parameter.** The browser cannot set headers on a WebSocket handshake, so the token can appear in access logs. In production the fix is a short-lived single-use ticket fetched over REST first.
4. **No authorization model.** Any signed-in user can post a job, refer to any job, and comment on any referral. There are no organizations, owners, or roles.
5. **No pagination.** The job feed is capped at the 100 most recent listings, and referrals and comments load in full.
6. **Referral status is stored but not editable** through the API. The column and its CHECK constraint exist for the obvious next step.
7. **The development JWT secret is committed** as a default. Set `JWT_SECRET` before deploying anywhere real, and never commit production credentials.

## Troubleshooting

### The backend cannot connect to PostgreSQL

Check that the container is running, then read its logs:

```bash
docker compose ps
```

```bash
docker compose logs postgres
```

### The feed says it cannot reach the API

Confirm the API is up:

```bash
curl -i http://localhost:8090/healthz
```

### Everything loads but nothing updates live

Look at the connection indicator in the header. If it shows reconnecting, the API probably restarted, and the socket recovers on its own within a couple of seconds.

### Requests fail with a CORS error

The web app is not on port 3001. Either run it on 3001, or set `CORS_ORIGIN` on the backend to the port you are using.

### Port 5434 is already in use

Change the host-side PostgreSQL port in `docker-compose.yml` and update `DATABASE_URL` to match.

### The database is in a bad state

Reset it and let migrations re-run on the next backend start. This deletes all data:

```bash
docker compose exec -T postgres psql -U referhub -d referhub -c 'DROP SCHEMA public CASCADE; CREATE SCHEMA public;'
```
