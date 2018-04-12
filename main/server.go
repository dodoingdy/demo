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
	fmt.Println(r.Method)
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		body_str := string(body)
		fmt.Println(body_str)
		var username Name
		if err := json.Unmarshal(body, &username); err == nil {
			fmt.Println("name:", username)
			username := username.Name
            fmt.Println("username", username)
			result := pghandler.UserInsert(username)
			fmt.Println(result)
			w.Write(result)
		} else {
			fmt.Println(err)
		}
	} else {
		result := pghandler.GetAllUsers()
		fmt.Println(result)
		w.Write(result)
	}
}

func find_uid(str string) []string{
	reg := regexp.MustCompile(`\d+`)
	result := reg.FindAllString(str, -1)
	return result
}

func getAllRelasHandler(w http.ResponseWriter, r*http.Request) {
	fmt.Println(r.URL.Path)
	uid_req := find_uid(r.URL.Path)
	if uid_req != nil{
		uid_str := uid_req[0]
		uid, err := strconv.Atoi(uid_str)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println(uid)
		result := pghandler.GetALLRelationships(uid)
		w.Write(result)
	} else {
		info := "uid's format is wrong!"
		err := []byte(info)
		w.Write(err)
	}
}

func newRelaHandler(w http.ResponseWriter, r*http.Request) {
	uid_req := find_uid(r.URL.Path)
	var uid int
	var user_id int
	if uid_req != nil{
		uid_str, user_id_str := uid_req[0], uid_req[1]
		var err error
		uid, err = strconv.Atoi(uid_str)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		user_id, err = strconv.Atoi(user_id_str)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println("first step:", uid, user_id)
	} else {
		info := "uid's format is wrong!"
		err := []byte(info)
		w.Write(err)
	}
	fmt.Println("second step:", uid, user_id)
	body, _ := ioutil.ReadAll(r.Body)
	body_str := string(body)
	fmt.Println(body_str)
	var state State
	if err := json.Unmarshal(body, &state); err == nil {
		fmt.Println("state:", state)
		state := state.State
		fmt.Println(uid, user_id, state)
		result := pghandler.NewRelationship(uid, user_id, state, "relationship")
		w.Write(result)
	} else {
		fmt.Println(err)
	}
}

func main() {
	pghandler.SqlOpen()
	r := mux.NewRouter()
	r.HandleFunc("/users", usershandler)
	r.HandleFunc("/users/{user_id}/relationships", getAllRelasHandler)
	r.HandleFunc("/users/{user_id}/relationships/{other_user_id}", newRelaHandler)
	http.Handle("/", r)
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	//defer pghandler.SqlClose()
}
