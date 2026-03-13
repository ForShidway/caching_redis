# Redis Caching Demo — Go + Memurai

Sistem REST API Katalog Produk yang mendemonstrasikan Redis caching dengan **track method**, **latency measurement**, dan **clean architecture** menggunakan Go.

---

## Prasyarat

- [Go](https://golang.org/dl/) 1.21+
- [Memurai](https://www.memurai.com/) (Redis untuk Windows) — pastikan sudah berjalan di `localhost:6379`

---

## Cara Menjalankan

```bash
# 1. Pastikan Memurai sudah running di background

# 2. Install dependencies
go mod tidy

# 3. Jalankan server
go run ./cmd/main.go
```

---

## Endpoint API

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/products` | Daftar semua produk |
| `GET` | `/products/{id}` | Detail produk (**demo CACHE HIT/MISS**) |
| `DELETE` | `/cache/{id}` | Invalidasi cache produk |
| `GET` | `/stats` | Statistik cache HIT/MISS/hit rate |
| `GET` | `/health` | Health check server |

---

## Demo Perbedaan Latency

Buka **dua terminal** setelah server berjalan:

```bash
# Terminal 2 — Request PERTAMA: MISS (~200ms)
curl http://localhost:8080/products/1

# Terminal 2 — Request KEDUA: HIT (<1ms) ← drastis lebih cepat!
curl http://localhost:8080/products/1

# Lihat statistik
curl http://localhost:8080/stats

# Invalidasi cache, lalu request lagi → MISS kembali
curl -X DELETE http://localhost:8080/cache/1
curl http://localhost:8080/products/1
```

**Log yang muncul di Terminal 1 (server):**
```
[CACHE MISS] key=product:1 | mengambil dari database...
[CACHE SET]  key=product:1 | data disimpan ke Redis
[LATENCY] GET    /products/1                    → 201.3ms

[CACHE HIT]  key=product:1
[LATENCY] GET    /products/1                    → 0.8ms
```

---

## Arsitektur

```
cmd/main.go          → Entry point, dependency injection
config/config.go     → Konfigurasi terpusat (Redis addr, TTL, port)
domain/
  product.go         → Model data
  interfaces.go      → Kontrak antar layer (interface)
cache/
  redis_cache.go     → Redis + Track Method (atomic HIT/MISS counter)
repository/
  product_repo.go    → Simulasi database (in-memory + slow query sim)
service/
  product_service.go → Cache-Aside Pattern (core logic)
handler/
  product_handler.go → HTTP handler (parse request, tulis response)
middleware/
  latency.go         → Ukur & cetak latency setiap request
```

---

## Penjelasan Fitur Utama

### Track Method
Setiap `Get()` ke Redis mencatat HIT atau MISS via atomic counter (thread-safe). Lihat hasilnya di `GET /stats`.

### Latency Setting
- `CacheTTL`: durasi data hidup di Redis (default: 30 detik)
- `SlowQuerySim`: delay simulasi query DB (default: 200ms)
- Latency tiap request dicetak otomatis oleh `LatencyMiddleware`

### Clean Architecture
Setiap layer hanya berkomunikasi lewat **interface** → mudah diganti, mudah ditest.

---

## 💬 Riwayat Prompt Permintaan

Berikut adalah seluruh prompt/permintaan yang diajukan selama proses pembuatan sistem ini karena sistem ini dibuat menggunakan Agentic AI Claude Sonet:

---

**Prompt 1 — Permintaan Awal & Pemilihan Use Case**
> *"saya ingin menguji catching dengan menggunakan redis dan saya juga sudah mendownload mamurai untuk menggunakan redis. saya juga ingin agar system ini dibangun dengan bahasa golang dan memakai track method, makai setingan latensi, dan clean code. apakah anda bisa berikan kasus sistem yang cocok?"*

---

**Prompt 2 — Pertanyaan tentang Konsep Sistem**
> *"saya mau bertanya, manakah penjelasan mengenai fungsi catching dan redisnya dan golangnya disini, dan gimana cara kerja sistem ini nantinya menggunakan redis golang, track method dan clean arsitektur ini?"*

---

**Prompt 3 — Pertanyaan tentang Cara Melihat Perbedaan Kecepatan**
> *"kapan perbedaan kecepatannya dapat dilihat dan dengan cara apa dan bagaimana"*

---

**Prompt 4 — Perintah Mulai Implementasi**
> *"okay sekarang lanjut bangun lagi ya"*

---

**Prompt 5 — Perintah Menyimpan Dokumentasi ke Folder doc**
> *"saya ingin agar file task, implementation plan, dan walkthrough ini disimpan dalam folder baru dengan nama doc apakah bisa?"*

---

**Prompt 6 — Perintah Menambahkan Riwayat Prompt ke README**
> *"saya ingin di dalam file readme.md itu ditambahkan prompt permintaan saya kepada anda tadi ya semua prompt permintaan saya yang berkaitan dengan sistem ini"*
