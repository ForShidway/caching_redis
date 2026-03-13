package main

import (
	"context"
	"log"
	"net/http"

	"redis_caching_demo/cache"
	"redis_caching_demo/config"
	"redis_caching_demo/handler"
	"redis_caching_demo/middleware"
	"redis_caching_demo/repository"
	"redis_caching_demo/service"
)

func main() {
	// ─── 1. LOAD KONFIGURASI ─────────────────────────────────────────────
	// Semua pengaturan (Redis addr, TTL, port) terpusat di satu tempat.
	cfg := config.DefaultConfig()

	log.Println("=================================================")
	log.Println("  Redis Caching Demo — Go + Memurai")
	log.Println("=================================================")
	log.Printf("  Redis Address  : %s", cfg.RedisAddr)
	log.Printf("  Cache TTL      : %v", cfg.CacheTTL)
	log.Printf("  DB Sim Delay   : %v (simulasi query lambat)", cfg.SlowQuerySim)
	log.Printf("  Server Port    : %s", cfg.ServerPort)
	log.Println("=================================================")

	// ─── 2. INISIALISASI CACHE (Redis) ───────────────────────────────────
	// Dependency Injection: cache dibuat sekali dan disuntikkan ke service.
	redisCache := cache.NewRedisCache(cfg)

	// Validasi koneksi ke Memurai sebelum server mulai
	if err := redisCache.Ping(context.Background()); err != nil {
		log.Fatalf("[ERROR] Tidak bisa terhubung ke Redis/Memurai: %v\n"+
			"Pastikan Memurai sudah berjalan di %s", err, cfg.RedisAddr)
	}
	log.Println("[OK] Terhubung ke Redis/Memurai")

	// ─── 3. INISIALISASI REPOSITORY (Simulasi DB) ────────────────────────
	// Dependency Injection: repository dibuat sekali dan disuntikkan ke service.
	productRepo := repository.NewProductRepository(cfg)

	// ─── 4. INISIALISASI SERVICE (Business Logic) ────────────────────────
	// Service menerima interface (bukan concrete type) → clean architecture.
	productSvc := service.NewProductService(productRepo, redisCache)

	// ─── 5. INISIALISASI HANDLER (HTTP Layer) ────────────────────────────
	productHandler := handler.NewProductHandler(productSvc)

	// ─── 6. SETUP ROUTER ─────────────────────────────────────────────────
	mux := http.NewServeMux()
	productHandler.RegisterRoutes(mux)

	// ─── 7. SETUP MIDDLEWARE ─────────────────────────────────────────────
	// LatencyMiddleware membungkus semua handler → mengukur durasi setiap request.
	// Ini adalah tempat di mana CACHE HIT vs MISS terlihat di log.
	serverHandler := middleware.LatencyMiddleware(mux)

	// ─── 8. JALANKAN SERVER ──────────────────────────────────────────────
	log.Printf("\n[SERVER] Berjalan di http://localhost%s", cfg.ServerPort)
	log.Println("\nEndpoint yang tersedia:")
	log.Println("  GET    /products          — daftar semua produk")
	log.Println("  GET    /products/{id}     — detail produk (demo cache HIT/MISS)")
	log.Println("  DELETE /cache/{id}        — hapus cache produk")
	log.Println("  GET    /stats             — statistik cache HIT/MISS/hit_rate")
	log.Println("  GET    /health            — health check server\n")

	server := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: serverHandler,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("[ERROR] Server gagal berjalan: %v", err)
	}
}
