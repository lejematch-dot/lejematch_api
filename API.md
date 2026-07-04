# Lejematch API

Base URL: `/api/v1`

Authentication is done via a **Bearer token** in the `Authorization` header:
```
Authorization: Bearer <token>
```

---

## Auth

### Login
`POST /auth/login`

**Public**

**Request**
```json
{
  "email": "string",
  "password": "string"
}
```

**Response `200`**
```json
{
  "token": "string",
  "userID": 1
}
```

**Errors**
| Status | Reason |
|--------|--------|
| `400` | Invalid request body |
| `401` | Invalid email or password |
| `500` | Internal error |

---

## Users

### Create user
`POST /users`

**Public**

**Request**
```json
{
  "FirstName": "string",
  "LastName": "string",
  "Email": "string",
  "Phone": "string",
  "Password": "string",
  "City": "string",
  "ImageURL": "string"
}
```

**Response `201`**
```json
{
  "id": 1,
  "createdAt": "2024-01-01T00:00:00Z"
}
```

**Errors**
| Status | Reason |
|--------|--------|
| `400` | Invalid request body |
| `409` | Email or phone already in use |
| `500` | Internal error |

---

### Get user
`GET /users/:id`

**JWT required** — own record or admin only

**Response `200`**
```json
{
  "ID": 1,
  "CreatedAt": "2024-01-01T00:00:00Z",
  "UpdatedAt": "2024-01-01T00:00:00Z",
  "FirstName": "string",
  "LastName": "string",
  "Email": "string",
  "Phone": "string",
  "IsAdmin": false,
  "IsActive": true
}
```

**Errors**
| Status | Reason |
|--------|--------|
| `400` | Invalid ID |
| `403` | Forbidden |
| `404` | User not found |
| `500` | Internal error |

---

### Update user
`PATCH /users/:id`

**JWT required** — own record or admin only

All fields are optional.

**Request**
```json
{
  "FirstName": "string",
  "LastName": "string",
  "Email": "string",
  "Phone": "string"
}
```

**Response `204`** — no body

**Errors**
| Status | Reason |
|--------|--------|
| `400` | Invalid ID or request body |
| `403` | Forbidden |
| `409` | Email or phone already in use |
| `500` | Internal error |

---

### Delete user
`DELETE /users/:id`

**JWT required** — own record or admin only

Also deletes the associated profile (cascade).

**Response `200`**
```json
1
```

**Errors**
| Status | Reason |
|--------|--------|
| `400` | Invalid ID |
| `403` | Forbidden |
| `404` | User not found |
| `500` | Internal error |

---

### Update password
`PUT /users/:id/password`

**JWT required** — own record or admin only

**Request**
```json
{
  "CurrentPassword": "string",
  "NewPassword": "string"
}
```

**Response `204`** — no body

**Errors**
| Status | Reason |
|--------|--------|
| `400` | Invalid ID or request body |
| `403` | Forbidden |
| `404` | User not found |
| `422` | Current password is wrong |
| `422` | New password too short |
| `500` | Internal error |

---

## Profiles

### Get profile
`GET /users/:id/profile`

**Public**

**Response `200`**
```json
{
  "displayName": "string",
  "bio": "string",
  "city": "string",
  "imageURL": "string"
}
```

**Errors**
| Status | Reason |
|--------|--------|
| `400` | Invalid ID |
| `404` | Profile not found |
| `500` | Internal error |

---

### Update profile
`PATCH /users/:id/profile`

**JWT required** — own record or admin only

All fields are optional.

**Request**
```json
{
  "DisplayName": "string",
  "Bio": "string",
  "City": "string",
  "ImageURL": "string"
}
```

**Response `204`** — no body

**Errors**
| Status | Reason |
|--------|--------|
| `400` | Invalid ID or request body |
| `403` | Forbidden |
| `500` | Internal error |

---

## Listings

### List listings
`GET /listings`

**Public**

**Query parameters**
| Param | Type | Description |
|-------|------|-------------|
| `page` | int | Page number (default: 1) |
| `city` | string | Filter by city |
| `roomType` | string | `private`, `shared`, or `apartment` |
| `minPrice` | int | Minimum price (DKK/month) |
| `maxPrice` | int | Maximum price (DKK/month) |

**Response `200`**
```json
{
  "data": [
    {
      "ID": 1,
      "CreatedAt": "2024-01-01T00:00:00Z",
      "UserID": 1,
      "Title": "string",
      "Description": "string",
      "Price": 5000,
      "City": "string",
      "Zip": "string",
      "Area": "string",
      "RoomType": "private",
      "Status": "active",
      "AvailableFrom": "2024-06-01",
      "Images": ["https://..."]
    }
  ],
  "page": 1,
  "pageSize": 20,
  "total": 42,
  "totalPages": 3
}
```

**Errors**
| Status | Reason |
|--------|--------|
| `500` | Internal error |

---

### List listings by user
`GET /users/:id/listings`

**Public**

Returns all listings (all statuses) for the given user, ordered newest first.

**Response `200`**
```json
[
  {
    "ID": 1,
    "CreatedAt": "2024-01-01T00:00:00Z",
    "UserID": 1,
    "Title": "string",
    "Description": "string",
    "Price": 5000,
    "City": "string",
    "Zip": "string",
    "Area": "string",
    "RoomType": "private",
    "Status": "active",
    "AvailableFrom": "2024-06-01",
    "Images": ["https://..."]
  }
]
```

**Errors**
| Status | Reason |
|--------|--------|
| `400` | Invalid ID |
| `500` | Internal error |

---

### Get listing
`GET /listings/:id`

**Public**

**Response `200`**
```json
{
  "ID": 1,
  "CreatedAt": "2024-01-01T00:00:00Z",
  "UserID": 1,
  "Title": "string",
  "Description": "string",
  "Price": 5000,
  "City": "string",
  "Zip": "string",
  "Area": "string",
  "RoomType": "private",
  "Status": "active",
  "AvailableFrom": "2024-06-01",
  "Images": ["https://..."]
}
```

**Errors**
| Status | Reason |
|--------|--------|
| `400` | Invalid ID |
| `404` | Listing not found |
| `500` | Internal error |

---

### Create listing
`POST /listings`

**JWT required**

**Request**
```json
{
  "Title": "string",
  "Description": "string",
  "Price": 5000,
  "City": "string",
  "Zip": "string",
  "Area": "string",
  "RoomType": "private",
  "AvailableFrom": "2024-06-01",
  "Images": ["https://..."]
}
```

`RoomType` must be one of: `private`, `shared`, `apartment`

**Response `201`**
```json
{
  "id": 1,
  "createdAt": "2024-01-01T00:00:00Z"
}
```

**Errors**
| Status | Reason |
|--------|--------|
| `400` | Invalid request body |
| `500` | Internal error |

---

### Update listing
`PATCH /listings/:id`

**JWT required** — own listing or admin only

All fields are optional.

**Request**
```json
{
  "Title": "string",
  "Description": "string",
  "Price": 5000,
  "City": "string",
  "Zip": "string",
  "Area": "string",
  "RoomType": "private",
  "Status": "rented",
  "AvailableFrom": "2024-06-01",
  "Images": ["https://..."]
}
```

`Status` must be one of: `active`, `rented`, `archived`

**Response `204`** — no body

**Errors**
| Status | Reason |
|--------|--------|
| `400` | Invalid ID or request body |
| `403` | Forbidden |
| `404` | Listing not found |
| `500` | Internal error |

---

### Delete listing
`DELETE /listings/:id`

**JWT required** — own listing or admin only

**Response `204`** — no body

**Errors**
| Status | Reason |
|--------|--------|
| `400` | Invalid ID |
| `403` | Forbidden |
| `404` | Listing not found |
| `500` | Internal error |

---

## Health

### Health check
`GET /health`

**Public**

**Response `200`**
```json
{
  "status": "ok"
}
```
