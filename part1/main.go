package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type ResponseMsg struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type Item struct {
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

type database map[string]Item

func main() {
	db := database{
		"cpu_amd":    Item{Name: "cpu_amd", Price: 24999, Quantity: 12},
		"cpu_intel":  Item{Name: "cpu_intel", Price: 20999, Quantity: 24},
		"gpu_amd":    Item{Name: "gpu_amd", Price: 47999, Quantity: 40},
		"gpu_nvidia": Item{Name: "gpu_nvidia", Price: 89999, Quantity: 5},
	}

	http.HandleFunc("/", db.list)
	http.HandleFunc("/show", db.show)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)

	log.Fatal(http.ListenAndServe("localhost:4000", nil))
}

func (db database) list(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	json.NewEncoder(w).Encode(db)
}

func (db database) show(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	item := r.URL.Query().Get("item")

	if _, exists := db[item]; !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseMsg{"item not found", http.StatusNotFound})
		return
	}

	json.NewEncoder(w).Encode(db[item])
}

func (db database) create(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	name := r.URL.Query().Get("name")
	price, err := strconv.Atoi(r.URL.Query().Get("price"))

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ResponseMsg{"invalid price", http.StatusUnprocessableEntity})
		return
	}

	quantity, err := strconv.Atoi(r.URL.Query().Get("quantity"))

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ResponseMsg{"invalid quantity", http.StatusUnprocessableEntity})
		return
	}

	if _, exists := db[name]; exists {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(ResponseMsg{fmt.Sprintf("item %s already exists", name), http.StatusUnprocessableEntity})
		return
	}

	db[name] = Item{Name: name, Price: price, Quantity: quantity}
	w.WriteHeader(http.StatusCreated)
}

func (db database) update(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	name := r.URL.Query().Get("item")

	if _, exists := db[name]; !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseMsg{"item not found", http.StatusNotFound})
		return
	}

	item := db[name]

	price := r.URL.Query().Get("price")
	if price != "" {
		price, err := strconv.Atoi(r.URL.Query().Get("price"))

		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ResponseMsg{"invalid price", http.StatusUnprocessableEntity})
			return
		}

		item.Price = price
	}

	quantity := r.URL.Query().Get("quantity")
	if quantity != "" {
		quantity, err := strconv.Atoi(r.URL.Query().Get("quantity"))

		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ResponseMsg{"invalid quantity", http.StatusUnprocessableEntity})
			return
		}

		item.Quantity = quantity
	}

	db[item.Name] = item

	json.NewEncoder(w).Encode(item)
}

func (db database) delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	item := r.URL.Query().Get("item")

	if _, exists := db[item]; !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ResponseMsg{"item not found", http.StatusNotFound})
		return
	}

	delete(db, item)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(ResponseMsg{"item deleted", http.StatusAccepted})
}
