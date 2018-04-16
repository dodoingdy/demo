package main

import (
	"io/ioutil"
	"net/http"
	"log"
	"encoding/json"
	"regexp"
	"demo/pghandler"
	"github.com/gorilla/mux"
	"strconv"
)

type Name struct {
	Name string `json:name`
}

type State struct {
	State string `json:state`
}

func usershandler(w http.ResponseWriter, r*http.Request) {
	var Users pghandler.Users
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		bodyStr := string(body)
		log.Println(bodyStr)
		var username Name
		if err := json.Unmarshal(body, &username); err == nil {
			username := username.Name
			log.Println("username", username)
			result := Users.UserInsert(username)
			w.Write(result)
		} else {
			log.Fatal(err)
		}
	} else if r.Method == "GET" {
		result := Users.GetAllUsers()
		w.Write(result)
	} else {
		info := "404 Not Found"
		err := []byte(info)
		w.Write(err)
	}
}

func getAllRelsHandler(w http.ResponseWriter, r*http.Request) {
	uidStr := mux.Vars(r)["user_id"]
	var Rel pghandler.Relationships
	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		log.Fatal(err)
		info := "ERROR:Uid is not number"
		errInfo := []byte(info)
		w.Write(errInfo)
	} else {
		result := Rel.GetALLRelationships(uid)
		w.Write(result)
	}
}

func newRelHandler(w http.ResponseWriter, r*http.Request) {
	var Rel pghandler.Relationships
	uidStr := mux.Vars(r)["user_id"]
	userIdStr := mux.Vars(r)["other_user_id"]

	uid, err := strconv.Atoi(uidStr)
	if err != nil {
		log.Fatal(err)
		info := "ERROR:Uid is not number"
		errInfo := []byte(info)
		w.Write(errInfo)
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		log.Fatal(err)
		info := "ERROR:Uid is not number"
		errInfo := []byte(info)
		w.Write(errInfo)
	}
	
	body, _ := ioutil.ReadAll(r.Body)
	bodyStr := string(body)
	log.Println(bodyStr)
	var state State
	if err := json.Unmarshal(body, &state); err == nil {
		state := state.State
		log.Println(uid, userId, state)
		result := Rel.NewRelationship(uid, userId, state, "relationship")
		w.Write(result)
	} else {
		log.Println(err)
	}
}

func main() {
	pghandler.SqlOpen()
	r := mux.NewRouter()
	r.HandleFunc("/users", usershandler).Methods("GET", "POST")
	r.HandleFunc("/users/{user_id}/relationships", getAllRelsHandler).Methods("GET")
	r.HandleFunc("/users/{user_id}/relationships/{other_user_id}", newRelHandler).Methods("PUT")
	http.Handle("/", r)
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	defer pghandler.SqlClose()
}
