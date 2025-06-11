# 🧳 Job Portal REST API (Golang, Gin, JWT, Redis, Docker)

A production-grade Job Portal REST API built in **Golang** using **Gin**, **GORM**, **JWT with RSA public/private keys**, and **Redis** for caching. The project follows clean architecture with middleware, dependency injection, structured logging via **Zerolog**, and unit testing using mocks. Docker is used for containerization, and `encoding/json` is used for data serialization.

## 🚀 Features

- User Signup/Login with JWT authentication (RSA-based)
- Company and Job management with protected routes
- Password reset via email-based flow
- Redis caching for performance
- Middleware for authentication and logging
- Structured logging using Zerolog
- Dockerized for easy deployment
- Unit testing with mocking

## 🔐 Authentication

- JWT Auth using RSA **private/public keys**
- Auth middleware protects all sensitive endpoints

## 📦 API Endpoints

### ✅ Public

| Method | Endpoint         | Description                     |
|--------|------------------|---------------------------------|
| GET    | `/check`         | Auth test route                 |
| POST   | `/signup`        | Register a new user             |
| POST   | `/login`         | Login and get JWT               |
| POST   | `/forget`        | Request password reset          |
| POST   | `/password`      | Set new password                |

### 🔒 Protected (JWT Required)

| Method | Endpoint                              | Description                          |
|--------|----------------------------------------|--------------------------------------|
| POST   | `/createCompany`                      | Create a new company                 |
| GET    | `/getallcompanies`                    | Get all companies                    |
| GET    | `/getacompany/:cid`                   | Get company by ID                    |
| POST   | `/companies/:cid`                     | Post a job under a company           |
| GET    | `/jobs/:CompanyId`                    | Get all jobs under a specific company|
| GET    | `/jobs`                               | Get all jobs                         |
| GET    | `/jobs/jid`                           | Get job by job ID                    |
| POST   | `/process/applications`               | Process job applications             |

## 🧪 Tech Stack

- **Golang**
- **Gin** – Web framework
- **GORM** – ORM for PostgreSQL
- **JWT (RS256)** – Auth with public/private key pair
- **Redis** – Caching layer
- **Zerolog** – Fast structured logging
- **Docker** – Containerization
- **encoding/json** – Data encoding
- **Testify + Mocks** – Unit testing

## 🧰 Getting Started

```bash
# Clone the repository
git clone https://github.com/your-username/jobportal-rest-api.git
cd jobportal-rest-api

# Copy and configure your environment variables
cp .env.example to .env

# Generate RSA keys
openssl genrsa -out private.key 2048
openssl rsa -in private.key -pubout -out public.key

# Run the application
 go run .\cmd\job-portal-api\main.go
