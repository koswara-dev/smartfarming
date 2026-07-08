# SmartFarming Backend API

This repository contains the backend RESTful API service for the SmartFarming fullstack application.

---

## 🚀 Tech Stack

*   **Programming Language**: Go (Golang)
*   **Web Framework**: [Gin Gonic](https://github.com/gin-gonic/gin)
*   **Database ORM**: [GORM](https://gorm.io/)
*   **SQL Database**: PostgreSQL (Stores user records, categories, and articles metadata)
*   **NoSQL Database**: MongoDB (Available client connection for telemetry and sensor data storage)
*   **Cache & Sesi blacklisting**: Redis (Handles temporary signup session caches, OTP verification, and JWT blacklisting on logout)
*   **Object Storage**: MinIO (Stores user profile photos and article image assets)
*   **Live Reload**: [Air](https://github.com/air-verse/air) (For hot compilation during development)
*   **API Documentation**: [Swagger UI / Swag](https://github.com/swaggo/swag)
*   **Logging**: Custom Daily Rotating File Logger (Rotates every 24 hours, retains logs for 30 days)

---

## 📁 Project Architecture

The codebase follows a modular clean-architecture pattern dividing responsibilities cleanly into layers:
```text
backend/
├── config/         # App Config, Database connections (Postgres, Mongo, Redis, MinIO), Seeders
├── docs/           # Auto-generated Swagger documentation
├── dto/            # Data Transfer Objects for requests and responses
├── handler/        # Controllers mapping HTTP requests to business logic
├── middleware/     # Auth checks, IDOR middleware, CORS settings, Role-based constraints (RBAC)
├── model/          # GORM database schemas and base audit hooks
├── repository/     # Database read/write queries
├── routes/         # API Endpoint routing setups
├── service/        # Core business operations, calculations, and validators
└── tests/          # External unit and integration security test suites
```

---

## 🔑 Core Features

1.  **2-Stage OTP Registration**: Users initiate registration; a 6-digit OTP is sent via Redis. The account is created in PostgreSQL only after verifying the code.
2.  **JWT Authentication & Logout**: Symmetrical tokens are validated on protected routes. Logout invalidates active tokens using a Redis-backed blacklist.
3.  **Role-Based Access Control (RBAC)**: Supports roles: `admin`, `operator`, `user`. Categories and Articles are publicly readable (`GET`) but modifications (`POST`/`PUT`/`DELETE`) are restricted to administrative roles.
4.  **Automatic Auditing**: Incorporates a GORM hook setup mapping creator (`CreatedBy`), updater (`UpdatedBy`), and deleter (`DeletedBy`) UUIDs automatically from authenticated contexts.
5.  **MinIO File Upload**: Profile photos and article images are validation-checked (max 5MB, whitelisted image MIMEs, unique UUID name generator) and uploaded directly to MinIO.
6.  **OWASP Top 10 Protection**: Built-in mitigations for SQLi, XSS, SSRF, session hijackings, and Broken Access Control (verified by integration tests).

---

## 🛠️ Getting Started

### 1. Prerequisites
Ensure you have the following services installed/running:
*   Go (v1.20 or newer)
*   PostgreSQL
*   MongoDB
*   Redis
*   MinIO

### 2. Configure Environment
Create/edit the environment files `.env.development` or `.env.production` in the root backend directory:
```env
PORT=8081
ENV=development

DB_HOST=localhost
DB_PORT=5434
DB_USER=postgres
DB_PASSWORD=secret45
DB_NAME=smartfarmingdb
DB_SSLMODE=disable

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

MONGO_URI=mongodb://localhost:27017
MONGO_DB=smartfarmingdb

JWT_SECRET=supersecretjwtkeythatisreallylongandsecurefordevelopment

MINIO_ENDPOINT=172.30.54.28:9000
MINIO_ACCESS_KEY=admin
MINIO_SECRET_KEY=Devjc@1324
MINIO_USE_SSL=false
MINIO_BUCKET_NAME=smartfarming
```

### 3. Run the Server
To run in hot-reload mode:
```bash
air
```
Or run directly:
```bash
go run main.go
```

### 4. Interactive API Documentation
Open your browser and navigate to:
`http://localhost:8081/swagger/index.html`

---

## 🧪 Running Security & Unit Tests

To execute the unit tests, RBAC access integration tests, and OWASP Top 10 checks:
```bash
go test ./... -v
```
