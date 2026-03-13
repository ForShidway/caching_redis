package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"redis_caching_demo/domain"
)

// productService adalah implementasi dari domain.ProductService.
// Layer ini menerapkan Cache-Aside Pattern:
//   1. Cek cache Redis terlebih dahulu
//   2. Jika HIT → return data dari cache (cepat)
//   3. Jika MISS → ambil dari repository (lambat) → simpan ke cache → return
type productService struct {
	repo  domain.ProductRepository // akses ke "database"
	cache domain.CacheProvider     // akses ke Redis
}

// NewProductService membuat instance service baru.
// Menerima interface (bukan concrete type) → mudah di-mock untuk testing.
func NewProductService(repo domain.ProductRepository, cache domain.CacheProvider) domain.ProductService {
	return &productService{
		repo:  repo,
		cache: cache,
	}
}

// cacheKeyProduct menghasilkan key Redis yang konsisten untuk satu produk.
// Contoh: "product:1", "product:42"
func cacheKeyProduct(id int) string {
	return fmt.Sprintf("product:%d", id)
}

// cacheKeyAllProducts adalah key Redis untuk daftar semua produk.
const cacheKeyAllProducts = "products:all"

// GetProductByID mengambil data produk dengan Cache-Aside Pattern.
// Ini adalah inti dari demonstrasi CACHE HIT vs MISS.
func (s *productService) GetProductByID(ctx context.Context, id int) (*domain.Product, error) {
	key := cacheKeyProduct(id)

	// === LANGKAH 1: Cek Redis Cache ===
	cached, err := s.cache.Get(ctx, key)
	if err == nil {
		// CACHE HIT → deserialize JSON dan return langsung
		log.Printf("[CACHE HIT]  key=%s", key)

		var product domain.Product
		if jsonErr := json.Unmarshal([]byte(cached), &product); jsonErr != nil {
			return nil, fmt.Errorf("gagal parse cache: %w", jsonErr)
		}
		return &product, nil
	}

	// === LANGKAH 2: Cache MISS → Ambil dari Repository (DB) ===
	log.Printf("[CACHE MISS] key=%s | mengambil dari database...", key)

	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// === LANGKAH 3: Simpan ke Cache agar request berikutnya HIT ===
	jsonData, jsonErr := json.Marshal(product)
	if jsonErr == nil {
		if setErr := s.cache.Set(ctx, key, string(jsonData)); setErr != nil {
			// Gagal simpan ke cache tidak menghentikan proses (non-fatal)
			log.Printf("[CACHE WARN] gagal menyimpan ke cache: %v", setErr)
		} else {
			log.Printf("[CACHE SET]  key=%s | data disimpan ke Redis", key)
		}
	}

	return product, nil
}

// GetAllProducts mengambil semua produk dengan cache pada list keseluruhan.
func (s *productService) GetAllProducts(ctx context.Context) ([]*domain.Product, error) {
	key := cacheKeyAllProducts

	// Cek cache
	cached, err := s.cache.Get(ctx, key)
	if err == nil {
		log.Printf("[CACHE HIT]  key=%s", key)

		var products []*domain.Product
		if jsonErr := json.Unmarshal([]byte(cached), &products); jsonErr != nil {
			return nil, fmt.Errorf("gagal parse cache: %w", jsonErr)
		}
		return products, nil
	}

	// MISS → ambil dari repository
	log.Printf("[CACHE MISS] key=%s | mengambil dari database...", key)

	products, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Simpan ke cache
	jsonData, jsonErr := json.Marshal(products)
	if jsonErr == nil {
		if setErr := s.cache.Set(ctx, key, string(jsonData)); setErr != nil {
			log.Printf("[CACHE WARN] gagal menyimpan ke cache: %v", setErr)
		} else {
			log.Printf("[CACHE SET]  key=%s | data disimpan ke Redis", key)
		}
	}

	return products, nil
}

// InvalidateProductCache menghapus cache produk tertentu dari Redis.
// Dipanggil saat data produk diubah agar cache tidak stale (usang).
func (s *productService) InvalidateProductCache(ctx context.Context, id int) error {
	key := cacheKeyProduct(id)
	log.Printf("[CACHE DEL]  key=%s | cache dihapus", key)
	return s.cache.Delete(ctx, key)
}

// GetCacheStats mendelegasikan ke cache provider untuk mengambil statistik HIT/MISS.
func (s *productService) GetCacheStats() *domain.CacheStats {
	return s.cache.GetStats()
}
