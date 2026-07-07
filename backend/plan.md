# Rencana Implementasi Backend Go (Golang) Gin GORM - SmartFarming

Rencana ini mencakup arsitektur, teknologi, struktur folder, dan langkah-langkah pembuatan backend SmartFarming menggunakan **Golang, Gin, dan GORM** dengan database **PostgreSQL, MongoDB, dan Redis**.

---

## 1. Teknologi & Pustaka (Dependencies)
- **Framework HTTP**: `github.com/gin-gonic/gin`
- **ORM (PostgreSQL)**: `gorm.io/gorm` & `gorm.io/driver/postgres`
- **MongoDB Driver**: `go.mongodb.org/mongo-driver/mongo`
- **Redis Driver**: `github.com/redis/go-redis/v9`
- **ID UUID**: `github.com/google/uuid`
- **Hash Password**: `golang.org/x/crypto/bcrypt`
- **Token Auth JWT**: `github.com/golang-jwt/jwt/v5`
- **Environment Management**: `github.com/joho/godotenv`
- **API Documentation**: `github.com/swaggo/swag`, `github.com/swaggo/gin-swagger`, `github.com/swaggo/files`

---

## 2. Struktur Proyek (Project Structure)
Struktur folder mengikuti pola Clean Architecture atau Layered Architecture:
```text
backend/
├── main.go               # Entry point aplikasi
├── config/
│   ├── config.go             # Load env & setup struct konfigurasi
│   ├── postgres.go           # Koneksi & auto-migrate PostgreSQL
│   ├── mongodb.go            # Koneksi MongoDB
│   └── redis.go              # Koneksi Redis
├── model/
│   ├── base_model.go         # Model dasar (ID UUID, Audit fields)
│   └── user.go               # Entity database User
├── dto/
│   ├── auth_dto.go           # Request & Response untuk Login/Register/Me
│   ├── user_dto.go           # Request & Response untuk User CRUD
│   └── common_dto.go         # Struct pagination & generic response
├── repository/
│   └── user_repository.go    # Query DB ke PostgreSQL (GORM)
├── service/
│   ├── auth_service.go       # Logika bisnis register, login, jwt
│   └── user_service.go       # Logika bisnis CRUD user & pagination
├── handler/
│   ├── auth_handler.go       # Controller endpoint Auth
│   └── user_handler.go       # Controller endpoint User
├── middleware/
│   ├── auth_middleware.go    # Validasi JWT & inject user ke context
│   └── idor_middleware.go    # Proteksi IDOR (Insecure Direct Object Reference)
├── routes/
│   └── routes.go             # Registrasi semua route Gin
├── docs/                     # Swagger API docs (Auto-generated)
├── .env.development          # Environment local development
├── .env.production           # Environment production
├── go.mod
└── go.sum
```

---

## 3. Detail Implementasi Fitur Utama

### A. Base Model UUID & Audit Fields
`BaseModel` akan memiliki hook otomatis dari GORM untuk mengisi audit fields (`CreatedBy`, `UpdatedBy`, `DeletedBy`) dari user context:
```go
type BaseModel struct {
    ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    CreatedAt time.Time      `json:"createdAt"`
    UpdatedAt time.Time      `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
    CreatedBy *uuid.UUID     `gorm:"type:uuid" json:"createdBy,omitempty"`
    UpdatedBy *uuid.UUID     `gorm:"type:uuid" json:"updatedBy,omitempty"`
    DeletedBy *uuid.UUID     `gorm:"type:uuid" json:"deletedBy,omitempty"`
}
```

### B. Proteksi IDOR (Insecure Direct Object Reference)
- Setiap request yang mencoba mengakses `/users/:id` akan diverifikasi melalui middleware atau service layer.
- Aturan akses:
  - User biasa hanya boleh mengakses/mengubah data milik **mereka sendiri** (`current_user_id == target_user_id`).
  - User dengan role `Admin` diperbolehkan untuk mengakses/mengubah data user lain.

### C. DTO dengan Pagination & Search
- Request pagination dikirim melalui query param: `page` (default 1), `limit` (default 10), dan `search` (untuk mencari berdasarkan nama/email).
- Response pagination mengembalikan metadata:
  ```json
  {
    "data": [...],
    "meta": {
      "currentPage": 1,
      "limit": 10,
      "totalRecords": 45,
      "totalPages": 5
    }
  }
  ```

---

## 4. Rencana Langkah Kerja (Roadmap)
1. **Inisialisasi Project**: `go mod init`, install library yang dibutuhkan.
2. **Configuration & Database Connections**: Buat setup `.env`, konfigurasi GORM PostgreSQL, MongoDB, dan Redis.
3. **Database Model & Migrations**: Buat `BaseModel` dan model `User`, lalu jalankan auto-migrate PostgreSQL.
4. **Implementasi Repository**: Query SQL CRUD dan pagination.
5. **Implementasi Service**: Bisnis logik registrasi, enkripsi password dengan bcrypt, generate JWT token, dll.
6. **Implementasi Middleware**: JWT Auth Validator & IDOR Protection middleware.
7. **Implementasi Handler & Router**: Setup Gin Router, register endpoint `/register`, `/login`, `/me`, dan `/users`.
8. **Integrasi Swagger**: Install swag CLI, anotasi API, compile docs.
9. **Verifikasi & Pengujian**: Pengujian unit test dan testing API endpoint (Postman/Curl).
