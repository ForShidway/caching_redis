package domain

import "time"

// Product adalah model utama yang merepresentasikan data produk.
// Dipakai di seluruh layer (repository, service, handler).
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// CacheStats menyimpan statistik penggunaan cache secara keseluruhan.
// Dipakai oleh endpoint GET /stats untuk melaporkan efektivitas cache.
type CacheStats struct {
	Hits         int64   `json:"hits"`           // jumlah cache HIT
	Misses       int64   `json:"misses"`         // jumlah cache MISS
	TotalRequest int64   `json:"total_requests"` // total request ke cache
	HitRate      string  `json:"hit_rate"`       // persentase HIT (misal: "83.33%")
}
