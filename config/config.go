package config

import "time"

// Config menyimpan semua konfigurasi aplikasi.
// Dipakai di semua layer tanpa import silang (clean code).
type Config struct {
	// Alamat Redis/Memurai
	RedisAddr string

	// Password Redis (kosong jika tidak ada)
	RedisPassword string

	// Database Redis (default: 0)
	RedisDB int

	// TTL (Time-To-Live) data di cache Redis
	// Setelah durasi ini, key otomatis terhapus dari Redis
	CacheTTL time.Duration

	// Simulasi delay query database yang lambat
	// Agar perbedaan CACHE HIT vs MISS terlihat jelas di log latency
	SlowQuerySim time.Duration

	// Port server HTTP
	ServerPort string
}

// DefaultConfig mengembalikan konfigurasi default yang siap pakai.
func DefaultConfig() *Config {
	return &Config{
		RedisAddr:     "localhost:6379",
		RedisPassword: "",
		RedisDB:       0,
		CacheTTL:      30 * time.Second,  // data expired setelah 30 detik
		SlowQuerySim:  200 * time.Millisecond, // simulasi DB lambat 200ms
		ServerPort:    ":8080",
	}
}
