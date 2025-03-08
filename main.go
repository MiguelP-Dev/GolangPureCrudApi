package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

var (
	users  = make(map[int]User)
	nextID = 1
	mu     sync.Mutex
)

func main() {
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userHandler)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	response := APIResponse{
		Success: false,
		Error:   message,
	}
	writeJSON(w, status, response)
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

	switch r.Method {
	case http.MethodGet:
		getUser(w, r, id)
	case http.MethodPut:
		updateUser(w, r, id)
	case http.MethodDelete:
		deleteUser(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func listUsers(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	var userList []User
	for _, user := range users {
		userList = append(userList, user)
	}
	mu.Unlock()

	response := APIResponse{
		Success: true,
		Data:    userList,
	}
	writeJSON(w, http.StatusOK, response)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	mu.Lock()
	user.ID = nextID
	nextID++
	users[user.ID] = user
	mu.Unlock()

	response := APIResponse{
		Success: true,
		Data:    user,
	}
	writeJSON(w, http.StatusCreated, response)
}

func getUser(w http.ResponseWriter, r *http.Request, id int) {
	mu.Lock()
	user, exists := users[id]
	mu.Unlock()

	if !exists {
		writeError(w, http.StatusNotFound, "User not found")
		return
	}

	response := APIResponse{
		Success: true,
		Data:    user,
	}
	writeJSON(w, http.StatusOK, response)
}

func updateUser(w http.ResponseWriter, r *http.Request, id int) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	mu.Lock()
	_, exists := users[id]
	if !exists {
		mu.Unlock()
		writeError(w, http.StatusNotFound, "User not found")
		return
	}
	user.ID = id
	users[id] = user
	mu.Unlock()

	response := APIResponse{
		Success: true,
		Data:    user,
	}
	writeJSON(w, http.StatusOK, response)
}

func deleteUser(w http.ResponseWriter, r *http.Request, id int) {
	mu.Lock()
	_, exists := users[id]
	if !exists {
		mu.Unlock()
		writeError(w, http.StatusNotFound, "User not found")
		return
	}
	delete(users, id)
	mu.Unlock()

	response := APIResponse{
		Success: true,
	}
	writeJSON(w, http.StatusOK, response)
}
