package handlers

import (
	"Reto05/shortener"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Método no permitido"}`, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		LongURL string `json:"long_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.LongURL) == "" {
		http.Error(w, `{"error": "URL inválida"}`, http.StatusBadRequest)
		return
	}

	parsed, err := url.ParseRequestURI(req.LongURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		http.Error(w, `{"error": "URL no válida o esquema no permitido"}`, http.StatusBadRequest)
		return
	}

	code, err := shortener.GenerateShortCode(req.LongURL)
	if err != nil {
		http.Error(w, `{"error": "Error al generar código corto"}`, http.StatusInternalServerError)
		return
	}

	resp := struct {
		ShortURL string `json:"short_url"`
	}{
		ShortURL: "http://localhost:8080/" + code,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/")
	if code == "" {
		http.NotFound(w, r)
		return
	}

	longURL, err := shortener.Get(code)
	if err != nil {
		http.Error(w, `{"error": "No se encontró la URL"}`, http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusMovedPermanently)
}
