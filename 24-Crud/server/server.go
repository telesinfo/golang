package server

import (
	"crud/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	ID    uint32 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// AddUser add new user
func AddUser(w http.ResponseWriter, r *http.Request) {

	// read body content
	bodyRequest, error := ioutil.ReadAll(r.Body)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Body request error!"))
		return
	}

	// convert body content to struct
	var u user
	if error = json.Unmarshal(bodyRequest, &u); error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error converting user to struct"))
		return
	}

	// connect to db
	db, error := database.Connect()
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error connecting to DB"))
		return
	}
	defer db.Close()

	// create statement
	statement, error := db.Prepare("insert into users (name, email) values (?, ?)")
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error creating statement"))
		return
	}
	defer statement.Close()

	// insert user
	ins, error := statement.Exec(u.Name, u.Email)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error creating user"))
		return
	}

	// get new user id
	id, error := ins.LastInsertId()
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to get new user id"))
		return
	}

	// return new id
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("New user created! ID: %d", id)))

}

// GetUsers get all users
func GetUsers(w http.ResponseWriter, r *http.Request) {

	// connect db
	db, error := database.Connect()
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error connecting to DB"))
		return
	}
	defer db.Close()

	// get users
	rows, error := db.Query("select * from users")
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error select users"))
		return
	}
	defer rows.Close()

	//
	var users []user
	for rows.Next() {
		var u user
		if error := rows.Scan(&u.ID, &u.Name, &u.Email); error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error get users"))
			return
		}

		users = append(users, u)
	}

	// return users
	w.WriteHeader(http.StatusOK)
	if error := json.NewEncoder(w).Encode(users); error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error return users"))
		return
	}

}

// GetUser get user with ID
func GetUser(w http.ResponseWriter, r *http.Request) {

	// get parameters in the request
	parameter := mux.Vars(r)

	ID, error := strconv.ParseUint(parameter["id"], 10, 32)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error convert parameter ID"))
		return

	}

	// connect db
	db, error := database.Connect()
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error connecting to DB"))
		return
	}
	defer db.Close()

	// get user with ID
	row, error := db.Query("select * from users where id = ?", ID)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error get user with ID"))
		return
	}

	// convert user to struct
	var u user
	if row.Next() {
		if error := row.Scan(&u.ID, &u.Name, &u.Email); error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error convert user to struct"))
			return
		}
	}

	//
	w.WriteHeader(http.StatusOK)
	if error := json.NewEncoder(w).Encode(u); error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error convert user to json"))
		return
	}

}

// UpdateUser update user
func UpdateUser(w http.ResponseWriter, r *http.Request) {

	// get parameters in the request
	parameter := mux.Vars(r)
	ID, error := strconv.ParseUint(parameter["id"], 10, 32)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error convert parameter ID"))
		return

	}

	// read body content
	bodyRequest, error := ioutil.ReadAll(r.Body)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Body request error!"))
		return
	}

	// convert body content to struct
	var u user
	if error = json.Unmarshal(bodyRequest, &u); error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error converting user to struct"))
		return
	}

	// connect db
	db, error := database.Connect()
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error connecting to DB"))
		return
	}
	defer db.Close()

	// create statement
	statement, error := db.Prepare("update users set name = ?, email = ? where id = ?")
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error prepare statement"))
		return
	}
	defer statement.Close()

	// execute statement
	if _, error := statement.Exec(u.Name, u.Email, ID); error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to update user"))
		return
	}

	w.WriteHeader(http.StatusOK)

}

// DeleteUser delete user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// get parameters in the request
	parameter := mux.Vars(r)

	ID, error := strconv.ParseUint(parameter["id"], 10, 32)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error convert parameter ID"))
		return

	}

	// connect db
	db, error := database.Connect()
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error connecting to DB"))
		return
	}
	defer db.Close()

	// create statement
	statement, error := db.Prepare("delete from users where id = ?")
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error prepare statement"))
		return
	}
	defer statement.Close()

	// execute statement
	if _, error := statement.Exec(ID); error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error to delete user"))
		return
	}

	w.WriteHeader(http.StatusOK)

}
