package main

import (
	"Reto05/handlers" // Importamos el paquete donde están los controladores de las rutas
	"fmt"             // Paquete para imprimir mensajes en la consola
	"log"             // Paquete para mostrar mensajes de error si el servidor falla
	"net/http"        // Paquete para crear el servidor web y definir las rutas (endpoints)
)

func main() { // Configuramos las rutas que el servidor va a manejar

	http.HandleFunc("/health", handlers.HealthHandler) // Verifica que el servidor esté activo

	http.HandleFunc("/shorten", handlers.ShortenHandler) // Recibe una URL larga y genera un código corto

	http.HandleFunc("/", handlers.RedirectHandler) // Nos redirige a la URL original usando el código corto

	fmt.Println("Servidor iniciado en el puerto :8080") // Imprime un mensaje en consola indicando que el servidor está funcionando

	log.Fatal(http.ListenAndServe(":8080", nil)) // Inicia el servidor en el puerto 8080 y detiene el proceso si hay algún error

}
