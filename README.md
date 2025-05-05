# TestTask - User Enrichment API

🚀 Test case for Junior Golang Developer: REST API for enriching user data (age, gender, nationality) with saving to PostgreSQL.

## 📋 Description

The service takes the user's full name, accesses an open API to get the estimated age, gender, and nationality, and stores the enriched data in the database.

Supports:
- CRUD operations
- Filtering by age, gender, nationality
- Pagination
- Swagger documentation

---

## 🛠 Technologies

- Go 1.22
- PostgreSQL
- GORM (ORM)
- Swaggo (Swagger)
- Standard net/http (router)

---

## 📦 Setup

1. Clone a repository:

```bash
git clone https://github.com/terfo1/TestTask.git
cd TestTask
```

2. Configure environment variables:

Create file `.env`:

```env
DBurl=postgres://postgres:password@localhost:5432/users_db?sslmode=disable
```

3. Start migrations:

```bash
go run cmd/main.go # AutoMigrate выполнит миграции
```

---

## 🚀 Run

```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080`.

### 📖 Swagger UI

Available at:

```
http://localhost:8080/swagger/index.html
```

---

## 📚 API Endpoints

| Method | Endpoint        | Description                   |
|-------|-----------------|----------------------------|
| GET   | `/user`         | Get users (filters + pagination) |
| POST  | `/createuser`   | Create a user       |
| PUT   | `/updateuser`   | Update user      |
| DELETE| `/deleteuser`   | Delete user       |

---

## 📝 Request example:

**User Creation:**

```http
POST /createuser
Content-Type: application/json

{
  "name": "Dmitriy",
  "surname": "Ushakov"
}
```

Answer:

```json
{
  "id": 1,
  "name": "Dmitriy",
  "surname": "Ushakov",
  "age": 32,
  "gender": "male",
  "nationality": "RU"
}
```
