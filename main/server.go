package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
	"encoding/json"
	"regexp"
	"demo/pghandler"
	"github.com/gorilla/mux"
	"strconv"
	"os"
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
		body_str := string(body)
		log.Println(body_str)
		var username Name
		if err := json.Unmarshal(body, &username); err == nil {
			username := username.Name
			log.Println("username", username)
			result := Users.UserInsert(username)
			w.Write(result)
		} else {
			fmt.Println(err)
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

func findUid(str string) []string{
	reg := regexp.MustCompile(`\d+`)
	result := reg.FindAllString(str, -1)
	return result
}

func getAllRelsHandler(w http.ResponseWriter, r*http.Request) {
	var Rel pghandler.Relationships
	log.Println(r.URL.Path)
	uid_req := findUid(r.URL.Path)
	if uid_req != nil{
		uid_str := uid_req[0]
		uid, err := strconv.Atoi(uid_str)
		if err != nil {
			log.Println(err)
			os.Exit(2)
		}
		result := Rel.GetALLRelationships(uid)
		w.Write(result)
	} else {
		info := "uid's format is wrong!"
		err := []byte(info)
		w.Write(err)
	}
}

func newRelHandler(w http.ResponseWriter, r*http.Request) {
	var Rel pghandler.Relationships
	uid_req := findUid(r.URL.Path)
	var uid int
	var user_id int
	if uid_req != nil{
		uid_str, user_id_str := uid_req[0], uid_req[1]
		var err error
		uid, err = strconv.Atoi(uid_str)
		if err != nil {
			log.Println(err)
			os.Exit(2)
		}
		user_id, err = strconv.Atoi(user_id_str)
		if err != nil {
			log.Println(err)
			os.Exit(2)
		}
	} else {
		info := "uid's format is wrong!"
		err := []byte(info)
		log.Println(err)
		w.Write(err)
	}
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	log.Println(body_str)
	var state State
	if err := json.Unmarshal(body, &state); err == nil {
		state := state.State
		log.Println(uid, user_id, state)
		result := Rel.NewRelationship(uid, user_id, state, "relationship")
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
