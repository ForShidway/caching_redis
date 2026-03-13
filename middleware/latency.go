package middleware

import (
	"log"
	"net/http"
	"time"
)

// LatencyMiddleware adalah HTTP middleware yang mengukur dan mencatat
// durasi (latency) setiap request yang masuk ke server.
//
// Cara kerja:
//   1. Catat waktu mulai sebelum handler dijalankan
//   2. Teruskan request ke handler berikutnya
//   3. Setelah handler selesai, hitung selisih waktu
//   4. Cetak ke log dengan format yang mudah dibaca
//
// Ini adalah tempat di mana perbedaan CACHE HIT (~0.7ms) vs
// CACHE MISS (~202ms) bisa dilihat langsung di terminal.
func LatencyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mulai ukur waktu
		start := time.Now()

		// Jalankan handler selanjutnya (bisa berupa handler lain atau route handler)
		next.ServeHTTP(w, r)

		// Hitung durasi setelah handler selesai
		duration := time.Since(start)

		// Log latency dengan format yang informatif
		// Contoh output:
		//   [LATENCY] GET /products/1          → 201.4ms   ← CACHE MISS
		//   [LATENCY] GET /products/1          → 0.8ms     ← CACHE HIT
		//   [LATENCY] GET /stats               → 0.1ms
		log.Printf("[LATENCY] %-6s %-30s → %v",
			r.Method,
			r.URL.Path,
			duration.Round(time.Microsecond),
		)
	})
}
