# FolksArt IAM Identity Governance Suite - Go & MySQL Back-End

This is the production-ready high-performance Golang backend microservice for the **FolksArt IAM Identity Governance Suite**.

It implements secure storage for corporate identity directories (Subject Principals) and operational audit trails on top of a relational **MySQL RDBMS**.

---

## 🛠️ Technology Stack & Architecture

- **Runtime**: Go (Golang) 1.21+
- **RDBMS**: MySQL 5.7 / 8.0+
- **Routing**: Lightweight native Go `net/http` router (No third-party packages required, ensuring high speed and zero dependencies)
- **Database Driver**: Official `github.com/go-sql-driver/mysql`
- **CORS Support**: Custom middleware enabled to support integration with multi-domain clients.
- **Deployments**: Fully multi-stage dockerized execution (optimizes runtime payload size to ~12MB)

---

## 📂 Design Structure

```
/backend-golang
│
├── config/
│   └── db.go             # Verifies connection & structures DB pooling
├── models/
│   └── models.go         # Strongly typed model representations of User and Audit
├── handlers/
│   ├── user.go           # Subject Principal enrollments, attributes, delete & CSV
│   └── audit.go          # Querying audit trails and manual system entries
├── routes/
│   └── routes.go         # Endpoint map and HTTP option cors controller queries
├── schema.sql            # Script to run tables and enterprise seeds 
├── go.mod                # Go module descriptors file
├── Dockerfile            # Compile & compact multistage run commands
└── README.md             # Standard technical documentation
```

---

## 💾 Relational Data Schema Map

The table structures are designed as follows:

### 1. `iam_users` (Subject Principals Table)
* `id` (VARCHAR 50, PRIMARY KEY) - Secure client Reference ID
* `name` (VARCHAR 100) - Real corporate name
* `username` (VARCHAR 50, UNIQUE) - LDAP/ActiveDirectory matching username
* `email` (VARCHAR 100) - Corporate electronic address
* `phone` (VARCHAR 30) - Enterprise phone index
* `role` (VARCHAR 50) - Evaluates access: `Administrator` | `Security Officer` | `End User`
* `status` (VARCHAR 50) - Status index: `Active` | `Banned` | `Deactivated`
* `kyc_status` (VARCHAR 50) - Compliance: `Pending` | `Verified` | `Failed` | `Suspicious`
* `department` (VARCHAR 100) - E.g. Security Operations, Legal, Engineering
* `risk_score` (INT) - Score calculated by threat analyzer vectors (0 - 100)
* `mfa_enabled` (BOOLEAN) - Guard shield indicator
* `created_at` (TIMESTAMP) - Verification timestamp

### 2. `audit_logs` (Administrative Logs Table)
* `id` (VARCHAR 50, PRIMARY KEY) - Unique security event UUID
* `timestamp` (TIMESTAMP) - Time of operator action
* `actor` (VARCHAR 50) - Username representing operator
* `action` (TEXT) - Description of state modification
* `target` (VARCHAR 100) - Affected Identity/System principal
* `severity` (VARCHAR 30) - Severity index: `Critical` | `High` | `Medium` | `Low`

---

## ⚙️ Configuration (Environment Variables)

Build and run parameters can be customised using environment parameters:

| Variable | Default Value | Description |
|---|---|---|
| `PORT` | `8080` | Network port for listening |
| `DB_USER` | `root` | MySQL access username |
| `DB_PASSWORD` | `secret` | MySQL access password |
| `DB_HOST` | `127.0.0.1` | MySQL Database IP/Hostname |
| `DB_PORT` | `3306` | MySQL Port |
| `DB_NAME` | `folksart_iam` | Target Database name |

---

## 🚀 Speed-Start Instructions

### Run Standalone
1. **Prepare MySQL Database**
   Log into your MySQL server and run the script located in `schema.sql`:
   ```bash
   mysql -u root -p < schema.sql
   ```

2. **Launch Back-End Server**
   Navigate to `/backend-golang` and run:
   ```bash
   export DB_USER="root"
   export DB_PASSWORD="your_password"
   export DB_HOST="localhost"
   export DB_PORT="3306"
   export DB_NAME="folksart_iam"
   export PORT="8080"
   
   go run main.go
   ```

### Run using Docker Container
Build and start the container directly:
```bash
docker build -t folksart-iam-backend .
docker run -d -p 8080:8080 --name iam-backend folksart-iam-backend
```

---

## 📡 API Testing Examples

### 1. List Users with Search and Filters
```bash
curl -X GET "http://localhost:8080/api/v1/users?search=Connor&role=Administrator"
```

### 2. Enroll a New Identity (Subject Principal)
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Kyle Reese",
    "email": "kreese@cyberdyne.org",
    "role": "Security Officer",
    "department": "Security Operations",
    "operator": "sarah_connor"
  }'
```

### 3. Patch Identity Attributes & KYC Risk score
```bash
curl -X PATCH http://localhost:8080/api/v1/users/usr-8821 \
  -H "Content-Type: application/json" \
  -d '{
    "status": "Active",
    "riskScore": 15,
    "kycStatus": "Verified",
    "mfaEnabled": true,
    "operator": "marcus_wright"
  }'
```

### 4. Download Complete Directory as Secure CSV
```bash
curl -o identity_export.csv "http://localhost:8080/api/v1/export/csv"
```

### 5. Fetch Audit Trail Console logs
```bash
curl -X GET "http://localhost:8080/api/v1/audit-logs?severity=Critical"
```
