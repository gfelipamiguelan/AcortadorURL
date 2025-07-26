// Almacenamiento concurrente en memoria de URLs largas y códigos cortos
package shortener

import (
	"errors" // Paquete para manejar errores
	"sync"   // Paquete para manejar concurrencia
)

var (
	Store = make(map[string]string)
	Mutex = sync.RWMutex{}
)

// Guarda una URL larga con su código corto
// Utiliza un mutex para asegurar que el acceso al mapa sea seguro en concurrencia
func Save(shortCode, longURL string) {
	Mutex.Lock()
	defer Mutex.Unlock()
	Store[shortCode] = longURL
}

// Obtiene la URL larga a partir de su código corto
func Get(shortCode string) (string, error) {
	Mutex.RLock() // Asegura que el acceso al mapa sea seguro en concurrencia
	defer Mutex.RUnlock()

	longURL, exists := Store[shortCode]
	if !exists {
		return "", errors.New("código no encontrado")
	}
	return longURL, nil
}

// Inicializa (o reinicia) el almacenamiento en memoria
func InitStore() {
	Mutex.Lock()
	defer Mutex.Unlock()
	Store = make(map[string]string)
}
