# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Lejematch (lejematch.dk) is a Danish roommate-finding platform. Users can list rooms (or entire apartments) they have available, and people searching for a place to live can browse and find a match. The core flow is: someone with a spare room lists it → someone looking for a room finds and connects with them.

## Commands

```bash
# Run the application (loads dev.env)
go run main.go -env dev

# Build
go build -o lejematch ./

# Run all tests
go test ./...

# Run a single package's tests
go test -v ./internal/security

# Vet
go vet ./...

# Start the PostgreSQL database
docker-compose up -d
```

The `-env` flag determines which `.env` file to load (e.g., `-env dev` loads `dev.env`). The app panics on startup if the env file is not found.

Database migrations run automatically via GORM AutoMigrate on startup.

## Architecture

**Framework:** Gofiber v2 (Express.js-like). **ORM:** GORM with PostgreSQL. **Auth:** JWT (HS256, 24h expiry) + Argon2id password hashing.

### Layered structure

```
main.go / cmd/app.go         → Fiber server setup, middleware, route registration
api/v1/routes/               → Route declarations (setup.go wires everything)
api/v1/handlers/             → HTTP handlers (thin layer: parse, call service, respond)
api/auth/JWTmiddleware.go    → JWT Bearer token validation middleware
internal/services/           → Business logic (userService.go, listingService.go, jwt.go)
internal/database/repo/      → Repository layer (generic.go + specialized repos)
internal/database/models/    → GORM models (User, Profile, Listing)
internal/security/           → Argon2id password hashing/verification
config/config.go             → Env-file based config loading
```

### Data model

- `User` — core account record (email unique, phone unique, IsAdmin, IsActive)
- `Profile` — 1:1 with User (cascade delete), holds display info (DisplayName, Bio, City, ImageURL)
- `Listing` — room/apartment listing posted by a User. Fields: Title, Description, Price (DKK/month), City, Area (street/neighbourhood — no house number), RoomType (`private`/`shared`/`apartment`), Status (`active`/`rented`/`archived`), AvailableFrom (ISO date string), Images (JSON array of URL strings stored as `jsonb`)
- User creation always creates a Profile in the same transaction via `UserService`.

### Repository pattern

`internal/database/repo/generic.go` provides a `GenericRepo[T]` using Go generics. Specialized repos (`UsersRepo`, `ProfilesRepo`, `ListingsRepo`) embed it and add query-specific methods. `ListingsRepo.FindFiltered` supports offset pagination and filtering by city, price range, and room type.

### API routes

```
GET    /health
POST   /api/v1/auth/login

POST   /api/v1/users                   (public — create user + profile)
GET    /api/v1/users/:id               (JWT required)
PATCH  /api/v1/users/:id               (JWT required)
DELETE /api/v1/users/:id               (JWT required)
PUT    /api/v1/users/:id/password      (JWT required)
GET    /api/v1/users/:id/profile       (public)
PATCH  /api/v1/users/:id/profile       (JWT required)

GET    /api/v1/listings                (public — paginated, filterable)
GET    /api/v1/listings/:id            (public)
POST   /api/v1/listings                (JWT required)
PATCH  /api/v1/listings/:id            (JWT required — own listing or admin)
DELETE /api/v1/listings/:id            (JWT required — own listing or admin)
```

JWT claims carry `UserID`, `Email`, `IsAdmin`, `IsActive`. Handlers enforce that non-admins can only access their own records. For listings, ownership is checked in the service layer via `ErrNotOwner`.

## Environment variables (`dev.env`)

| Variable | Purpose |
|---|---|
| `ENV` | Environment name |
| `API_PORT` | Fiber listen port (default 3000) |
| `DATABASE_URL` | GORM PostgreSQL DSN |
| `JWT_SECRET` | HMAC signing secret |