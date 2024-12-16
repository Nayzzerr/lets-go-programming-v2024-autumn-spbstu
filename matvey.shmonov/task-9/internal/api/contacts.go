package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Koshsky/task-9/internal/database"

	"github.com/gorilla/mux"
)

type API struct {
	ContactManager *database.ContactManager
}

func (api *API) GetContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	contactList, err := api.ContactManager.GetContacts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(contactList)
}

func (api *API) GetContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	contact, err := api.ContactManager.GetContact(id)
	if err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(contact)
}

func (api *API) CreateContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	name := r.URL.Query().Get("name")
	phone := r.URL.Query().Get("phone")

	if name == "" || phone == "" {
		http.Error(w, "name and phone are required", http.StatusBadRequest)
		return
	}

	id, err := api.ContactManager.CreateContact(name, phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contact := database.Contact{ID: id, Name: name, Phone: phone}
	json.NewEncoder(w).Encode(contact)
}

func (api *API) UpdateContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	name := r.URL.Query().Get("name")
	phone := r.URL.Query().Get("phone")

	if name == "" && phone == "" {
		http.Error(w, "At least one of name or phone must be provided", http.StatusBadRequest)
		return
	}

	if err := api.ContactManager.UpdateContact(id, name, phone); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	contact, _ := api.ContactManager.GetContact(id)
	json.NewEncoder(w).Encode(contact)
}

func (api *API) DeleteContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := api.ContactManager.DeleteContact(id); err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func RegisterRoutes(r *mux.Router, cm *database.ContactManager) {
	api := &API{ContactManager: cm}
	r.HandleFunc("/contacts", api.GetContacts).Methods("GET")
	r.HandleFunc("/contacts/{id:[0-9]+}", api.GetContact).Methods("GET")
	r.HandleFunc("/contacts", api.CreateContact).Methods("POST")
	r.HandleFunc("/contacts/{id:[0-9]+}", api.UpdateContact).Methods("PUT")
	r.HandleFunc("/contacts/{id:[0-9]+}", api.DeleteContact).Methods("DELETE")
}
