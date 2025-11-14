package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"sync"

	"github.com/a1eksandre/pack-calculator/internal/calculator"
)

// Server holds pack sizes and exposes HTTP handlers.
type Server struct {
	mu        sync.RWMutex
	packSizes []int
}

func NewServer(defaultPackSizes []int) *Server {
	sizes := append([]int(nil), defaultPackSizes...)
	return &Server{
		packSizes: sizes,
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/pack-sizes", s.handlePackSizes)
	mux.HandleFunc("/api/calculate", s.handleCalculate)

	return loggingMiddleware(mux)
}

func (s *Server) handlePackSizes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetPackSizes(w, r)
	case http.MethodPut:
		s.handleSetPackSizes(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleGetPackSizes(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	sizes := append([]int(nil), s.packSizes...)
	s.mu.RUnlock()

	resp := packSizesResponse{PackSizes: sizes}
	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleSetPackSizes(w http.ResponseWriter, r *http.Request) {
	var req packSizesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if len(req.PackSizes) == 0 {
		http.Error(w, "packSizes must not be empty", http.StatusBadRequest)
		return
	}

	seen := make(map[int]struct{})
	sizes := make([]int, 0, len(req.PackSizes))
	for _, v := range req.PackSizes {
		if v <= 0 {
			http.Error(w, "packSizes must contain positive numbers", http.StatusBadRequest)
			return
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		sizes = append(sizes, v)
	}
	sort.Ints(sizes)

	s.mu.Lock()
	s.packSizes = sizes
	s.mu.Unlock()

	resp := packSizesResponse{PackSizes: sizes}
	writeJSON(w, http.StatusOK, resp)
}

// handleCalculate returns pack counts for a given number of items.
func (s *Server) handleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req calculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}
	if req.Items <= 0 {
		http.Error(w, "items must be > 0", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	sizes := append([]int(nil), s.packSizes...)
	s.mu.RUnlock()

	sol, err := calculator.CalculatePacks(req.Items, sizes)
	if err != nil {
		http.Error(w, "cannot calculate packs: "+err.Error(), http.StatusBadRequest)
		return
	}

	entries := make([]packEntry, 0, len(sol.Packs))
	for pack, qty := range sol.Packs {
		entries = append(entries, packEntry{Pack: pack, Quantity: qty})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Pack < entries[j].Pack
	})

	resp := calculateResponse{
		Items:      req.Items,
		PackSizes:  sizes,
		Solution:   entries,
		TotalItems: sol.TotalItems,
		ExtraItems: sol.ExtraItems,
	}

	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON error: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
