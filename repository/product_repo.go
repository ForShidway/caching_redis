package repository

import (
	"context"
	"fmt"
	"time"

	"redis_caching_demo/config"
	"redis_caching_demo/domain"
)

// productRepo adalah implementasi dari domain.ProductRepository.
// Mensimulasikan database dengan data in-memory (slice/map).
// Dalam proyek nyata, ini akan berisi query ke MySQL/PostgreSQL.
type productRepo struct {
	db  map[int]*domain.Product // simulasi tabel database
	cfg *config.Config          // untuk mengambil SlowQuerySim duration
}

// NewProductRepository membuat instance repository baru dengan data dummy.
// Data produk diinisialisasi di sini (simulasi seed database).
func NewProductRepository(cfg *config.Config) domain.ProductRepository {
	now := time.Now()

	// Data dummy — simulasi isi tabel 'products' di database
	db := map[int]*domain.Product{
		1: {
			ID: 1, Name: "Laptop Pro X", Category: "Electronics",
			Price: 15000000, Stock: 10,
			Description: "Laptop performa tinggi untuk developer",
			CreatedAt:   now,
		},
		2: {
			ID: 2, Name: "Mechanical Keyboard RGB", Category: "Accessories",
			Price: 850000, Stock: 50,
			Description: "Keyboard mekanikal dengan backlight RGB",
			CreatedAt:   now,
		},
		3: {
			ID: 3, Name: "4K Monitor Ultrawide", Category: "Electronics",
			Price: 7500000, Stock: 5,
			Description: "Monitor 34 inci resolusi 4K untuk produktivitas",
			CreatedAt:   now,
		},
		4: {
			ID: 4, Name: "Wireless Mouse Ergonomic", Category: "Accessories",
			Price: 450000, Stock: 80,
			Description: "Mouse wireless ergonomis untuk kerja seharian",
			CreatedAt:   now,
		},
		5: {
			ID: 5, Name: "USB-C Hub 7-in-1", Category: "Accessories",
			Price: 320000, Stock: 30,
			Description: "Hub USB-C dengan port HDMI, USB-A, SD Card",
			CreatedAt:   now,
		},
	}

	return &productRepo{db: db, cfg: cfg}
}

// FindByID mencari produk berdasarkan ID.
// SlowQuerySim: delay artifisial untuk mensimulasikan query DB yang lambat.
// Ini membuat perbedaan CACHE HIT vs MISS terlihat dramatis di log latency.
func (r *productRepo) FindByID(ctx context.Context, id int) (*domain.Product, error) {
	// Simulasi operasi database yang membutuhkan waktu (baca dari disk)
	time.Sleep(r.cfg.SlowQuerySim)

	product, exists := r.db[id]
	if !exists {
		return nil, fmt.Errorf("produk dengan ID %d tidak ditemukan", id)
	}

	return product, nil
}

// FindAll mengambil semua produk dari "database".
// Juga menggunakan SlowQuerySim karena join atau full scan biasanya lebih lambat.
func (r *productRepo) FindAll(ctx context.Context) ([]*domain.Product, error) {
	// Simulasi full table scan
	time.Sleep(r.cfg.SlowQuerySim)

	products := make([]*domain.Product, 0, len(r.db))
	for _, p := range r.db {
		products = append(products, p)
	}

	return products, nil
}
