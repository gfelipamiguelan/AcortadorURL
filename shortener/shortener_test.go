// Contiene las pruebas de la lógica del acortador
package shortener

import (
	"sync"    // Paquete para manejar concurrencia
	"testing" // Paquete para pruebas unitarias
)

// Prueba si el almacenamiento en memoria funciona correctamente
func TestConcurrentSaveAndGet(t *testing.T) {
	InitStore() // Inicializa el almacenamiento en memoria
	url := "https://www.example.com"
	code := "abc123"

	var wg sync.WaitGroup   // Permite esperar que múltiples goroutines terminen
	const concurrency = 100 // Define el número de goroutines concurrentes

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Save(code, url)
			retrieved, err := Get(code)
			if err != nil { //
				t.Errorf("Error getting code: %v", err)
			}
			if retrieved != url {
				t.Errorf("Expected %s, got %s", url, retrieved)
			}
		}()
	}

	wg.Wait()
}

func TestGenerateShortCodeLength(t *testing.T) {
	InitStore()

	url := "https://www.testlength.com"
	code, err := GenerateShortCode(url)
	if err != nil {
		t.Fatalf("Error generating short code: %v", err)
	}
	if len(code) < 6 || len(code) > 8 {
		t.Errorf("Expected code length 6-8, got %d: %s", len(code), code)
	}
}
