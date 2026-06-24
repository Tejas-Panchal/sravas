# Sravas

A YouTube clone built to learn DevOps end-to-end: 11 microservices, React frontend, Docker, Terraform, Kubernetes, Helm, CI/CD.

## Architecture

```
┌──────────┐    ┌──────────────┐
│  React   │───▶│ API Gateway  │
│ Frontend │    │   :8000      │
└──────────┘    └──────┬───────┘
                       │ routes to services
       ┌───────────────┼─────────────┐
       │               │             │
┌──────▼──┐   ┌────────▼┐   ┌────────▼─┐
│ Auth    │   │ Upload  │   │  Stream  │
│ :3001   │   │ :3003   │   │  :3004   │
└─────────┘   └─────────┘   └──────────┘
┌─────────┐   ┌─────────┐   ┌──────────┐
│ User    │   │ Search  │   │  Comment │
│ :3002   │   │ :3005   │   │  :3006   │
└─────────┘   └─────────┘   └──────────┘
┌───────────┐   ┌─────────────┐
│ Analytics │   │ Notification│
│  :3007    │   │  :3008      │
└───────────┘   └─────────────┘
┌────────────────┐
│ Scheduler      │
│ (cron only)    │
└────────────────┘
```

## Stack

| Layer           | Technology                                  |
| --------------- | ------------------------------------------- |
| Frontend        | React 18, Vite, Tailwind CSS, Redux Toolkit |
| Backend         | Go (chi/v5), Python (FastAPI)               |
| Database        | PostgreSQL, MongoDB, Elasticsearch          |
| Cache           | Redis                                       |
| Queue           | Apache Kafka                                |
| Container       | Docker, docker-compose                      |
| Infra (planned) | Terraform, AWS EKS, Helm                    |
| CI/CD (planned) | GitHub Actions                              |

## Services

| Service              | Port | Language | Description                                       |
| -------------------- | ---- | -------- | ------------------------------------------------- |
| api-gateway          | 8000 | Go       | Reverse proxy, CORS, logging                      |
| auth-service         | 3001 | Go       | Register, login, JWT                              |
| user-service         | 3002 | Go       | Profile, subscriptions                            |
| upload-service       | 3003 | Go       | Multipart upload, validation                      |
| streaming-service    | 3004 | Go       | HLS manifest, segment serving                     |
| search-service       | 3005 | Go       | Inverted index, autocomplete                      |
| comment-service      | 3006 | Go       | CRUD, nested replies, likes                       |
| analytics-service    | 3007 | Go       | Event tracking, trending                          |
| notification-service | 3008 | Go       | In-app notifications                              |
| scheduler-service    | —    | Go       | Cron jobs (uploads, trending, cleanup, analytics) |

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Git

### Run locally

```bash
git clone https://github.com/Tejas-Panchal/sravas.git
cd sravas
cp .env.example .env
docker compose up --build
```

Open **http://localhost** for the frontend.
The API Gateway is available at **http://localhost:8000**.

### Run frontend without Docker

```bash
cd frontend
npm install
npm run dev
```

Vite dev server runs on **http://localhost:5173** and proxies `/api` to `:8000`.

### Run a Go service without Docker

```bash
cd services/<service-name>
go run ./cmd/server
```

## Project Structure

```
sravas/
├── frontend/           # React app (Vite + Tailwind)
│   ├── src/            #   28 source files
│   ├── Dockerfile
│   └── nginx.conf
├── services/
│   ├── api-gateway/          # Go — reverse proxy
│   ├── auth-service/         # Go — JWT auth
│   ├── user-service/         # Go — user CRUD
│   ├── upload-service/       # Go — file upload
│   ├── streaming-service/    # Go — HLS streaming
│   ├── search-service/       # Go — full-text search
│   ├── comment-service/      # Go — comments
│   ├── analytics-service/    # Go — analytics
│   ├── notification-service/ # Go — notifications
│   └── scheduler-service/    # Go — cron jobs
├── infrastructure/     # Terraform (planned)
├── monitoring/         # Prometheus/Grafana (planned)
├── tests/              # Unit, integration, e2e, load
├── docker-compose.yml
└─ .env.example
```
