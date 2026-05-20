# StepUp AI

![Go](https://img.shields.io/badge/Go-1.25-blue)
![gRPC](https://img.shields.io/badge/gRPC-enabled-green)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17-blue)
![Redis](https://img.shields.io/badge/Redis-7-red)
![NATS](https://img.shields.io/badge/NATS-2.10-orange)
![License](https://img.shields.io/badge/license-MIT-green)

University admission helper. Microservice-based backend written in Go.

# StepUp AI

University admission helper. Microservice-based backend written in Go.

## What it does

Helps students apply to foreign universities — analyzes admission chances, finds matching grants, generates preparation roadmaps and reviews essays using OpenAI.

## Stack

Go 1.25, gRPC, PostgreSQL, Redis, NATS, OpenAI API.

## Services

- **api-gateway** (port 8080) — REST API, JWT auth, routes requests to services via gRPC
- **user-service** (port 9001) — auth, profiles, password reset emails
- **university-service** (port 9002) — university and grant search
- **ai-service** (port 9003) — OpenAI integration for analysis and essay review

NATS is used for events between services (for example, user-service publishes "user.registered" and ai-service listens to it).

## Endpoints

Auth: register, login, logout, refresh, forgot-password, reset-password
Profile: get, update
Universities: search, details, save, list saved
Grants: search, save, list saved
AI: analyze admission, generate roadmap, review essay, match grants, history

Total: 20+ endpoints.

## How to run

You need PostgreSQL, Redis and NATS running locally.

```bash
psql -U postgres -c "CREATE DATABASE user_db;"
psql -U postgres -c "CREATE DATABASE university_db;"
psql -U postgres -c "CREATE DATABASE ai_db;"

nats-server &
redis-server &
```

Then run each service in its own terminal:

```bash
cd user-service && go run cmd/main.go
cd university-service && go run cmd/main.go
cd ai-service && go run cmd/main.go
cd api-gateway && go run cmd/main.go
```

Environment variables you'll need:
- `DATABASE_URL` for each service
- `JWT_SECRET` for user-service and api-gateway
- `OPENAI_API_KEY` for ai-service
- `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS` for user-service email
- `NATS_URL`, `REDIS_URL` where applicable

## Tests

```bash
cd user-service && go test ./...
cd university-service && go test ./...
```

Integration tests in university-service need a running PostgreSQL.

## Team

- Asylbek — api-gateway, user-service
- Ansar — ai-service
- Arai — university-service


Grant database-та жоқ (seed жасалмаған шығар немесе search фильтрі бос мән қабылдамайды). Grant-сыз жалғастырамыз — оның endpoint-тары аз. 

Енді **толық Postman/curl guide** жасап беремін. README-ге қою үшін.

Сақта осы UUID-дар:

```
USER_ID:     04d4f64c-b9cc-4364-859b-3bd91c2c38d2
MIT_ID:      749a812f-ff64-480f-878b-bb2a882319e8
STANFORD_ID: 16350e8c-2af2-40fa-80a3-fa5173208856
TOKEN:       (login кезінде алған token, $TOKEN-да тұр)
```

---

# TESTING GUIDE — 36 ENDPOINT



UUID

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"essay@test.com","password":"password123"}' \
  | python3 -c "import sys,json; print(json.load(sys.stdin)['access_token'])")

echo $TOKEN | cut -d. -f2 | base64 -d 2>/dev/null
```

---

# 🔵 USER SERVICE — localhost:9001

Proto: `proto/user/user.proto`

### 1. RegisterUser
```
Method: user.UserService/RegisterUser
{
  "email": "newuser@test.com",
  "password": "password123",
  "full_name": "New User"
}
```

### 2. LoginUser
```
Method: user.UserService/LoginUser
{
  "email": "newuser@test.com",
  "password": "password123"
}
```

### 3. LogoutUser
```
Method: user.UserService/LogoutUser
{
  "access_token": "PASTE_ACCESS_TOKEN_FROM_LOGIN"
}
```

### 4. RefreshToken
```
Method: user.UserService/RefreshToken
{
  "refresh_token": "PASTE_REFRESH_TOKEN_FROM_LOGIN"
}
```
⚠️ Refresh token — random hex string (eyJ... емес!). Login response-нан refresh_token field-ын ал.

### 5. GetUserProfile
```
Method: user.UserService/GetUserProfile
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2"
}
```

### 6. UpdateUserProfile
```
Method: user.UserService/UpdateUserProfile
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2",
  "gpa": 3.8,
  "sat_score": 1450,
  "ielts_score": 7.5,
  "country": "Kazakhstan"
}
```

### 7. SendPasswordResetEmail
```
Method: user.UserService/SendPasswordResetEmail
{
  "email": "essay@test.com"
}
```

### 8. ResetPassword
```
Method: user.UserService/ResetPassword
{
  "token": "RESET_TOKEN",
  "new_password": "newpassword123"
}
```

### 9. DeleteUser
```
Method: user.UserService/DeleteUser
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2"
}
```
⚠️ Осы шақырғаннан кейін user өшеді — басқа тестке кедергі болуы мүмкін.

### 10. GetUserByID
```
Method: user.UserService/GetUserByID
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2"
}
```

### 11. ChangePassword
```
Method: user.UserService/ChangePassword
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2",
  "old_password": "password123",
  "new_password": "newpass456"
}
```

### 12. VerifyEmail
```
Method: user.UserService/VerifyEmail
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2",
  "verification_code": "123456"
}
```

---

# 🟡 UNIVERSITY SERVICE — localhost:9002

Proto: `proto/university/university.proto`

### 1. SearchUniversities
```
Method: university.UniversityService/SearchUniversities
{
  "country": "USA",
  "university_type": "",
  "min_acceptance_rate": 0,
  "max_acceptance_rate": 100
}
```

### 2. GetUniversityDetails
```
Method: university.UniversityService/GetUniversityDetails
{
  "university_id": "749a812f-ff64-480f-878b-bb2a882319e8"
}
```

### 3. SaveUniversity
```
Method: university.UniversityService/SaveUniversity
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2",
  "university_id": "749a812f-ff64-480f-878b-bb2a882319e8"
}
```

### 4. GetSavedUniversities
```
Method: university.UniversityService/GetSavedUniversities
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2"
}
```

### 5. SearchGrants
```
Method: university.UniversityService/SearchGrants
{
  "country": "USA"
}
```

### 6. SaveGrant
```
Method: university.UniversityService/SaveGrant
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2",
  "grant_id": "GRANT_UUID"
}
```

### 7. GetSavedGrants
```
Method: university.UniversityService/GetSavedGrants
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2"
}
```

### 8. RemoveSavedUniversity
```
Method: university.UniversityService/RemoveSavedUniversity
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2",
  "university_id": "749a812f-ff64-480f-878b-bb2a882319e8"
}
```

### 9. RemoveSavedGrant
```
Method: university.UniversityService/RemoveSavedGrant
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2",
  "grant_id": "GRANT_UUID"
}
```

### 10. GetGrantDetails
```
Method: university.UniversityService/GetGrantDetails
{
  "grant_id": "GRANT_UUID"
}
```

### 11. ListUniversitiesByCategory
```
Method: university.UniversityService/ListUniversitiesByCategory
{
  "category": "reach"
}
```
Other options: `"target"`, `"safety"`

### 12. GetUniversityStatistics
```
Method: university.UniversityService/GetUniversityStatistics
{}
```

---

# 🟢 AI SERVICE — localhost:9003

Proto: `proto/ai/ai.proto`

### 1. AnalyzeAdmissionChances
```
Method: ai.AIService/AnalyzeAdmissionChances
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2",
  "university_id": "749a812f-ff64-480f-878b-bb2a882319e8",
  "gpa": 3.8,
  "sat_score": 1450,
  "ielts_score": 7.5
}
```

### 2. GenerateRoadmap
```
Method: ai.AIService/GenerateRoadmap
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2",
  "target_university_id": "749a812f-ff64-480f-878b-bb2a882319e8",
  "months_until_application": 12
}
```

### 3. ReviewEssay
```
Method: ai.AIService/ReviewEssay
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2",
  "essay_text": "My passion for computer science started in high school when I built my first program. I have always been fascinated by how technology can solve real-world problems. I want to study at MIT because of their world-class research in AI and machine learning. Their interdisciplinary approach perfectly aligns with my goal of combining computer science with neuroscience to understand human cognition.",
  "university_name": "MIT",
  "program_name": "Computer Science",
  "word_limit": 500
}
```

### 4. MatchGrants
```
Method: ai.AIService/MatchGrants
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2",
  "country": "Kazakhstan",
  "gpa": 3.9,
  "achievements": ["olympiad winner", "research paper published", "volunteer work"]
}
```

### 5. GetAnalysisHistory
```
Method: ai.AIService/GetAnalysisHistory
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2"
}
```

### 6. GetRoadmapByUserID
```
Method: ai.AIService/GetRoadmapByUserID
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2"
}
```

### 7. GetEssayReviewsByUserID
```
Method: ai.AIService/GetEssayReviewsByUserID
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2"
}
```

### 8. GetEssayReviewByID
```
Method: ai.AIService/GetEssayReviewByID
{
  "review_id": "REVIEW_UUID_FROM_REVIEWESSAY_RESPONSE"
}
```

### 9. DeleteAnalysis
```
Method: ai.AIService/DeleteAnalysis
{
  "analysis_id": "ANALYSIS_UUID_FROM_ANALYZE_RESPONSE"
}
```

### 10. CompareUniversities
```
Method: ai.AIService/CompareUniversities
{
  "university_id_one": "749a812f-ff64-480f-878b-bb2a882319e8",
  "university_id_two": "16350e8c-2af2-40fa-80a3-fa5173208856"
}
```

### 11. GetRecommendedUniversities
```
Method: ai.AIService/GetRecommendedUniversities
{
  "user_id": "04d4f64c-b9cc-4364-859b-3bd91c2c38d2",
  "gpa": 3.8
}
```

### 12. GetAdmissionAnalysisByID
```
Method: ai.AIService/GetAdmissionAnalysisByID
{
  "analysis_id": "ANALYSIS_UUID_FROM_ANALYZE_RESPONSE"
}
```

