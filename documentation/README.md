# File Vault

A secure file storage and management system with user authentication, file upload/download, access control, and metadata tracking.

## Features

- 🔐 User authentication (sign up, login, logout)
- 📁 Upload, download, and delete files
- 👀 File visibility toggling (public/private)
- 📝 Metadata tracking (file name, size, uploader, references)
- 🗄️ PostgreSQL database for users and file metadata
- 💻 React frontend with protected routes
- 🚀 REST API built with Go (Gin framework)

## Tech Stack

- **Backend:** Go (Gin) + GORM (Postgres ORM)
- **Frontend:** React + TypeScript + Vite
- **Database:** PostgreSQL
- **Storage:** Local disk (can be extended to S3, GCS, etc.)
- **Auth:** Cookie-based sessions (can be extended to JWT)

## Project Structure

```

file-vault/
│
├── backend/             # Go REST API
│   ├── cmd/             # Main entry point
│   ├── internal/        # Application logic
│   │   ├── api/         # HTTP handlers
│   │   ├── db/          # Database connection
│   │   ├── models/      # GORM models
│   │   ├── repository/  # DB access
│   │   └── services/    # Business logic
│   ├── migrations/      # Database migrations
│   └── go.mod
│
├── frontend/            # React + Vite app
│   ├── src/
│   │   ├── api/         # API calls
│   │   ├── components/  # Reusable components
│   │   ├── context/     # Auth context
│   │   ├── pages/       # App pages
│   │   └── routes/      # Route definitions
│   │   ├── styles/         
│   │   ├── hooks/         
│   └── package.json
│
└── docker-compose.yml   # Dev setup with Postgres + services

````

## Setup Instructions

### 1. Prerequisites
- Node.js (>=18) + npm
- Go (>=1.25)
- PostgreSQL

### 2. Clone the Repository

### 3. Start Backend

```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

### 4. Run Database (Postgres)
connect local postgres server via credentials (.env)

Apply migrations:

```bash
psql -h <host> -U <username> -d <database> -f <migration_file.sql>

```

### 5. Start Frontend

```bash
cd frontend
npm install
npm run dev
```

Visit: [http://localhost:5173](http://localhost:5173)

---

## API Endpoints (Backend)

| Method | Endpoint                    | Description           |
| ------ | --------------------------- | --------------------- |
| POST   | `/api/signup`               | Register new user     |
| POST   | `/api/login`                | Login user            |
| POST   | `/api/logout`               | Logout user           |
| GET    | `/api/me`                   | Get current user info |
| POST   | `/api/upload`                | Upload file           |
| GET    | `/api/files`                | List all files        |
| GET    | `/api/files/:id/download`            | Download file         |
| DELETE | `/api/files/:id/delete`            | Delete file           |
| PATCH  | `/api/files/:id/visibility` | Toggle visibility     |
| GET  | `/api/storage-stats` | user storage info    |


---

## Future Improvements

* ✅ GraphQL architecture for scalable apis
* ✅ File previews (images, PDFs)
* ✅ Role-based access control (admin, user)
* ✅ Cloud storage support (AWS S3, GCS, Azure Blob)

---
