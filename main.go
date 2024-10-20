package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"os"
	"github.com/gorilla/handlers"
)

type Member struct {
	ID          string `json:"ID"`
	Name        string `json:"name"`
	Gender      string `json:"gender"`
	BirthPlace  string `json:"birth_place"`
	BirthDate   int64  `json:"birth_date"`
	Phone       int64  `json:"phone"`
	Kelurahan   string `json:"kelurahan"`
	Kecamatan   string `json:"kecamatan"`
	Job         string `json:"job"`
	RT          int    `json:"rt"`
	RW          int    `json:"rw"`
	Address     string `json:"address"`
	Status      int    `json:"status"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

var members = make(map[string]Member)

func GetMember(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if member, ok := members[id]; ok {
		jsonResponse(w, http.StatusOK, "success_ok", member)
		return
	}
	jsonResponse(w, http.StatusNotFound, "member_not_found", nil)
}

func GetMembers(w http.ResponseWriter, r *http.Request) {
	var memberList []Member
	for _, member := range members {
		memberList = append(memberList, member)
	}
	jsonResponse(w, http.StatusOK, "success_ok", memberList)
}

func CreateMember(w http.ResponseWriter, r *http.Request) {
	var member Member
	err := json.NewDecoder(r.Body).Decode(&member)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	members[member.ID] = member
	jsonResponse(w, http.StatusCreated, "member_created", member)
}

func UpdateMember(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var updatedMember Member
	err := json.NewDecoder(r.Body).Decode(&updatedMember)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := members[id]; ok {
		updatedMember.ID = id
		members[id] = updatedMember
		jsonResponse(w, http.StatusOK, "member_updated", updatedMember)
		return
	}
	jsonResponse(w, http.StatusNotFound, "member_not_found", nil)
}

func DeleteMember(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	if _, ok := members[id]; ok {
		delete(members, id)
		jsonResponse(w, http.StatusOK, "member_deleted", nil)
		return
	}
	jsonResponse(w, http.StatusNotFound, "member_not_found", nil)
}

func jsonResponse(w http.ResponseWriter, code int, message string, data interface{}) {
	response := map[string]interface{}{
		"code":    code,
		"message": message,
		"data":    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/members/{id}", GetMember).Methods("GET")
	r.HandleFunc("/api/v1/members", GetMembers).Methods("GET")
	r.HandleFunc("/api/v1/members", CreateMember).Methods("POST")
	r.HandleFunc("/api/v1/members/{id}", UpdateMember).Methods("PUT")
	r.HandleFunc("/api/v1/members/{id}", DeleteMember).Methods("DELETE")

	h := handlers.CORS(handlers.AllowedOrigins([]string{"*"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}), handlers.AllowedHeaders([]string{"Content-Type"}))(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), h))
}

