# Redis Caching System — Go + Memurai (Redis)

## 🎯 Use Case: REST API Katalog Produk dengan Redis Caching

Sistem ini mensimulasikan backend e-commerce ringan di mana data produk diambil dari **"database"** (simulasi dengan in-memory store). Redis (Memurai) digunakan sebagai **layer cache** untuk mempercepat response time. Sistem mencatat **cache HIT/MISS**, mengukur **latency** setiap request, dan dibangun dengan **clean architecture** (domain → repository → service → handler).

### Mengapa Katalog Produk?
- Data produk dibaca jauh lebih sering daripada ditulis (read-heavy) → ideal untuk caching
- Mudah didemonstrasikan secara langsung via HTTP
- Bisa terlihat jelas perbedaan latency antara cache HIT vs MISS

---

## Struktur Project

```
x:/topik_khusus_p2/
├── cmd/
│   └── main.go                  # Entry point
├── config/
│   └── config.go                # Config (Redis addr, TTL, port)
├── domain/
│   ├── product.go               # Model Product
│   └── interfaces.go            # Interface Repository & Cache
├── repository/
│   └── product_repo.go          # Simulasi DB (in-memory)
├── cache/
│   └── redis_cache.go           # Redis cache + Track Method (HIT/MISS/latency)
├── service/
│   └── product_service.go       # Business logic + cache-aside pattern
├── handler/
│   └── product_handler.go       # HTTP handler (net/http)
├── middleware/
│   └── latency.go               # Middleware pengukur latency per request
├── doc/
│   ├── task.md                  # Checklist pengerjaan
│   ├── implementation_plan.md   # Rencana implementasi
│   └── walkthrough.md           # Panduan menjalankan & menguji
├── go.mod
└── README.md
```

---

## Proposed Changes

### Domain Layer

#### [NEW] domain/product.go
Model `Product` dan `CacheStats`.

#### [NEW] domain/interfaces.go
Interface `ProductRepository` dan `CacheProvider` agar tiap layer bisa diganti (clean code / DI).

---

### Cache Layer

#### [NEW] cache/redis_cache.go
- Koneksi ke Memurai via `github.com/redis/go-redis/v9`
- Method: `Get`, `Set`, `Delete`, `Ping`
- **Track Method**: Setiap `Get` mencatat apakah HIT atau MISS ke `CacheTracker` (struct dengan counter atomic)
- **Latency Setting**: TTL dikonfigurasi dari `config`, duration operasi cache dicatat

---

### Repository Layer

#### [NEW] repository/product_repo.go
Simulasi database (map in-memory). Delay artifisial (`time.Sleep`) untuk mensimulasikan query DB lambat, sehingga perbedaan cache HIT vs MISS terlihat nyata.

---

### Service Layer

#### [NEW] service/product_service.go
Implementasi **Cache-Aside Pattern**:
1. Cek cache → jika HIT, return data dari Redis
2. Jika MISS → ambil dari repository → simpan ke cache → return data

---

### Handler & Middleware

#### [NEW] handler/product_handler.go
HTTP handler dengan endpoint:
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/products` | List semua produk |
| GET | `/products/{id}` | Detail produk (demonstrasi cache) |
| DELETE | `/cache/{id}` | Invalidasi cache produk |
| GET | `/stats` | Statistik cache (HIT/MISS count) |

#### [NEW] middleware/latency.go
Middleware yang mengukur dan mencetak latency tiap HTTP request ke log.

---

### Config & Entry Point

#### [NEW] config/config.go
```go
type Config struct {
    RedisAddr    string        // "localhost:6379"
    CacheTTL     time.Duration // default: 30s
    SlowQuerySim time.Duration // simulasi DB: 200ms
    ServerPort   string        // ":8080"
}
```

#### [NEW] cmd/main.go
Inisialisasi semua layer dengan dependency injection manual.

---

## Verification Plan

### Automated Tests

> [!NOTE]
> Tidak ada existing tests. Verifikasi dilakukan via HTTP request manual menggunakan `curl` atau browser.

### Manual Verification (Step-by-Step)

**Prerequisite**: Memurai harus running di `localhost:6379`.

**1. Jalankan server:**
```bash
cd x:/topik_khusus_p2
go run ./cmd/main.go
```

**2. Test MISS (pertama kali, data dari DB):**
```bash
curl http://localhost:8080/products/1
```
→ Log akan menampilkan `[CACHE MISS]` + latency ~200ms

**3. Test HIT (kedua kali, data dari Redis):**
```bash
curl http://localhost:8080/products/1
```
→ Log akan menampilkan `[CACHE HIT]` + latency <5ms

**4. Lihat statistik cache:**
```bash
curl http://localhost:8080/stats
```
→ Response JSON: `{"hits": 1, "misses": 1, "hit_rate": "50.00%"}`

**5. Invalidasi cache lalu hit lagi:**
```bash
curl -X DELETE http://localhost:8080/cache/1
curl http://localhost:8080/products/1
```
→ Log: `[CACHE MISS]` kembali muncul
