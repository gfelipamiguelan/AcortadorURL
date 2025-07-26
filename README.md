# Reto 5 - Acortador de URLs en Go

Este proyecto implementa un servicio básico de acortamiento de URLs utilizando Go. Permite generar una URL corta para cualquier URL larga ingresada, y redirecciona correctamente a la original cuando se accede a la corta.

---

## Generación de códigos cortos

El código corto se genera usando el siguiente algoritmo:

1. Se concatena la URL larga con la marca de tiempo actual y un número aleatorio.
2. Se genera un `hash` MD5 del resultado.
3. Se selecciona una longitud aleatoria entre **6 y 8 caracteres**.
4. Se toma esa cantidad de caracteres desde el hash (convertido a base 36).
5. Se verifica si el código ya existe. Si hay colisión, se reintenta hasta 5 veces.

Esto asegura que:
- Cada código corto sea único, aleatorio y difícil de predecir.
- El rango de longitud cumple el requisito de ser entre 6 y 8 caracteres.

---

## Manejo de colisiones

Para evitar conflictos con códigos ya generados:

- El sistema intenta hasta 5 veces generar un código único.
- En cada intento, el `hash` base cambia por el uso de `time.Now().UnixNano()` y `rand.Int()`.
- Si después de 5 intentos no se logra generar un código único, se devuelve un error.

---

## Redirección: ¿301 o 307?

Se ha optado por usar el código de estado **HTTP 301 (Moved Permanently)** en lugar de 307 por las siguientes razones:

- La URL acortada se considera **permanente**: siempre redirigirá a la misma URL larga.
- El código 301 es ideal para enlaces permanentes y permite que navegadores y motores de búsqueda lo cacheen.
- 307 es útil cuando el método HTTP no debe cambiar, pero en este caso no es necesario mantener el método (usamos `GET` para redirigir).

---

## Ejecución del proyecto

1. Ejecutar el servidor:
   ```bash
   go run main.go
   ```

2. Crear una URL corta:
   ```bash
   curl -X POST http://localhost:8080/shorten \
        -H "Content-Type: application/json" \
        -d '{"url": "https://www.example.com"}'
   ```

3. Acceder desde el navegador:
   ```
   http://localhost:8080/{codigo}
   ```

---

## Ejecutar pruebas

El proyecto cuenta con pruebas automatizadas para verificar el correcto funcionamiento del acortador y sus endpoints.

### Instrucciones:

1. Abre una terminal y navega a la raíz del proyecto.
2. Ejecuta el siguiente comando:
   ```bash
   go test ./... -v
   ```

### Qué se prueba:

- Validación de requests válidos e inválidos (`POST /shorten`).
- Redirecciones correctas (`GET /{code}`).
- Colisiones y generación de códigos únicos.
- Comportamiento concurrente del almacenamiento en memoria.

--------------------------------------
--------------------------------------


Se utilizó una estructura modular para separar responsabilidades:

handlers/: controla las rutas y la lógica HTTP.

shortener/: contiene la lógica de negocio (generación de códigos y almacenamiento).
Esto facilita el mantenimiento, las pruebas unitraias y la escalabilidad futura del sistema.

## Estructura del proyecto

```
Reto05/
├── main.go                        # Punto de entrada del servidor
├── go.mod                        # Dependencias del proyecto
├── handlers/
│   ├── http.go                   # Handlers HTTP (shorten, redirect, health)
│   └── http_test.go              # Pruebas HTTP de los endpoints
├── shortener/
│   ├── service.go                # Lógica para generación de códigos cortos
│   ├── store.go                  # Almacenamiento concurrente en memoria
│   └── shortener_test.go         # Pruebas de la lógica del acortador
└── README.md                     # Explicación técnica del proyecto
```
