# 💸 Loan Service - Go REST API

A scalable, clean-architecture-based REST API for managing loan states (proposed, approved, invested, disbursed) — built with **Golang**, **Docker**, and **hot reload via Air** for development.

---

## 📦 Environment & Tech Stack

- **Go**: 1.21+
- **PostgreSQL**: 15+
- **Docker + Docker Compose**
- **Air** (hot reload)
- **Makefile**: workflow automation
- **logrus**: structured logging
- **Viper**: config management
- **GORM**: ORM for Go
- **Testify**: for testing

---

## 📂 Project Structure

```
loan-service/
├── cmd/                  # Entry point
├── config/               # Load .env.* config
├── controllers/          # HTTP handlers
├── routes/               # Router setup
├── middleware/           # Auth middleware (API key)
├── models/               # Database models
├── dto/                  # DTOs for request/response
├── services/             # Business logic layer
├── repository/           # Data access abstraction
├── utils/                # Helper functions (email, file, logger)
├── Dockerfile            # For production container
├── docker-compose.yml    # For production
├── docker-compose.dev.yml # For development w/ hot reload
├── .air.toml             # Hot reload config
└── Makefile              # CLI workflow
```

---

## ⚙️ Features

- ✅ Clean project structure (modular)
- ✅ `.env` config loading
- ✅ RESTful endpoint `/loans`
- ✅ Health check `/health`
- ✅ API key authentication via `X-API-KEY`
- ✅ Dockerized (production-ready)
- ✅ Development hot reload with `air`
- ✅ Makefile for easy command access

---

## 🔐 Environment Setup

API_SECRET=your-secure-api-key
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=loan_service

---

## 🚀 How to Run

### 📦 Production

> Build & run production container

```bash
make prod-build
```

Test API:

```bash
curl http://localhost:8080/health
```

### 🔁 Development (Auto reload)

> Watch for file changes and rebuild app automatically

```bash
make dev-build
```

When editing any `.go` file, the service auto-restarts.

### 🛠️ Makefile Commands

| Command            | Description                             |
|-------------------|-----------------------------------------|
| `make dev-build`  | Dev mode + rebuild + hot reload         |
| `make dev-logs`   | View dev logs                           |
| `make dev-stop`   | Stop dev containers + volumes           |
| `make prod-build` | Rebuild and run production container    |
| `make prod-logs`  | View production logs                    |
| `make prod-stop`  | Stop production containers              |

---

## 🧪 API Endpoints

### 1. Health Check
**GET** `/health`

### 2. Create Loan
**POST** `/loans`

**Headers:**
`X-API-KEY: {your_api_key}`

**Request Body:**
```json
{
  "borrower_id": "BORR-001",
  "principal": 1500000,
  "rate": 12,
  "roi": 5
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Loan created",
  "data": { ... }
}
```

### 3. Approve Loan
**PUT** `/loans/{id}/approve`

**Form Data:**
- `field_validator_id`: EMP-001
- `proof_image`: (jpg/jpeg)

### 4. Invest in Loan
**POST** `/loans/{id}/invest`

**JSON:**
```json
{
  "investor_id": "INV-001",
  "amount": 1000000
}
```

### 5. Disburse Loan
**PUT** `/loans/{id}/disburse`

**Form Data:**
- `field_officer_id`: EMP-002
- `signed_agreement`: (jpg/jpeg)
- `disburse_notes`: (optional)

---

## 🔐 API Key Security

All endpoints except `/health` require:

```
X-API-KEY: your_secret_key_here
```

Set in `.env.dev` or `.env.prod`:

```
API_SECRET=your_very_strong_secret_key
```

---

## 🧪 Testing

Run tests:

```bash
make dev-test
```

Coverage report generated using `go test -cover`.

---

## ✨ Author

Created by **Erick Estrada** with ❤️ in Golang.