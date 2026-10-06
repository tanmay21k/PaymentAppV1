# Stack

- go
- postgresql

### Tables

bank - userId, balance
user - firstname, lastname, username, password, paymentPin, profilePicture
transaction - userId, type, amount , balance, counterpartyId, transferId   , status

### middleware

- auth middleware using jwt

### routes

/balance
/recent
/transfer

"/signup"
"/signin"
"/pin"
"/profile-picture",
"/bulk"
"/recipient/:id
"/me"


```
backend/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   │
│   ├── database/
│   │   └── postgres.go
│   │
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── logger.go
│   │   └── recovery.go
│   │
│   ├── auth/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── token.go
│   │
│   ├── user/
│   │   ├── handler.go - How does HTTP talk to my application?
│   │   ├── service.go - What business rules should happen?
│   │   ├── model.go - How do I get/store data?
│   │   ├── types.go - What is the domain object?
│   │   └── repository.go - What data does this operation accept/return?
│   │
│   └── server/
│       └── server.go
│
├── db/
│   ├── migrations/
│   │   ├── 000001_create_users.up.sql
│   │   └── 000001_create_users.down.sql
│   │
│   └── queries/
│       └── users.sql
│
├── api/
│   └── ...
│
├── tests/
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
└── .env.example

```