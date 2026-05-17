# StepUp AI — API Documentation

Base URL: `http://localhost:8080`

## Authentication

### POST /auth/register

Register a new user.

**Request body:**
\`\`\`json
{
  "email": "user@example.com",
  "password": "securepassword123",
  "full_name": "John Doe"
}
\`\`\`

**Response 201:**
\`\`\`json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "abc123..."
}
\`\`\`

**Errors:**
- 400 — Invalid request body
- 409 — User already exists

---

### POST /auth/login

Login with email and password.

**Request body:**
\`\`\`json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
\`\`\`

**Response 200:**
\`\`\`json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "abc123..."
}
\`\`\`

---

### POST /auth/refresh

Get a new access token using a refresh token.

**Request body:**
\`\`\`json
{
  "refresh_token": "abc123..."
}
\`\`\`

---

### POST /auth/logout

Logout the current user. Requires Authorization header.

---

### POST /auth/forgot-password

Send a password reset email.

**Request body:**
\`\`\`json
{ "email": "user@example.com" }
\`\`\`

---

### POST /auth/reset-password

Reset password using a token from email.

---

## Profile

All profile endpoints require `Authorization: Bearer <token>`.

### GET /profile

Get the current user profile.

### PUT /profile

Update profile.

**Request body:**
\`\`\`json
{
  "gpa": 3.9,
  "sat_score": 1500,
  "ielts_score": 8.0,
  "country": "Kazakhstan"
}
\`\`\`

---

## Universities

### GET /universities?country=USA

Search universities with optional filters.

### GET /universities/:id

Get details of a specific university.

### POST /universities/save

Save a university to favorites.

### GET /universities/saved

Get all saved universities.

---

## Grants

### GET /grants

Search grants.

### POST /grants/save

Save a grant.

### GET /grants/saved

Get saved grants.

---

## AI

### POST /ai/analyze

Analyze admission chances.

**Request body:**
\`\`\`json
{
  "university_id": "uuid",
  "gpa": 3.8,
  "sat_score": 1450,
  "ielts_score": 7.5
}
\`\`\`

### POST /ai/roadmap

Generate a preparation roadmap.

### POST /ai/essay

Review an application essay.

### POST /ai/grants/match

Match grants to user profile.

### GET /ai/history

Get analysis history.