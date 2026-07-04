# TODO

## Security

- [ ] **CORS** — restrict `AllowOrigins` from `"*"` to specific domain(s) (e.g. `https://lejematch.dk`) before production (`cmd/app.go`)
- [ ] **ImageURL validation** — validate that `ImageURL` on profile and `Images` on listings are valid URLs (`internal/services/`)
- [ ] **Email format validation** — add application-level email format check on user create/update (`internal/services/userService.go`)
- [ ] **Minimum password length** — consider raising from 8 to 12+ characters (`internal/security/password.go`)
- [ ] **Server header** — remove or obscure the version/build info in `ServerHeader` (`cmd/app.go`)
- [ ] **Rate limiting** — re-enable commented-out login (10/min) and signup (5/hr) limiters before production (`cmd/app.go`)

## Testing

### Tier 1 — Unit tests (no DB required)
- [ ] `internal/services/jwt_test.go` — generate/parse token, expiry, malformed token, wrong secret
- [ ] `internal/services/userService_test.go` — `titleCase()`, `isDuplicateKeyError()`, input normalization
- [ ] `internal/services/listingService_test.go` — `ErrNotOwner`, admin bypass, status defaults to "active"
- [ ] `internal/database/models/listing_test.go` — `StringSlice` JSON `Value`/`Scan` round-trip

### Tier 2 — Handler tests (Fiber `app.Test()`, mocked services)
- [ ] `api/auth/JWTmiddleware_test.go` — valid/missing/expired/malformed token → correct status codes
- [ ] `api/v1/handlers/login/login_test.go` — success, wrong password, nonexistent email, missing body
- [ ] `api/v1/handlers/users/create_test.go` — 201 success, 409 duplicate, 400 invalid body
- [ ] `api/v1/handlers/listings/list_test.go` — pagination defaults, filter params, response shape

### Tier 3 — Integration tests (requires PostgreSQL)
- [ ] `internal/database/repo/users_repo_test.go` — CRUD, password field omission, unique constraints
- [ ] `internal/database/repo/listings_repo_test.go` — `FindFiltered` with all filter combos, promoted sorting
