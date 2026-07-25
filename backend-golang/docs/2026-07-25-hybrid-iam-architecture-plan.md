# Rencana Arsitektur & Roadmap: Hybrid IAM Governance (Workforce & Customer CIAM)

**Tanggal Pembaharuan**: 2026-07-25  
**Dokumen Referensi**: `docs/2026-07-25-hybrid-iam-architecture-plan.md`  
**Tujuan**: Merancang transformasi layanan **`folksart-be`** dari platform pengamanan identitas dasar menjadi **Hybrid IAM Suite** yang melayani pengguna internal (Workforce/Employees) sekaligus pengguna eksternal (Customer/Consumers), lengkap dengan analisis kelemahan sistem dan protokol pelacakan progres pengerjaan (*execution checkpoint protocol*).

---

## 🚩 PROTOKOL CHECKPOINT & EXECUTION FLAG (Antisipasi Token Habis / Sesi Rejeksi)

> [!IMPORTANT]
> **Panduan Resume untuk Agen AI & Developer Berikutnya**  
> Jika dalam pembicaraan berikutnya token batas (*token context window*) habis atau sesi teregenerasi, periksa **Tabel Execution Checkpoint** dan **Daftar Rencana Implementasi** di dokumen ini.  
> - **Jangan menulis ulang** kode yang sudah berstatus **`[x] COMPLETED`**.  
> - Lanjutkan secara eksklusif dari item pertama yang berstatus **`[/] IN PROGRESS`** atau **`[ ] PENDING`**.
> - Setelah menyelesaikan tugas baru, perbarui flag status di file ini menggunakan tool `replace_file_content` agar riwayat konotasi sistem senantiasa tersinkronisasi.

### Tabel Execution Checkpoint (Status Riwayat Layanan)

| Kode Tahap | Nama Paket / Fitur | Status Pengerjaan | Catatan & Referensi Verifikasi |
| :---: | :--- | :---: | :--- |
| **THP-00** | Refactoring JSON Tags ke `camelCase` | **[x] COMPLETED** | Seluruh package `dto` dan `security` telah terkonfirmasi mengadopsi standar camelCase. |
| **THP-01** | Quick Wins & Build Failure Resolution | **[x] COMPLETED** | Perbaikan impor `os` tak terpakai di `config/db.go` & pembuatan generator `docs/docs.go`. Modul berganti ke `folksart-be/backend-golang`. |
| **THP-02** | Security Hardening & RBAC Route Protection | **[x] COMPLETED** | Penghapusan hardcoded secrets di `jwt.go`, implementasi `JWTAuth` & `RBAC("Administrator")` middleware pada rute API di `routes.go`. |
| **THP-03** | Usecase Testing Suite & Global Rule Validation | **[x] COMPLETED** | Unit test di `usecases/`, `security/`, dan `middleware/`. Validasi Aturan Global (pencatatan *Alert/Audit Notification* pada tiap modifikasi *Update/Delete* database) terpenuhi. |
| **THP-04** | **Transformasi Skema Hybrid (Workforce vs CIAM)** | **[x] COMPLETED** | Pembuatan pemisahan `UserType`, endpoint *Self-Service Sign-up* untuk Customer, dan pendaftaran profil otomatis. |
| **THP-05** | **Real Database Persistence & Unit of Work (TxManager)** | **[ ] PENDING** | Pembuatan skema SQL migrasi (DDL) dan refaktorisasi usecases *mock* (Role, Policy, KYC, dll.) ke integrasi SQL nyata berlandaskan transaksi ACID. |
| **THP-06** | **Automating Global Rule Compliance & Async Notification Worker** | **[ ] PENDING** | Implementasi *Repository Decorator* untuk otomatisasi Alerting Event pada tiap query Update/Delete via Async Event Bus. |
| **THP-07** | **Standardisasi Respons API untuk Integrasi React Frontend** | **[ ] PENDING** | Penambahan obyek metadata paginasi (*Paging Meta*) dan Katalog Kode Error Bisnis terstruktur. |

---

## 🔍 ANALISIS KELEMAHAN (WEAKNESS REVIEW) SISTEM SAAT INI

Berdasarkan evaluasi arsitektur Clean Architecture pada versi service yang terverifikasi (pasca Tahap 3), terdapat 4 kelemahan fundamental yang dapat menjadi batu sandungan bagi target **Hybrid IAM (Internal & External Users)**:

### 1. Kelemahan Model & Isolasi Identitas (Belum Dukung Hybrid User Type)
- **Kondisi**: Entitas [user.go](file:///Users/a2275/Artech/golang/folksart-be/backend-golang/internal/domain/user.go#L8-L21) (`IAMUser`) saat ini berorientasi tunggal pada karyawan internal/workforce (mengandung field wajib `Department` dan `RiskScore`, serta rute pendaftaran `POST /users` yang dikawal ketat oleh RBAC `"Administrator"`).
- **Dampak**: Pengguna eksternal / kustomer (CIAM) yang ingin mendaftar ke aplikasi tidak memiliki jalur mandiri (*Self-Service Sign-Up*). Selain itu, bercampurnya tipe pengguna tanpa pembeda (*flag isolation*) seperti `UserType: "workforce" | "customer"` berisiko membuat kustomer umum terkena kebijakan tata kelola ketat yang semestinya hanya untuk staf kantor.

### 2. Kelemahan Persistensi & Arsitektur Data ("Stub / Mock Architecture")
- **Kondisi**: Hanya modul `User` dan `Audit` yang terkoneksi ke driver database `*sql.DB`. Modul kritis lainnya seperti `RoleUsecase`, `PolicyUsecase`, `AccessRequest`, `KYC`, `Notification`, dan `Report` hanya berisi string *dummy mock return* dan sama sekali tidak tersimpan di database.
- **Dampak**: Layanan belum bersiap untuk lingkungan produksi sesungguhnya. Tanpa tabel relasi nyata (`user_roles`, `policies_rules`), keputusan otorisasi dan peninjauan regulasi KYC bersifat hampa (*stateless mock*).

### 3. Ketiadaan Transaksionalitas (ACID) & Kerentanan Manual pada Aturan Global
- **Kondisi**: Meskipun seluruh tes pada `UserUsecase` telah memvalidasi kepatuhan pada **Aturan Global (*Global Rules*): "Setiap query Update & Delete wajib mengirimkan Notifikasi/Alert"**, implementasi saat ini dipaksakan lewat penulisan kode berurutan (manual) pada Usecase: memanggil `.Update()` di repo satu, lalu memanggil `.Store()` di `auditRepo`.
- **Dampak**:
  - **Data Tidak Konsisten**: Jika penyimpanan log audit atau pengiriman notifikasi gagal, data utama yang tersimpan tidak di-*rollback* karena ketiadaan mekanisme *Unit-of-Work / DB Transaction* (`sql.Tx`).
  - **Human Error / Omission**: Developer di masa depan yang menambah fitur di usecase lain berisiko lupa menyisipkan perintah simpan ke audit log, sehingga melanggar aturan sistem global.

### 4. Kesiapan Integrasi ke Frontend React (Missing Pagination Meta & Structured Errors)
- **Kondisi**: Respons baku pada [response.go](file:///Users/a2275/Artech/golang/folksart-be/backend-golang/httputil/response.go) mengembalikan paket obyek tunggal atau array merta.
- **Dampak**: Jika dipasangkan dengan antarmuka **React Clean Architecture**, modul presentasi antarmuka akan kesulitan menyewa komponen *Pagination Table* karena tidak tersedianya metadata paginasi (seperti `currentPage`, `totalItems`, `totalPages`). Selain itu, penanganan galat masih mengandalkan string pesan teks mentah alih-alih kode katalog konsisten (seperti `ERR_USER_BANNED` atau `ERR_KYC_REQUIRED`).

---

## 🗺️ RENCANA IMPLEMENTASI DETAILED (ROADMAP)

Berikut adalah panduan eksekusi teknis untuk menyelesaikan kelemahan di atas sekaligus mendedikasikan aplikasi sebagai **Hybrid Workforce & Customer IAM Engine**.

### FASE 4: Transformasi Skema Hybrid (Workforce vs CIAM Integration)

#### 4.1 Modifikasi Entitas Domain
- [x] Update [internal/domain/user.go](file:///Users/a2275/Artech/golang/folksart-be/backend-golang/internal/domain/user.go) untuk menyisipkan atribut pendukung hybrid:
  ```go
  type IAMUser struct {
      // ... field eksisting ...
      UserType   string `json:"userType"`   // Nilai: "workforce" atau "customer"
      TenantID   string `json:"tenantId"`   // Untuk pemisahan instansi atau grup kustomer
      IsVerified bool   `json:"isVerified"` // Validasi instan KYC/Email
  }
  ```

#### 4.2 Pembuatan Jalur Customer Onboarding (CIAM Self-Service Register)
- [x] Buat metode baru di antarmuka `AuthUsecase` dan `AuthHandler`: `RegisterCustomer(ctx, name, email, password, phone)`.
- [x] Buka rute publik di [routes/routes.go](file:///Users/a2275/Artech/golang/folksart-be/backend-golang/routes/routes.go) pada grup `/auth`:
  ```go
  auth.Post("/register", hc.AuthHandler.RegisterCustomer)
  ```
- [x] Pastikan alur `RegisterCustomer` memberikan `UserType = "customer"`, `Role = "VerifiedCustomer"`, dan menetapkan `KYCStatus = "Pending"` untuk dilanjutkan pada alur verifikasi identitas Mandiri.

---

### FASE 5: Real Database Persistence & Unit of Work (TxManager)

#### 5.1 Skema Migrasi Database (DDL)
- [ ] Buat direktori dan file spesifikasi migrasi database pada `migrations/0001_initial_iam_schema.sql` dengan kelenturan struktur:
  - `users` (dengan tambahan kolom `user_type`, `tenant_id`, `is_verified`).
  - `roles`, `permissions`, dan relasi m-2-m `role_permissions` & `user_roles`.
  - `policies` dan `access_requests` untuk menampung approval tiket akses internal.
  - `kyc_records` untuk mengawasi status upload identitas kustomer eksternal.

#### 5.2 Implementasi Transaction Manager / Unit of Work
- [ ] Buat abstraksi pengelola transaksi di [internal/ctxutil/tx.go](file:///Users/a2275/Artech/golang/folksart-be/backend-golang/internal/ctxutil) berbasis `*sql.DB` dan `*sql.Tx`:
  ```go
  type TxManager interface {
      WithTx(ctx context.Context, fn func(txCtx context.Context) error) error
  }
  ```
- [ ] Implementasikan `WithTx` pada `UserUsecase` saat melakukan pendaftaran ataupun penghapusan agar operasi pencatatan profil utama dan audit trail tergembok dalam 1 transaksi atomik (ACID): jika satu gagal, maka *auto-rollback*.

---

### FASE 6: Otomatisasi Aturan Global (Repository Decorator) & Async Notification

#### 6.1 Interceptive Repository Decorator (Kepatuhan Global Rules)
- [ ] Buat *wrapper/decorator pattern* pada lapisan repository di [internal/repositories/audit_decorator.go](file:///Users/a2275/Artech/golang/folksart-be/backend-golang/internal/repositories) untuk membungkus repository utama (seperti `UserRepository`, `RoleRepository`, `PolicyRepository`).
- [ ] Atur agar setiap kali fungsi `.Update()` atau `.Delete()` dipicu melalui decorator ini, sistem secara **OTOMatis & PASIF** meneruskan *Query Action Trigger* ke service audit/alerting tanpa campur tangan kode di level *Usecase*. Hal ini menjamin **100% kepatuhan pada Aturan Global tanpa risiko kelupaan dari developer**.

#### 6.2 Async Background Worker untuk Notifikasi Real-time
- [ ] Hubungkan modul `NotificationUsecase` dengan *In-Memory Event Bus (Go Channel Worker)* atau driver queue (seperti *Asynq/Redis*).
- [ ] Ketika kustomer menyerahkan pengajuan KYC, atau akun staf di-*Decommission*, worker latar belakang segera mendistribusikan email/webhook ke sistem antarmuka pemantauan operator tanpa menahan laju latensi HTTP handler.

---

### FASE 7: Standardisasi Respons untuk React Frontend

#### 7.1 Struktur Paginasi Baku (Paging Meta)
- [ ] Perluas struktur `Response` di [httputil/response.go](file:///Users/a2275/Artech/golang/folksart-be/backend-golang/httputil/response.go) untuk menghamparkan metadata tabel frontend:
  ```go
  type PaginationMeta struct {
      CurrentPage int `json:"currentPage"`
      Limit       int `json:"limit"`
      TotalItems  int `json:"totalItems"`
      TotalPages  int `json:"totalPages"`
  }

  type PaginatedResponse struct {
      Code       string          `json:"code"`
      Message    string          `json:"message"`
      Data       interface{}     `json:"data"`
      Pagination *PaginationMeta `json:"pagination,omitempty"`
  }
  ```

#### 7.2 Katalog Kode Error Bisnis (Frontend React Toast Alignment)
- [ ] Daftarkan konstanta kode error terstandarisasi di [errs/errors.go](file:///Users/a2275/Artech/golang/folksart-be/backend-golang/errs) agar frontend React dapat me-render UI Alert/Toast khusus:
  - `ERR_USER_BANNED` (Akun ditangguhkan, arahkan ke Customer Care).
  - `ERR_KYC_REQUIRED` (Batas akses customer umum, arahkan ke layar unggah identitas KYC).
  - `ERR_MFA_CHALLENGE` (Risiko tinggi terdeteksi, minta input OTP MFA).
  - `ERR_INSUFFICIENT_PERM` (Penolakan akses admin ops console, tampilkan banner 403).

---

## 📝 Catatan Penutupan & Instruksi Lanjutan
Dokumen ini merupakan panduan konkrit bagi siapapun yang menangani perbaikan lanjutan service ini. Ketika pembicaraan selanjutnya dilaksananakan, pastikan untuk memeriksa status flag di **Tabel Execution Checkpoint**, pilihlah item **PENDING (`[ ]`)** berurut, lakukan pengodean sejelas mungkin, dan perbarui centang tabel sesudahnya.
