package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	_ "github.com/lib/pq"
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

var db *sql.DB

func initDB() {
	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	var err error
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")
}

func GetMember(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var member Member
	err := db.QueryRow(`SELECT * FROM members WHERE id = $1`, id).Scan(&member.ID, &member.Name, &member.Gender, &member.BirthPlace, &member.BirthDate, &member.Phone, &member.Kelurahan, &member.Kecamatan, &member.Job, &member.RT, &member.RW, &member.Address, &member.Status, &member.CreatedAt, &member.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			jsonResponse(w, http.StatusNotFound, "member_not_found", nil)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, "success_ok", member)
}

func GetMembers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT * FROM members`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var member Member
		err := rows.Scan(&member.ID, &member.Name, &member.Gender, &member.BirthPlace, &member.BirthDate, &member.Phone, &member.Kelurahan, &member.Kecamatan, &member.Job, &member.RT, &member.RW, &member.Address, &member.Status, &member.CreatedAt, &member.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		members = append(members, member)
	}
	jsonResponse(w, http.StatusOK, "success_ok", members)
}

func CreateMember(w http.ResponseWriter, r *http.Request) {
	var member Member
	err := json.NewDecoder(r.Body).Decode(&member)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO members (id, name, gender, birth_place, birth_date, phone, kelurahan, kecamatan, job, rt, rw, address, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	_, err = db.Exec(query, member.ID, member.Name, member.Gender, member.BirthPlace, member.BirthDate, member.Phone, member.Kelurahan, member.Kecamatan, member.Job, member.RT, member.RW, member.Address, member.Status, member.CreatedAt, member.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

	query := `UPDATE members SET name=$1, gender=$2, birth_place=$3, birth_date=$4, phone=$5, kelurahan=$6, kecamatan=$7, job=$8, rt=$9, rw=$10, address=$11, status=$12, created_at=$13, updated_at=$14 WHERE id=$15`
	result, err := db.Exec(query, updatedMember.Name, updatedMember.Gender, updatedMember.BirthPlace, updatedMember.BirthDate, updatedMember.Phone, updatedMember.Kelurahan, updatedMember.Kecamatan, updatedMember.Job, updatedMember.RT, updatedMember.RW, updatedMember.Address, updatedMember.Status, updatedMember.CreatedAt, updatedMember.UpdatedAt, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	numRows, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if numRows == 0 {
		jsonResponse(w, http.StatusNotFound, "member_not_found", nil)
		return
	}
	jsonResponse(w, http.StatusOK, "member_updated", updatedMember)
}

func DeleteMember(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	result, err := db.Exec(`DELETE FROM members WHERE id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	numRows, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if numRows == 0 {
		jsonResponse(w, http.StatusNotFound, "member_not_found", nil)
		return
	}
	jsonResponse(w, http.StatusOK, "member_deleted", nil)
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

// Handler function for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	initDB()

	routes := mux.NewRouter()
	routes.HandleFunc("/api/v1/members/{id}", GetMember).Methods("GET")
	routes.HandleFunc("/api/v1/members", GetMembers).Methods("GET")
	routes.HandleFunc("/api/v1/members", CreateMember).Methods("POST")
	routes.HandleFunc("/api/v1/members/{id}", UpdateMember).Methods("PUT")
	routes.HandleFunc("/api/v1/members/{id}", DeleteMember).Methods("DELETE")

	h := handlers.CORS(handlers.AllowedOrigins([]string{"*"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}), handlers.AllowedHeaders([]string{"Content-Type"}))(routes)

	h.ServeHTTP(w, r)
}
