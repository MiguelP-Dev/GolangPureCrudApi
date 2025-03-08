# API CRUD en Go Puro

## Descripción

Este proyecto implementa una API CRUD básica usando Go sin frameworks externos. Incluye manejo de usuarios con operaciones Crear, Leer, Actualizar y Eliminar.

## Características

- API RESTful
- Manejo de JSON
- Respuestas consistentes
- Control de concurrencia
- Manejo de errores estructurado

## Estructura de Respuesta API

```json
{
    "success": true,
    "data": {},
    "error": ""
}
```

## Métodos de Prueba

### Usando cURL

#### Crear Usuario (POST)

```bash
curl -X POST http://localhost:8080/users \
    -H "Content-Type: application/json" \
    -d '{"name":"Juan Pérez","email":"juan.perez@ejemplo.com"}'
```

#### Listar Usuarios (GET)

```bash
curl -X GET http://localhost:8080/users
```

#### Obtener Usuario (GET)

```bash
curl -X GET http://localhost:8080/users/1
```

#### Actualizar Usuario (PUT)

```bash
curl -X PUT http://localhost:8080/users/1 \
    -H "Content-Type: application/json" \
    -d '{"name":"Ana López","email":"ana.lopez@ejemplo.com"}'
```

#### Eliminar Usuario (DELETE)

```bash
curl -X DELETE http://localhost:8080/users/1
```

### Usando Postman

1. **Crear Usuario (POST)**
     - URL: `http://localhost:8080/users`
     - Método: `POST`
     - Cabeceras: `Content-Type: application/json`
     - Cuerpo:

         ```json
         {
             "name": "Juan Pérez",
             "email": "juan.perez@ejemplo.com"
         }
         ```

2. **Listar Usuarios (GET)**
     - URL: `http://localhost:8080/users`
     - Método: `GET`

3. **Obtener Usuario (GET)**
     - URL: `http://localhost:8080/users/1`
     - Método: `GET`

4. **Actualizar Usuario (PUT)**
     - URL: `http://localhost:8080/users/1`
     - Método: `PUT`
     - Cabeceras: `Content-Type: application/json`
     - Cuerpo:

         ```json
         {
             "name": "Ana López",
             "email": "ana.lopez@ejemplo.com"
         }
         ```

5. **Eliminar Usuario (DELETE)**
     - URL: `http://localhost:8080/users/1`
     - Método: `DELETE`

### Usando Thunder Client (Extensión de VS Code)

1. **Crear Usuario**
     - Nueva Petición
     - Método: `POST`
     - URL: `http://localhost:8080/users`
     - Cabeceras: `Content-Type: application/json`
     - Cuerpo (JSON):

         ```json
         {
             "name": "Juan Pérez",
             "email": "juan.perez@ejemplo.com"
         }
         ```

2. **Listar Usuarios**
     - Nueva Petición
     - Método: `GET`
     - URL: `http://localhost:8080/users`

3. **Obtener Usuario**
     - Nueva Petición
     - Método: `GET`
     - URL: `http://localhost:8080/users/1`

4. **Actualizar Usuario**
     - Nueva Petición
     - Método: `PUT`
     - URL: `http://localhost:8080/users/1`
     - Cabeceras: `Content-Type: application/json`
     - Cuerpo (JSON):

         ```json
         {
             "name": "Ana López",
             "email": "ana.lopez@ejemplo.com"
         }
         ```

5. **Eliminar Usuario**
     - Nueva Petición
     - Método: `DELETE`
     - URL: `http://localhost:8080/users/1`

## Implementación del Código

### Estructuras Principales

```go
// Estructura de Usuario
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// Estructura de Respuesta API
type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

// Variables globales
var (
    users  = make(map[int]User)
    nextID = 1
    mu     sync.Mutex
)
```

### Manejadores HTTP

```go
func main() {
    http.HandleFunc("/users", usersHandler)
    http.HandleFunc("/users/", userHandler)
    log.Println("Server started at :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        listUsers(w, r)
    case http.MethodPost:
        createUser(w, r)
    default:
        writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
    }
}

func userHandler(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Path[len("/users/"):])
    if err != nil {
        writeError(w, http.StatusBadRequest, "Invalid user ID")
        return
    }
    // ... [resto de la implentación del handler]
}
```

## Mejoras Implementadas en el CRUD

**Estructura de Respuesta Consistente**

```go
type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}
```

**Manejo de Errores Mejorado**

```go
func writeError(w http.ResponseWriter, status int, message string) {
    response := APIResponse{
        Success: false,
        Error:   message,
    }
    writeJSON(w, status, response)
}
```

**Respuestas JSON Estandarizadas**

```go
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        log.Printf("Error encoding response: %v", err)
    }
}
```

**Control de Concurrencia con Mutex**

```go
var (
    users  = make(map[int]User)
    nextID = 1
    mu     sync.Mutex
)
```

## Ejecución del Proyecto

1. Clonar el repositorio
2. Navegar al directorio del proyecto
3. Ejecutar:

```bash
go run main.go
```

- El servidor estará disponible en `http://localhost:8080`

## Requisitos

- Go 1.23.1 o superior
- Ninguna dependencia externa

## Notas

- La API utiliza almacenamiento en memoria
- Los datos se pierden al reiniciar el servidor
- Implementación thread-safe con mutex
