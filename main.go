package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis"
)

type pingResponse struct {
	Status string `json:"status"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Println(path)

	/* client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}) */

	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Use "redis" instead of "localhost"
		Password: "",
		DB:       0,
	})

	switch r.Method {
	case http.MethodGet:
		if path == "/" {
			fmt.Println("SWITCH:", path)
			homepage := "This is a sample JSON homepage for the microservice."
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(homepage)
		} else if path == "/ping" {
			fmt.Println("SWITCH:", path)
			response := pingResponse{Status: "working"}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else if strings.HasPrefix(path, "/tenant/") {
			tenantName := strings.TrimPrefix(path, "/tenant/")
			fmt.Println("SWITCH:", path)
			value, err := client.Get(tenantName).Result()
			fmt.Println(value)
			if err == redis.Nil {
				http.Error(w, "Tenant not found", http.StatusNotFound)
			} else if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, value)
			}
		} else {
			http.Error(w, "Invalid request path", http.StatusBadRequest)
		}

	case http.MethodPost:
		if strings.HasPrefix(path, "/") {
			tenantName := strings.TrimPrefix(path, "/")
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusBadRequest)
				return
			}

			err = client.Set(tenantName, string(body), 0).Err()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				fmt.Fprintf(w, "POST request for tenant %s, saved JSON data\n", tenantName)
			}
		} else {
			http.Error(w, "Invalid request path", http.StatusBadRequest)
		}

	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
