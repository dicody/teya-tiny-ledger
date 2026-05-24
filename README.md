# tiny-ledger 🏦

A simple ledger service built with Go. Records deposits and withdrawals, derives balances from transaction history, and
provides a REST API.

## Assumptions

- **Amounts in cents** — All monetary values are `int64` cents (e.g., `$10.50` → `1050`) to avoid floating-point
  precision issues.
- **Balance is derived** — No stored balance. Balance is always computed by summing transactions — the "ledger way".
- **Accounts auto-created** — Accounts are created implicitly on the first deposit. No explicit account creation
  endpoint.
- **Negative balance rejected** — Withdrawals that would make the balance negative return `400 Bad Request`.
- **No authentication** — All endpoints are open.
- **In-memory storage** — Data is lost when the server stops.
- **Concurrency / Race Conditions** — To keep the solution simple and aligned with the "no atomic operations" constraint, application-level locking (like optimistic concurrency control) is not implemented. While the in-memory store uses a mutex to prevent map panics, a TOCTOU (Time-Of-Check to Time-Of-Use) race condition could technically occur during concurrent withdrawals.

## Getting Started

### Prerequisites 

- Go 1.26+ (uses method-based routing in `net/http`)

### Run the server

```bash
make run
```

The server starts on `http://localhost:8080`.

### Run tests

```bash
make test
```

### Build binary

```bash
make build
./bin/tiny-ledger
```

## API

All amounts are in **cents** (integer).

### Record a deposit

```bash
curl -X POST http://localhost:8080/api/v1/accounts/dmitry/transactions \
  -H 'Content-Type: application/json' \
  -d '{"type":"deposit","amount":5000}'
```

Response (`201 Created`):

```json
{
  "id": "a1b2c3d4e5f6a7b8",
  "account_id": "dmitry",
  "type": "deposit",
  "amount": 5000,
  "created_at": "2026-05-24T12:00:00Z"
}
```

### Record a withdrawal

```bash
curl -X POST http://localhost:8080/api/v1/accounts/dmitry/transactions \
  -H 'Content-Type: application/json' \
  -d '{"type":"withdrawal","amount":2000}'
```

Response (`201 Created`):

```json
{
  "id": "b2c3d4e5f6a7b8c9",
  "account_id": "dmitry",
  "type": "withdrawal",
  "amount": 2000,
  "created_at": "2026-05-24T12:01:00Z"
}
```

### Get balance

```bash
curl http://localhost:8080/api/v1/accounts/dmitry/balance
```

Response (`200 OK`):

```json
{
  "account_id": "dmitry",
  "balance": 3000
}
```

### Get transaction history

```bash
curl http://localhost:8080/api/v1/accounts/dmitry/transactions
```

Response (`200 OK`):

```json
{
  "account_id": "dmitry",
  "transactions": [
    {"id": "a1b2c3d4e5f6a7b8", "account_id": "dmitry", "type": "deposit", "amount": 5000, "created_at": "2026-05-24T12:00:00Z"},
    {"id": "b2c3d4e5f6a7b8c9", "account_id": "dmitry", "type": "withdrawal", "amount": 2000, "created_at": "2026-05-24T12:01:00Z"}
  ]
}
```

### Response codes

| Status | Meaning                                                          |
|--------|------------------------------------------------------------------|
| `201`  | Transaction created successfully                                 |
| `200`  | Balance or transaction history returned                          |
| `400`  | Invalid amount, invalid type, insufficient funds, malformed JSON |

## Project Structure

```
tiny-ledger/
├── cmd/ledger/main.go           # Entry point
├── internal/
│   ├── domain/
│   │   ├── model.go             # Transaction, CalculateBalance, validation
│   │   └── model_test.go
│   ├── app/
│   │   ├── service.go           # Ledger & Repository interfaces
│   │   ├── ledger.go            # Ledger implementation
│   │   └── ledger_test.go
│   ├── storage/memstore.go      # In-memory repository
│   └── handler/http/
│       ├── http.go              # Router, writeJSON helper
│       ├── balance.go           # GET balance handler
│       ├── transactions.go      # GET/POST transactions handlers
│       └── http_test.go
├── docs/ledger_requirements.md
├── Makefile
└── README.md
```
