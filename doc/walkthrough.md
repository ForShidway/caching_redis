# Walkthrough — Redis Caching Demo (Go + Memurai)

## Status: ✅ Selesai — Build & Vet Berhasil

## Yang Dibangun

| File | Deskripsi |
|------|-----------|
| `config/config.go` | Konfigurasi terpusat (Redis addr, TTL, delay) |
| `domain/product.go` | Model `Product` dan `CacheStats` |
| `domain/interfaces.go` | Interface untuk semua layer (clean arch) |
| `cache/redis_cache.go` | Redis + **Track Method** (atomic HIT/MISS counter) |
| `repository/product_repo.go` | Simulasi DB in-memory + slow query 200ms |
| `service/product_service.go` | **Cache-Aside Pattern** |
| `handler/product_handler.go` | 5 HTTP endpoint |
| `middleware/latency.go` | **Latency measurement** tiap request |
| `cmd/main.go` | Entry point + dependency injection |

## Verifikasi Build

```bash
go build ./...   # ✅ sukses, 0 error
go vet ./...     # ✅ sukses, 0 warning
```

---

## Cara Menjalankan & Menguji

### 1. Pastikan Memurai Running
Buka Memurai atau jalankan via command:
```bash
memurai
```

### 2. Jalankan Server Go
```bash
cd X:\topik_khusus_p2
go run ./cmd/main.go
```

Output yang muncul:
```
=================================================
  Redis Caching Demo — Go + Memurai
=================================================
  Redis Address  : localhost:6379
  Cache TTL      : 30s
  DB Sim Delay   : 200ms
  Server Port    : :8080
=================================================
[OK] Terhubung ke Redis/Memurai
[SERVER] Berjalan di http://localhost:8080
```

### 3. Uji CACHE MISS (Request Pertama — ~200ms)
```bash
curl http://localhost:8080/products/1
```
**Log Server:**
```
[CACHE MISS] key=product:1 | mengambil dari database...
[CACHE SET]  key=product:1 | data disimpan ke Redis
[LATENCY] GET    /products/1   → 201.3ms  ← lambat
```

### 4. Uji CACHE HIT (Request Kedua — <1ms)
```bash
curl http://localhost:8080/products/1
```
**Log Server:**
```
[CACHE HIT]  key=product:1
[LATENCY] GET    /products/1   → 0.8ms   ← super cepat!
```

### 5. Lihat Statistik Cache
```bash
curl http://localhost:8080/stats
```
```json
{
  "hits": 1,
  "misses": 1,
  "total_requests": 2,
  "hit_rate": "50.00%"
}
```

### 6. Invalidasi Cache → MISS Lagi
```bash
curl -X DELETE http://localhost:8080/cache/1
curl http://localhost:8080/products/1   # kembali MISS ~200ms
```

### 7. Semua Produk
```bash
curl http://localhost:8080/products
```
