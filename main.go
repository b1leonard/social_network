package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var database *sql.DB

// Users Comment
type Users struct {
	Users []User `json:"users"`
}

// User Comment
type User struct {
	ID    int    "json:id"
	Name  string "json:username"
	Email string "json:email"
	First string "json:first"
	Last  string "json:last"
}

// CreateResponse comment
type CreateResponse struct {
	Error     string "json:error"
	ErrorCode int    "json:code"
}

// ErrMsg comment
type ErrMsg struct {
	ErrCode    int
	StatusCode int
	Msg        string
}

// UserCreate Comment
func UserCreate(w http.ResponseWriter, r *http.Request) {

	NewUser := User{}
	NewUser.Name = r.FormValue("user")
	NewUser.Email = r.FormValue("email")
	NewUser.First = r.FormValue("first")
	NewUser.Last = r.FormValue("last")
	output, err := json.Marshal(NewUser)
	fmt.Println(string(output))
	if err != nil {
		fmt.Println("Something went wrong")
	}

	Response := CreateResponse{}
	sql := "INSERT INTO users set user_nickname='" + NewUser.Name + "', user_first='" + NewUser.First + "', user_last='" + NewUser.Last + "', user_email='" + NewUser.Email + "'"
	q, err := database.Exec(sql)
	if err != nil {
		errorMessage, errorCode := dbErrorParser(err.Error())
		fmt.Println(errorMessage)
		error, httpCode, msg := ErrorMessages(errorCode)
		Response.Error = msg
		Response.ErrorCode = error
		fmt.Println(httpCode)
	}
	fmt.Println(q)
	createOutput, _ := json.Marshal(Response)
	fmt.Fprintln(w, string(createOutput))
}

//UsersRetrieve Comment
func UsersRetrieve(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Pragma", "no-cache")

	rows, _ := database.Query("select * from users LIMIT 10")
	Response := Users{}
	for rows.Next() {
		user := User{}
		rows.Scan(&user.ID, &user.Name, &user.First, &user.Last, &user.Email)

		Response.Users = append(Response.Users, user)
	}
	output, _ := json.Marshal(Response)
	fmt.Fprintln(w, string(output))
}

// ErrorMessages comment
func ErrorMessages(err int64) (int, int, string) {
	errorMessage := ""
	statusCode := 200
	errorCode := 0
	switch err {
	case 1062:
		errorMessage = "Duplicate entry"
		errorCode = 10
		statusCode = 409
	}
	return errorCode, statusCode, errorMessage
}

func dbErrorParser(err string) (string, int64) {
	Parts := strings.Split(err, ":")
	errorMessage := Parts[1]
	Code := strings.Split(Parts[0], "Error ")
	errorCode, _ := strconv.ParseInt(Code[1], 10, 32)
	return errorMessage, errorCode
}

func main() {
	mysqlString := "golang@localhost" + ":" + "123" + "@tcp(" + "localhost" + ":3306)/social_network"
	db, err := sql.Open("mysql", mysqlString)
	if err != nil {

	}
	database = db
	routes := mux.NewRouter()
	routes.HandleFunc("/api/users", UserCreate).Methods("POST")
	routes.HandleFunc("/api/users", UsersRetrieve).Methods("GET")
	http.Handle("/", routes)
	http.ListenAndServe(":8080", nil)
}
