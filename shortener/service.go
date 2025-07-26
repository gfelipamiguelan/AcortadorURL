// Lógica para la generación de códigos cortos
package shortener

import (
	"crypto/md5" // Paquete para generar un hash MD5
	"fmt"        // Paquete para formatear cadenas
	"math/rand"  // Paquete para generar números aleatorios
	"strings"    // Paquete para manipular cadenas
	"time"       // Paquete para manejar el tiempo
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Genera un código corto a partir del hash MD5, utiliza una longitud determinada
func generateCode(hash []byte, length int) string {
	var code strings.Builder
	for _, b := range hash {
		code.WriteByte(charset[int(b)%len(charset)])
		if code.Len() == length {
			break
		}
	}
	return code.String()
}

// Genera un código corto único para una URL larga
// Intenta generar el código hasta 5 veces
func GenerateShortCode(longURL string) (string, error) {
	const maximoIntentos = 5

	for i := 0; i < maximoIntentos; i++ {
		// Genera un hash a partir del tiempo + URL
		seed := fmt.Sprintf("%d:%s:%d", time.Now().UnixNano(), longURL, rand.Int())
		hash := md5.Sum([]byte(seed))

		// Longitud aleatoria entre 6 y 8
		length := rand.Intn(3) + 6

		code := generateCode(hash[:], length)

		// Verifica si ya existe
		if _, err := Get(code); err != nil {
			Save(code, longURL)
			return code, nil
		}
	}

	return "", fmt.Errorf("no se pudo generar un código único después de %d intentos", maximoIntentos)
}
