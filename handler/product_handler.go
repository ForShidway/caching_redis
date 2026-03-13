package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"redis_caching_demo/domain"
)

// ProductHandler menangani semua HTTP request yang berkaitan dengan produk.
// Hanya bertugas: parse request → panggil service → tulis response.
// Tidak ada business logic di sini (clean separation).
type ProductHandler struct {
	service domain.ProductService
}

// NewProductHandler membuat instance handler baru.
func NewProductHandler(service domain.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// writeJSON adalah helper untuk menulis response JSON dengan status code.
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeError adalah helper untuk menulis response error dalam format JSON.
func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]string{"error": message})
}

// RegisterRoutes mendaftarkan semua route ke ServeMux.
// Dipisahkan dari main agar handler bisa di-test tanpa menjalankan server.
func (h *ProductHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/products", h.handleProducts)
	mux.HandleFunc("/products/", h.handleProductByID)
	mux.HandleFunc("/cache/", h.handleCacheInvalidation)
	mux.HandleFunc("/stats", h.handleStats)
	mux.HandleFunc("/health", h.handleHealth)
}

// handleProducts menangani GET /products
// Mengambil semua produk (dengan caching).
func (h *ProductHandler) handleProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method tidak diizinkan")
		return
	}

	products, err := h.service.GetAllProducts(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  products,
		"total": len(products),
	})
}

// handleProductByID menangani GET /products/{id}
// Ini adalah endpoint utama demonstrasi CACHE HIT vs MISS.
func (h *ProductHandler) handleProductByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method tidak diizinkan")
		return
	}

	// Ekstrak ID dari path "/products/1" → "1"
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "ID produk diperlukan")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID produk harus berupa angka")
		return
	}

	product, err := h.service.GetProductByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data": product,
	})
}

// handleCacheInvalidation menangani DELETE /cache/{id}
// Menghapus cache produk tertentu dari Redis.
func (h *ProductHandler) handleCacheInvalidation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, "gunakan method DELETE")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/cache/")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "ID produk diperlukan")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID produk harus berupa angka")
		return
	}

	if err := h.service.InvalidateProductCache(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "cache produk berhasil dihapus",
		"id":      idStr,
	})
}

// handleStats menangani GET /stats
// Menampilkan statistik CACHE HIT, MISS, dan hit rate.
func (h *ProductHandler) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method tidak diizinkan")
		return
	}

	stats := h.service.GetCacheStats()
	writeJSON(w, http.StatusOK, stats)
}

// handleHealth menangani GET /health
// Endpoint sederhana untuk mengecek apakah server berjalan.
func (h *ProductHandler) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "OK",
		"service": "Redis Caching Demo",
	})
}
