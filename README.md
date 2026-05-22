# LinkVault

> A side project to learn new tech — a full-stack bookmark manager with auto-metadata scraping, JWT authentication, and AWS DynamoDB persistence.

[![Go Version](https://img.shields.io/badge/Go-1.25-blue)](https://go.dev)
[![React](https://img.shields.io/badge/React-19-61DAFB?logo=react)](https://react.dev)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## Table of Contents

- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Features](#features)
- [Project Structure](#project-structure)
- [Quick Start](#quick-start)
- [Environment Variables](#environment-variables)
- [API Reference](#api-reference)
- [DynamoDB Design](#dynamodb-design)
- [Deployment](#deployment)
- [Screenshots](#screenshots)
- [Contributing](#contributing)

---

## Overview

**LinkVault** is a modern bookmark management application that lets you save, organize, and retrieve your favorite links. It automatically scrapes webpage titles so you don't have to type them manually, supports tagging, and uses a clean single-table DynamoDB design for efficient data access.

This project was built as a learning exercise to explore:
- Go with the Gin web framework
- AWS SDK v2 and DynamoDB single-table design
- JWT authentication flow
- Docker containerization & Caddy reverse proxy
- React 19 + TypeScript + Vite frontend tooling

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| **Backend** | Go 1.25, Gin, AWS SDK Go v2 |
| **Database** | Amazon DynamoDB (single-table design) |
| **Auth** | JWT (golang-jwt/jwt/v5), bcrypt |
| **Frontend** | React 19, TypeScript, Vite |
| **DevOps** | Docker, Docker Compose, Caddy |
| **Utilities** | UUID generation, URL title scraper |

---

## Features

- **User Authentication** — Register and login with email/password. JWT tokens with 72-hour expiry.
- **Bookmark CRUD** — Create, list, and delete bookmarks.
- **Auto Title Scraping** — Automatically fetches the `<title>` from a URL if left blank.
- **Tagging System** — Organize bookmarks with custom tags.
- **Single-Table DynamoDB** — Efficient PK/SK pattern for users and bookmarks.
- **CORS Ready** — Pre-configured for frontend integration.
- **Dockerized** — One-command local setup with DynamoDB Local.
- **Caddy Reverse Proxy** — Production-ready HTTPS and routing.

---

## Project Structure

```
LinkVault/
├── cmd/server/           # Application entry point
│   └── main.go
├── internal/
│   ├── config/           # Environment configuration
│   ├── db/               # DynamoDB client setup
│   ├── handler/          # HTTP handlers (auth, bookmark, health)
│   ├── middleware/       # JWT authentication middleware
│   ├── model/            # Data models (User, Bookmark)
│   ├── router/           # Route definitions & CORS
│   └── scraper/          # URL title fetcher
├── frontend/             # React + TypeScript + Vite SPA
│   ├── src/
│   ├── index.html
│   └── package.json
├── aws/                  # AWS deployment configurations
├── docs/                 # Architecture documentation
├── docker-compose.yml    # Local orchestration
├── Dockerfile            # API container image
├── Dockerfile.caddy      # Caddy container image
├── Caddyfile             # Caddy reverse proxy config
└── .env.example          # Environment variable template
```

---

## Quick Start

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- [Go 1.25+](https://go.dev/dl/) (for local development)
- [Node.js 20+](https://nodejs.org/) (for frontend development)

### 1. Clone the repository

```bash
git clone https://github.com/Sayyedarham/LinkVault.git
cd LinkVault
```

### 2. Set up environment variables

```bash
cp .env.example .env
```

Edit `.env` with your values (see [Environment Variables](#environment-variables)).

### 3. Run with Docker Compose

```bash
docker-compose up --build
```

This starts three services:
- **API** at `http://localhost:8080`
- **DynamoDB Local** at `http://localhost:8000`
- **Caddy** reverse proxy at `http://localhost:80`

### 4. Run frontend locally (optional)

```bash
cd frontend
npm install
npm run dev
```

The frontend will be available at `http://localhost:5173`.

### 5. Run backend locally (optional)

```bash
go run cmd/server/main.go
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `:8080` | API server port |
| `AWS_REGION` | `us-east-1` | AWS region |
| `DYNAMO_ENDPOINT` | `http://dynamodb-local:8000` | DynamoDB endpoint (local) |
| `TABLE_NAME` | `LinkTable` | DynamoDB table name |
| `JWT_SECRET` | `dev-secret-change-in-prod` | JWT signing secret |

> **Production Tip:** Always change `JWT_SECRET` to a cryptographically secure random string in production.

---

## API Reference

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/auth/register` | Register a new user |
| `POST` | `/api/v1/auth/login` | Login and receive JWT |

**Register / Login Request Body:**
```json
{
  "email": "user@example.com",
  "password": "min6chars"
}
```

**Success Response:**
```json
{
  "token": "eyJhbG...",
  "user": {
    "id": "uuid",
    "email": "user@example.com"
  }
}
```

### Bookmarks *(Requires Bearer Token)*

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/bookmarks` | Create a new bookmark |
| `GET` | `/api/v1/bookmarks` | List all user bookmarks |
| `DELETE` | `/api/v1/bookmarks/:id` | Delete a bookmark |

**Create Bookmark Request Body:**
```json
{
  "url": "https://example.com",
  "title": "Optional Title",
  "tags": ["go", "web"]
}
```

**List Bookmarks Response:**
```json
{
  "bookmarks": [
    {
      "id": "uuid",
      "url": "https://example.com",
      "title": "Example Domain",
      "tags": ["demo"],
      "created_at": "2026-05-22T12:00:00Z",
      "updated_at": "2026-05-22T12:00:00Z"
    }
  ],
  "count": 1
}
```

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/health` | Service health status |

---

## DynamoDB Design

LinkVault uses a **single-table design** with composite keys:

| Entity | PK | SK | Attributes |
|--------|----|----|-----------|
| **User** | `USER#<email>` | `PROFILE` | `userId`, `email`, `passwordHash` |
| **Bookmark** | `USER#<userId>` | `BOOKMARK#<bookmarkId>` | `id`, `url`, `title`, `tags`, `created_at`, `updated_at` |

This pattern enables efficient querying:
- Fetch user profile by PK + SK
- List all bookmarks for a user with a single `Query` on PK with `begins_with(SK, "BOOKMARK#")`

See [`docs/dynamodb-design.md`](docs/dynamodb-design.md) for detailed schema documentation.

---

## Deployment

### AWS Deployment

The `aws/` directory contains infrastructure-as-code and deployment scripts for AWS services. Configure your AWS credentials and run the provided scripts to deploy DynamoDB tables and optionally ECS/Fargate for the API.

### Production Checklist

- [ ] Change `JWT_SECRET` to a secure random string
- [ ] Use AWS DynamoDB instead of DynamoDB Local
- [ ] Enable HTTPS via Caddy or AWS ALB
- [ ] Set up CI/CD pipeline (GitHub Actions recommended)
- [ ] Configure CloudFront for frontend static assets
- [ ] Add monitoring (CloudWatch / Prometheus)

---

## Screenshots

*Coming soon*

---

## Contributing

Contributions are welcome! This is a learning project, so feel free to open issues or PRs for:
- Frontend UI improvements
- Additional bookmark features (folders, search, import/export)
- Testing coverage
- Infrastructure enhancements

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## License

This project is open source and available under the [MIT License](LICENSE).

---

<p align="center">Built with ☕ by <a href="https://github.com/Sayyedarham">@Sayyedarham</a></p>
