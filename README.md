# Go Gin Clean Architecture Boilerplate (Feature-Based CSR + GORM + gRPC)

Boilerplate backend **Golang + Gin + GORM + gRPC** dengan pendekatan **Feature-Based Clean Architecture (CSR: Controller – Service – Repository)**, diadaptasi dari arsitektur [ifanferdi/POS-nextjs-express](https://github.com/ifanferdi/POS-nextjs-express/tree/api/apps/api).

## Tech Stack

| Kategori | Teknologi |
| :--- | :--- |
| Language | Go 1.24+ |
| HTTP Framework | Gin (`github.com/gin-gonic/gin`) |
| RPC Framework | gRPC (`google.golang.org/grpc`) |
| ORM & Database | GORM (`gorm.io/gorm`) + PostgreSQL |
| Validation | `go-playground/validator/v10` (terintegrasi dengan Gin binding) |
| Config | `joho/godotenv` |
| Agent Workflow | Ponytail (`AGENTS.md` + `.agents/skills/`) |

---

## Struktur Direktori

```text
go-gin-boilerplate/
├── .agents/
│   └── skills/                                           # Agent skills (ponytail & go-gin-feature)
├── cmd/
│   └── api/
│       └── main.go                                       # Entrypoint HTTP + gRPC & CLI flags (--migrate, --drop, --fresh, --seed)
├── docs/
│   └── proto/
│       ├── category.proto                                # Kontrak Protobuf gRPC CategoryService
│       └── product.proto                                 # Kontrak Protobuf gRPC ProductService
├── internal/
│   ├── bootstrap/
│   │   └── bootstrap.go                                  # Wiring Gin, gRPC Server/Client, & registrasi semua Feature
│   ├── config/
│   │   └── config.go                                     # Konfigurasi environment (App, Database, GRPC)
│   ├── constants/
│   │   └── http_status.go                                # Konstanta HTTP status & pesan standar
│   ├── helpers/
│   │   ├── error.go                                      # Custom AppError (BadRequest, NotFound, Conflict, dll.)
│   │   └── paginate.go                                   # Generic Paginate[T] & ParseOrderBy
│   ├── middlewares/
│   │   ├── error_handler.go                              # Global Error & Validation Handler Middleware
│   │   └── not_found_handler.go                          # 404 Route Not Found Handler
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── postgres.go                               # Koneksi GORM PostgreSQL
│   │   │   ├── migrations/
│   │   │   │   ├── migration.go                          # Batch-based Migration Manager (Migrate, Rollback, DropAll, Fresh)
│   │   │   │   ├── 20261027000001_create_categories_table.go
│   │   │   │   └── 20261027000002_create_products_and_pivot_table.go
│   │   │   └── seeders/
│   │   │       ├── seeder.go                             # Runner semua seeder
│   │   │       ├── category_seeder.go                    # Seeder data awal Category
│   │   │       └── product_seeder.go                     # Seeder data awal Product & pivot product_has_categories
│   │   └── grpc/
│   │       ├── server.go                                 # gRPC Server lifecycle & reflection
│   │       └── connection.go                             # gRPC Client Connection helper (grpc.NewClient)
│   └── features/
│       ├── category/                                     # Feature Category
│       │   ├── _services/
│       │   │   └── category_service_create.go            # Extended business logic pembuatan Category
│       │   ├── category_entity.go                        # GORM Model Category (table: product_categories)
│       │   ├── category_dto.go                           # DTO Request & Validation tags
│       │   ├── category_repository.go                    # Interface & Implementasi GORM Repository
│       │   ├── category_service.go                       # Interface & Service utama
│       │   ├── category_controller.go                    # Gin HTTP Controller
│       │   ├── category_grpc.go                          # gRPC Server Handler
│       │   └── category_routes.go                        # Registrasi HTTP & gRPC Routes
│       └── product/                                      # Feature Product + Pivot ProductHasCategory
│           ├── _services/
│           │   └── product_service_create.go             # Extended business logic pembuatan Product + validasi SKU/Barcode/Categories
│           ├── product_entity.go                         # GORM Model Product & ProductHasCategory (table: product_has_categories)
│           ├── product_dto.go                            # DTO Request (termasuk categoryIds: []uint)
│           ├── product_repository.go                     # GORM Repository (Preload Categories & transaksi pivot)
│           ├── product_service.go                        # Interface & Service utama
│           ├── product_controller.go                     # Gin HTTP Controller
│           ├── product_grpc.go                           # gRPC Server Handler
│           ├── product_routes.go                         # Registrasi HTTP & gRPC Routes
│           └── product_test.go                           # Unit/Integration Test HTTP, Pivot Many-to-Many, & gRPC
├── AGENTS.md
├── opencode.json
├── .env.example
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── go.mod
```

---

## Cara Menjalankan Project

### 1. Persiapan Environment & Database
```bash
cp .env.example .env
docker compose up -d
```

### 2. Migration & Seeder (Ekuivalen `bun run db:*` di POS-nextjs-express)
```bash
# Jalankan migrasi tabel (Up)
go run ./cmd/api --migrate
# atau: make db-migrate

# Jalankan seeder data awal
go run ./cmd/api --seed
# atau: make db-seed

# Reset database (Drop + Migrate + Seed sekaligus)
go run ./cmd/api --fresh --seed
# atau: make db-fresh-seed

# Rollback batch migrasi terakhir
go run ./cmd/api --rollback
# atau: make db-rollback
```

### 3. Menjalankan Server (HTTP `:8000` & gRPC `:50051`)
```bash
go run ./cmd/api
# atau: make dev
```

---

## Kontrak Endpoint (`Category` & `Product` CRUD + Pivot `product_has_categories`)

### HTTP REST API (`/api/v1/categories` & `/api/v1/products`)

| Method | Endpoint | Deskripsi |
| :--- | :--- | :--- |
| `GET` | `/api/v1/categories` | List kategori + pagination (`?page=1&limit=10&q=coffee&productId=1`) |
| `POST` | `/api/v1/categories/get` | Alternatif list kategori via JSON body |
| `GET` | `/api/v1/categories/:id` | Detail kategori berdasarkan ID |
| `POST` | `/api/v1/categories` | Buat kategori baru (didelegasikan ke `category/_services/category_service_create.go`) |
| `PUT` | `/api/v1/categories/:id` | Update kategori berdasarkan ID |
| `DELETE` | `/api/v1/categories/:id` | Hapus kategori berdasarkan ID |
| `GET` | `/api/v1/products` | List produk + preload `categories` (`?page=1&limit=10&q=americano&categoryId=1`) |
| `POST` | `/api/v1/products/get` | Alternatif list produk via JSON body |
| `GET` | `/api/v1/products/:id` | Detail produk + preload `categories` berdasarkan ID |
| `POST` | `/api/v1/products` | Buat produk baru + attach `categoryIds` ke tabel pivot `product_has_categories` |
| `PUT` | `/api/v1/products/:id` | Update produk + sinkronisasi pivot `product_has_categories` |
| `DELETE` | `/api/v1/products/:id` | Soft-delete produk berdasarkan ID |

### gRPC Services (pada `:50051`)
- `category.v1.CategoryService`: `FindAll`, `FindByID`, `Create`, `Update`, `Delete`
- `product.v1.ProductService`: `FindAll`, `FindByID`, `Create`, `Update`, `Delete`
