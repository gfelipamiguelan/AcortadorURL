package handlers

import (
	"Reto05/shortener"  // Inicializa y guarda datos
	"bytes"             // Paquete que prepara los datos JSON
	"encoding/json"     // Paquete para codificar y decodificar datos JSON
	"net/http"          // Paquete para manejar solicitudes HTTP
	"net/http/httptest" // Paquete que simula las solicitudes HTTP para pruebas
	"testing"           // Paquete para pruebas unitarias
)

// Prueba si el endpoint "/shorten" responde correctamente a una solicitud válida
func TestShortenHandler_ValidRequest(t *testing.T) {

	shortener.InitStore() // Inicializa el almacenamiento en memoria

	reqBody := map[string]string{ // Prepara el cuerpo de la solicitud con la URL larga
		"long_url": "https://www.example.com",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// Crea una solicitud HTTP POST con el cuerpo JSON
	// y establece el tipo de contenido como JSON
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder() // Crea un grabador que captura la respuesta HTTP

	ShortenHandler(w, req) // Llama al manejador ShortenHandler  para ejecutar la lógica de acortamiento

	// Captura la respuesta simulada
	resp := w.Result()
	defer resp.Body.Close()

	// Verifica que el código de estado sea 200 OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verifica que la respuesta JSON tenga una clave short_url
	var res map[string]string
	json.NewDecoder(resp.Body).Decode(&res)

	if _, ok := res["short_url"]; !ok {
		t.Error("Expected short_url in response")
	}
}

// Prueba si el endpoint "/shorten" responde correctamente a una solicitud inválida
// Debería responder con error 400
func TestShortenHandler_InvalidRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	ShortenHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d", resp.StatusCode)
	}
}

// Prueba que una URL acortada redirija correctamente a la URL original
// Verifica que el código de estado sea 301 o 307 (redirección permanente o temporal)
func TestRedirectHandler_ValidCode(t *testing.T) {
	shortener.InitStore()

	// Guarda manualmente la URL
	code := "abc123"
	url := "https://www.ejemplo.com"
	shortener.Save(code, url)

	// Simula una visita a "/abc123"
	req := httptest.NewRequest(http.MethodGet, "/"+code, nil)
	w := httptest.NewRecorder()

	RedirectHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	// Verifica que el código de estado sea 301 o 307 (redirección permanente o temporal)
	if resp.StatusCode != http.StatusMovedPermanently && resp.StatusCode != http.StatusTemporaryRedirect {
		t.Errorf("Expected 301 or 307, got %d", resp.StatusCode)
	}

	loc := resp.Header.Get("Location")
	if loc != url {
		t.Errorf("Expected redirect to %s, got %s", url, loc)
	}
}

// Simula el acceso a un código inválido.
// Devuelve un 404
func TestRedirectHandler_InvalidCode(t *testing.T) {
	shortener.InitStore()

	req := httptest.NewRequest(http.MethodGet, "/invalidcode", nil)
	w := httptest.NewRecorder()

	RedirectHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404 Not Found, got %d", resp.StatusCode)
	}
}
