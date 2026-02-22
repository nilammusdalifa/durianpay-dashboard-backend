# Backend — DurianPay Payment Dashboard

REST API built in Go using the [chi](https://github.com/go-chi/chi) router, JWT authentication, and in-memory storage. The API contract is defined in `openapi.yaml` and validated at runtime using `oapi-codegen`.

---

## Tech Stack

| Library | Purpose |
|---------|---------|
| `go-chi/chi` | HTTP router |
| `oapi-codegen` | OpenAPI code generation + request validation middleware |
| `golang-jwt/jwt` | JWT creation and validation |
| `golang.org/x/crypto/bcrypt` | Password hashing |
| `joho/godotenv` | `.env` file loading |

---

## Project Structure

```
backend/
├── main.go                          # Entry point — wires everything together, seeds data
├── env.sample                       # Example environment variables
├── go.mod / go.sum                  # Go module files
├── Makefile                         # Dev commands
├── Dockerfile                       # Container build
├── internal/
│   ├── api/
│   │   └── api_handler.go           # Implements the OpenAPI ServerInterface
│   ├── config/
│   │   └── env.go                   # Reads environment variables
│   ├── entity/
│   │   ├── error.go                 # AppError type and constructors
│   │   ├── payment.go               # Payment struct and PaymentFilter
│   │   └── user.go                  # User struct
│   ├── module/
│   │   ├── auth/
│   │   │   ├── handler/auth.go      # HTTP handler — decodes request, calls usecase
│   │   │   ├── repository/user.go   # In-memory user store (thread-safe)
│   │   │   └── usecase/auth.go      # Login logic — bcrypt verify + JWT sign
│   │   └── payment/
│   │       ├── handler/payment.go   # HTTP handler — validates JWT, calls usecase
│   │       ├── repository/payment.go# In-memory payment store (thread-safe, filterable)
│   │       └── usecase/payment.go   # Business logic — status validation
│   ├── openapigen/
│   │   └── openapi.gen.go           # Auto-generated from openapi.yaml (do not edit)
│   ├── service/http/
│   │   └── server.go                # HTTP server setup, CORS, OpenAPI validation middleware
│   └── transport/
│       └── jsonerror.go             # Error → JSON HTTP response writer
└── script/gen-secret/
    └── main.go                      # Helper to generate a JWT secret
```

---

## Setup

**1. Copy the environment file:**

```bash
cp env.sample .env
```

**2. Edit `.env` if needed** (defaults work out of the box):

```env
HTTP_ADDR=:8080
OPENAPIYAML_LOCATION=../openapi.yaml
JWT_SECRET=your-very-secret-key
JWT_EXPIRED=24h
```

**3. Download dependencies:**

```bash
go mod download
```

**4. Run:**

```bash
go run main.go
```

You should see:
```
✓ seeded 2 users
✓ seeded 50 payments
🚀 Server starting on :8080 (in-memory store)
listening on :8080
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_ADDR` | `:8080` | Address the server listens on |
| `OPENAPIYAML_LOCATION` | `../openapi.yaml` | Path to the OpenAPI spec file |
| `JWT_SECRET` | `dev-secret-replace-me` | Secret key used to sign JWT tokens |
| `JWT_EXPIRED` | `24h` | How long a JWT token is valid |

---

## API Endpoints

### POST `/dashboard/v1/auth/login`

No authentication required.

**Request:**
```json
{
  "email": "cs@test.com",
  "password": "password"
}
```

**Response `200`:**
```json
{
  "email": "cs@test.com",
  "role": "cs",
  "token": "<jwt-token>"
}
```

**Response `401`:**
```json
{
  "code": 401,
  "message": "invalid credentials"
}
```

---

### GET `/dashboard/v1/payments`

Requires `Authorization: Bearer <token>` header.

**Query parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `status` | string | Filter by `completed`, `processing`, or `failed` |
| `id` | string | Filter by exact payment ID (e.g. `PAY-00001`) |
| `sort` | string | Sort field. Prefix with `-` for descending. e.g. `-created_at`, `amount` |

**Response `200`:**
```json
{
  "payments": [
    {
      "id": "PAY-00001",
      "merchant": "Tokopedia",
      "status": "completed",
      "amount": "150000.00",
      "created_at": "2026-02-10T10:30:00Z"
    }
  ]
}
```

**Response `401`:**
```json
{
  "code": 401,
  "message": "Unauthenticated: missing or invalid token"
}
```

---

## Example curl Commands

**Login:**
```bash
curl -X POST http://localhost:8080/dashboard/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"cs@test.com","password":"password"}'
```

**Get all payments:**
```bash
curl http://localhost:8080/dashboard/v1/payments \
  -H "Authorization: Bearer <token>"
```

**Filter by status:**
```bash
curl "http://localhost:8080/dashboard/v1/payments?status=completed" \
  -H "Authorization: Bearer <token>"
```

**Sort by newest:**
```bash
curl "http://localhost:8080/dashboard/v1/payments?sort=-created_at" \
  -H "Authorization: Bearer <token>"
```

**One-liner (login + fetch):**
```bash
TOKEN=$(curl -s -X POST http://localhost:8080/dashboard/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"cs@test.com","password":"password"}' \
  | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

curl http://localhost:8080/dashboard/v1/payments \
  -H "Authorization: Bearer $TOKEN"
```

---

## Running Tests

```bash
go test ./... -v
```

Tests cover the payment usecase:
- Listing all payments with no filter
- Filtering by valid status (`completed`, `processing`, `failed`)
- Rejecting an invalid status value

---

## Build for Production

```bash
go build -o bin/server main.go
./bin/server
```

---

## Generating a Secure JWT Secret

```bash
go run ./script/gen-secret/main.go
```

Copy the output into your `.env` as `JWT_SECRET`.

---

## Regenerating OpenAPI Code

If you modify `openapi.yaml`, regenerate the server code with:

```bash
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
oapi-codegen -generate "types,chi-server,spec" -package openapigen \
  -o internal/openapigen/openapi.gen.go ../openapi.yaml
```