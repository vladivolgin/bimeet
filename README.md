# Bimeet — Technical Documentation

Bimeet is an event planning and coordination platform. It lets organizers create events and manage participants, money collections, polls, item assignments, carpools, and links — all in one place.

The platform consists of three independently deployed applications:

| Application | Role | Domain |
|---|---|---|
| **Backend** | REST API server | (internal, no public domain) |
| **Frontend** | React SPA | `app.bimeet.pro` |
| **Landing** | Marketing site | `bimeet.pro`, `www.bimeet.pro` |

---

## Cloud Infrastructure

```mermaid
graph TD
    Internet((Internet))

    subgraph VPS["VPS — Docker Bridge Network (app-net)"]
        Traefik["Traefik v3.3\n(reverse proxy + TLS)\n:80 → :443 redirect\nLet's Encrypt ACME"]

        subgraph Apps["Application Containers"]
            Landing["Landing\nnginx:alpine\nbimeet.pro / www.bimeet.pro"]
            Frontend["Frontend\nnginx:alpine\napp.bimeet.pro"]
            Backend["Backend\nalpine:3.20\nGo binary :8080"]
        end

        subgraph Data["Data Containers"]
            Postgres["PostgreSQL 16\nalpine\nbimeet DB"]
            MinIO["MinIO\nS3-compatible\nobject storage"]
        end

        subgraph Volumes["Named Volumes"]
            PGVol[("postgres_data")]
            MinioVol[("minio_data")]
            AcmeVol[("acme_data")]
        end
    end

    GHCR["GHCR\nghcr.io\nDocker image registry"]

    Internet --> Traefik
    Traefik -->|"Host: bimeet.pro\nwww.bimeet.pro"| Landing
    Traefik -->|"Host: app.bimeet.pro"| Frontend
    Traefik -->|"Host: s3.bimeet.pro"| MinIO
    Frontend -->|"proxy /api/*"| Backend
    Backend --> Postgres
    Backend --> MinIO
    Postgres --- PGVol
    MinIO --- MinioVol
    Traefik --- AcmeVol
    GHCR -.->|"docker pull on deploy"| Backend
    GHCR -.->|"docker pull on deploy"| Frontend
    GHCR -.->|"docker pull on deploy"| Landing
```

### Infrastructure Components

| Component | Image | Purpose |
|---|---|---|
| Traefik | `traefik:v3.3` | Reverse proxy, TLS termination, Let's Encrypt certificates via HTTP-01 challenge |
| PostgreSQL | `postgres:16-alpine` | Primary relational database |
| MinIO | `minio/minio:latest` | S3-compatible object storage for user avatars |
| GHCR | GitHub Container Registry | Stores built Docker images (`bimeet-backend`, `bimeet-frontend`, `bimeet-landing`) |

All services run on a single Docker bridge network `app-net`. Three named volumes provide persistent storage: `postgres_data`, `minio_data`, and `acme_data` (TLS certificates).

---

## Deployment Pipeline

```mermaid
flowchart TD
    Push["git push → master"] --> GHA["GitHub Actions\ndeploy.yml\nubuntu-latest runner"]

    GHA --> Checkout["actions/checkout@v4"]
    Checkout --> Login["docker/login-action@v3\nLogin to GHCR\nusing GITHUB_TOKEN"]

    Login --> BuildBackend["docker/build-push-action@v5\nBuild backend\ncontext: ./backend\nghcr.io/OWNER/bimeet-backend:latest"]
    Login --> BuildFrontend["docker/build-push-action@v5\nBuild frontend\ncontext: ./frontend\nARG: VITE_API_URL (secret)\nghcr.io/OWNER/bimeet-frontend:latest"]
    Login --> BuildLanding["docker/build-push-action@v5\nBuild landing\ncontext: ./landing\nARG: NEXT_PUBLIC_APP_URL (secret)\nghcr.io/OWNER/bimeet-landing:latest"]

    BuildBackend --> SSH
    BuildFrontend --> SSH
    BuildLanding --> SSH

    SSH["appleboy/ssh-action@v1\nSSH into VPS\n(SSH_HOST / SSH_USER / SSH_KEY secrets)"]

    SSH --> Pull["docker compose -f docker-compose.prod.yml pull\n(pulls all 3 new images from GHCR)"]
    Pull --> Up["docker compose -f docker-compose.prod.yml up -d\n(recreates changed containers)"]
    Up --> Prune["docker image prune -f\n(removes dangling old images)"]
```

### GitHub Actions Secrets Required

| Secret | Used by | Description |
|---|---|---|
| `VITE_API_URL` | Frontend build | Backend API base URL injected at build time |
| `NEXT_PUBLIC_APP_URL` | Landing build | App URL for CTA links |
| `SSH_HOST` | Deploy step | VPS IP or hostname |
| `SSH_USER` | Deploy step | SSH username on VPS |
| `SSH_KEY` | Deploy step | SSH private key |
| `GITHUB_TOKEN` | GHCR login | Auto-provided by GitHub Actions |

---

## Applications

### Backend

A stateless REST API server implementing a strict three-layer architecture: **handler → service → repository → PostgreSQL**.

**Responsibilities:**
- JWT-based authentication (register / login)
- Event lifecycle management (create, update, delete, complete, invite)
- Participant management with status tracking (invited / confirmed / declined)
- Money collections with contribution and payment confirmation flows
- Polls with multi-option voting
- Item checklists with user assignment
- Carpools with seat and passenger tracking
- Business event links management
- Real-time-like notifications (generated on mutations, dispatched via goroutines)
- Event changelog recording on field updates
- User profile management with avatar upload to S3

**Stack:**

| Dependency | Version | Role |
|---|---|---|
| Go | 1.26 | Language runtime |
| `go-chi/chi` | v5.1.0 | HTTP router and middleware |
| `go-chi/cors` | v1.2.1 | CORS middleware |
| `jackc/pgx` | v5.5.5 | PostgreSQL driver and connection pool |
| `golang-jwt/jwt` | v5.2.1 | JWT generation and validation |
| `google/uuid` | v1.6.0 | UUID generation |
| `joho/godotenv` | v1.5.1 | `.env` file loading |
| `aws/aws-sdk-go-v2` | v1.41.3 | S3-compatible storage client (MinIO / AWS S3) |
| `golang.org/x/crypto` | v0.24.0 | bcrypt password hashing |

**Build (multi-stage Dockerfile):**

```
Stage 1 — golang:1.26-alpine
  CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

Stage 2 — alpine:3.20
  ca-certificates + tzdata
  EXPOSE 8080
  CMD ["./server"]
```

**Environment Variables:**

| Variable | Example | Description |
|---|---|---|
| `PORT` | `8080` | HTTP listen port |
| `DSN` | `postgres://bimeet:bimeet@localhost:5432/bimeet` | PostgreSQL connection string |
| `JWT_SECRET` | `change-me` | HMAC secret for JWT signing |
| `JWT_EXP_HOURS` | `72` | Token expiry in hours |
| `SMTP_HOST` | `localhost` | SMTP server host |
| `SMTP_PORT` | `1025` | SMTP server port |
| `SMTP_FROM` | `Bimeet <noreply@bimeet.local>` | Sender address |
| `SMTP_USER` / `SMTP_PASS` | — | SMTP credentials (optional) |
| `FRONTEND_URL` | `http://localhost:5173` | Used in email links |
| `AWS_REGION` | `us-east-1` | S3 region |
| `AWS_ACCESS_KEY_ID` | `minioadmin` | S3 access key |
| `AWS_SECRET_ACCESS_KEY` | `minioadmin` | S3 secret key |
| `S3_BUCKET` | `avatars` | Bucket name |
| `S3_ENDPOINT` | `http://localhost:9000` | Custom endpoint for MinIO; omit for real AWS S3 |
| `S3_PUBLIC_BASE_URL` | `http://localhost:9000/avatars` | Public base URL for serving uploaded files |

**API Endpoints:**

```
Public
  POST   /api/auth/register
  POST   /api/auth/login
  GET    /api/events/invite/{token}

Protected (Authorization: Bearer <token> required)
  GET    /api/auth/me
  PUT    /api/auth/me
  GET    /api/auth/me/stats
  POST   /api/auth/me/avatar
  DELETE /api/auth/me/avatar

  GET    /api/events
  POST   /api/events
  GET    /api/events/public
  POST   /api/events/invite/{token}       — join by invite link
  GET    /api/events/{id}
  PUT    /api/events/{id}
  DELETE /api/events/{id}
  POST   /api/events/{id}/complete
  POST   /api/events/{id}/join            — join public event

  POST   /api/events/{id}/participants
  PATCH  /api/events/{id}/participants/{userId}

  GET    /api/events/{id}/collections
  POST   /api/events/{id}/collections
  GET    /api/events/{id}/collections/summary
  DELETE /api/events/{id}/collections/{collectionId}
  POST   /api/events/{id}/collections/{collectionId}/contribute
  POST   /api/events/{id}/collections/{collectionId}/contributions/{contribId}/confirm
  POST   /api/events/{id}/collections/{collectionId}/contributions/{contribId}/reject
  POST   /api/events/{id}/collections/{collectionId}/contributions/mark-paid

  GET    /api/events/{id}/polls
  POST   /api/events/{id}/polls
  POST   /api/events/{id}/polls/{pollId}/vote

  GET    /api/events/{id}/items
  POST   /api/events/{id}/items
  PATCH  /api/events/{id}/items/{itemId}

  GET    /api/events/{id}/carpools
  POST   /api/events/{id}/carpools
  POST   /api/events/{id}/carpools/{carpoolId}/join

  GET    /api/events/{id}/links
  POST   /api/events/{id}/links
  DELETE /api/events/{id}/links/{linkId}

  GET    /api/notifications
  PATCH  /api/notifications/{id}/read
  POST   /api/notifications/read-all
  DELETE /api/notifications/{id}
  DELETE /api/notifications
```

---

### Frontend

A single-page application following **Feature-Sliced Design** (FSD). Nginx serves the static build and proxies all `/api/` requests to the backend container on the same Docker network.

**Responsibilities:**
- Event discovery (public events list, invite link flow)
- Full event detail page: participants, collections, polls, items, carpools, links, changelog
- Auth flows: register, login, JWT stored in `localStorage`
- User profile with avatar upload
- Notification center
- Dark/light mode via Chakra UI semantic tokens

**Stack:**

| Dependency | Version | Role |
|---|---|---|
| React | 19.2 | UI library |
| TypeScript | ~5.9 | Type safety |
| Vite | 8.0 | Dev server and bundler |
| `@chakra-ui/react` | ^2.10 | Component library (custom theme with semantic tokens) |
| `@tanstack/react-query` | ^5.90 | Server state, caching, invalidation |
| `@tanstack/react-table` | ^8.21 | Table primitives |
| `react-router-dom` | ^7.13 | Client-side routing |
| `i18next` + `react-i18next` | ^26 / ^17 | Internationalisation |
| `framer-motion` | ^12.37 | Animations |
| `react-icons` | ^5.6 | Icon set |
| `@emotion/react` + `@emotion/styled` | ^11.14 | CSS-in-JS (required by Chakra UI) |
| pnpm | — | Package manager |

**FSD Layer Structure:**

```
src/
  app/        providers (Chakra, Query, Auth), router, theme
  pages/      route-level components
  widgets/    Navbar layout, EventCard
  features/   auth, event-manage, participants, collections, polls, items, carpools, links
  entities/   event, user, collection, poll, item, carpool, notification — types + API + queries
  shared/     apiFetch client (reads JWT from localStorage), formatDate utilities
```

**Build (multi-stage Dockerfile):**

```
Stage 1 — node:20-alpine
  pnpm install --frozen-lockfile
  ARG VITE_API_URL (injected at build time)
  pnpm run build → /app/dist

Stage 2 — nginx:alpine
  COPY dist → /usr/share/nginx/html
  COPY nginx.conf (proxies /api/ to backend:8080, SPA fallback)
  EXPOSE 80
```

**Nginx config highlights:**
- `location /api/` → `proxy_pass http://backend:8080` — backend is resolved via Docker DNS on `app-net`
- `location /` → `try_files $uri $uri/ /index.html` — SPA fallback for React Router

---

### Landing

A static marketing site built with Next.js static export and served by Nginx.

**Stack:**

| Dependency | Version | Role |
|---|---|---|
| Next.js | 16.2 | Framework (static export via `next build` → `out/`) |
| React | 19.2 | UI library |
| TypeScript | ^5 | Type safety |
| Tailwind CSS | ^4 | Utility-first CSS |
| PostCSS | — | Tailwind CSS processing |

**Build (multi-stage Dockerfile):**

```
Stage 1 — node:20-alpine
  npm ci
  ARG NEXT_PUBLIC_APP_URL (injected at build time)
  npm run build → /app/out (static HTML/CSS/JS)

Stage 2 — nginx:alpine
  COPY out → /usr/share/nginx/html
  EXPOSE 80
```

---

## Database Schema

PostgreSQL 16. Migrations run automatically at backend startup via embedded `*.up.sql` files tracked in the `schema_migrations` table.

```mermaid
erDiagram
    users {
        uuid id PK
        text name
        text email
        text password_hash
        timestamptz created_at
    }
    events {
        uuid id PK
        text title
        text description
        timestamptz date_start
        timestamptz date_end
        text location
        uuid organizer_id FK
        uuid invite_token
        timestamptz created_at
        timestamptz updated_at
    }
    event_participants {
        uuid id PK
        uuid event_id FK
        uuid user_id FK
        text status
    }
    collections {
        uuid id PK
        uuid event_id FK
        text title
        numeric target_amount
        uuid created_by FK
        timestamptz created_at
    }
    collection_contributions {
        uuid id PK
        uuid collection_id FK
        uuid user_id FK
        boolean paid
        timestamptz paid_at
    }
    polls {
        uuid id PK
        uuid event_id FK
        text question
        uuid created_by FK
        timestamptz created_at
    }
    poll_options {
        uuid id PK
        uuid poll_id FK
        text label
    }
    poll_votes {
        uuid poll_option_id FK
        uuid user_id FK
    }
    items {
        uuid id PK
        uuid event_id FK
        text name
        uuid assigned_to FK
    }
    carpools {
        uuid id PK
        uuid event_id FK
        uuid driver_id FK
        int seats_available
        text departure_point
    }
    carpool_passengers {
        uuid carpool_id FK
        uuid user_id FK
    }
    event_change_logs {
        uuid id PK
        uuid event_id FK
        uuid changed_by FK
        text field_name
        text old_value
        text new_value
        timestamptz changed_at
    }
    notifications {
        uuid id PK
        uuid user_id FK
        uuid event_id FK
        text type
        text message
        boolean is_read
        timestamptz created_at
    }

    users ||--o{ events : organizer
    users ||--o{ event_participants : participant
    events ||--o{ event_participants : has
    events ||--o{ collections : has
    events ||--o{ polls : has
    events ||--o{ items : has
    events ||--o{ carpools : has
    events ||--o{ event_change_logs : tracks
    events ||--o{ notifications : triggers
    users ||--o{ notifications : receives
    collections ||--o{ collection_contributions : has
    polls ||--o{ poll_options : has
    poll_options ||--o{ poll_votes : receives
    users ||--o{ poll_votes : casts
    carpools ||--o{ carpool_passengers : carries
    users ||--o{ carpool_passengers : joins
```

---

## Local Development

### Prerequisites
- Docker + Docker Compose (for services)
- Go 1.26+
- Node.js 20+ and pnpm

### Start infrastructure services

```bash
cd backend
docker compose up -d   # starts PostgreSQL :5432, MinIO :9000/:9001, MailHog :1025/:8025
```

### Backend

```bash
cd backend
cp .env.example .env   # edit as needed
go run ./cmd/server    # starts on :8080, runs migrations automatically
```

### Frontend

```bash
cd frontend
pnpm install
pnpm run dev           # starts Vite dev server on :5173
```

### Landing

```bash
cd landing
npm install
npm run dev            # starts Next.js dev server on :3000
```

###Local service URLs

| Service | URL |
|---|---|
| Frontend dev | `http://localhost:5173` |
| Landing dev | `http://localhost:3000` |
| Backend API | `http://localhost:8080` |
| MinIO console | `http://localhost:9001` |
| MailHog UI | `http://localhost:8025` |
