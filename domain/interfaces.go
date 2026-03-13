package domain

import "context"

// ProductRepository adalah kontrak untuk mengakses data produk dari storage (DB).
// Dengan interface ini, repository bisa diganti ke MySQL/PostgreSQL tanpa mengubah service.
type ProductRepository interface {
	// FindByID mencari satu produk berdasarkan ID.
	FindByID(ctx context.Context, id int) (*Product, error)

	// FindAll mengambil semua produk.
	FindAll(ctx context.Context) ([]*Product, error)
}

// CacheProvider adalah kontrak untuk operasi cache.
// Dengan interface ini, Redis bisa diganti ke Memcached/in-memory tanpa mengubah service.
type CacheProvider interface {
	// Get mengambil nilai dari cache berdasarkan key.
	// Mengembalikan error jika key tidak ditemukan (cache MISS).
	Get(ctx context.Context, key string) (string, error)

	// Set menyimpan nilai ke cache dengan TTL tertentu.
	Set(ctx context.Context, key string, value string) error

	// Delete menghapus key dari cache (cache invalidation).
	Delete(ctx context.Context, key string) error

	// Ping mengecek koneksi ke Redis server.
	Ping(ctx context.Context) error

	// GetStats mengembalikan statistik cache HIT/MISS.
	GetStats() *CacheStats
}

// ProductService adalah kontrak untuk business logic produk.
type ProductService interface {
	// GetProductByID mengambil produk (dari cache atau DB).
	GetProductByID(ctx context.Context, id int) (*Product, error)

	// GetAllProducts mengambil semua produk.
	GetAllProducts(ctx context.Context) ([]*Product, error)

	// InvalidateProductCache menghapus cache produk tertentu.
	InvalidateProductCache(ctx context.Context, id int) error

	// GetCacheStats mengembalikan statistik penggunaan cache.
	GetCacheStats() *CacheStats
}
